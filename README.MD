# go-esnsi

[![Go Reference](https://pkg.go.dev/badge/github.com/ofstudio/go-esnsi.svg)](https://pkg.go.dev/github.com/ofstudio/go-esnsi)
[![Go Report Card](https://goreportcard.com/badge/github.com/ofstudio/go-esnsi)](https://goreportcard.com/report/github.com/ofstudio/go-esnsi)

Библиотека для работы со справочниками (классификаторами) Единой системы нормативно-справочной информации (ЕСНСИ)
на языке Go.

## Установка

```bash
go get github.com/ofstudio/go-esnsi
```

## Пакет esnsi

Основной пакет предоставляет инструменты для работы с простыми классификаторами ЕСНСИ в формате XML ЦНСИ.

### Основные возможности:
- Автоматическое декодирование XML простых классификаторов ЕСНСИ в типизированные структуры Go
- Поддержка строковых и целочисленных атрибутов записей
- Валидация типов полей структуры записи
- Поддержка пользовательских обработчиков записей

### Пример базового использования:

```go
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
```

## Пакет classifiers

Пакет `classifiers` содержит готовые к использованию реализации популярных классификаторов ЕСНСИ с дополнительными индексами для быстрого поиска записей.

### Реализованные справочники:

- **[SFR_CO](https://esnsi.gosuslugi.ru/classifiers/10991/)** — Клиентские службы СФР 
- **[classifierOkato](https://esnsi.gosuslugi.ru/classifiers/16270/)** — Общероссийский классификатор объектов административно-территориального деления (ОКАТО)

### Пример использования справочника СФР:

```go
package main

import (
    "fmt"
    "log"
    "os"

    "github.com/ofstudio/go-esnsi/classifiers"
)

func main() {
    // Открываем файл с данными классификатора СФР
    file, err := os.Open("SFR_CO_55UTF-8.xml")
    if err != nil {
        log.Fatal(err)
    }
    defer file.Close()

    // Создаем классификатор из XML-файла
    sfr, err := classifiers.NewSfr(file)
    if err != nil {
        log.Fatal(err)
    }

    // Выводим информацию о классификаторе
    fmt.Printf("Классификатор: %s\n", sfr.Name)
    fmt.Printf("Версия: %d\n", sfr.Version)
    fmt.Printf("Всего записей: %d\n", len(sfr.Records))

    // Поиск службы по коду ОКАТО
    if service := sfr.ByOkato["77401"]; service != nil {
        fmt.Printf("Найдена служба: %s\n", service.ToSfrName)
        fmt.Printf("Адрес: %s\n", service.Address)
        fmt.Printf("Телефон: %s\n", service.Phone)
    }

    // Поиск всех служб в регионе по 5-значному коду ОКАТО
    if services := sfr.ByOkato5["77000"]; len(services) > 0 {
        fmt.Printf("\nНайдено %d служб в регионе:\n", len(services))
        for _, service := range services {
            fmt.Printf("- %s\n", service.ToSfrName)
        }
    }
}
```

## Документация ЕСНСИ

Методические рекомендации по работе с ЕСНСИ доступны в документе `doc/Методические рекомендации ЕСНСИ.docx`.

## Лицензия

Распространяется по лицензии MIT. Более подробная информация в файле LICENSE.