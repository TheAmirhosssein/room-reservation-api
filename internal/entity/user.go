package entity

import "gorm.io/gorm"

type User struct {
	gorm.Model
	FullName     string
	MobileNumber string `gorm:"unique"`
}

func NewUser(fullName, mobileNumber string) User {
	return User{FullName: fullName, MobileNumber: mobileNumber}
}
