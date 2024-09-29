package models

import "github.com/TheAmirhosssein/room-reservation-api/internal/entity"

type (
	City struct {
		Title string `json:"title" binding:"required"`
	}
	CityResponse struct {
		Id         uint   `json:"id"`
		Title      string `json:"title"`
		StateId    uint   `json:"state_id"`
		StateTitle string `json:"state_title"`
	}
)

func NewCityResponse(city entity.City) CityResponse {
	return CityResponse{
		Id:         city.ID,
		Title:      city.Title,
		StateId:    city.StateID,
		StateTitle: city.State.Title,
	}
}
