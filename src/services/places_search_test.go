package services

import (
	"testing"

	"github.com/championswimmer/api.midpoint.place/src/config"
	"github.com/championswimmer/api.midpoint.place/src/dto"
	"github.com/championswimmer/api.midpoint.place/src/utils/applogger"
	"github.com/samber/lo"
	"github.com/stretchr/testify/assert"
)

func TestPlaceSearchService_NearbyPlaces(t *testing.T) {
	placeSearchService := NewPlaceSearchService()

	places, err := placeSearchService.NearbyPlaces(dto.Location{
		Latitude:  28.6139,
		Longitude: 77.2090,
	}, 1000, config.PlaceTypePark)

	assert.NoError(t, err)
	assert.NotEmpty(t, places)
	assert.Equal(t, 3, len(places))
	lo.ForEach(places, func(place dto.Place, _ int) {
		applogger.Info("place", place)
	})
}
