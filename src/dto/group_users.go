package dto

// GroupUserJoinRequest represents the request to add a user to a group
type GroupUserJoinRequest struct {
	UserID    uint    `json:"user_id" validate:"required"`
	GroupID   string  `json:"group_id" validate:"required"`
	Latitude  float64 `json:"latitude" validate:"required"`
	Longitude float64 `json:"longitude" validate:"required"`
}

// GroupUserLeaveRequest represents the request to remove a user from a group
type GroupUserLeaveRequest struct {
	UserID  uint   `json:"user_id" validate:"required"`
	GroupID string `json:"group_id" validate:"required"`
}

// GroupUserResponse represents the response for group user operations
type GroupUserResponse struct {
	UserID    uint    `json:"user_id"`
	GroupID   string  `json:"group_id"`
	Latitude  float64 `json:"latitude"`
	Longitude float64 `json:"longitude"`
}
