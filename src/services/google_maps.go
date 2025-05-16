package services

import (
	"sync"

	"github.com/championswimmer/api.midpoint.place/src/config"
	"github.com/samber/lo"
	"googlemaps.github.io/maps"
)

var googleMapsClient *maps.Client
var googleMapsClientOnce sync.Once

func GetGoogleMapsClient() *maps.Client {

	googleMapsClientOnce.Do(func() {
		googleMapsClient = lo.Must(maps.NewClient(maps.WithAPIKey(config.GoogleMapsAPIKey)))
	})

	return googleMapsClient
}
