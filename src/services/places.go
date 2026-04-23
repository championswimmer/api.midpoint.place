package services

import (
	"context"
	"sync"

	places "cloud.google.com/go/maps/places/apiv1"
	placespb "cloud.google.com/go/maps/places/apiv1/placespb"
	"github.com/championswimmer/api.midpoint.place/src/config"
	"github.com/championswimmer/api.midpoint.place/src/utils/applogger"
	"github.com/samber/lo"
	"google.golang.org/api/option"
)

type PlacesClientInterface interface {
	SearchNearby(ctx context.Context, req *placespb.SearchNearbyRequest, opts ...interface{}) (*placespb.SearchNearbyResponse, error)
	Close() error
}

type RealPlacesClient struct {
	client *places.Client
}

func (r *RealPlacesClient) SearchNearby(ctx context.Context, req *placespb.SearchNearbyRequest, opts ...interface{}) (*placespb.SearchNearbyResponse, error) {
	return r.client.SearchNearby(ctx, req)
}

func (r *RealPlacesClient) Close() error {
	return r.client.Close()
}

var googlePlacesClient PlacesClientInterface
var googlePlacesClientOnce sync.Once

func GetGooglePlacesClient() PlacesClientInterface {
	googlePlacesClientOnce.Do(func() {
		if config.Env == "production" || !config.UseMockPlaces {
			if config.GoogleMapsAPIKey == "" {
				if config.Env == "production" {
					applogger.Fatal("Google Maps API key is required in production")
				}
				applogger.Warn("No Google Maps API key provided, using mock Places client")
				googlePlacesClient = NewMockPlacesClient()
				return
			}
			applogger.Info("Using real Places client")
			clientOpts := []option.ClientOption{option.WithAPIKey(config.GoogleMapsAPIKey)}
			realClient := lo.Must(places.NewClient(context.Background(), clientOpts...))
			googlePlacesClient = &RealPlacesClient{client: realClient}
		} else {
			applogger.Info("Using mock Places client (USE_MOCK_PLACES=true)")
			googlePlacesClient = NewMockPlacesClient()
		}
	})

	return googlePlacesClient
}
