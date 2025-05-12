package dto

import "github.com/championswimmer/api.midpoint.place/src/config"

// GroupUserJoinRequest represents the request to add a user to a group
type GroupUserJoinRequest struct {
	Location
}

// GroupUserResponse represents the response for group user operations
type GroupUserResponse struct {
	UserID    string               `json:"user_id"`
	GroupID   string               `json:"group_id"`
	Latitude  float64              `json:"latitude"`
	Longitude float64              `json:"longitude"`
	Role      config.GroupUserRole `json:"role"`
}
