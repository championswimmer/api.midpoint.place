package validators

import (
	"github.com/championswimmer/api.midpoint.place/src/dto"
	"github.com/gofiber/fiber/v2"
)

var mandatoryUserDtoFieldsError = &ValidationError{
	status:  fiber.StatusUnprocessableEntity,
	message: "Username and password are required",
}

func ValidateCreateUserRequest(dto *dto.CreateUserRequest) *ValidationError {
	if dto.Username == "" || dto.Password == "" {
		return mandatoryUserDtoFieldsError
	}
	return nil
}

func ValidateLoginUserRequest(dto *dto.LoginUserRequest) *ValidationError {
	if dto.Username == "" || dto.Password == "" {
		return mandatoryUserDtoFieldsError
	}
	return nil
}

func ValidateLocation(location dto.Location) *ValidationError {
	if location.Latitude < -90 || location.Latitude > 90 {
		return &ValidationError{
			status:  fiber.StatusBadRequest,
			message: "Invalid latitude value. Must be between -90 and 90",
		}
	}
	if location.Longitude < -180 || location.Longitude > 180 {
		return &ValidationError{
			status:  fiber.StatusBadRequest,
			message: "Invalid longitude value. Must be between -180 and 180",
		}
	}
	return nil
}
