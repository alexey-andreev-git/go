package models

// Struct document models
type Document struct {
	Name   string `json:"name"`
	Number string `json:"email"`
	Issued string `json:"issued"`
	Expiry string `json:"expiry"`
}

type DocumentList struct {
	Items map[int]Document `json:"email_items"`
}
