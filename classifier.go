package esnsi

// Classifier - простой классификатор ЕСНСИ
type Classifier[T any] struct {
	Name    string
	Code    string
	UID     string
	Version int
	Records []T
}
