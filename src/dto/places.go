package dto

import "github.com/championswimmer/api.midpoint.place/src/config"

type Place struct {
	Location
	Id      string           `json:"id" validate:"required"`
	Name    string           `json:"name" validate:"required"`
	Address string           `json:"address" validate:"required"`
	Type    config.PlaceType `json:"type" validate:"required"`
	Rating  float64          `json:"rating" validate:"required,min=1.0,max=5.0"`
	MapURI  string           `json:"map_uri" validate:"required"`
}
