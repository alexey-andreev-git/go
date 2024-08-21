package models

import (
	address "what-to.com/internal/models/address"
	email "what-to.com/internal/models/email"
	phone "what-to.com/internal/models/phone"
)

// Struct person models
type Person struct {
	FirstName   string               `json:"first_name" database:"first_name"`
	MiddleName  string               `json:"middle_name" database:"middle_name"`
	LastName    string               `json:"last_name" database:"last_name"`
	EmailsList  *email.EmaislList    `json:"email_list" database:"email_list"`
	PhonesList  *phone.PhonesList    `json:"phone_list" database:"phone_list"`
	AddressList *address.AddressList `json:"address_list" database:"address_list"`
}

type PersonsList struct {
	Items map[int]Person `json:"person_items"`
}
