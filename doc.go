/*
Package esnsi предоставляет инструменты для работы с простыми классификаторами ЕСНСИ
(Единая система нормативно-справочной информации) в формате XML ЦНСИ.

Основные возможности:
  - Автоматическое декодирование XML простых классификаторов ЕСНСИ в типизированные структуры Go на основе структурных тегов
  - Поддержка строковых и целочисленных атрибутов записей
  - Валидация типов полей структуры записи в соответствии с типами атрибутов в XML
  - Поддержка пользовательских обработчиков записей

# Основные типы

Classifier[T] - простой классификатор ЕСНСИ, содержащий метаданные и записи классификатора.

Decoder[T] - декодер для преобразования XML формата ЦНСИ в структуры Go.

# Пример базового использования

	package main

	import (
		"fmt"
		"os"
		"github.com/ofstudio/go-esnsi"
	)

	// Определяем структуру записи классификатора
	type RegionRecord struct {
		Code   string `esnsi:"RegionCode"`
		Name   string `esnsi:"RegionName"`
		Type   int    `esnsi:"RegionType"`
	}

	func main() {
		// Открываем XML файл классификатора
		file, err := os.Open("regions.xml")
		if err != nil {
			panic(err)
		}
		defer file.Close()

		// Создаем декодер и классификатор
		decoder := esnsi.NewDecoder[RegionRecord](file)
		classifier := &esnsi.Classifier[RegionRecord]{}

		// Декодируем XML в классификатор
		if err := decoder.Decode(classifier); err != nil {
			panic(err)
		}

		// Используем данные
		fmt.Printf("Классификатор: %s (%s)\n", classifier.Name, classifier.Code)
		fmt.Printf("Версия: %d\n", classifier.Version)
		fmt.Printf("Загружено записей: %d\n", len(classifier.Records))

		// Обрабатываем записи
		for _, record := range classifier.Records {
			fmt.Printf("Регион %s: %s (тип %d)\n", record.Code, record.Name, record.Type)
		}
	}

# Пример с пользовательским обработчиком

Обработчики позволяют выполнять дополнительную валидацию, трансформацию данных или фильтрацию записей:

	func main() {
		file, err := os.Open("regions.xml")
		if err != nil {
			panic(err)
		}
		defer file.Close()

		var validRecords []RegionRecord
		decoder := esnsi.NewDecoder[RegionRecord](file).WithHandler(func(rec *RegionRecord) error {
			// Валидируем запись
			if rec.Code == "" {
				return fmt.Errorf("пустой код региона")
			}

			// Преобразуем данные
			rec.Name = strings.ToUpper(rec.Name)

			// Фильтруем записи (сохраняем только регионы определенного типа)
			if rec.Type == 1 {
				validRecords = append(validRecords, *rec)
			}

			return nil
		})

		classifier := &esnsi.Classifier[RegionRecord]{}
		if err := decoder.Decode(classifier); err != nil {
			panic(err)
		}

		fmt.Printf("Обработано %d валидных записей\n", len(validRecords))
	}

# Требования к структуре записи

Структура записи должна быть struct с полями, помеченными тегом `esnsi`, который должен соответствовать имени атрибута в XML классификаторе:

	type MyRecord struct {
		ID     int    `esnsi:"RecordID"`      // integer-attribute в XML
		Name   string `esnsi:"RecordName"`    // string-attribute или text-attribute в XML
		Active int   `esnsi:"IsActive"`       // integer-attribute (0/1) в XML

		// Поля без тега esnsi игнорируются при декодировании
		CustomField string
	}

Поддерживаемые типы полей:
  - string - для string-attribute и text-attribute
  - int - для integer-attribute

Декодер автоматически проверяет соответствие типов полей структуры типам атрибутов в XML
и возвращает ошибку при несоответствии.

# Специализированные классификаторы

Пакет classifiers содержит готовые реализации для конкретных классификаторов ЕСНСИ.
*/
package esnsi
