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

func TestGroupsRoute_UpdateGroup(t *testing.T) {
	// First create a new user
	createdUser := tests.TestUtil_CreateUser(t, "testuser145", "testpassword145")

	// Create a new group
	reqBody := []byte(`{
		"name": "Test Group 02",
		"type": "public",
		"radius": 800
	}`)

	createGroupReq := httptest.NewRequest("POST", "/v1/groups", bytes.NewBuffer(reqBody))
	createGroupReq.Header.Set("Content-Type", "application/json")
	createGroupReq.Header.Set("Authorization", "Bearer "+createdUser.Token)

	createGroupResp := lo.Must(tests.App.Test(createGroupReq, -1))
	assert.Equal(t, fiber.StatusCreated, createGroupResp.StatusCode)

	createGroupRespBody := lo.Must(io.ReadAll(createGroupResp.Body))
	var createGroupRespData dto.GroupResponse
	err := json.Unmarshal(createGroupRespBody, &createGroupRespData)
	assert.NoError(t, err)
	assert.Equal(t, "Test Group 02", createGroupRespData.Name)
	assert.Equal(t, config.GroupTypePublic, createGroupRespData.Type)
	assert.Equal(t, 800, createGroupRespData.Radius)

	groupID := createGroupRespData.ID
	groupCode := createGroupRespData.Code

	testcases := []struct {
		name           string
		groupIDOrCode  string
		updateBody     []byte
		expectedStatus int
		checkResponse  func(t *testing.T, body []byte)
	}{
		{
			name:          "successful group update using group id",
			groupIDOrCode: groupID,
			updateBody: []byte(`{
				"name": "Updated Group",
				"radius": 1000
			}`),
			expectedStatus: fiber.StatusAccepted,
			checkResponse: func(t *testing.T, body []byte) {
				var updateResp dto.GroupResponse
				err := json.Unmarshal(body, &updateResp)
				assert.NoError(t, err)
				assert.Equal(t, "Updated Group", updateResp.Name)
				assert.Equal(t, 1000, updateResp.Radius)
			},
		},
		{
			name:          "successful group update using group code",
			groupIDOrCode: groupCode,
			updateBody: []byte(`{
				"name": "Updated Group 2",
				"radius": 1200
			}`),
			expectedStatus: fiber.StatusAccepted,
			checkResponse: func(t *testing.T, body []byte) {
				var updateResp dto.GroupResponse
				err := json.Unmarshal(body, &updateResp)
				assert.NoError(t, err)
				assert.Equal(t, "Updated Group 2", updateResp.Name)
				assert.Equal(t, 1200, updateResp.Radius)
			},
		},
	}

	for _, tc := range testcases {
		t.Run(tc.name, func(t *testing.T) {
			updateReq := httptest.NewRequest("PATCH", "/v1/groups/"+tc.groupIDOrCode, bytes.NewBuffer(tc.updateBody))
			updateReq.Header.Set("Content-Type", "application/json")
			updateReq.Header.Set("Authorization", "Bearer "+createdUser.Token)

			updateResp := lo.Must(tests.App.Test(updateReq, -1))
			assert.Equal(t, tc.expectedStatus, updateResp.StatusCode)

			body := []byte{}
			if updateResp.Body != nil {
				body = lo.Must(io.ReadAll(updateResp.Body))
			}
			tc.checkResponse(t, body)
		})
	}
}
