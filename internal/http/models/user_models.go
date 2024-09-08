package models

import (
	"time"

	"github.com/TheAmirhosssein/room-reservation-api/internal/entity"
)

type (
	Authenticate struct {
		MobileNumber string `json:"mobile_number" binding:"required"`
	}

	Token struct {
		MobileNumber string `json:"mobile_number" binding:"required"`
		Code         string `json:"code" binding:"required"`
	}
)

type UserResponse struct {
	Id           uint      `json:"id"`
	MobileNumber string    `json:"mobile_number"`
	FullName     string    `json:"full_name"`
	JoinedAt     time.Time `json:"joined_at"`
}

func NewUserResponse(userEntity entity.User) UserResponse {
	return UserResponse{
		Id:           userEntity.ID,
		MobileNumber: userEntity.MobileNumber,
		FullName:     userEntity.FullName,
		JoinedAt:     userEntity.CreatedAt,
	}
}
