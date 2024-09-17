package entity

import "gorm.io/gorm"

const (
	AdminRole   string = "Admin"
	SupportRole string = "Support"
	UserRole    string = "User"
)

type User struct {
	gorm.Model
	FullName     string
	MobileNumber string `gorm:"unique"`
	Role         string
}

func NewUser(fullName, mobileNumber, role string) User {
	return User{FullName: fullName, MobileNumber: mobileNumber}
}
