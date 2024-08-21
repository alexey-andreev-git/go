package models

import (
	address "what-to.com/internal/models/address"
	email "what-to.com/internal/models/email"
	phone "what-to.com/internal/models/phone"
)

// Struct company models
type Company struct {
	Name        string               `json:"name"`
	EmailsList  *email.EmaislList    `json:"email_list"`
	PhonesList  *phone.PhonesList    `json:"phone_list"`
	AddressList *address.AddressList `json:"address_list"`
}

type CompaniesList struct {
	Items map[int]Company `json:"company_items"`
}
