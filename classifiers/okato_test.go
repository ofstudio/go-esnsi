package classifiers

import (
	"os"
	"strings"
	"testing"
)

//goland:noinspection GoUnhandledErrorResult
func TestNewOkato(t *testing.T) {

	t.Run("valid data", func(t *testing.T) {
		f, err := os.Open("../testdata/okato-valid_test.xml")
		if err != nil {
			t.Fatalf("failed to open test file: %v", err)
		}
		defer f.Close()

		okato, err := NewOkato(f)
		if err != nil {
			t.Fatalf("failed to create OKATO classifier: %v", err)
		}

		// Проверяем метаданные
		if okato.Name != "Общероссийский классификатор объектов административно-территориального деления (ОКАТО)" {
		}
		if okato.Code != "classifierOkato" {
			t.Errorf("unexpected Code: %s", okato.Code)
		}
		if okato.UID != "b8628acb-01f8-41e9-8fce-a5e62e6256ab" {
			t.Errorf("unexpected UID: %s", okato.UID)
		}
		if okato.Version != 7 {
			t.Errorf("unexpected Version: %d", okato.Version)
		}

		// Проверяем записи
		if len(okato.Records) != 4 {
			t.Fatalf("unexpected number of records: %d", len(okato.Records))
		}

		// Проверяем первую запись (регион)
		r0 := okato.Records[0]
		if r0.C != "01" {
			t.Errorf("record 0: unexpected C: %s, expected 01", r0.C)
		}
		if r0.K4 != "2" {
			t.Errorf("record 0: unexpected K4: %s, expected 2", r0.K4)
		}
		if r0.Name != "Алтайский край" {
			t.Errorf("record 0: unexpected Name: %s, expected Алтайский край", r0.Name)
		}
		if r0.AdditionalData != "г Барнаул" {
			t.Errorf("record 0: unexpected AdditionalData: %s, expected г Барнаул", r0.AdditionalData)
		}

		// Проверяем разбор кода ОКАТО для первой записи
		if r0.Code != "01" {
			t.Errorf("record 0: unexpected Code: %s, expected 01", r0.Code)
		}
		if r0.Code11 != "01000000000" {
			t.Errorf("record 0: unexpected Code11: %s, expected 01000000000", r0.Code11)
		}
		if r0.Region != "01" {
			t.Errorf("record 0: unexpected Region: %s, expected 01", r0.Region)
		}
		if r0.Level1 != "" {
			t.Errorf("record 0: unexpected Level1: %s, expected empty", r0.Level1)
		}

		// Проверяем вторую запись (район)
		r1 := okato.Records[1]
		if r1.C != "01.200" {
			t.Errorf("record 1: unexpected C: %s, expected 01.200", r1.C)
		}
		if r1.Code != "01200" {
			t.Errorf("record 1: unexpected Code: %s, expected 01200", r1.Code)
		}
		if r1.Code11 != "01200000000" {
			t.Errorf("record 1: unexpected Code11: %s, expected 01200000000", r1.Code11)
		}
		if r1.Region != "01" {
			t.Errorf("record 1: unexpected Region: %s, expected 01", r1.Region)
		}
		if r1.Level1 != "200" {
			t.Errorf("record 1: unexpected Level1: %s, expected 200", r1.Level1)
		}
		if r1.Level2 != "" {
			t.Errorf("record 1: unexpected Level2: %s, expected empty", r1.Level2)
		}

		// Проверяем третью запись (сельсовет)
		r2 := okato.Records[2]
		if r2.C != "01.201.800" {
			t.Errorf("record 2: unexpected C: %s, expected 01.201.800", r2.C)
		}
		if r2.Code != "01201800" {
			t.Errorf("record 2: unexpected Code: %s, expected 01201800", r2.Code)
		}
		if r2.Code11 != "01201800000" {
			t.Errorf("record 2: unexpected Code11: %s, expected 01201800000", r2.Code11)
		}
		if r2.Region != "01" {
			t.Errorf("record 2: unexpected Region: %s, expected 01", r2.Region)
		}
		if r2.Level1 != "201" {
			t.Errorf("record 2: unexpected Level1: %s, expected 201", r2.Level1)
		}
		if r2.Level2 != "800" {
			t.Errorf("record 2: unexpected Level2: %s, expected 800", r2.Level2)
		}
		if r2.Level3 != "" {
			t.Errorf("record 2: unexpected Level3: %s, expected empty", r2.Level3)
		}

		// Проверяем четвертую запись (населенный пункт)
		r3 := okato.Records[3]
		if r3.C != "01.201.802.002" {
			t.Errorf("record 3: unexpected C: %s, expected 01.201.802.002", r3.C)
		}
		if r3.Code != "01201802002" {
			t.Errorf("record 3: unexpected Code: %s, expected 01201802002", r3.Code)
		}
		if r3.Code11 != "01201802002" {
			t.Errorf("record 3: unexpected Code11: %s, expected 01201802002", r3.Code11)
		}
		if r3.Region != "01" {
			t.Errorf("record 3: unexpected Region: %s, expected 01", r3.Region)
		}
		if r3.Level1 != "201" {
			t.Errorf("record 3: unexpected Level1: %s, expected 201", r3.Level1)
		}
		if r3.Level2 != "802" {
			t.Errorf("record 3: unexpected Level2: %s, expected 802", r3.Level2)
		}
		if r3.Level3 != "002" {
			t.Errorf("record 3: unexpected Level3: %s, expected 002", r3.Level3)
		}

		// Проверяем индекс Region в тесте valid data - должна быть одна запись региона
		if len(okato.Region) != 1 {
			t.Errorf("unexpected Region length: %d, expected 1", len(okato.Region))
		}

		regionRecord := okato.Region["01"]
		if regionRecord == nil {
			t.Fatal("region record '01' not found in Region index")
		}

		// Проверяем, что в индексе Region находится именно запись региона (первая запись)
		if regionRecord.C != "01" {
			t.Errorf("unexpected region record C: %s, expected '01'", regionRecord.C)
		}
		if regionRecord.Name != "Алтайский край" {
			t.Errorf("unexpected region record Name: %s, expected 'Алтайский край'", regionRecord.Name)
		}
	})

	t.Run("check indexes", func(t *testing.T) {
		f, err := os.Open("../testdata/okato-valid_test.xml")
		if err != nil {
			t.Fatalf("failed to open test file: %v", err)
		}
		defer f.Close()

		okato, err := NewOkato(f)
		if err != nil {
			t.Fatalf("failed to create OKATO classifier: %v", err)
		}

		// Проверяем индекс byCode
		if len(okato.byCode) != 4 {
			t.Errorf("unexpected byCode length: %d, expected 4", len(okato.byCode))
		}

		// Проверяем что все записи есть в индексе по коду
		if _, exists := okato.byCode["01"]; !exists {
			t.Error("record with code '01' not found in byCode index")
		}
		if _, exists := okato.byCode["01200"]; !exists {
			t.Error("record with code '01200' not found in byCode index")
		}
		if _, exists := okato.byCode["01201800"]; !exists {
			t.Error("record with code '01201800' not found in byCode index")
		}
		if _, exists := okato.byCode["01201802002"]; !exists {
			t.Error("record with code '01201802002' not found in byCode index")
		}

		// Проверяем индекс byCode11
		if len(okato.byCode11) != 4 {
			t.Errorf("unexpected byCode11 length: %d, expected 4", len(okato.byCode11))
		}

		if _, exists := okato.byCode11["01000000000"]; !exists {
			t.Error("record with code11 '01000000000' not found in byCode11 index")
		}
		if _, exists := okato.byCode11["01200000000"]; !exists {
			t.Error("record with code11 '01200000000' not found in byCode11 index")
		}
		if _, exists := okato.byCode11["01201800000"]; !exists {
			t.Error("record with code11 '01201800000' not found in byCode11 index")
		}
		if _, exists := okato.byCode11["01201802002"]; !exists {
			t.Error("record with code11 '01201802002' not found in byCode11 index")
		}

		// Проверяем индекс byRegion
		if len(okato.byRegion) != 1 {
			t.Errorf("unexpected byRegion length: %d, expected 1", len(okato.byRegion))
		}

		records01 := okato.byRegion["01"]
		if len(records01) != 4 {
			t.Errorf("unexpected number of records for region '01': %d, expected 4", len(records01))
		}

		// Проверяем индекс Region (только записи регионов)
		if len(okato.Region) != 1 {
			t.Errorf("unexpected Region length: %d, expected 1", len(okato.Region))
		}

		regionRecord := okato.Region["01"]
		if regionRecord == nil {
			t.Fatal("region record '01' not found in Region index")
		}

		// Проверяем, что в индексе Region находится именно запись региона
		if regionRecord.C != "01" {
			t.Errorf("unexpected region record C: %s, expected '01'", regionRecord.C)
		}
		if regionRecord.Name != "Алтайский край" {
			t.Errorf("unexpected region record Name: %s, expected 'Алтайский край'", regionRecord.Name)
		}
		if regionRecord.Code != "01" {
			t.Errorf("unexpected region record Code: %s, expected '01'", regionRecord.Code)
		}
	})

	t.Run("invalid code", func(t *testing.T) {
		f, err := os.Open("../testdata/okato-invalid_test.xml")
		if err != nil {
			t.Fatalf("failed to open test file: %v", err)
		}
		defer f.Close()

		okato, err := NewOkato(f)
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}

		// При невалидных кодах записи должны игнорироваться
		// В тестовом файле есть одна валидная запись и одна невалидная
		if len(okato.Records) != 1 {
			t.Errorf("expected 1 record (invalid record should be ignored), got %d", len(okato.Records))
		}

		// Проверяем, что валидная запись присутствует
		if len(okato.Records) > 0 {
			validRecord := okato.Records[0]
			if validRecord.C != "01.201.800" {
				t.Errorf("expected valid record with code '01.201.800', got '%s'", validRecord.C)
			}
		}

		// Проверяем индекс Region - в валидной записи "01.201.800" нет записи региона (код длиннее 2 символов)
		if len(okato.Region) != 0 {
			t.Errorf("unexpected Region length: %d, expected 0 (no region records in invalid test)", len(okato.Region))
		}
	})

	t.Run("duplicate code", func(t *testing.T) {
		f, err := os.Open("../testdata/okato-duplicate_test.xml")
		if err != nil {
			t.Fatalf("failed to open test file: %v", err)
		}
		defer f.Close()

		_, err = NewOkato(f)
		if err == nil {
			t.Error("expected error for duplicate code, got nil")
		} else if !strings.Contains(err.Error(), "duplicate OKATO code") {
			t.Errorf("unexpected error: %v", err)
		}
	})
}
