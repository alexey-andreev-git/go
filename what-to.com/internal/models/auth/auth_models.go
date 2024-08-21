package models

import (
	"gorm.io/gorm"
	person "what-to.com/internal/models/person"
)

// Struct auth models
type User struct {
	gorm.Model
	Name     string        `json:"name" gorm:"size:255"`
	Password string        `json:"password" gorm:"size:255"`
	Person   person.Person `json:"person" gorm:"constraint:OnUpdate:CASCADE,OnDelete:SET NULL;"`
}

type UsersList struct {
	Items map[int]User `json:"auth_user_data_items"`
}
