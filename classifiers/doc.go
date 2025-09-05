/*
Package classifiers предоставляет готовые к использованию справочники ЕСНСИ.

# Справочники

  - Sfr - Клиентские службы СФР
  - Okato - Общероссийский классификатор объектов административно-территориального деления (ОКАТО)

# Пример использования справочника Sfr

	package main

	import (
		"fmt"
		"log"
		"os"

		"github.com/ofstudio/go-esnsi/classifiers"
	)

	func main() {
		// Открываем файл с данными классификатора
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

		// Поиск службы по точному коду ОКАТО
		if service := sfr.ByOkato["77401"]; service != nil {
			fmt.Printf("Найдена служба: %s\n", service.ToSfrName)
			fmt.Printf("Адрес: %s\n", service.Address)
		}

		// Перебор всех записей
		for _, record := range sfr.Records {
			fmt.Printf("Код СФР: %s, Название: %s\n",
				record.ToSfrCode, record.ToSfrName)
		}
	}
*/
package classifiers
