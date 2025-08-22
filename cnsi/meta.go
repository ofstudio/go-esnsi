package cnsi

// Meta represents a classifier metadata
type Meta struct {
	Name    string `xml:"name,attr"`
	Code    string `xml:"code,attr"`
	UID     string `xml:"uid,attr"`
	Version int    `xml:"version,attr"`
}
