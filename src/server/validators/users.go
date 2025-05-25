package validators

import (
	"regexp"

	"github.com/championswimmer/api.midpoint.place/src/dto"
	"github.com/gofiber/fiber/v2"
)

var mandatoryUserDtoFieldsError = &ValidationError{
	status:  fiber.StatusUnprocessableEntity,
	message: "Email and password are required",
}

var invalidEmailFormatError = &ValidationError{
	status:  fiber.StatusUnprocessableEntity,
	message: "Email format incorrect",
}

var invalidDisplayNameLengthError = &ValidationError{
	status:  fiber.StatusUnprocessableEntity,
	message: "Display name must be between 3 and 25 characters long",
}

// Basic email regex, can be improved for more strict validation
var emailRegex = regexp.MustCompile(`^[a-zA-Z0-9._%+\'-]+@[a-zA-Z0-9.-]+\.[a-zA-Z]{2,}$`)

func ValidateCreateUserRequest(dto *dto.CreateUserRequest) *ValidationError {
	if dto.Email == "" || dto.Password == "" {
		return mandatoryUserDtoFieldsError
	}
	if !emailRegex.MatchString(dto.Email) {
		return invalidEmailFormatError
	}
	if len(dto.DisplayName) < 3 || len(dto.DisplayName) > 25 {
		return invalidDisplayNameLengthError
	}
	return nil
}

func ValidateLoginUserRequest(dto *dto.LoginUserRequest) *ValidationError {
	if dto.Email == "" || dto.Password == "" {
		return mandatoryUserDtoFieldsError
	}
	if !emailRegex.MatchString(dto.Email) {
		return invalidEmailFormatError
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
