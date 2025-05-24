package e2e

import (
	"bytes"
	"encoding/json"
	"io"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/championswimmer/api.midpoint.place/src/dto"
	"github.com/championswimmer/api.midpoint.place/tests"
	"github.com/gofiber/fiber/v2"
	"github.com/samber/lo"
	"github.com/stretchr/testify/assert"
)

func TestGroupMidpointUpdate(t *testing.T) {
	user1 := tests.TestUtil_CreateUser(t, "testuser1101", "testpassword1101")
	user2 := tests.TestUtil_CreateUser(t, "testuser2101", "testpassword2101")

	group1 := tests.TestUtil_CreateGroup(t, user1.Token, "Test Group 1101")

	testcases := []struct {
		name           string
		groupID        string
		userToken      string
		requestBody    *dto.GroupUserJoinRequest
		expectedStatus int
		checkResponse  func(t *testing.T, body []byte)
	}{
		{
			name:      "successful user1 joins group1",
			groupID:   group1.ID,
			userToken: user1.Token,
			requestBody: &dto.GroupUserJoinRequest{
				Location: dto.Location{
					Latitude:  51.5051821,
					Longitude: -0.2160895,
				},
			},
			expectedStatus: fiber.StatusAccepted,
			checkResponse: func(t *testing.T, body []byte) {
				var response dto.GroupResponse
				err := json.Unmarshal(body, &response)
				assert.NoError(t, err)
				assert.Equal(t, group1.ID, response.ID)
				assert.Equal(t, 51.5051821, response.MidpointLatitude)
				assert.Equal(t, -0.2160895, response.MidpointLongitude)
			},
		},
		{
			name:      "successful user2 joins group1",
			groupID:   group1.ID,
			userToken: user2.Token,
			requestBody: &dto.GroupUserJoinRequest{
				Location: dto.Location{
					Latitude:  51.4974653,
					Longitude: -0.1536909,
				},
			},
			expectedStatus: fiber.StatusAccepted,
			checkResponse: func(t *testing.T, body []byte) {
				var response dto.GroupResponse
				err := json.Unmarshal(body, &response)
				assert.NoError(t, err)
				assert.Equal(t, group1.ID, response.ID)
				assert.Equal(t, 51.5013237, response.MidpointLatitude)
				assert.Equal(t, -0.1848902, response.MidpointLongitude)

				assert.NotZero(t, len(response.Members))
				member1 := response.Members[0]
				assert.NotNil(t, member1.Username)
			},
		},
	}

	for _, tc := range testcases {
		t.Run(tc.name, func(t *testing.T) {
			jsonBody, err := json.Marshal(tc.requestBody)
			assert.NoError(t, err)
			req := httptest.NewRequest(fiber.MethodPut, "/v1/groups/"+tc.groupID+"/join", bytes.NewBuffer(jsonBody))
			req.Header.Set("Content-Type", "application/json")
			req.Header.Set("Authorization", "Bearer "+tc.userToken)
			resp := lo.Must(tests.App.Test(req, -1))
			assert.Equal(t, tc.expectedStatus, resp.StatusCode)

			// sleep for 20ms
			time.Sleep(20 * time.Millisecond)

			// fetch group details to check new midpoint

			req = httptest.NewRequest(fiber.MethodGet, "/v1/groups/"+tc.groupID+"?includeUsers=true", nil)
			req.Header.Set("Authorization", "Bearer "+tc.userToken)
			resp = lo.Must(tests.App.Test(req, -1))
			assert.Equal(t, fiber.StatusOK, resp.StatusCode)
			body := lo.Must(io.ReadAll(resp.Body))
			tc.checkResponse(t, body)
		})
	}
}
