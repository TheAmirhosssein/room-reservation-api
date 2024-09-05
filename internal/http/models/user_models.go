package models

type (
	Authenticate struct {
		MobileNumber string `json:"mobile_number" binding:"required"`
	}

	Token struct {
		MobileNumber string `json:"mobile_number" binding:"required"`
		Code         string `json:"code" binding:"required"`
	}
)
