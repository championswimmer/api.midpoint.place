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

func TestGroupsRoute_MembershipOperations(t *testing.T) {
	// Create two users
	user1 := tests.TestUtil_CreateUser(t, "testuser101", "testpassword101")
	user2 := tests.TestUtil_CreateUser(t, "testuser201", "testpassword201")

	// Create two groups
	createGroup1Req := []byte(`{
		"name": "Group 1",
		"radius": 100
	}`)

	createGroup2Req := []byte(`{
		"name": "Group 2",
		"radius": 100
	}`)

	var group1, group2 *dto.GroupResponse

	req := httptest.NewRequest("POST", "/v1/groups", bytes.NewBuffer(createGroup1Req))
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+user1.Token)

	resp := lo.Must(tests.App.Test(req, -1))
	assert.Equal(t, fiber.StatusCreated, resp.StatusCode)
	body := lo.Must(io.ReadAll(resp.Body))
	err := json.Unmarshal(body, &group1)
	assert.NoError(t, err)

	req = httptest.NewRequest("POST", "/v1/groups", bytes.NewBuffer(createGroup2Req))
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+user1.Token)

	resp = lo.Must(tests.App.Test(req, -1))
	assert.Equal(t, fiber.StatusCreated, resp.StatusCode)
	body = lo.Must(io.ReadAll(resp.Body))
	err = json.Unmarshal(body, &group2)
	assert.NoError(t, err)

	testcases := []struct {
		name           string
		groupID        string
		requestBody    []byte
		method         string
		expectedStatus int
		checkResponse  func(t *testing.T, body []byte)
	}{
		{
			name:    "join group with group id",
			groupID: group1.ID,
			requestBody: []byte(`{
				"latitude": 12.971645,
				"longitude": 77.594562
			}`),
			method:         "PUT",
			expectedStatus: fiber.StatusAccepted,
			checkResponse: func(t *testing.T, body []byte) {
				var groupUserResp dto.GroupUserResponse
				err := json.Unmarshal(body, &groupUserResp)
				assert.NoError(t, err)
				assert.Equal(t, group1.ID, groupUserResp.GroupID)
				assert.Equal(t, user2.Id, groupUserResp.UserID)
			},
		},
		{
			name:    "join group with group code",
			groupID: group2.Code,
			requestBody: []byte(`{
				"latitude": 12.971645,
				"longitude": 77.594562
			}`),
			method:         "PUT",
			expectedStatus: fiber.StatusAccepted,
			checkResponse: func(t *testing.T, body []byte) {
				var groupUserResp dto.GroupUserResponse
				err := json.Unmarshal(body, &groupUserResp)
				assert.NoError(t, err)
				assert.Equal(t, group2.ID, groupUserResp.GroupID)
				assert.Equal(t, user2.Id, groupUserResp.UserID)
			},
		},
		{
			name:           "leave group with group id",
			groupID:        group1.ID,
			requestBody:    []byte(``),
			method:         "DELETE",
			expectedStatus: fiber.StatusAccepted,
			checkResponse: func(t *testing.T, body []byte) {

			},
		},
	}

	for _, tc := range testcases {
		t.Run(tc.name, func(t *testing.T) {
			req := httptest.NewRequest(tc.method, "/v1/groups/"+tc.groupID+"/join", bytes.NewBuffer(tc.requestBody))
			req.Header.Set("Content-Type", "application/json")
			req.Header.Set("Authorization", "Bearer "+user2.Token)

			resp := lo.Must(tests.App.Test(req, -1))
			assert.Equal(t, tc.expectedStatus, resp.StatusCode)
			body := lo.Must(io.ReadAll(resp.Body))
			tc.checkResponse(t, body)
		})
	}

}
