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
