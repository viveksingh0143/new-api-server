package dto

type IError struct {
	Field   string `json:"field,omitempty"`
	Tag     string `json:"tag,omitempty"`
	Message string `json:"message"`
	Value   string `json:"value,omitempty"`
}
