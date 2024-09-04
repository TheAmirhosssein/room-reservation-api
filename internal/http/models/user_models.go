package models

type Authenticate struct {
	MobileNumber string `json:"mobile_number" binding:"required"`
}
