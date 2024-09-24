package models

import "github.com/TheAmirhosssein/room-reservation-api/internal/entity"

type (
	State struct {
		Title string `json:"title" binding:"required"`
	}
	StateResponse struct {
		Id    uint   `json:"id"`
		Title string `json:"title"`
	}
)

func NewStateResponse(state entity.State) StateResponse {
	return StateResponse{
		Id:    state.ID,
		Title: state.Title,
	}
}
