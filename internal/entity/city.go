package entity

import "gorm.io/gorm"

type City struct {
	gorm.Model
	Title string
	State State
}

func NewCity(title string, state State) City {
	return City{Title: title, State: state}
}
