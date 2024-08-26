package entity

import "gorm.io/gorm"

type User struct {
	gorm.Model
	FullName     string
	MobileNumber string `gorm:"unique"`
	IsVerified   bool   `gorm:"default:false"`
}
