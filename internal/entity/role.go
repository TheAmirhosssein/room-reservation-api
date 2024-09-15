package entity

import "gorm.io/gorm"

const (
	AdminRole   string = "admin"
	SupportRole string = "support"
	UserRole    string = "user"
)

type Role struct {
	gorm.Model
	Name string `gorm:"unique"`
}

func NewRole(Name string) Role {
	return Role{Name: Name}
}
