package classifiers

import (
	"os"
	"strings"
	"testing"
)

//goland:noinspection GoUnhandledErrorResult
func TestNewSfr(t *testing.T) {
	t.Run("reader is nil", func(t *testing.T) {
		_, err := NewSfr(nil)
		if err == nil {
			t.Error("expected error, got nil")
		} else if !strings.Contains(err.Error(), "reader is nil") {
			t.Errorf("unexpected error: %v", err)
		}
	})

	t.Run("valid data", func(t *testing.T) {
		f, err := os.Open("../testdata/sfr-valid_test.xml")
		if err != nil {
			t.Fatalf("failed to open test file: %v", err)
		}
		defer f.Close()

		sfr, err := NewSfr(f)
		if err != nil {
			t.Fatalf("failed to create SFR classifier: %v", err)
		}

		// Проверяем метаданные
		if sfr.Name != "Клиентские службы СФР" {
			t.Errorf("unexpected Name: %s", sfr.Name)
		}
		if sfr.Code != "SFR_CO" {
			t.Errorf("unexpected Code: %s", sfr.Code)
		}
		if sfr.UID != "2fad55ce-854c-4044-873e-7ae806cc94cb" {
			t.Errorf("unexpected UID: %s", sfr.UID)
		}
		if sfr.Version != 55 {
			t.Errorf("unexpected Version: %d", sfr.Version)
		}

		// Проверяем записи
		if len(sfr.Records) != 2 {
			t.Fatalf("unexpected number of records: %d", len(sfr.Records))
		}

		// Проверяем первую запись
		record0 := sfr.Records[0]
		if record0.COID != 1831 {
			t.Errorf("record 0: unexpected COID: %d, expected 1831", record0.COID)
		}
		if record0.ToSfrCode != "210" {
			t.Errorf("record 0: unexpected ToSfrCode: %s, expected 210", record0.ToSfrCode)
		}
		if record0.RegionName != "Республика Татарстан" {
			t.Errorf("record 0: unexpected RegionName: %s, expected Республика Татарстан", record0.RegionName)
		}
		if record0.OKATO != "92430000000" {
			t.Errorf("record 0: unexpected OKATO: %s, expected 92430000000", record0.OKATO)
		}
		if record0.OKATOArea != "92430, 92432" {
			t.Errorf("record 0: unexpected OKATOArea: %s, expected '92430, 92432'", record0.OKATOArea)
		}

		// Проверяем разбор OKATOAreas
		expectedAreas0 := []string{"92430", "92432"}
		if len(record0.OKATOAreas) != len(expectedAreas0) {
			t.Fatalf("record 0: unexpected OKATOAreas length: %d, expected %d", len(record0.OKATOAreas), len(expectedAreas0))
		}
		for i, area := range expectedAreas0 {
			if record0.OKATOAreas[i] != area {
				t.Errorf("record 0: OKATOAreas[%d]: %s, expected %s", i, record0.OKATOAreas[i], area)
			}
		}

		// Проверяем вторую запись
		record1 := sfr.Records[1]
		if record1.COID != 949 {
			t.Errorf("record 1: unexpected COID: %d, expected 949", record1.COID)
		}
		if record1.ToSfrCode != "201" {
			t.Errorf("record 1: unexpected ToSfrCode: %s, expected 201", record1.ToSfrCode)
		}
		if record1.RegionName != "Москва" {
			t.Errorf("record 1: unexpected RegionName: %s, expected Москва", record1.RegionName)
		}
		if record1.OKATO != "45277592000" {
			t.Errorf("record 1: unexpected OKATO: %s, expected 45277592000", record1.OKATO)
		}

		// Проверяем разбор OKATOAreas для второй записи
		expectedAreas1 := []string{"45277592"}
		if len(record1.OKATOAreas) != len(expectedAreas1) {
			t.Fatalf("record 1: unexpected OKATOAreas length: %d, expected %d", len(record1.OKATOAreas), len(expectedAreas1))
		}
		for i, area := range expectedAreas1 {
			if record1.OKATOAreas[i] != area {
				t.Errorf("record 1: OKATOAreas[%d]: %s, expected %s", i, record1.OKATOAreas[i], area)
			}
		}
	})

	t.Run("check indexes", func(t *testing.T) {
		f, err := os.Open("../testdata/sfr-valid_test.xml")
		if err != nil {
			t.Fatalf("failed to open test file: %v", err)
		}
		defer f.Close()

		sfr, err := NewSfr(f)
		if err != nil {
			t.Fatalf("failed to create SFR classifier: %v", err)
		}

		// Проверяем индекс ByOkato
		if len(sfr.ByOkato) != 3 {
			t.Errorf("unexpected ByOkato length: %d, expected 3", len(sfr.ByOkato))
		}

		// Проверяем что записи есть в индексах по оригинальным кодам
		if _, exists := sfr.ByOkato["92430"]; !exists {
			t.Error("record with OKATO '92430' not found in ByOkato index")
		}
		if _, exists := sfr.ByOkato["92432"]; !exists {
			t.Error("record with OKATO '92432' not found in ByOkato index")
		}
		if _, exists := sfr.ByOkato["45277592"]; !exists {
			t.Error("record with OKATO '45277592' not found in ByOkato index")
		}

		// Проверяем индекс ByOkato11
		if len(sfr.ByOkato11) != 3 {
			t.Errorf("unexpected ByOkato11 length: %d, expected 3", len(sfr.ByOkato11))
		}

		if _, exists := sfr.ByOkato11["92430000000"]; !exists {
			t.Error("record with OKATO11 '92430000000' not found in ByOkato11 index")
		}
		if _, exists := sfr.ByOkato11["92432000000"]; !exists {
			t.Error("record with OKATO11 '92432000000' not found in ByOkato11 index")
		}
		if _, exists := sfr.ByOkato11["45277592000"]; !exists {
			t.Error("record with OKATO11 '45277592000' not found in ByOkato11 index")
		}

		// Проверяем индекс ByOkato8 - теперь ожидаем 3 записи
		if len(sfr.ByOkato8) != 3 {
			t.Errorf("unexpected ByOkato8 length: %d, expected 3", len(sfr.ByOkato8))
		}

		records8_92430 := sfr.ByOkato8["92430000"]
		if len(records8_92430) != 1 {
			t.Errorf("unexpected number of records for OKATO8 '92430000': %d, expected 1", len(records8_92430))
		}

		records8_92432 := sfr.ByOkato8["92432000"]
		if len(records8_92432) != 1 {
			t.Errorf("unexpected number of records for OKATO8 '92432000': %d, expected 1", len(records8_92432))
		}

		records8_45277 := sfr.ByOkato8["45277592"]
		if len(records8_45277) != 1 {
			t.Errorf("unexpected number of records for OKATO8 '45277592': %d, expected 1", len(records8_45277))
		}

		// Проверяем индекс ByOkato5 - теперь ожидаем 3 записи
		if len(sfr.ByOkato5) != 3 {
			t.Errorf("unexpected ByOkato5 length: %d, expected 3", len(sfr.ByOkato5))
		}

		records5_92430 := sfr.ByOkato5["92430"]
		if len(records5_92430) != 1 {
			t.Errorf("unexpected number of records for OKATO5 '92430': %d, expected 1", len(records5_92430))
		}

		records5_92432 := sfr.ByOkato5["92432"]
		if len(records5_92432) != 1 {
			t.Errorf("unexpected number of records for OKATO5 '92432': %d, expected 1", len(records5_92432))
		}

		records5_45277 := sfr.ByOkato5["45277"]
		if len(records5_45277) != 1 {
			t.Errorf("unexpected number of records for OKATO5 '45277': %d, expected 1", len(records5_45277))
		}

		// Проверяем индекс ByOkato2 - у нас 2 различных региона: 92 и 45
		if len(sfr.ByOkato2) != 2 {
			t.Errorf("unexpected ByOkato2 length: %d, expected 2", len(sfr.ByOkato2))
		}

		records2_92 := sfr.ByOkato2["92"]
		if len(records2_92) != 2 {
			t.Errorf("unexpected number of records for OKATO2 '92': %d, expected 2", len(records2_92))
		}

		records2_45 := sfr.ByOkato2["45"]
		if len(records2_45) != 1 {
			t.Errorf("unexpected number of records for OKATO2 '45': %d, expected 1", len(records2_45))
		}
	})

	t.Run("invalid OKATO area", func(t *testing.T) {
		f, err := os.Open("../testdata/sfr-okato-area-invalid_test.xml")
		if err != nil {
			t.Fatalf("failed to open test file: %v", err)
		}
		defer f.Close()

		_, err = NewSfr(f)
		if err == nil {
			t.Error("expected error, got nil")
		} else if !strings.Contains(err.Error(), "invalid OKATO") {
			t.Errorf("unexpected error: %v", err)
		}
	})

	t.Run("empty OKATO area", func(t *testing.T) {
		f, err := os.Open("../testdata/sfr-okato-area-empty_test.xml")
		if err != nil {
			t.Fatalf("failed to open test file: %v", err)
		}
		defer f.Close()

		sfr, err := NewSfr(f)
		if err != nil {
			t.Fatalf("failed to create SFR classifier: %v", err)
		}

		// Проверяем, что записи создались корректно
		if len(sfr.Records) != 2 {
			t.Fatalf("unexpected number of records: %d, expected 2", len(sfr.Records))
		}

		// Проверяем первую запись - у неё есть валидный код "123456" и пустое значение после запятой
		record0 := sfr.Records[0]
		if record0.OKATOArea != "123456," {
			t.Errorf("record 0: unexpected OKATOArea: %s, expected '123456,'", record0.OKATOArea)
		}
		// Проверяем, что в OKATOAreas попал только валидный код, пустое значение игнорируется
		expectedAreas0 := []string{"123456"}
		if len(record0.OKATOAreas) != len(expectedAreas0) {
			t.Fatalf("record 0: unexpected OKATOAreas length: %d, expected %d", len(record0.OKATOAreas), len(expectedAreas0))
		}
		for i, area := range expectedAreas0 {
			if record0.OKATOAreas[i] != area {
				t.Errorf("record 0: OKATOAreas[%d]: %s, expected %s", i, record0.OKATOAreas[i], area)
			}
		}

		// Проверяем вторую запись - у неё только пробел в OKATOArea (пустое значение)
		record1 := sfr.Records[1]
		if record1.OKATOArea != " " {
			t.Errorf("record 1: unexpected OKATOArea: %s, expected ' '", record1.OKATOArea)
		}
		// Проверяем, что OKATOAreas пустой для второй записи (пробел игнорируется)
		if len(record1.OKATOAreas) != 0 {
			t.Fatalf("record 1: unexpected OKATOAreas length: %d, expected 0", len(record1.OKATOAreas))
		}

		// Проверяем индексы - должен быть только один валидный код "123456"
		if len(sfr.ByOkato) != 1 {
			t.Errorf("unexpected ByOkato length: %d, expected 1", len(sfr.ByOkato))
		}

		if _, exists := sfr.ByOkato["123456"]; !exists {
			t.Error("record with OKATO '123456' not found in ByOkato index")
		}

		// Проверяем индекс ByOkato11 - должен содержать только один полный код
		if len(sfr.ByOkato11) != 1 {
			t.Errorf("unexpected ByOkato11 length: %d, expected 1", len(sfr.ByOkato11))
		}

		if _, exists := sfr.ByOkato11["12345600000"]; !exists {
			t.Error("record with OKATO11 '12345600000' not found in ByOkato11 index")
		}

		// Проверяем индекс ByOkato8 - должен содержать только один 8-значный код
		records11_123456 := sfr.ByOkato11["12345600000"]
		if len(records11_123456) != 1 {
			t.Errorf("unexpected number of records for OKATO11 '12345600000': %d, expected 1", len(records11_123456))
		}

		records8_123456 := sfr.ByOkato8["12345600"]
		if len(records8_123456) != 1 {
			t.Errorf("unexpected number of records for OKATO8 '12345600': %d, expected 1", len(records8_123456))
		}

		// Проверяем индекс ByOkato5 - должен содержать только один 5-значный код
		if len(sfr.ByOkato5) != 1 {
			t.Errorf("unexpected ByOkato5 length: %d, expected 1", len(sfr.ByOkato5))
		}

		records5_12345 := sfr.ByOkato5["12345"]
		if len(records5_12345) != 1 {
			t.Errorf("unexpected number of records for OKATO5 '12345': %d, expected 1", len(records5_12345))
		}

		// Проверяем индекс ByOkato2 - должен содержать только один 2-значный код
		if len(sfr.ByOkato2) != 1 {
			t.Errorf("unexpected ByOkato2 length: %d, expected 1", len(sfr.ByOkato2))
		}

		records2_12 := sfr.ByOkato2["12"]
		if len(records2_12) != 1 {
			t.Errorf("unexpected number of records for OKATO2 '12': %d, expected 1", len(records2_12))
		}
	})
}
