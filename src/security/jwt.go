package security

import (
	"errors"
	"strconv"
	"time"

	"github.com/championswimmer/api.midpoint.place/src/config"
	"github.com/championswimmer/api.midpoint.place/src/db/models"
	"github.com/championswimmer/api.midpoint.place/src/utils/applogger"
	"github.com/golang-jwt/jwt/v5"
	"github.com/samber/lo"
)

var SigningKey = config.JWTSigningKey
var SigningMethod = jwt.SigningMethodHS256
var ExpirationDays = config.JWTExpirationDays

// CreateJWTFromUser creates a JWT token from a user
func CreateJWTFromUser(user *models.User) string {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"sub": strconv.FormatUint(uint64(user.ID), 10),
		"iat": jwt.NewNumericDate(time.Now()),
		"exp": jwt.NewNumericDate(time.Now().AddDate(0, 0, ExpirationDays)),
	})

	applogger.Info(lo.Must(token.SignedString([]byte(SigningKey))))

	return lo.Must(token.SignedString([]byte(SigningKey)))
}

// ValidateAndParseJWT validates and parses a JWT token
// and returns a user object with the user ID (this is not a fully populated user object)
func ValidateAndParseJWT(tokenString string) (*models.User, error) {

	claims := jwt.MapClaims{}

	token, err := jwt.ParseWithClaims(tokenString, claims, func(token *jwt.Token) (interface{}, error) {
		return []byte(SigningKey), nil
	})

	if err != nil {
		return nil, err
	}

	if !token.Valid {
		return nil, errors.New("invalid token") // TODO: create custom error later
	}

	userID, err := strconv.ParseUint(claims["sub"].(string), 10, 64)
	if err != nil {
		return nil, err // TODO: create custom error later
	}

	return &models.User{
		ID: uint(userID),
	}, nil
}
