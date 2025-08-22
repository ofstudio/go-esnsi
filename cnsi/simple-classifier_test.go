package cnsi

import (
	"bytes"
	"errors"
	"os"
	"reflect"
	"strings"
	"testing"
)

// TestNewSimpleClassifier tests the NewSimpleClassifier function
func TestNewSimpleClassifier(t *testing.T) {
	// Test successful case
	t.Run("Valid XML file", func(t *testing.T) {
		r, err := os.Open("testdata/SFR_CO_test.xml")
		if err != nil {
			t.Fatalf("Failed to open test file: %v", err)
		}
		//goland:noinspection ALL
		defer r.Close()

		classifier, err := NewSimpleClassifier(r)
		if err != nil {
			t.Fatalf("Failed to create SimpleClassifier: %v", err)
		}

		if classifier == nil {
			t.Fatal("Expected non-nil classifier")
		}
	})

	t.Run("Invalid XML", func(t *testing.T) {
		invalidXML := bytes.NewBufferString("<invalid>xml</")
		_, err := NewSimpleClassifier(invalidXML)
		if err == nil {
			t.Fatal("Expected error for invalid XML, got nil")
		}
	})
}

// TestSimpleClassifier_Unmarshal tests the Unmarshal method
func TestSimpleClassifier_Unmarshal(t *testing.T) {
	// Setup classifier from test file
	r, err := os.Open("testdata/SFR_CO_test.xml")
	if err != nil {
		t.Fatalf("Failed to open test file: %v", err)
	}
	//goland:noinspection ALL
	defer r.Close()

	classifier, err := NewSimpleClassifier(r)
	if err != nil {
		t.Fatalf("Failed to create SimpleClassifier: %v", err)
	}

	// Test successful unmarshalling
	t.Run("Successful unmarshalling", func(t *testing.T) {
		var items []Item
		if err := classifier.Unmarshal(&items, nil); err != nil {
			t.Fatalf("Failed to unmarshal: %v", err)
		}

		if len(items) != 4 {
			t.Fatalf("Expected 4 items, got %d", len(items))
		}

		// Validate first item's fields
		expected := Item{
			COID:               1831,
			AutoKey:            "SFR_CO_1831_1",
			ToSfrCode:          "210",
			RegionCode:         "013",
			RegionName:         "Республика Татарстан",
			DivisionCode:       "401",
			OfficeCode:         "401",
			OfficeDistrictName: "г. Набережные Челны",
			ToSfrName:          "Клиентская служба СФР в г. Набережные Челны (пр-кт Мира)",
			OfficeType:         2,
			Predecessor:        1,
			Address:            "423812, Республика Татарстан, г. Набережные Челны, пр-кт Мира, д 24Е (7/20)",
			OKATO:              "92430000000",
			OKATOArea:          "92430",
			OKTMO:              "92730000001",
			OKTMOArea:          "92730",
			Email:              "tatarstan@16.sfr.gov.ru",
			Phone:              "8 (800) 100-00-01",
			Latitude:           "55.7310453",
			Longitude:          "52.3967788",
			WorkingTime:        "Пн-Чт: с 8:00 до 17:00, Пт: с 8:00 до 15:45, без перерыва",
			UTC:                "UTC+03:00",
			MSK:                0,
			TOFSS:              "1600",
		}

		if !reflect.DeepEqual(items[0], expected) {
			t.Fatalf("Item data mismatch.\nExpected: %+v\nGot: %+v", expected, items[0])
		}
	})

	// Test with callback function
	t.Run("Unmarshal with callback", func(t *testing.T) {
		var items []Item
		callCount := 0

		err := classifier.Unmarshal(&items, func(index int) error {
			callCount++
			return nil
		})

		if err != nil {
			t.Fatalf("Failed to unmarshal with callback: %v", err)
		}

		if callCount != 4 {
			t.Fatalf("Expected callback to be called 4 times, got %d", callCount)
		}
	})

	// Test error cases
	t.Run("Non-pointer destination", func(t *testing.T) {
		var items []Item
		err := classifier.Unmarshal(items, nil)
		if err == nil {
			t.Fatal("Expected error for non-pointer destination, got nil")
		}
	})

	t.Run("Non-slice destination", func(t *testing.T) {
		var item Item
		err := classifier.Unmarshal(&item, nil)
		if err == nil {
			t.Fatal("Expected error for non-slice destination, got nil")
		}
	})

	t.Run("Non-struct slice elements", func(t *testing.T) {
		var items []string
		err := classifier.Unmarshal(&items, nil)
		if err == nil {
			t.Fatal("Expected error for non-struct slice elements, got nil")
		}
	})

	t.Run("Callback returning error", func(t *testing.T) {
		var items []Item
		expectedErr := errors.New("expected error")

		err := classifier.Unmarshal(&items, func(index int) error {
			return expectedErr
		})

		if !errors.Is(err, expectedErr) {
			t.Fatal("Expected error from callback, got nil")
		}
	})
}

// TestStructureWithInvalidFields tests unmarshaling to a struct with invalid field tags
func TestStructureWithInvalidFields(t *testing.T) {
	r, err := os.Open("testdata/SFR_CO_test.xml")
	if err != nil {
		t.Fatalf("Failed to open test file: %v", err)
	}
	//goland:noinspection ALL
	defer r.Close()

	classifier, err := NewSimpleClassifier(r)
	if err != nil {
		t.Fatalf("Failed to create SimpleClassifier: %v", err)
	}

	t.Run("Unknown attribute", func(t *testing.T) {
		type InvalidItem struct {
			COID    int    `esnsi:"COID"`
			Invalid string `esnsi:"NonExistentAttribute"` // This attribute doesn't exist
		}

		var items []InvalidItem
		err := classifier.Unmarshal(&items, nil)
		if err == nil || !strings.Contains(err.Error(), "not found in classifier") {
			t.Fatalf("Expected 'not found in classifier' error, got: %v", err)
		}
	})

	t.Run("Type mismatch", func(t *testing.T) {
		type InvalidTypeItem struct {
			COID string `esnsi:"COID"` // Should be int, not string
		}

		var items []InvalidTypeItem
		err := classifier.Unmarshal(&items, nil)
		if err == nil || !strings.Contains(err.Error(), "expected") {
			t.Fatalf("Expected type mismatch error, got: %v", err)
		}
	})
}

// TestSimpleClassifierProperties tests the properties and metadata
func TestSimpleClassifierProperties(t *testing.T) {
	r, err := os.Open("testdata/SFR_CO_test.xml")
	if err != nil {
		t.Fatalf("Failed to open test file: %v", err)
	}
	//goland:noinspection ALL
	defer r.Close()

	classifier, err := NewSimpleClassifier(r)
	if err != nil {
		t.Fatalf("Failed to create SimpleClassifier: %v", err)
	}

	// Test metadata values
	if classifier.Properties.Meta.Code != "SFR_CO" {
		t.Errorf("Expected code 'SFR_CO', got '%s'", classifier.Properties.Meta.Code)
	}

	if classifier.Properties.Meta.Name == "" {
		t.Error("Expected non-empty name")
	}

	// Check that attributes are properly loaded
	attributeCount := len(classifier.Properties.Attributes.IntegerAttributes) +
		len(classifier.Properties.Attributes.StringAttributes) +
		len(classifier.Properties.Attributes.TextAttributes)

	if attributeCount == 0 {
		t.Error("No attributes found in classifier")
	}
}

type Item struct {
	COID               int    `esnsi:"COID"`               // Уникальный идентификатор записи в классификаторе. Пример: 1831
	AutoKey            string `esnsi:"autokey"`            // Уникальный идентификатор для записи. Пример: SFR_CO_1831_1
	ToSfrCode          string `esnsi:"ToSfrCode"`          // Код тероргана СФР. Пример: 210
	RegionCode         string `esnsi:"RegionCode"`         // Код региона. Пример: 013
	RegionName         string `esnsi:"RegionName"`         // Наименование региона. Пример: Республика Татарстан
	DivisionCode       string `esnsi:"DivisionCode"`       // Код подразделения (необязательный). Пример: 401
	OfficeCode         string `esnsi:"OfficeCode"`         // Код офиса. Пример: 401
	OfficeDistrictName string `esnsi:"OfficeDistrictName"` // Наименование района офиса (необязательный). Пример: "г. Набережные Челны"
	ToSfrName          string `esnsi:"ToSfrName"`          // Наименование тероргана СФР (необязательный). Пример: "Клиентская служба СФР в г. Набережные Челны (пр-кт Мира)"
	OfficeType         int    `esnsi:"OfficeType"`         // ? Тип офиса. Пример: 2 (Клиентская служба)
	Predecessor        int    `esnsi:"Predecessor"`        // ? Идентификатор предшественника. Пример: 1
	Address            string `esnsi:"Address"`            // Адрес офиса. Пример: 423812, г. Набережные Челны, пр-кт Мира, д. 16
	OKATO              string `esnsi:"OKATO"`              // Код ОКАТО клиентской службы. Пример: 92430000000
	OKATOArea          string `esnsi:"OKATO_Area"`         // ОКАТО обслуживаемой территории (необязательный). Пример: 92430
	OKTMO              string `esnsi:"OKTMO"`              // Код ОКТМО клиентской службы. Пример: 92730000001
	OKTMOArea          string `esnsi:"OKTMO_Area"`         // ОКТМО обслуживаемой территории (необязательный). Пример: 92730
	Email              string `esnsi:"Email"`              // Электронная почта
	Phone              string `esnsi:"Phone"`              // Телефон. Пример: 8 (800) 100-00-01
	Latitude           string `esnsi:"Latitude"`           // Широта. Пример: 55.7310453
	Longitude          string `esnsi:"Longitude"`          // Долгота. Пример: 52.3967788
	WorkingTime        string `esnsi:"WorkingTime"`        // Время работы. Пример: "Пн-Пт: с 8:00 до 17:00, Сб: с 8:00 до 15:45, Вс - выходной"
	UTC                string `esnsi:"UTC"`                // Часовой пояс. Пример: "UTC+03:00"
	MSK                int    `esnsi:"MSK"`                // Смещение относительно московского времени в часах. Пример: 0
	TOFSS              string `esnsi:"TOFSS"`              // Код ФСС (необязательный). Пример: 1600
}
