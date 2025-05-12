package e2e

import (
	"bytes"
	"encoding/json"
	"io"
	"net/http/httptest"
	"testing"

	"github.com/championswimmer/api.midpoint.place/src/dto"
	"github.com/championswimmer/api.midpoint.place/tests"
	"github.com/gofiber/fiber/v2"
	"github.com/samber/lo"
	"github.com/stretchr/testify/assert"
)

func TestUsersRoute_LoginUser(t *testing.T) {
	// First register a user
	tests.TestUtil_CreateUser(t, "testuser111", "testpassword111")

	// Test cases
	testCases := []struct {
		name           string
		requestBody    []byte
		expectedStatus int
		checkResponse  func(t *testing.T, body []byte)
	}{
		{
			name:           "successful login",
			requestBody:    []byte(`{"username": "testuser111", "password": "testpassword111"}`),
			expectedStatus: fiber.StatusOK,
			checkResponse: func(t *testing.T, body []byte) {
				var userResp dto.UserResponse
				err := json.Unmarshal(body, &userResp)
				assert.NoError(t, err)
				assert.Equal(t, "testuser111", userResp.Username)
				assert.NotEmpty(t, userResp.Token)
				assert.NotEmpty(t, userResp.Id)
			},
		},
		{
			name:           "wrong password",
			requestBody:    []byte(`{"username": "testuser111", "password": "wrongpassword222"}`),
			expectedStatus: fiber.StatusUnauthorized,
			checkResponse: func(t *testing.T, body []byte) {
				var errResp dto.ErrorResponse
				err := json.Unmarshal(body, &errResp)
				assert.NoError(t, err)
				assert.Equal(t, "Invalid password", errResp.Message)
			},
		},
		{
			name:           "user not found",
			requestBody:    []byte(`{"username": "nonexistentuser", "password": "testpassword"}`),
			expectedStatus: fiber.StatusNotFound,
			checkResponse: func(t *testing.T, body []byte) {
				var errResp dto.ErrorResponse
				err := json.Unmarshal(body, &errResp)
				assert.NoError(t, err)
				assert.Equal(t, "User not found", errResp.Message)
			},
		},
		{
			name:           "invalid request body",
			requestBody:    []byte(`{"invalid": "json"`),
			expectedStatus: fiber.StatusBadRequest,
			checkResponse: func(t *testing.T, body []byte) {
				var errResp dto.ErrorResponse
				err := json.Unmarshal(body, &errResp)
				assert.NoError(t, err)
				assert.Contains(t, errResp.Message, "not valid")
			},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			req := httptest.NewRequest("POST", "/v1/users/login", bytes.NewBuffer(tc.requestBody))
			req.Header.Set("Content-Type", "application/json")

			resp := lo.Must(tests.App.Test(req, -1))
			assert.Equal(t, tc.expectedStatus, resp.StatusCode)
			body := []byte{}
			if resp.Body != nil {
				body, _ = io.ReadAll(resp.Body)
			}
			tc.checkResponse(t, body)
		})
	}
}
