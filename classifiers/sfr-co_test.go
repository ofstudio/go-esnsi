package classifiers

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func TestNewSfrCo(t *testing.T) {
	t.Run("successful parsing of valid file", func(t *testing.T) {
		// Arrange
		f, err := os.Open(filepath.Join("testdata", "SFR_CO_valid_test.xml"))
		if err != nil {
			t.Fatalf("Failed to open test file: %v", err)
		}
		defer f.Close()

		// Act
		sfrCo, err := NewSfrCo(f)

		// Assert
		if err != nil {
			t.Fatalf("Expected no error, got: %v", err)
		}
		if sfrCo == nil {
			t.Fatal("Expected non-nil SfrCo")
		}

		// Check classifier metadata
		if sfrCo.Properties.Meta.Code != "SFR_CO" {
			t.Errorf("Expected Code to be SFR_CO, got %s", sfrCo.Properties.Meta.Code)
		}
		if sfrCo.Properties.Meta.Name != "Клиентские службы СФР" {
			t.Errorf("Expected Name to be 'Клиентские службы СФР', got %s", sfrCo.Properties.Meta.Name)
		}
		if sfrCo.Properties.Meta.UID == "" {
			t.Error("Expected non-empty UID")
		}

		// Check records were loaded
		if len(sfrCo.Records) == 0 {
			t.Fatal("Expected records to be loaded")
		}

		// Check specific record content
		record := sfrCo.Records[0]
		if record.COID == 0 {
			t.Error("Expected non-zero COID")
		}
		if record.AutoKey == "" {
			t.Error("Expected non-empty AutoKey")
		}
		if record.ToSfrCode == "" {
			t.Error("Expected non-empty ToSfrCode")
		}
		if record.RegionName == "" {
			t.Error("Expected non-empty RegionName")
		}

		// Check OKATO areas were parsed correctly
		if len(record.OKATOAreas) == 0 {
			t.Error("Expected non-empty OKATOAreas")
		}

		// Check index was built properly
		for _, okato := range record.OKATOAreas {
			indexedRecord, exists := sfrCo.ByOKATO[okato]
			if !exists {
				t.Errorf("Record should be indexed by OKATO %s", okato)
			}
			if exists && indexedRecord.COID != record.COID {
				t.Errorf("Expected indexed record COID to be %d, got %d", record.COID, indexedRecord.COID)
			}
		}
	})

	t.Run("invalid OKATO", func(t *testing.T) {
		// Arrange
		f, err := os.Open(filepath.Join("testdata", "SFR_CO_OKATOArea_Invalid_test.xml"))
		if err != nil {
			t.Fatalf("Failed to open test file: %v", err)
		}
		defer f.Close()

		// Act
		_, err = NewSfrCo(f)

		// Assert
		if err == nil {
			t.Error("Expected error, got nil")
		}
		if err != nil && !strings.Contains(err.Error(), "invalid OKATO") {
			t.Errorf("Expected error message to contain 'invalid OKATO', got: %s", err.Error())
		}
	})

	t.Run("duplicate OKATO code", func(t *testing.T) {
		// Arrange
		f, err := os.Open(filepath.Join("testdata", "SFR_CO_OKATOArea_Duplicate_test.xml"))
		if err != nil {
			t.Fatalf("Failed to open test file: %v", err)
		}
		defer f.Close()

		// Act
		_, err = NewSfrCo(f)

		// Assert
		if err == nil {
			t.Error("Expected error, got nil")
		}
		if err != nil && !strings.Contains(err.Error(), "duplicate OKATO") {
			t.Errorf("Expected error message to contain 'duplicate OKATO', got: %s", err.Error())
		}
	})

	t.Run("invalid OKATO empty", func(t *testing.T) {
		// Arrange
		f, err := os.Open(filepath.Join("testdata", "SFR_CO_OKATOArea_Empty_test.xml"))
		if err != nil {
			t.Fatalf("Failed to open test file: %v", err)
		}
		defer f.Close()

		// Act
		_, err = NewSfrCo(f)

		// Assert
		if err == nil {
			t.Error("Expected error, got nil")
		}
		if err != nil && !strings.Contains(err.Error(), "invalid OKATO") {
			t.Errorf("Expected error message to contain 'invalid OKATO', got: %s", err.Error())
		}
	})

}
