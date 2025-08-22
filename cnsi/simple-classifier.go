package cnsi

import (
	"encoding/xml"
	"fmt"
	"io"
	"reflect"
)

// SimpleClassifier представляет собой простой XML-классификатор формата ЦНСИ.
type SimpleClassifier struct {
	XMLName    xml.Name                   `xml:"document"`
	Properties SimpleClassifierProperties `xml:"simple-classifier"`
	Records    []Record                   `xml:"data>record"`
	PFRF       struct {
		Records []Record `xml:"record"`
	}
}

// SimpleClassifierProperties содержит свойства классификатора:
// метаданные и описание полей.
type SimpleClassifierProperties struct {
	Meta
	Attributes
}

// NewSimpleClassifier создает новый экземпляр SimpleClassifier из ридерa.
func NewSimpleClassifier(r io.Reader) (*SimpleClassifier, error) {
	decoder := xml.NewDecoder(r)
	var classifier SimpleClassifier
	if err := decoder.Decode(&classifier); err != nil {
		return nil, fmt.Errorf("error decoding xml: %v", err)
	}
	return &classifier, nil
}

// Unmarshal распаковывает записи из SimpleClassifier в структуры данных.
// dst должен быть указателем на слайс структур, которые соответствуют атрибутам классификатора.
// Если callback не nil, он будет вызван для каждого индекса после добавления элемента в
// dst, что позволяет выполнять дополнительные действия или обработку.
func (c *SimpleClassifier) Unmarshal(dst interface{}, callback func(index int) error) error {
	// Проверяем, что dst является указателем на слайс структур
	dstVal := reflect.ValueOf(dst)
	if dstVal.Kind() != reflect.Ptr ||
		dstVal.Elem().Kind() != reflect.Slice ||
		dstVal.Elem().Type().Elem().Kind() != reflect.Struct {
		return fmt.Errorf("destination must be a pointer to a slice of structs, got %s", dstVal.Kind())
	}
	structElem := dstVal.Elem().Type().Elem()

	// Создаем индексы для атрибутов
	attrNameToUID, attrUIDToKind := c.Properties.indexes()

	// Создаем индекс UID->индекс полей структуры
	fieldUIDToIndex := make(map[string]int)
	for i := 0; i < structElem.NumField(); i++ {
		field := structElem.Field(i)
		attrName := field.Tag.Get("esnsi")
		if attrName == "" {
			continue
		}

		uid, ok := attrNameToUID[attrName]
		if !ok {
			return fmt.Errorf("attribute %s not found in classifier", attrName)
		}

		kind, ok := attrUIDToKind[uid]
		if !ok {
			return fmt.Errorf("attribute %s with UID %s not found in classifier", attrName, uid)
		}

		if field.Type.Kind() != kind {
			return fmt.Errorf("field %s has type %s, expected %s for attribute %s with UID %s",
				field.Name, field.Type.Kind(), kind, attrName, uid)
		}

		fieldUIDToIndex[uid] = i
	}

	// Обходим записи в классификаторе и заполняем структуры.
	for i, record := range c.Records {

		// Создаем новый элемент структуры
		newElem := reflect.New(structElem).Elem()

		// Заполняем поля структуры значениями атрибутов из записи
		for _, attrValue := range record.AttributeValues {
			fieldIndex, exists := fieldUIDToIndex[attrValue.AttributeRef]
			if !exists {
				continue // Пропускаем атрибуты, которые не соответствуют полям целевой структуры
			}

			field := newElem.Field(fieldIndex)

			switch field.Kind() {
			case reflect.Int:
				if attrValue.IntegerValue == nil {
					return fmt.Errorf("integer value for attribute %s is nil", attrValue.AttributeRef)
				}
				field.SetInt(int64(attrValue.IntegerValue.Value))

			case reflect.String:
				if attrValue.TextValue == nil && attrValue.StringValue == nil {
					return fmt.Errorf("text or string values for attribute %s is nil", attrValue.AttributeRef)
				}
				if attrValue.TextValue != nil {
					field.SetString(attrValue.TextValue.Value)
				} else {
					field.SetString(attrValue.StringValue.Value)
				}

			default:
				return fmt.Errorf("unsupported field type %s for attribute %s", field.Kind(), attrValue.AttributeRef)
			}
		}

		// Добавляем новый элемент в слайс
		dstVal.Elem().Set(reflect.Append(dstVal.Elem(), newElem))

		// Вызываем callback, если он задан
		if callback != nil {
			if err := callback(i); err != nil {
				return fmt.Errorf("callback error at index %d: %w", i, err)
			}
		}
	}

	return nil
}
