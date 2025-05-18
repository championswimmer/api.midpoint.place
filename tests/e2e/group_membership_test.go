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

func TestGroupsRoute_MembershipOperations(t *testing.T) {
	// Create two users
	user1 := tests.TestUtil_CreateUser(t, "testuser101", "testpassword101")
	user2 := tests.TestUtil_CreateUser(t, "testuser201", "testpassword201")

	group1 := tests.TestUtil_CreateGroup(t, user1.Token, "Test Group 1")
	group2 := tests.TestUtil_CreateGroup(t, user1.Token, "Test Group 2")

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
				assert.Equal(t, user2.ID, groupUserResp.UserID)
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
				assert.Equal(t, user2.ID, groupUserResp.UserID)
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
			// sleep for 20ms
			time.Sleep(20 * time.Millisecond)
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

func TestGroupsRoute_AddUsersToGroupAndVerify(t *testing.T) {
	// Create two users
	user1 := tests.TestUtil_CreateUser(t, "testuser301", "testpassword301")
	user2 := tests.TestUtil_CreateUser(t, "testuser401", "testpassword401")

	// Create a group
	group := tests.TestUtil_CreateGroup(t, user1.Token, "Test Group 3")

	// Add user2 to the group
	joinRequestBody := []byte(`{
		"latitude": 12.971645,
		"longitude": 77.594562
	}`)
	req := httptest.NewRequest("PUT", "/v1/groups/"+group.ID+"/join", bytes.NewBuffer(joinRequestBody))
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+user2.Token)

	resp := lo.Must(tests.App.Test(req, -1))
	assert.Equal(t, fiber.StatusAccepted, resp.StatusCode)

	// Verify that user2 is included in the group members
	req = httptest.NewRequest("GET", "/v1/groups/"+group.ID+"?includeUsers=true", nil)
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+user1.Token)

	resp = lo.Must(tests.App.Test(req, -1))
	assert.Equal(t, fiber.StatusOK, resp.StatusCode)

	var groupResp dto.GroupResponse
	body := lo.Must(io.ReadAll(resp.Body))
	err := json.Unmarshal(body, &groupResp)
	assert.NoError(t, err)

	assert.Len(t, groupResp.Members, 1)
	assert.Equal(t, user2.ID, groupResp.Members[0].UserID)
}
