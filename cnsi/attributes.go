package cnsi

import "reflect"

type Attributes struct {
	StringAttributes  []StringAttribute  `xml:"string-attribute"`
	TextAttributes    []StringAttribute  `xml:"text-attribute"`
	IntegerAttributes []IntegerAttribute `xml:"integer-attribute"`
}

func (a Attributes) indexes() (map[string]string, map[string]reflect.Kind) {
	// Create a map to hold the UID for each attribute name
	nameToUID := make(map[string]string)
	UIDToKind := make(map[string]reflect.Kind)

	// Add string attributes
	for _, attr := range a.StringAttributes {
		nameToUID[attr.Name] = attr.UID
		UIDToKind[attr.UID] = reflect.String
	}

	// Add text attributes
	for _, attr := range a.TextAttributes {
		nameToUID[attr.Name] = attr.UID
		UIDToKind[attr.UID] = reflect.String
	}

	// Add integer attributes
	for _, attr := range a.IntegerAttributes {
		nameToUID[attr.Name] = attr.UID
		UIDToKind[attr.UID] = reflect.Int
	}

	return nameToUID, UIDToKind
}

// StringAttribute represents a string and text attribute definition
type StringAttribute struct {
	UID  string `xml:"uid,attr"`
	Name string `xml:"name,attr"`
}

// IntegerAttribute represents an integer attribute definition
type IntegerAttribute struct {
	UID  string `xml:"uid,attr"`
	Name string `xml:"name,attr"`
}
