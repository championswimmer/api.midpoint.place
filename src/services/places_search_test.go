package services

import (
	"context"
	"testing"

	"github.com/googleapis/gax-go/v2"
	placespb "cloud.google.com/go/maps/places/apiv1/placespb"
	"github.com/championswimmer/api.midpoint.place/src/config"
	"github.com/championswimmer/api.midpoint.place/src/dto"
	"github.com/stretchr/testify/assert"
	"google.golang.org/genproto/googleapis/type/latlng"
	"google.golang.org/genproto/googleapis/type/localized_text"
	"google.golang.org/grpc/metadata"
)

// fakePlacesClient implements placesNearbyClient for testing
type fakePlacesClient struct {
	// Store the last request received for assertions
	lastRequest *placespb.SearchNearbyRequest
	lastContext context.Context
}

func (f *fakePlacesClient) SearchNearby(ctx context.Context, req *placespb.SearchNearbyRequest, opts ...gax.CallOption) (*placespb.SearchNearbyResponse, error) {
	f.lastRequest = req
	f.lastContext = ctx

	// Return a fake response with 3 places
	return &placespb.SearchNearbyResponse{
		Places: []*placespb.Place{
			{
				Id: "place1",
				DisplayName: &localized_text.LocalizedText{
					Text: "Test Park 1",
				},
				Location: &latlng.LatLng{
					Latitude:  28.6140,
					Longitude: 77.2091,
				},
				ShortFormattedAddress: "Test Address 1",
				GoogleMapsUri:         "https://maps.google.com/place1",
				Rating:                4.5,
			},
			{
				Id: "place2",
				DisplayName: &localized_text.LocalizedText{
					Text: "Test Park 2",
				},
				Location: &latlng.LatLng{
					Latitude:  28.6141,
					Longitude: 77.2092,
				},
				ShortFormattedAddress: "Test Address 2",
				GoogleMapsUri:         "https://maps.google.com/place2",
				Rating:                4.0,
			},
			{
				Id: "place3",
				DisplayName: &localized_text.LocalizedText{
					Text: "Test Park 3",
				},
				Location: &latlng.LatLng{
					Latitude:  28.6142,
					Longitude: 77.2093,
				},
				ShortFormattedAddress: "Test Address 3",
				GoogleMapsUri:         "https://maps.google.com/place3",
				Rating:                3.5,
			},
		},
	}, nil
}

func TestPlaceSearchService_NearbyPlaces(t *testing.T) {
	fakeClient := &fakePlacesClient{}
	placeSearchService := NewPlaceSearchServiceWithClient(fakeClient)

	testLocation := dto.Location{
		Latitude:  28.6139,
		Longitude: 77.2090,
	}
	testRadius := 1000
	testPlaceType := config.PlaceTypePark

	places, err := placeSearchService.NearbyPlaces(testLocation, testRadius, testPlaceType)

	// Verify no error
	assert.NoError(t, err)

	// Verify we got 3 places
	assert.NotEmpty(t, places)
	assert.Equal(t, 3, len(places))

	// Verify the request was properly constructed
	assert.NotNil(t, fakeClient.lastRequest)
	assert.Equal(t, int32(3), fakeClient.lastRequest.MaxResultCount)
	assert.Equal(t, []string{"park", "garden"}, fakeClient.lastRequest.IncludedTypes)

	// Verify location restriction was set correctly
	circle := fakeClient.lastRequest.GetLocationRestriction().GetCircle()
	assert.NotNil(t, circle)
	assert.Equal(t, testLocation.Latitude, circle.Center.Latitude)
	assert.Equal(t, testLocation.Longitude, circle.Center.Longitude)
	assert.Equal(t, float64(testRadius), circle.Radius)

	// Verify x-goog-fieldmask header was set in context
	md, ok := metadata.FromOutgoingContext(fakeClient.lastContext)
	assert.True(t, ok)
	fieldMask := md.Get("x-goog-fieldmask")
	assert.NotEmpty(t, fieldMask)
	assert.Equal(t, fieldsToRequest, fieldMask[0])

	// Verify returned places have correct type and data
	for _, place := range places {
		assert.Equal(t, testPlaceType, place.Type)
		assert.NotEmpty(t, place.Id)
		assert.NotEmpty(t, place.Name)
	}
}
