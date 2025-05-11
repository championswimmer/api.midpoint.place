package security

import (
	"testing"
	"time"

	"github.com/championswimmer/api.midpoint.place/src/db/models"
	"github.com/golang-jwt/jwt/v5"
	"github.com/stretchr/testify/assert"
)

func TestCreateJWTFromUser(t *testing.T) {
	// Test case 1: Create JWT for a valid user
	user := &models.User{ID: 123}
	token := CreateJWTFromUser(user)

	// Verify token is not empty
	assert.NotEmpty(t, token)

	// Parse the token to verify its contents
	claims := jwt.MapClaims{}
	parsedToken, err := jwt.ParseWithClaims(token, claims, func(token *jwt.Token) (interface{}, error) {
		return []byte(SigningKey), nil
	})

	assert.NoError(t, err)
	assert.True(t, parsedToken.Valid)
	assert.Equal(t, "123", claims["sub"])
	assert.NotNil(t, claims["iat"])
	assert.NotNil(t, claims["exp"])

	// Verify expiration time
	expTime := time.Unix(int64(claims["exp"].(float64)), 0)
	expectedExpTime := time.Now().AddDate(0, 0, ExpirationDays)
	assert.True(t, expTime.Sub(expectedExpTime) < time.Hour) // Allow 1 hour difference for test execution time
}

func TestValidateAndParseJWT(t *testing.T) {
	// Test case 1: Valid token
	user := &models.User{ID: 456}
	validToken := CreateJWTFromUser(user)
	parsedUser, err := ValidateAndParseJWT(validToken)
	assert.NoError(t, err)
	assert.Equal(t, user.ID, parsedUser.ID)

	// Test case 2: Invalid token (malformed)
	invalidToken := "invalid.token.string"
	_, err = ValidateAndParseJWT(invalidToken)
	assert.Error(t, err)

	// Test case 3: Expired token
	expiredToken := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"sub": "789",
		"iat": jwt.NewNumericDate(time.Now().AddDate(0, 0, -1)), // issued 1 day ago
		"exp": jwt.NewNumericDate(time.Now().AddDate(0, 0, -1)), // expired 1 day ago
	})
	expiredTokenString, _ := expiredToken.SignedString([]byte(SigningKey))
	_, err = ValidateAndParseJWT(expiredTokenString)
	assert.Error(t, err)

	// Test case 4: Token with invalid signature
	invalidSignatureToken := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"sub": "789",
		"iat": jwt.NewNumericDate(time.Now()),
		"exp": jwt.NewNumericDate(time.Now().AddDate(0, 0, ExpirationDays)),
	})
	invalidSignatureTokenString, _ := invalidSignatureToken.SignedString([]byte("wrong-signing-key"))
	_, err = ValidateAndParseJWT(invalidSignatureTokenString)
	assert.Error(t, err)
}
