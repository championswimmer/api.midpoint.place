package controllers

import (
	"testing"

	"github.com/championswimmer/api.midpoint.place/src/config"
	"github.com/championswimmer/api.midpoint.place/src/db/models"
	"github.com/championswimmer/api.midpoint.place/src/dto"
	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

func setupGroupsControllerTestDB(t *testing.T) *gorm.DB {
	t.Helper()
	db, err := gorm.Open(sqlite.Open("file::memory:?cache=shared"), &gorm.Config{})
	require.NoError(t, err)
	require.NoError(t, db.AutoMigrate(&models.User{}, &models.Group{}, &models.GroupUser{}))
	return db
}

func createGroupFixture(t *testing.T, db *gorm.DB, groupType config.GroupType) models.Group {
	t.Helper()
	creator := models.User{
		Email:       uuid.NewString() + "@test.com",
		DisplayName: "creator",
		Password:    "password",
	}
	require.NoError(t, db.Create(&creator).Error)

	group := models.Group{
		ID:        uuid.NewString(),
		CreatorID: creator.ID,
		Name:      "test-group",
		Code:      uuid.NewString()[:10],
		Secret:    "123456",
		Type:      groupType,
		Radius:    1000,
	}
	require.NoError(t, db.Create(&group).Error)

	return group
}

func createGroupMembers(t *testing.T, db *gorm.DB, groupID string, count int) {
	t.Helper()
	for i := 0; i < count; i++ {
		user := models.User{
			Email:       uuid.NewString() + "@test.com",
			DisplayName: "member",
			Password:    "password",
		}
		require.NoError(t, db.Create(&user).Error)
		require.NoError(t, db.Create(&models.GroupUser{
			UserID:    user.ID,
			GroupID:   groupID,
			Latitude:  12.9716,
			Longitude: 77.5946,
		}).Error)
	}
}

func TestUpdateGroup_AllowsOpeningPrivacy_WhenMemberCountIsAtMostOne(t *testing.T) {
	db := setupGroupsControllerTestDB(t)
	controller := &GroupsController{db: db}
	group := createGroupFixture(t, db, config.GroupTypePrivate)
	createGroupMembers(t, db, group.ID, 1)

	resp, err := controller.UpdateGroup(group.ID, &dto.UpdateGroupRequest{Type: config.GroupTypePublic})
	require.NoError(t, err)
	assert.Equal(t, config.GroupTypePublic, resp.Type)
}

func TestUpdateGroup_AllowsOpeningPrivacy_WhenMemberCountIsZero(t *testing.T) {
	db := setupGroupsControllerTestDB(t)
	controller := &GroupsController{db: db}
	group := createGroupFixture(t, db, config.GroupTypePrivate)

	resp, err := controller.UpdateGroup(group.ID, &dto.UpdateGroupRequest{Type: config.GroupTypePublic})
	require.NoError(t, err)
	assert.Equal(t, config.GroupTypePublic, resp.Type)
}

func TestUpdateGroup_BlocksOpeningPrivacy_WhenMemberCountIsMoreThanOne(t *testing.T) {
	db := setupGroupsControllerTestDB(t)
	controller := &GroupsController{db: db}
	group := createGroupFixture(t, db, config.GroupTypePrivate)
	createGroupMembers(t, db, group.ID, 2)

	resp, err := controller.UpdateGroup(group.ID, &dto.UpdateGroupRequest{Type: config.GroupTypeProtected})
	require.Nil(t, resp)
	require.Error(t, err)
	fiberErr, ok := err.(*fiber.Error)
	require.True(t, ok)
	assert.Equal(t, fiber.StatusUnprocessableEntity, fiberErr.Code)
	assert.Contains(t, fiberErr.Message, "Cannot make group privacy more open")
}

func TestUpdateGroup_AllowsClosingPrivacy_WhenMemberCountIsMoreThanOne(t *testing.T) {
	db := setupGroupsControllerTestDB(t)
	controller := &GroupsController{db: db}
	group := createGroupFixture(t, db, config.GroupTypePublic)
	createGroupMembers(t, db, group.ID, 2)

	resp, err := controller.UpdateGroup(group.ID, &dto.UpdateGroupRequest{Type: config.GroupTypePrivate})
	require.NoError(t, err)
	assert.Equal(t, config.GroupTypePrivate, resp.Type)
}
