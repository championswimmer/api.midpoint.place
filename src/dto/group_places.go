package dto

import "github.com/championswimmer/api.midpoint.place/src/config"

// GroupPlacesAddRequest represents the request to add places to a group
type GroupPlacesAddRequest struct {
	Places []Place `json:"places" validate:"required,min=1"`
}

// GroupPlaceResponse represents a place in a group
type GroupPlaceResponse struct {
	ID        string           `json:"id"`
	GroupID   string           `json:"group_id"`
	PlaceID   string           `json:"place_id"`
	Name      string           `json:"name"`
	Address   string           `json:"address"`
	Type      config.PlaceType `json:"type"`
	Rating    float64          `json:"rating"`
	MapURI    string           `json:"map_uri"`
	Latitude  float64          `json:"latitude"`
	Longitude float64          `json:"longitude"`
}
