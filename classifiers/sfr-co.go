package classifiers

import (
	"fmt"
	"io"
	"regexp"
	"strings"

	"github.com/ofstudio/go-esnsi/cnsi"
)

// SfrCo представляет структуру для работы с классификатором SFR_CO "Клиентские службы СФР".
// https://esnsi.gosuslugi.ru/classifiers/10991/
type SfrCo struct {
	cnsi.SimpleClassifier
	ByOKATO map[string]*SfrCoRecord // Индекс по коду ОКАТО
	Records []SfrCoRecord           // Список записей классификатора
}

// NewSfrCo создает новый экземпляр SfrCo из ридерa.
func NewSfrCo(r io.Reader) (*SfrCo, error) {
	classifier, err := cnsi.NewSimpleClassifier(r)
	if err != nil {
		return nil, fmt.Errorf("failed to create simple clasifier : %v", err)
	}

	sfrCo := &SfrCo{
		SimpleClassifier: *classifier,
	}

	// Инициализируем индексы для быстрого доступа к записям
	sfrCo.ByOKATO = make(map[string]*SfrCoRecord)

	// Разбираем записи классификатора
	if err = sfrCo.Unmarshal(&sfrCo.Records, sfrCo.unmarshalHelper); err != nil {
		return nil, fmt.Errorf("failed to unmarshal records: %v", err)
	}

	return sfrCo, nil
}

// reOKATOPlain - регулярное выражение для проверки корректности кода ОКАТО.
var reOKATOPlain = regexp.MustCompile(`^\d{2,11}$`)

// unmarshalHelper выполняет дополнительную обработку каждой записи при разборе классификатора SFR_CO.
func (c *SfrCo) unmarshalHelper(index int) error {
	record := &c.Records[index]

	// Разбираем OKATOAreaOrig на массив строк
	record.OKATOAreas = strings.Split(record.OKATOAreaOrig, ",")

	// Добавляем запись в индекс по кодам ОКАТО
	for i := range record.OKATOAreas {
		// Проверяем корректность кода ОКАТО
		record.OKATOAreas[i] = strings.TrimSpace(record.OKATOAreas[i])
		if !reOKATOPlain.MatchString(record.OKATOAreas[i]) {
			return fmt.Errorf("invalid OKATO '%s' in record at index %d", record.OKATOAreas[i], index)
		}
		// Проверяем на дубликаты в индексе
		// Если запись с таким OKATO уже существует, возвращаем ошибку
		if _, exists := c.ByOKATO[record.OKATOAreas[i]]; exists {
			return fmt.Errorf("duplicate OKATO '%s' found in record at index %d", record.OKATOAreas[i], index)
		}
		c.ByOKATO[record.OKATOAreas[i]] = record
	}

	return nil
}

// SfrCoRecord представляет запись в классификаторе SFR_CO.
type SfrCoRecord struct {
	COID               int    `esnsi:"COID"`               // Уникальный идентификатор записи в классификаторе. Пример: 1831
	AutoKey            string `esnsi:"autokey"`            // Уникальный идентификатор для записи. Пример: SFR_CO_1831_1
	ToSfrCode          string `esnsi:"ToSfrCode"`          // Код тероргана СФР. Пример: 210
	RegionCode         string `esnsi:"RegionCode"`         // Код региона. Пример: 013
	RegionName         string `esnsi:"RegionName"`         // Наименование региона. Пример: Республика Татарстан
	DivisionCode       string `esnsi:"DivisionCode"`       // Код подразделения. Пример: 401
	OfficeCode         string `esnsi:"OfficeCode"`         // Код офиса. Пример: 401
	OfficeDistrictName string `esnsi:"OfficeDistrictName"` // Наименование района офиса. Пример: "г. Набережные Челны"
	ToSfrName          string `esnsi:"ToSfrName"`          // Наименование тероргана СФР. Пример: "Клиентская служба СФР в г. Набережные Челны (пр-кт Мира)"
	OfficeType         int    `esnsi:"OfficeType"`         // ? Тип офиса. Пример: 2
	Predecessor        int    `esnsi:"Predecessor"`        // ? Идентификатор предшественника. Пример: 1
	Address            string `esnsi:"Address"`            // Адрес офиса. Пример: 423812, г. Набережные Челны, пр-кт Мира, д. 16
	OKATO              string `esnsi:"OKATO"`              // Код ОКАТО клиентской службы. Пример: 92430000000
	OKATOAreaOrig      string `esnsi:"OKATO_Area"`         // ОКАТО обслуживаемой территории. Пример: 92430,92430 (может быть список через запятую)
	OKTMO              string `esnsi:"OKTMO"`              // Код ОКТМО клиентской службы. Пример: 92730000001
	OKTMOAreaOrig      string `esnsi:"OKTMO_Area"`         // ОКТМО обслуживаемой территории. Пример: 92730
	Email              string `esnsi:"Email"`              // Электронная почта
	Phone              string `esnsi:"Phone"`              // Телефон. Пример: 8 (800) 100-00-01
	Latitude           string `esnsi:"Latitude"`           // Широта. Пример: 55.7310453
	Longitude          string `esnsi:"Longitude"`          // Долгота. Пример: 52.3967788
	WorkingTime        string `esnsi:"WorkingTime"`        // Время работы. Пример: "Пн-Пт: с 8:00 до 17:00, Сб: с 8:00 до 15:45, Вс - выходной"
	UTC                string `esnsi:"UTC"`                // Часовой пояс. Пример: "UTC+03:00"
	MSK                int    `esnsi:"MSK"`                // Смещение относительно московского времени в часах. Пример: 0
	TOFSS              string `esnsi:"TOFSS"`              // Код ФСС. Пример: 1600

	OKATOAreas []string // ОКАТО обслуживаемых территорий. Пример: ["92430", "92431"]
}
