package e2e

import (
	"bytes"
	"encoding/json"
	"io"
	"net/http/httptest"
	"testing"

	"github.com/championswimmer/api.midpoint.place/src/config"
	"github.com/championswimmer/api.midpoint.place/src/dto"
	"github.com/championswimmer/api.midpoint.place/tests"
	"github.com/gofiber/fiber/v2"
	"github.com/samber/lo"
	"github.com/stretchr/testify/assert"
)

func TestGroupsRoute_CreateGroup(t *testing.T) {
	// First register a user
	createdUser := tests.TestUtil_CreateUser(t, "testuser222@test.com", "testpassword222")
	var createdGroup dto.GroupResponse

	testcases := []struct {
		name           string
		requestBody    []byte
		expectedStatus int
		checkResponse  func(t *testing.T, body []byte)
	}{
		{
			name: "successful group creation",
			requestBody: []byte(`{
				"name": "Test Group",
				"type": "public",
				"radius": 1200
			}`),
			expectedStatus: fiber.StatusCreated,
			checkResponse: func(t *testing.T, body []byte) {
				var groupResp dto.GroupResponse
				err := json.Unmarshal(body, &groupResp)
				assert.NoError(t, err)
				assert.Equal(t, "Test Group", groupResp.Name)
				assert.Equal(t, config.GroupTypePublic, groupResp.Type)
				assert.Equal(t, 1200, groupResp.Radius)
				assert.Equal(t, createdUser.ID, groupResp.Creator.ID)
				assert.Equal(t, createdUser.DisplayName, groupResp.Creator.DisplayName)
				createdGroup = groupResp
			},
		},
		{
			name: "missing required fields",
			requestBody: []byte(`{
				"radius": 1200
			}`),
			expectedStatus: fiber.StatusUnprocessableEntity,
			checkResponse: func(t *testing.T, body []byte) {
				var errResp dto.ErrorResponse
				err := json.Unmarshal(body, &errResp)
				assert.NoError(t, err)
				assert.Contains(t, errResp.Message, "Name must be")
			},
		},
	}

	for _, tc := range testcases {
		t.Run(tc.name, func(t *testing.T) {
			req := httptest.NewRequest("POST", "/v1/groups", bytes.NewBuffer(tc.requestBody))
			req.Header.Set("Content-Type", "application/json")
			req.Header.Set("Authorization", "Bearer "+createdUser.Token)

			resp := lo.Must(tests.App.Test(req, -1))
			assert.Equal(t, tc.expectedStatus, resp.StatusCode)
			body := []byte{}
			if resp.Body != nil {
				body = lo.Must(io.ReadAll(resp.Body))
			}
			tc.checkResponse(t, body)
		})
	}

	// Get the group
	t.Run("get group", func(t *testing.T) {
		req := httptest.NewRequest("GET", "/v1/groups/"+createdGroup.ID, nil)
		req.Header.Set("Content-Type", "application/json")
		req.Header.Set("Authorization", "Bearer "+createdUser.Token)

		resp := lo.Must(tests.App.Test(req, -1))
		assert.Equal(t, fiber.StatusOK, resp.StatusCode)
	})
}
