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

	UpdateUser struct {
		FullName string `json:"full_name" binding:"required"`
	}
)

type UserResponse struct {
	Id           uint      `json:"id"`
	MobileNumber string    `json:"mobile_number"`
	FullName     string    `json:"full_name"`
	JoinedAt     time.Time `json:"joined_at"`
	Role         string    `json:"role"`
}

func NewUserResponse(userEntity entity.User) UserResponse {
	return UserResponse{
		Id:           userEntity.ID,
		MobileNumber: userEntity.MobileNumber,
		FullName:     userEntity.FullName,
		JoinedAt:     userEntity.CreatedAt,
	}
}

func NewUserListResponse(userEntityList []entity.User) []UserResponse {
	var finalResponse []UserResponse
	for _, user := range userEntityList {
		finalResponse = append(finalResponse, NewUserResponse(user))
	}
	return finalResponse
}
