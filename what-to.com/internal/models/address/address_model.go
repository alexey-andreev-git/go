package models

import (
	"gorm.io/gorm"
)

type Address struct {
	gorm.Model
	Unit     string `json:"unit" gorm:"size:255"`
	Building string `json:"building" gorm:"size:255"`
	Sreet    string `json:"address_line" gorm:"size:255"`
	City     string `json:"city" gorm:"size:255"`
	State    string `json:"state" gorm:"size:255"`
	ZipCode  string `json:"zip_code" gorm:"size:20"`
}

type AddressList struct {
	Items map[int]Address `json:"address_items"`
}
