package services

import (
	"context"
	"os"
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
		// In test environment without GOOGLE_MAPS_API_KEY, skip initialization
		// The client will be nil and callers should handle this gracefully
		if os.Getenv("ENV") == "test" && config.GoogleMapsAPIKey == "" {
			return
		}
		clientOpts := []option.ClientOption{option.WithAPIKey(config.GoogleMapsAPIKey)}
		googlePlacesClient = lo.Must(places.NewClient(context.Background(), clientOpts...))
	})

	return googlePlacesClient
}
