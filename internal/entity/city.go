package entity

import "gorm.io/gorm"

type City struct {
	gorm.Model
	Title   string
	StateID uint
	State   State `gorm:"foreignKey:StateID;references:ID"`
}

func NewCity(title string, state State) City {
	return City{Title: title, State: state}
}
