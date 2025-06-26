package validators

import (
	"regexp"
	"strconv"

	"github.com/championswimmer/api.midpoint.place/src/config"
	"github.com/championswimmer/api.midpoint.place/src/dto"
	"github.com/gofiber/fiber/v2"
	"github.com/umahmood/haversine"
)

const MAX_DISTANCE_KM = 100

func validateName(name string) *ValidationError {
	if len(name) > 100 || len(name) < 1 {
		return &ValidationError{
			status:  fiber.StatusUnprocessableEntity,
			message: "Name must be between 1 and 100 characters",
		}
	}
	return nil
}

func validateSecret(secret string) *ValidationError {
	if secret == "" {
		return nil
	}
	if matched, _ := regexp.MatchString(`^\d{6}$`, secret); !matched {
		return &ValidationError{
			status:  fiber.StatusUnprocessableEntity,
			message: "Secret must be a 6-digit string or empty",
		}
	}
	return nil
}

func validateRadius(radius int) *ValidationError {
	if radius < 0 {
		return &ValidationError{
			status:  fiber.StatusUnprocessableEntity,
			message: "Radius must be a positive integer",
		}
	}
	return nil
}

func validatePlaceTypes(placeTypes []config.PlaceType) *ValidationError {
	if len(placeTypes) == 0 {
		return nil
	}
	if len(placeTypes) > 10 {
		return &ValidationError{
			status:  fiber.StatusUnprocessableEntity,
			message: "Cannot select more than 10 place types",
		}
    }
    // Validate each place type is valid and unique
    validTypes := make(map[config.PlaceType]bool)
    for _, validType := range config.AllPlaceTypes {
        validTypes[validType] = true
    }
    
    seenTypes := make(map[config.PlaceType]bool)
    for _, placeType := range placeTypes {
        if !validTypes[placeType] {
            return &ValidationError{
                status:  fiber.StatusUnprocessableEntity,
                message: "Invalid place type: " + string(placeType),
            }
        }
        if seenTypes[placeType] {
            return &ValidationError{
                status:  fiber.StatusUnprocessableEntity,
                message: "Duplicate place type: " + string(placeType),
            }
        }
        seenTypes[placeType] = true
    }
	return nil
}

func ValidateCreateGroupRequest(req *dto.CreateGroupRequest) *ValidationError {
	if err := validateName(req.Name); err != nil {
		return err
	}
	if err := validateSecret(req.Secret); err != nil {
		return err
	}
	if err := validateRadius(req.Radius); err != nil {
		return err
	}
	if err := validatePlaceTypes(req.PlaceTypes); err != nil {
		return err
	}
	// Add any other specific validations for CreateGroupRequest
	return nil
}

func ValidateUpdateGroupRequest(req *dto.UpdateGroupRequest) *ValidationError {
	if req.Name != "" { // Only validate if provided for update
		if err := validateName(req.Name); err != nil {
			return err
		}
	}
	if req.Secret != "" { // Only validate if provided for update, or if explicitly set to be cleared (empty string)
		if err := validateSecret(req.Secret); err != nil {
			return err
		}
	}
	if req.Radius > 0 { // Only validate if provided for update and positive
		if err := validateRadius(req.Radius); err != nil {
			return err
		}
	}
	if len(req.PlaceTypes) > 0 { // Only validate if provided for update
		if err := validatePlaceTypes(req.PlaceTypes); err != nil {
			return err
		}
	}
	// Add any other specific validations for UpdateGroupRequest
	return nil
}

// TODO: there are no e2e tests for this yet.
func ValidateLocationProximity(loc1 dto.Location, loc2 dto.Location) *ValidationError {
	coord1 := haversine.Coord{
		Lat: loc1.Latitude,
		Lon: loc1.Longitude,
	}
	coord2 := haversine.Coord{
		Lat: loc2.Latitude,
		Lon: loc2.Longitude,
	}
	_, kmDistance := haversine.Distance(coord1, coord2)

	if kmDistance > float64(MAX_DISTANCE_KM) {
		return &ValidationError{
			status:  fiber.StatusUnprocessableEntity,
			message: "Cannot join group from too far away. Limit is " + strconv.Itoa(MAX_DISTANCE_KM) + "km",
		}
	}
	return nil
}
