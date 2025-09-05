package esnsi

import (
	"encoding/xml"
	"fmt"
	"io"
	"reflect"
)

// Decoder - декодер классификаторов ЕСНСИ из XML формата ЦНСИ
type Decoder[T any] struct {
	r       io.Reader
	handler func(*T) error
}

// NewDecoder - создает новый декодер классификатора для типа записей T.
func NewDecoder[T any](r io.Reader) *Decoder[T] {
	return &Decoder[T]{r: r}
}

// WithHandler - задает обработчик для каждой записи.
// Если обработчик задан, записи не добавляются в слайс Records классификатора,
// а передаются в обработчик, который может их добавить самостоятельно.
// Если обработчик возвращает ошибку, разбор прерывается.
func (d *Decoder[T]) WithHandler(h func(*T) error) *Decoder[T] {
	d.handler = h
	return d
}

// Decode - выполняет разбор XML и возвращает классификатор.
func (d *Decoder[T]) Decode(c *Classifier[T]) error {
	// Проверка, что c не nil
	if c == nil {
		return fmt.Errorf("nil pointer passed")
	}
	// Проверка, что T - структура
	if typeOf := reflect.TypeOf(*new(T)); typeOf.Kind() != reflect.Struct {
		return fmt.Errorf("struct type expected, got %s", typeOf.Kind())
	}

	// Проверяем, что r не nil
	if d.r == nil {
		return fmt.Errorf("reader is nil")
	}

	// Читаем XML формата ЦНСИ
	doc := &cnsiDoc{}
	if err := xml.NewDecoder(d.r).Decode(doc); err != nil {
		return fmt.Errorf("failed to decode XML: %w", err)
	}

	// Разбираем документ и заполняем классификатор
	return d.unmarshal(c, doc)
}

func (d *Decoder[T]) unmarshal(c *Classifier[T], doc *cnsiDoc) error {
	// Создаем индексы атрибутов
	attrNameToRef := make(map[string]string)
	attrRefToKind := make(map[string]reflect.Kind)

	for _, attr := range doc.Meta.StringAttrs {
		attrNameToRef[attr.Name] = attr.UID
		attrRefToKind[attr.UID] = reflect.String
	}
	for _, attr := range doc.Meta.TextAttrs {
		attrNameToRef[attr.Name] = attr.UID
		attrRefToKind[attr.UID] = reflect.String
	}
	for _, attr := range doc.Meta.IntegerAttrs {
		attrNameToRef[attr.Name] = attr.UID
		attrRefToKind[attr.UID] = reflect.Int
	}

	// Создаем индекс полей структуры записи
	fieldRefToIndex := make(map[string]int)
	typeOf := reflect.TypeOf(*new(T))
	// Проходим по всем полям структуры записи
	// и проверяем, что для каждого поля с тегом esnsi
	// существует соответствующий атрибут в классификаторе
	// и что тип поля совпадает с типом атрибута
	for i := 0; i < typeOf.NumField(); i++ {
		field := typeOf.Field(i)

		// Пропускаем поля без тега esnsi
		attrName := field.Tag.Get("esnsi")
		if attrName == "" {
			continue
		}

		// Проверяем, что атрибут существует в классификаторе
		uid, ok := attrNameToRef[attrName]
		if !ok {
			return fmt.Errorf("attribute %s not found in classifier", attrName)
		}
		kind, ok := attrRefToKind[uid]
		if !ok {
			return fmt.Errorf("attribute %s with Ref %s not found in classifier", attrName, uid)
		}

		// Проверяем, что тип поля совпадает с типом атрибута
		if field.Type.Kind() != kind {
			return fmt.Errorf("field %s has type %s, expected %s for attribute %s with Ref %s",
				field.Name, field.Type.Kind(), kind, attrName, uid)
		}
		fieldRefToIndex[uid] = i
	}

	// Заполняем метаданные классификатора
	c.Name = doc.Meta.Name
	c.Code = doc.Meta.Code
	c.UID = doc.Meta.UID
	c.Version = doc.Meta.Version

	// Если обработчик не задан, инициализируем слайс записей
	if d.handler != nil {
		c.Records = make([]T, 0, len(doc.Records))
	}

	// Разбираем записи документа и добавляем их в классификатор
	for i, docRecord := range doc.Records {
		var record T
		val := reflect.ValueOf(&record).Elem()

		// Заполняем поля структуры записи
		for _, attrVal := range docRecord.AttrVals {
			fieldIndex, found := fieldRefToIndex[attrVal.AttrRef]
			if !found {
				continue // Поле не нужно сохранять
			}
			field := val.Field(fieldIndex)

			switch field.Kind() {
			case reflect.String:
				if attrVal.StringVal != nil {
					field.SetString(attrVal.StringVal.Val)
				} else if attrVal.TextVal != nil {
					field.SetString(attrVal.TextVal.Val)
				} else {
					return fmt.Errorf("string value for attribute %s is not set", attrVal.AttrRef)
				}
			case reflect.Int:
				if attrVal.IntegerVal != nil {
					field.SetInt(int64(attrVal.IntegerVal.Val))
				} else {
					return fmt.Errorf("integer value for attribute %s is not set", attrVal.AttrRef)
				}
			default:
				return fmt.Errorf("unsupported field type: %s", field.Kind())
			}
		}

		// Если обработчик не задан, просто добавляем запись в слайс
		if d.handler == nil {
			c.Records = append(c.Records, record)
			continue
		}

		// Иначе вызываем обработчик, при этом запись в слайс не добавляем
		// (она должна быть добавлена обработчиком)
		if err := d.handler(&record); err != nil {
			return fmt.Errorf("handler error at record %d: %w", i, err)
		}
	}

	return nil
}
