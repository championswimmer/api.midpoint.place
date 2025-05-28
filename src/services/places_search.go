package services

import (
	"context"

	placespb "cloud.google.com/go/maps/places/apiv1/placespb"
	"github.com/championswimmer/api.midpoint.place/src/config"
	"github.com/championswimmer/api.midpoint.place/src/dto"
	"github.com/championswimmer/api.midpoint.place/src/utils/applogger"
	"github.com/samber/lo"
	"google.golang.org/genproto/googleapis/type/latlng"
	"google.golang.org/grpc/metadata"
)

type PlaceSearchService struct {
	placesClient PlacesClientInterface
}

func NewPlaceSearchService() *PlaceSearchService {

	return &PlaceSearchService{
		placesClient: GetGooglePlacesClient(),
	}
}

const fieldsToRequest = "places.id,places.displayName,places.formattedAddress,places.googleMapsUri,places.primaryTypeDisplayName,places.rating,places.location,places.shortFormattedAddress"

func (s *PlaceSearchService) NearbyPlaces(location dto.Location, radius int, placeType config.PlaceType) ([]dto.Place, error) {

	ctx := metadata.AppendToOutgoingContext(context.Background(), "x-goog-fieldmask", fieldsToRequest)

	searchResp, err := s.placesClient.SearchNearby(
		ctx,
		&placespb.SearchNearbyRequest{
			IncludedTypes:  _getIncludedTypes(placeType),
			MaxResultCount: 3, // TODO: fetch from config
			LocationRestriction: &placespb.SearchNearbyRequest_LocationRestriction{
				Type: &placespb.SearchNearbyRequest_LocationRestriction_Circle{
					Circle: &placespb.Circle{
						Center: &latlng.LatLng{
							Latitude:  location.Latitude,
							Longitude: location.Longitude,
						},
						Radius: float64(radius),
					},
				},
			},
		},
	)
	if err != nil {
		applogger.Error("Error searching for nearby places", location, "with radius", radius, "and place type", placeType, err)
		return nil, err
	}

	places := lo.Map(searchResp.Places, func(place *placespb.Place, _ int) dto.Place {
		return _googlePlaceToPlaceDto(place, placeType)
	})

	return places, nil
}

func _getIncludedTypes(placeType config.PlaceType) []string {
	switch placeType {
	case config.PlaceTypeRestaurant:
		return []string{"restaurant", "bar_and_grill", "food_court"}
	case config.PlaceTypeBar:
		return []string{"bar", "pub"}
	case config.PlaceTypeCafe:
		return []string{"cafe", "coffee_shop"}
	case config.PlaceTypePark:
		return []string{"park", "garden"}
	}
	return []string{}
}

func _googlePlaceToPlaceDto(googlePlace *placespb.Place, placeType config.PlaceType) dto.Place {

	place := dto.Place{
		Location: dto.Location{
			Latitude:  googlePlace.Location.Latitude,
			Longitude: googlePlace.Location.Longitude,
		},
		Id:      googlePlace.Id,
		Name:    googlePlace.DisplayName.Text,
		Address: googlePlace.ShortFormattedAddress,
		MapURI:  googlePlace.GoogleMapsUri,
		Type:    placeType,
		Rating:  googlePlace.Rating,
	}
	return place
}
