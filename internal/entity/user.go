package entity

import "gorm.io/gorm"

type User struct {
	gorm.Model
	FullName     string `json:"full_name"`
	MobileNumber string `json:"mobile_number" binding:"required" gorm:"unique"`
}

func (user User) Table() string {
	return "users"
}
