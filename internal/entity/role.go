package entity

import "gorm.io/gorm"

type Role struct {
	gorm.Model
	Name string `gorm:"unique"`
}

func NewRole(Name string) Role {
	return Role{Name: Name}
}
