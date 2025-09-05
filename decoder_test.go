package esnsi

import (
	"fmt"
	"os"
	"strings"
	"testing"
)

//goland:noinspection GoUnhandledErrorResult
func TestDecoder_Decode(t *testing.T) {
	t.Run("wrong type", func(t *testing.T) {
		err := NewDecoder[int](nil).Decode(&Classifier[int]{})
		if err == nil {
			t.Error("expected error, got nil")
		} else if !strings.Contains(err.Error(), "struct type expected") {
			t.Errorf("unexpected error: %v", err)
		}
	})

	t.Run("reader is nil", func(t *testing.T) {
		err := NewDecoder[testRecord](nil).Decode(&Classifier[testRecord]{})
		if err == nil {
			t.Error("expected error, got nil")
		} else if !strings.Contains(err.Error(), "reader is nil") {
			t.Errorf("unexpected error: %v", err)
		}
	})

	t.Run("nil classifier", func(t *testing.T) {
		f, err := os.Open("testdata/decoder-valid_test.xml")
		if err != nil {
			t.Fatalf("failed to open test file: %v", err)
		}
		defer f.Close()

		err = NewDecoder[testRecord](f).Decode(nil)
		if err == nil {
			t.Error("expected error, got nil")
		} else if !strings.Contains(err.Error(), "nil pointer passed") {
			t.Errorf("unexpected error: %v", err)
		}
	})

	t.Run("valid data", func(t *testing.T) {
		f, err := os.Open("testdata/decoder-valid_test.xml")
		if err != nil {
			t.Fatalf("failed to open test file: %v", err)
		}
		defer f.Close()

		classifier := &Classifier[testRecord]{}
		err = NewDecoder[testRecord](f).
			WithHandler(func(rec *testRecord) error {
				// Просто добавляем запись в слайс Records
				classifier.Records = append(classifier.Records, *rec)
				return nil
			}).
			Decode(classifier)
		if err != nil {
			t.Fatalf("failed to decode classifier: %v", err)
		}

		// Проверяем метаданные
		if classifier.Name != "Тестовый классификатор" {
			t.Errorf("unexpected Name: %s", classifier.Name)
		}
		if classifier.Code != "TestClassifier" {
			t.Errorf("unexpected Code: %s", classifier.Code)
		}
		if classifier.UID != "2fad55ce-854c-4044-873e-7ae806cc94cb" {
			t.Errorf("unexpected UID: %s", classifier.UID)
		}
		if classifier.Version != 55 {
			t.Errorf("unexpected Version: %d", classifier.Version)
		}

		// Проверяем записи
		if len(classifier.Records) != 4 {
			t.Fatalf("unexpected number of records: %d", len(classifier.Records))
		}

		// Проверяем первую запись
		record0 := classifier.Records[0]
		if record0.ToSfrCode != "210" {
			t.Errorf("record 0: unexpected ToSfrCode: %s, expected 210", record0.ToSfrCode)
		}
		if record0.RegionName != "Республика Татарстан" {
			t.Errorf("record 0: unexpected RegionName: %s, expected Республика Татарстан", record0.RegionName)
		}
		if record0.OfficeType != 2 {
			t.Errorf("record 0: unexpected OfficeType: %d, expected 2", record0.OfficeType)
		}

		// Проверяем вторую запись
		r1 := classifier.Records[1]
		if r1.ToSfrCode != "201" {
			t.Errorf("record 1: unexpected ToSfrCode: %s, expected 201", r1.ToSfrCode)
		}
		if r1.RegionName != "Москва" {
			t.Errorf("record 1: unexpected RegionName: %s, expected Москва", r1.RegionName)
		}
		if r1.OfficeType != 2 {
			t.Errorf("record 1: unexpected OfficeType: %d, expected 2", r1.OfficeType)
		}

		// Проверяем третью запись
		r2 := classifier.Records[2]
		if r2.ToSfrCode != "059" {
			t.Errorf("record 2: unexpected ToSfrCode: %s, expected 059", r2.ToSfrCode)
		}
		if r2.RegionName != "Магаданская область" {
			t.Errorf("record 2: unexpected RegionName: %s, expected Магаданская область", r2.RegionName)
		}
		if r2.OfficeType != 2 {
			t.Errorf("record 2: unexpected OfficeType: %d, expected 2", r2.OfficeType)
		}

		// Проверяем четвертую запись
		r3 := classifier.Records[3]
		if r3.ToSfrCode != "041" {
			t.Errorf("record 3: unexpected ToSfrCode: %s, expected 041", r3.ToSfrCode)
		}
		if r3.RegionName != "Белгородская область" {
			t.Errorf("record 3: unexpected RegionName: %s, expected Белгородская область", r3.RegionName)
		}
		if r3.OfficeType != 2 {
			t.Errorf("record 3: unexpected OfficeType: %d, expected 2", r3.OfficeType)
		}
	})

	t.Run("wrong field type", func(t *testing.T) {
		f, err := os.Open("testdata/decoder-valid_test.xml")
		if err != nil {
			t.Fatalf("failed to open test file: %v", err)
		}
		defer f.Close()

		err = NewDecoder[testWrongFieldTypeRecord](f).Decode(&Classifier[testWrongFieldTypeRecord]{})
		if err == nil {
			t.Error("expected error, got nil")
		} else if !strings.Contains(err.Error(), "field OfficeType has type string, expected int") {
			t.Errorf("unexpected error: %v", err)
		}
	})

	t.Run("missing field", func(t *testing.T) {
		f, err := os.Open("testdata/decoder-missing-field_test.xml")
		if err != nil {
			t.Fatalf("failed to open test file: %v", err)
		}
		defer f.Close()

		err = NewDecoder[testRecord](f).Decode(&Classifier[testRecord]{})
		if err == nil {
			t.Error("expected error, got nil")
		} else if !strings.Contains(err.Error(), "attribute RegionName not found in classifier") {
			t.Errorf("unexpected error: %v", err)
		}
	})

	t.Run("with custom handler", func(t *testing.T) {
		f, err := os.Open("testdata/decoder-valid_test.xml")
		if err != nil {
			t.Fatalf("failed to open test file: %v", err)
		}
		defer f.Close()

		var processedRecords []testRecord
		decoder := NewDecoder[testRecord](f).WithHandler(func(rec *testRecord) error {
			// Добавляем пометку о том, что запись была обработана
			rec.TestField = "processed by handler"
			processedRecords = append(processedRecords, *rec)
			return nil
		})

		classifier := &Classifier[testRecord]{}
		err = decoder.Decode(classifier)
		if err != nil {
			t.Fatalf("failed to decode with handler: %v", err)
		}

		// При использовании handler'а записи не должны добавляться в classifier.Records
		if len(classifier.Records) != 0 {
			t.Errorf("expected 0 records in classifier, got %d", len(classifier.Records))
		}

		// Но должны быть обработаны handler'ом
		if len(processedRecords) != 4 {
			t.Fatalf("expected 4 processed records, got %d", len(processedRecords))
		}

		// Проверяем, что handler обработал записи
		for i, rec := range processedRecords {
			if rec.TestField != "processed by handler" {
				t.Errorf("record %d: expected TestField 'processed by handler', got '%s'", i, rec.TestField)
			}
		}
	})

	t.Run("handler error", func(t *testing.T) {
		f, err := os.Open("testdata/decoder-valid_test.xml")
		if err != nil {
			t.Fatalf("failed to open test file: %v", err)
		}
		defer f.Close()

		decoder := NewDecoder[testRecord](f).WithHandler(func(rec *testRecord) error {
			// Возвращаем ошибку при обработке второй записи
			if rec.ToSfrCode == "201" {
				return fmt.Errorf("handler error for record %s", rec.ToSfrCode)
			}
			return nil
		})

		classifier := &Classifier[testRecord]{}
		err = decoder.Decode(classifier)
		if err == nil {
			t.Error("expected error, got nil")
		} else if !strings.Contains(err.Error(), "handler error for record 201") {
			t.Errorf("unexpected error: %v", err)
		}
	})
}

// testRecord - валидная тестовая запись классификатора
type testRecord struct {
	ToSfrCode  string `esnsi:"ToSfrCode"`
	RegionName string `esnsi:"RegionName"`
	OfficeType int    `esnsi:"OfficeType"`
	TestField  string
}

// testWrongFieldTypeRecord - запись с неправильным типом поля: string вместо int
type testWrongFieldTypeRecord struct {
	OfficeType string `esnsi:"OfficeType"`
}
