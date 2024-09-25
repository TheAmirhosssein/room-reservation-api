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

func NewStateListResponse(states []entity.State) []StateResponse {
	var finalResponse []StateResponse
	for _, state := range states {
		finalResponse = append(finalResponse, NewStateResponse(state))
	}
	return finalResponse
}
