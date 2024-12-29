package models

// Struct auth models
// Struct auth models
//
//	type User struct {
//		DbModel
//		Name     string `json:"name" gorm:"size:255"`
//		Password string `json:"password" gorm:"size:255"`
//		Token    string `json:"token" gorm:"size:255"`
//		Person   Person `json:"person" gorm:"constraint:OnUpdate:CASCADE,OnDelete:SET NULL;foreignKey:ID"`
//	}
type (
	User struct {
		Name     string `json:"username" gorm:"size:255"`
		Password string `json:"password" gorm:"size:255"`
		Token    string `json:"token" gorm:"size:255"`
		// Person   Person `json:"person" gorm:"constraint:OnUpdate:CASCADE,OnDelete:SET NULL;"`
	}
	UsersList struct {
		Items map[int]User `json:"auth_user_data_items"`
	}
)
