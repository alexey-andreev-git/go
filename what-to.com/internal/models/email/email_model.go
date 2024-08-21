package models

// Struct email models
type Email struct {
	Name  string `json:"name"`
	Email string `json:"email"`
}

type EmaislList struct {
	Items map[int]Email `json:"email_items"`
}
