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

func TestWaitlistSignup(t *testing.T) {
	testCases := []struct {
		name           string
		requestBody    []byte
		expectedStatus int
		checkResponse  func(t *testing.T, body []byte)
	}{
		{
			name: "successful signup",
			requestBody: []byte(`{
				"email": "test@example.com"
			}`),
			expectedStatus: fiber.StatusCreated,
			checkResponse: func(t *testing.T, body []byte) {
				var signupResp dto.WaitlistSignupResponse
				err := json.Unmarshal(body, &signupResp)
				assert.NoError(t, err)
				assert.Equal(t, "Successfully added to waitlist", signupResp.Message)
			},
		},
		{
			name: "duplicate email",
			requestBody: []byte(`{
				"email": "test@example.com"
			}`),
			expectedStatus: fiber.StatusConflict,
			checkResponse: func(t *testing.T, body []byte) {
				var errResp dto.ErrorResponse
				err := json.Unmarshal(body, &errResp)
				assert.NoError(t, err)
				assert.Equal(t, "Email already signed up", errResp.Message)
			},
		},
		{
			name: "invalid email format",
			requestBody: []byte(`{
				"email": "invalid-email"
			}`),
			expectedStatus: fiber.StatusBadRequest,
			checkResponse: func(t *testing.T, body []byte) {
				var errResp dto.ErrorResponse
				err := json.Unmarshal(body, &errResp)
				assert.NoError(t, err)
				assert.Equal(t, "Invalid email format", errResp.Message)
			},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			req := httptest.NewRequest("POST", "/v1/waitlist/signup", bytes.NewBuffer(tc.requestBody))
			req.Header.Set("Content-Type", "application/json")

			resp := lo.Must(tests.App.Test(req, -1))
			assert.Equal(t, tc.expectedStatus, resp.StatusCode)
			body := []byte{}
			if resp.Body != nil {
				body = lo.Must(io.ReadAll(resp.Body))
			}
			tc.checkResponse(t, body)
		})
	}
}
