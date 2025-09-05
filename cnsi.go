package esnsi

import (
	"encoding/xml"
)

// cnsiDoc - структура для разбора XML формата ЦНСИ.
//
//	urn://x-artefacts-nsi-gov-ru/services/cnsi/2.0.0.0
type cnsiDoc struct {
	XMLName xml.Name     `xml:"document"`
	Meta    cnsiMeta     `xml:"simple-classifier"`
	Records []cnsiRecord `xml:"data>record"`
}

// cnsiMeta - метаданные классификатора
type cnsiMeta struct {
	Name         string            `xml:"name,attr"`
	Code         string            `xml:"code,attr"`
	UID          string            `xml:"uid,attr"`
	Version      int               `xml:"version,attr"`
	StringAttrs  []cnsiStringAttr  `xml:"string-attribute"`
	TextAttrs    []cnsiStringAttr  `xml:"text-attribute"`
	IntegerAttrs []cnsiIntegerAttr `xml:"integer-attribute"`
}

// cnsiStringAttr - строковый атрибут
type cnsiStringAttr struct {
	UID  string `xml:"uid,attr"`
	Name string `xml:"name,attr"`
}

// cnsiIntegerAttr - целочисленный атрибут
type cnsiIntegerAttr struct {
	UID  string `xml:"uid,attr"`
	Name string `xml:"name,attr"`
}

// cnsiRecord - запись классификатора
type cnsiRecord struct {
	UID      string        `xml:"uid,attr,omitempty"`
	AttrVals []cnsiAttrVal `xml:"attribute-value"`
}

// cnsiAttrVal - значение атрибута записи
type cnsiAttrVal struct {
	AttrRef    string          `xml:"attribute-ref,attr"`
	IntegerVal *cnsiIntegerVal `xml:"integer"`
	TextVal    *cnsiStringVal  `xml:"text"`
	StringVal  *cnsiStringVal  `xml:"string"`
}

// Вспомогательные структуры для значений атрибутов

type cnsiIntegerVal struct {
	Val int `xml:",chardata"`
}

type cnsiStringVal struct {
	Val string `xml:",chardata"`
}
