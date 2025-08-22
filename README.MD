# go-esnsi

[![Go Reference](https://pkg.go.dev/badge/github.com/ofstudio/go-esnsi.svg)](https://pkg.go.dev/github.com/ofstudio/go-esnsi)
[![Go Report Card](https://goreportcard.com/badge/github.com/ofstudio/go-esnsi)](https://goreportcard.com/report/github.com/ofstudio/go-esnsi)

Библиотека для работы с классификаторами Единой системы нормативно-справочной информации (ЕСНСИ)
на языке Go.

## Описание

Данный модуль предоставляет инструменты для работы с классификаторами ЕСНСИ, 
позволяя загружать и анализировать XML файлы, соответствующие стандартной XSD схеме ЕСНСИ. 
Модуль разделен на два основных пакета:

- `cnsi` - базовые структуры и разбор XML-файлов классификаторов формата ЦНСИ
- `classifiers` - специализированные реализации конкретных классификаторов

## Установка

```bash
go get github.com/ofstudio/go-esnsi
```


## Примеры использования

### Загрузка и использование простого классификатора

```go
package main

import (
    "fmt"
    "os"
    
    "github.com/ofstudio/go-esnsi/cnsi"
)

type MyRecord struct {
	Name    string `esnsi:"Наименование"`
	Address string `esnsi:"Адрес"`
	Value   int    `esnsi:"Значение"`
}

func main() {
    // Открываем файл с классификатором
    file, err := os.Open("classifier.xml")
    if err != nil {
        fmt.Printf("Ошибка открытия файла: %v\n", err)
        return
    }
    defer file.Close()
    
    // Загружаем классификатор
    classifier, err := cnsi.NewSimpleClassifier(r)
    if err != nil {
        fmt.Printf("Ошибка загрузки классификатора: %v\n", err)
        return
    }
	
    // Получаем записи из классификатора
	var records []MyRecord
	err := classifier.Unmarshal(&records, nil)
	if err != nil {
        fmt.Printf("Ошибка разбора классификатора: %v\n", err)
        return
    }

    // Выводим данные из классификатора
    for _, record := range records {
        fmt.Printf("Наименование: %s, Адрес: %s, Значение: %d\n", 
                   record.Name, record.Address, record.Value)
    }
}
```

### Работа с классификатором клиентских служб СФР

```go
package main

import (
    "fmt"
    "os"
    
    "github.com/ofstudio/go-esnsi/classifiers"
)

func main() {
    // Открываем файл с классификатором клиентских служб СФР
    file, err := os.Open("SFR_CO.xml")
    if err != nil {
        fmt.Printf("Ошибка открытия файла: %v\n", err)
        return
    }
    defer file.Close()
    
    // Загружаем классификатор
    sfrCo, err := classifiers.NewSfrCo(file)
    if err != nil {
        fmt.Printf("Ошибка загрузки классификатора СФР: %v\n", err)
        return
    }
    
    // Поиск клиентской службы по коду ОКАТО
    okato := "12345678901"
    office, exists := sfrCo.ByOKATO[okato]
    if exists {
        fmt.Printf("Найдена клиентская служба: %s\n", office.ToSfrName)
        fmt.Printf("Адрес: %s\n", office.Address)
        fmt.Printf("Телефон: %s\n", office.Phone)
        fmt.Printf("Почта: %s\n", office.Email)
    } else {
        fmt.Printf("Клиентская служба для кода ОКАТО %s не найдена\n", okato)
    }
    
    // Получение всех клиентских служб региона
    regionCode := "77"
    for _, record := range sfrCo.Records {
        if record.RegionCode == regionCode {
            fmt.Printf("Клиентская служба в регионе %s: %s\n", 
                       record.RegionName, record.ToSfrName)
        }
    }
}
```

## Документация

Документация доступна в коде через комментарии в формате godoc.

Более подробные методические рекомендации по работе с ЕСНСИ 
доступны в документе `cnsi/doc/Методические рекомендации ЕСНСИ.docx`.


## Лицензия

Распространяется по лицензии MIT. Более подробная информация в файле LICENSE.