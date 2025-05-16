package services

import (
	"context"
	"sync"

	places "cloud.google.com/go/maps/places/apiv1"
	"github.com/championswimmer/api.midpoint.place/src/config"
	"github.com/samber/lo"
	"google.golang.org/api/option"
)

var googlePlacesClient *places.Client
var googlePlacesClientOnce sync.Once

func GetGooglePlacesClient() *places.Client {

	googlePlacesClientOnce.Do(func() {
		clientOpts := []option.ClientOption{option.WithAPIKey(config.GoogleMapsAPIKey)}
		googlePlacesClient = lo.Must(places.NewClient(context.Background(), clientOpts...))
	})

	return googlePlacesClient
}
