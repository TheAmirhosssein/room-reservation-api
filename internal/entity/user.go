package entity

import "gorm.io/gorm"

type User struct {
	gorm.Model
	FullName     string
	MobileNumber string `gorm:"unique"`
}

func (user User) Table() string {
	return "users"
}
