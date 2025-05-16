package services

import "googlemaps.github.io/maps"

type GoogleMapsPlacesService struct {
	googleMapsClient *maps.Client
}

func NewGoogleMapsPlacesService() *GoogleMapsPlacesService {
	return &GoogleMapsPlacesService{googleMapsClient: GetGoogleMapsClient()}
}
