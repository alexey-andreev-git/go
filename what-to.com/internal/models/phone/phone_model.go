package models

// Struct phone models
type Phone struct {
	Name  string `json:"name"`
	Phone string `json:"phone"`
}

type PhonesList struct {
	Items map[int]Phone `json:"phone_items"`
}
