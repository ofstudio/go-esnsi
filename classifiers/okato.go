package classifiers

import (
	"fmt"
	"io"
	"regexp"
	"strings"

	"github.com/ofstudio/go-esnsi"
)

// Okato - классификатор "Общероссийский классификатор объектов
// административно-территориального деления (ОКАТО)".
//
// https://esnsi.gosuslugi.ru/classifiers/16270
type Okato struct {
	esnsi.Classifier[OkatoRecord]
	byCode   map[string]*OkatoRecord   // Индекс по коду ОКАТО
	byCode11 map[string]*OkatoRecord   // Индекс по полному коду ОКАТО 11 символов
	byRegion map[string][]*OkatoRecord // Индекс по региону (первые 2 символа кода ОКАТО)
	Region   map[string]*OkatoRecord   // Регионы
}

// NewOkato - создает новый классификатор Okato из XML-данных.
//
// Примечание: по состоянию на сентябрь 2025 в классификаторе
// содержатся невалидные коды ОКАТО, например ";;classifierOkato_75.249.550".
// Поэтому при разборе записи с невалидным кодом игнорируются и не добавляются.
func NewOkato(r io.Reader) (*Okato, error) {
	byCode := make(map[string]*OkatoRecord)
	byCode11 := make(map[string]*OkatoRecord)
	byRegion := make(map[string][]*OkatoRecord)
	Region := make(map[string]*OkatoRecord)

	var c esnsi.Classifier[OkatoRecord]
	if err := esnsi.NewDecoder[OkatoRecord](r).WithHandler(func(rec *OkatoRecord) error {
		// Разбираем код ОКАТО
		if err := rec.ParseCode(); err != nil {
			// Если код невалидный, пропускаем запись
			return nil
		}
		// Добавляем запись в индексы по коду ОКАТО
		if _, exists := byCode[rec.Code]; exists {
			return fmt.Errorf("duplicate OKATO code '%s'", rec.Code)
		}
		byCode[rec.Code] = rec
		byCode11[rec.Code11] = rec
		byRegion[rec.Region] = append(byRegion[rec.Region], rec)

		// Добавляем запись в список регионов, если это регион (код из 2 символов)
		if len(rec.Code) == 2 {
			Region[rec.Code] = rec
		}

		// Добавляем запись в классификатор
		c.Records = append(c.Records, *rec)
		return nil
	}).Decode(&c); err != nil {
		return nil, fmt.Errorf("error decoding: %w", err)
	}

	return &Okato{
		Classifier: c,
		byCode:     byCode,
		byCode11:   byCode11,
		byRegion:   byRegion,
		Region:     Region,
	}, nil
}

// reValidOkatoC - регулярное выражение для проверки корректности кода ОКАТО
// в формате с точками (например, "01.201.800" или "01.201.800.001").
var reValidOkatoC = regexp.MustCompile(`^\d{2}(\.\d{3}(\.\d{3}(\.\d{3})?)?)?$`)

// OkatoRecord - запись классификатора Okato.
type OkatoRecord struct {
	// Данные из файла справочника
	C              string `esnsi:"Код"`                   // Код, например "01.201.800"
	K4             string `esnsi:"КЧ"`                    // Неизвестно (контрольное число?). Пример: "1"
	Name           string `esnsi:"Наименование"`          // Наименование территориального объекта. Пример: "Сельсоветы Алейского р-на"
	AdditionalData string `esnsi:"Дополнительные данные"` // Дополнительные сведения. Как правило, название центрального населенного пункта. Пример: "с Толстая Дуброва"

	// Данные, полученные при разборе
	Code   string // Код ОКАТО, например "0120188"
	Code11 string // Код ОКАТО, полный 11-символьный, например "01201800000"
	Region string // Регион (первые два символа кода ОКАТО), например "01"
	Level1 string // Уровень 1: район/город (символы с 3 по 5 кода ОКАТО), например "201"
	Level2 string // Уровень 2: рабочий поселок/сельсовет (символы с 6 по 8 кода ОКАТО), например "800"
	Level3 string // Уровень 3: населенный пункт (символы с 9 по 11 кода ОКАТО), например "001"
}

// ParseCode - разбирает и валидирует поле C (код ОКАТО с точками),
// заполняет поля Code, Code11, Region, Level1, Level2, Level3.
func (r *OkatoRecord) ParseCode() error {
	// Проверяем корректность кода ОКАТО
	if !reValidOkatoC.MatchString(r.C) {
		return fmt.Errorf("invalid OKATO code '%s'", r.C)
	}

	// Убираем точки из кода
	r.Code = ""
	for _, ch := range r.C {
		if ch != '.' {
			r.Code += string(ch)
		}
	}

	// Заполняем поле Code11 (полный 11-символьный код)
	r.Code11 = r.Code + strings.Repeat("0", 11-len(r.Code))

	// Заполняем поля Region, Level1, Level2, Level3
	if len(r.Code) >= 2 {
		r.Region = r.Code[:2]
	}
	if len(r.Code) >= 5 {
		r.Level1 = r.Code[2:5]
	}
	if len(r.Code) >= 8 {
		r.Level2 = r.Code[5:8]
	}
	if len(r.Code) == 11 {
		r.Level3 = r.Code[8:11]
	}

	return nil
}
