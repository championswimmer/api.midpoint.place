package services

import (
	"context"

	placespb "cloud.google.com/go/maps/places/apiv1/placespb"
	"google.golang.org/genproto/googleapis/type/latlng"
)

type MockPlacesClient struct{}

func NewMockPlacesClient() *MockPlacesClient {
	return &MockPlacesClient{}
}

func (m *MockPlacesClient) SearchNearby(ctx context.Context, req *placespb.SearchNearbyRequest, opts ...interface{}) (*placespb.SearchNearbyResponse, error) {
	placeType := req.IncludedTypes[0]
	center := req.LocationRestriction.GetCircle().Center

	mockPlaces := []*placespb.Place{
		{
			Id:                    "mock_place_1",
			Name:                  "Mock " + placeType,
			DisplayName:           nil,
			ShortFormattedAddress: "123 Mock Street",
			GoogleMapsUri:         "https://maps.google.com/mock",
			Location: &latlng.LatLng{
				Latitude:  center.Latitude + 0.001,
				Longitude: center.Longitude + 0.001,
			},
			Rating: 4.5,
		},
		{
			Id:                    "mock_place_2",
			Name:                  "Another Mock " + placeType,
			DisplayName:           nil,
			ShortFormattedAddress: "456 Mock Avenue",
			GoogleMapsUri:         "https://maps.google.com/mock2",
			Location: &latlng.LatLng{
				Latitude:  center.Latitude - 0.001,
				Longitude: center.Longitude - 0.001,
			},
			Rating: 4.0,
		},
		{
			Id:                    "mock_place_3",
			Name:                  "Third Mock " + placeType,
			DisplayName:           nil,
			ShortFormattedAddress: "789 Mock Road",
			GoogleMapsUri:         "https://maps.google.com/mock3",
			Location: &latlng.LatLng{
				Latitude:  center.Latitude,
				Longitude: center.Longitude + 0.002,
			},
			Rating: 4.2,
		},
	}

	return &placespb.SearchNearbyResponse{
		Places: mockPlaces,
	}, nil
}

func (m *MockPlacesClient) Close() error {
	return nil
}
