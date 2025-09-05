package classifiers

import (
	"fmt"
	"io"
	"regexp"
	"strings"

	"github.com/ofstudio/go-esnsi"
)

// Sfr - классификатор SFR_CO "Клиентские службы СФР".
//
// https://esnsi.gosuslugi.ru/classifiers/10991/
//
// Примечание: один офис может обслуживать несколько территорий обслуживания ОКАТО.
// Также один код ОКАТО может обслуживаться несколькими офисами.
// В этом случае в индекс ByOkato попадает последний из обработанных офисов.
type Sfr struct {
	esnsi.Classifier[SfrRecord]
	ByOkato   map[string]*SfrRecord   // Индекс по коду ОКАТО как в исходном справочнике (например, "92430" или "92430000000")
	ByOkato11 map[string][]*SfrRecord // Индекс по полному коду ОКАТО 11 символов (например, "92430000000")
	ByOkato8  map[string][]*SfrRecord // Индекс по коду ОКАТО 8 символов (например, "92430000")
	ByOkato5  map[string][]*SfrRecord // Индекс по коду ОКАТО 5 символов (например, "92430")
	ByOkato2  map[string][]*SfrRecord // Индекс по коду ОКАТО 2 символа (например, "92")
}

// NewSfr - создает новый классификатор Sfr из XML-данных.
func NewSfr(r io.Reader) (*Sfr, error) {
	byOkato := make(map[string]*SfrRecord)
	byOkato11 := make(map[string][]*SfrRecord)
	byOkato8 := make(map[string][]*SfrRecord)
	byOkato5 := make(map[string][]*SfrRecord)
	byOkato2 := make(map[string][]*SfrRecord)

	var c esnsi.Classifier[SfrRecord]
	if err := esnsi.NewDecoder[SfrRecord](r).WithHandler(func(rec *SfrRecord) error {

		// Разбираем OKATOArea на массив строк
		areas := strings.Split(rec.OKATOArea, ",")

		// Добавляем запись в индексы по кодам ОКАТО
		for _, area := range areas {
			// Проверяем корректность кода ОКАТО
			area = strings.TrimSpace(area)
			// Если пустая строка, пропускаем
			if area == "" {
				continue
			}
			if !reOKATOPlain.MatchString(area) {
				return fmt.Errorf("invalid OKATO '%s'", area)
			}

			// Добавляем в список обслуживаемых территорий
			rec.OKATOAreas = append(rec.OKATOAreas, area)

			// Индекс по коду ОКАТО как в исходном справочнике
			byOkato[area] = rec
			// Индекс по полному код ОКАТО 11 символов
			okato11 := area + strings.Repeat("0", 11-len(area))
			// В индекс по полному коду ОКАТО 11 символов
			// добавляем только если такого кода еще нет
			if _, exists := byOkato11[okato11]; !exists {
				byOkato11[okato11] = append(byOkato11[okato11], rec)
			}
			// Индекс по коду ОКАТО 8 символов
			okato8 := okato11[:8]
			byOkato8[okato8] = append(byOkato8[okato8], rec)
			// Индекс по коду ОКАТО 5 символов
			okato5 := okato11[:5]
			byOkato5[okato5] = append(byOkato5[okato5], rec)
			// Индекс по коду ОКАТО 2 символа
			okato2 := okato11[:2]
			byOkato2[okato2] = append(byOkato2[okato2], rec)
		}

		// Добавляем запись в классификатор
		c.Records = append(c.Records, *rec)

		return nil
	}).Decode(&c); err != nil {
		return nil, fmt.Errorf("error decoding: %v", err)
	}

	return &Sfr{
		Classifier: c,
		ByOkato:    byOkato,
		ByOkato11:  byOkato11,
		ByOkato8:   byOkato8,
		ByOkato5:   byOkato5,
		ByOkato2:   byOkato2,
	}, nil

}

// reOKATOPlain - регулярное выражение для проверки корректности кода ОКАТО.
var reOKATOPlain = regexp.MustCompile(`^\d{2,11}$`)

// SfrRecord - запись в классификаторе Sfr.
type SfrRecord struct {
	// Данные из файла справочника
	COID               int    `esnsi:"COID"`               // Уникальный идентификатор записи в классификаторе. Пример: 1831
	AutoKey            string `esnsi:"autokey"`            // Уникальный идентификатор для записи. Пример: SFR_CO_1831_1
	ToSfrCode          string `esnsi:"ToSfrCode"`          // Код тероргана СФР. Пример: 210
	RegionCode         string `esnsi:"RegionCode"`         // Код региона. Пример: 013
	RegionName         string `esnsi:"RegionName"`         // Наименование региона. Пример: Республика Татарстан
	DivisionCode       string `esnsi:"DivisionCode"`       // Код подразделения. Пример: 401
	OfficeCode         string `esnsi:"OfficeCode"`         // Код офиса. Пример: 401
	OfficeDistrictName string `esnsi:"OfficeDistrictName"` // Наименование района офиса. Пример: "г. Набережные Челны"
	ToSfrName          string `esnsi:"ToSfrName"`          // Наименование тероргана СФР. Пример: "Клиентская служба СФР в г. Набережные Челны (пр-кт Мира)"
	OfficeType         int    `esnsi:"OfficeType"`         // (?) Тип офиса. Пример: 2
	Predecessor        int    `esnsi:"Predecessor"`        // (?) Идентификатор предшественника. Пример: 1
	Address            string `esnsi:"Address"`            // Адрес офиса. Пример: 423812, г. Набережные Челны, пр-кт Мира, д. 16
	OKATO              string `esnsi:"OKATO"`              // Код ОКАТО клиентской службы. Пример: 92430000000
	OKATOArea          string `esnsi:"OKATO_Area"`         // ОКАТО обслуживаемой территории. Пример: 92430,92431 (может быть список через запятую)
	OKTMO              string `esnsi:"OKTMO"`              // Код ОКТМО клиентской службы. Пример: 92730000001
	OKTMOArea          string `esnsi:"OKTMO_Area"`         // ОКТМО обслуживаемой территории. Пример: 92730,92430 (может быть список через запятую)
	Email              string `esnsi:"Email"`              // Электронная почта
	Phone              string `esnsi:"Phone"`              // Телефон. Пример: 8 (800) 100-00-01
	Latitude           string `esnsi:"Latitude"`           // Широта. Пример: 55.7310453
	Longitude          string `esnsi:"Longitude"`          // Долгота. Пример: 52.3967788
	WorkingTime        string `esnsi:"WorkingTime"`        // Время работы. Пример: "Пн-Пт: с 8:00 до 17:00, Сб: с 8:00 до 15:45, Вс - выходной"
	UTC                string `esnsi:"UTC"`                // Часовой пояс. Пример: "UTC+03:00"
	MSK                int    `esnsi:"MSK"`                // Смещение относительно московского времени в часах. Пример: 0
	TOFSS              string `esnsi:"TOFSS"`              // Код ФСС. Пример: 1600

	// Данные, полученные при разборе
	OKATOAreas []string // Список ОКАТО обслуживаемых территорий. Пример: ["92430", "92431"]
}
