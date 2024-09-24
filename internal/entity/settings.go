package entity

import "gorm.io/gorm"

type State struct {
	gorm.Model
	Title string
}

func NewState(title string) State {
	return State{Title: title}
}
