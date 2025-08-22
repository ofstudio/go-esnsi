package cnsi

// Record represents a data record for a classifier
type Record struct {
	UID             string           `xml:"uid,attr,omitempty"`
	AttributeValues []AttributeValue `xml:"attribute-value"`
}

// AttributeValue represents an attribute value with reference to its definition
type AttributeValue struct {
	AttributeRef string        `xml:"attribute-ref,attr"`
	IntegerValue *IntegerValue `xml:"integer"`
	TextValue    *StringValue  `xml:"text"`
	StringValue  *StringValue  `xml:"string"`
}

// IntegerValue wraps an integer value
type IntegerValue struct {
	Value int `xml:",chardata"`
}

// StringValue wraps a string and text values
type StringValue struct {
	Value string `xml:",chardata"`
}
