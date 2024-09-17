package entity

import "gorm.io/gorm"

const (
	AdminRole   = "Admin"
	SupportRole = "Support"
	UserRole    = "User"
)

type User struct {
	gorm.Model
	FullName     string
	MobileNumber string `gorm:"unique"`
	Role         string
}

func NewUser(fullName, mobileNumber string) User {
	return User{FullName: fullName, MobileNumber: mobileNumber}
}
