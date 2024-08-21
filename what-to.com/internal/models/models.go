package models

import (
	"database/sql"
	"time"
)

type DbModel struct {
	ID        uint `db:"id;primarykey"`
	CreatedAt time.Time
	UpdatedAt time.Time
	DeletedAt sql.NullTime `db:"index"`
}

type Address struct {
	DbModel
	Unit     string `json:"unit" gorm:"size:255" db:"address_unit"`
	Building string `json:"building" gorm:"size:255" db:"address_building"`
	Street   string `json:"address_line" gorm:"size:255" db:"address_line"`
	City     string `json:"city" gorm:"size:255" db:"address_city"`
	State    string `json:"state" gorm:"size:255" db:"address_state"`
	ZipCode  string `json:"zip_code" gorm:"size:20" db:"address_zip_code"`
}

type AddressList struct {
	Items map[int]Address `json:"address_items"`
}

// Struct auth models
type User struct {
	DbModel
	Name     string `json:"name" gorm:"size:255"`
	Password string `json:"password" gorm:"size:255"`
	Person   Person `json:"person" gorm:"constraint:OnUpdate:CASCADE,OnDelete:SET NULL;foreignKey:ID"`
}

type UsersList struct {
	Items []User `json:"user_data_items"`
}

// Struct company models
type Company struct {
	DbModel
	Name        string    `json:"name" gorm:"size:255"`
	AddressList []Address `gorm:"many2many:company_addresses"`
	EmailsList  []Email   `gorm:"many2many:company_emails"`
	PhonesList  []Phone   `gorm:"many2many:company_phone"`
}

// Join tables
type CompanyAddress struct {
	AddressID uint `json:"address_id"`
	CompanyID uint `json:"company_id"`
}

type CompanyEmail struct {
	EmailID   uint `json:"email_id"`
	CompanyID uint `json:"company_id"`
}

type CompanyPhone struct {
	PhoneID   uint `json:"phone_id"`
	CompanyID uint `json:"company_id"`
}

type CompaniesList struct {
	Items map[int]Company `json:"company_items"`
}

// Struct document models
type Document struct {
	Name   string `json:"name" gorm:"size:255"`
	Number string `json:"email" gorm:"size:255"`
	Issued string `json:"issued" gorm:"size:255"`
	Expiry string `json:"expiry" gorm:"size:255"`
}

type DocumentList struct {
	Items map[int]Document `json:"email_items"`
}

// Struct email models
type Email struct {
	DbModel
	Name  string `json:"name" gorm:"size:255"`
	Email string `json:"email" gorm:"size:255;uniqueIndex"`
}

type EmaislList struct {
	Items map[int]Email `json:"email_items"`
}

// Entity represents an entity in the database
type Entity struct {
	Id      int    `json:"id"`
	Name    string `json:"name"`
	Comment string `json:"comment"`
}

// EntityData represents the data associated with an entity in the database
type EntityData struct {
	Entity int `json:"entity"`
	Order  int `json:"order"`
	Value  int `json:"value"`
}

// EntityDataReference represents the reference data associated with an entity in the database
type EntityDataReference struct {
	Entity  int    `json:"entity"`
	Order   int    `json:"order"`
	Name    string `json:"name"`
	Type    string `json:"type"`
	Comment string `json:"comment"`
}

// Struct person models
type Person struct {
	DbModel
	FirstName   string    `json:"first_name" gorm:"size:255" db:"first_name"`
	MiddleName  string    `json:"middle_name" gorm:"size:255" db:"middle_name"`
	LastName    string    `json:"last_name" gorm:"size:255" db:"last_name"`
	AddressList []Address `gorm:"many2many:person_address" db:"address_list"`
	EmailsList  []Email   `gorm:"many2many:person_emails" db:"email_list"`
	PhonesList  []Phone   `gorm:"many2many:person_phones" db:"phone_list"`
}

type PersonAddress struct {
	AddressID uint `json:"address_id"`
	PersonID  uint `json:"person_id"`
}

type PersonEmail struct {
	EmailID  uint `json:"email_id"`
	PersonID uint `json:"person_id"`
}

type PersonPhone struct {
	PhoneID  uint `json:"phone_id"`
	PersonID uint `json:"person_id"`
}

type PersonsList struct {
	Items map[int]Person `json:"person_items"`
}

// Struct phone models
type Phone struct {
	DbModel
	Name  string `json:"name" gorm:"size:255"`
	Phone string `json:"phone" gorm:"size:255;uniqueIndex"`
}

type PhonesList struct {
	Items map[int]Phone `json:"phone_items"`
}
