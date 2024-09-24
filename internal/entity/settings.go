package entity

import "gorm.io/gorm"

type State struct {
	gorm.Model
	Title string
}

func NewString(title string) State {
	return State{Title: title}
}
