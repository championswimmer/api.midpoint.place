package dto

type CreateUserRequest struct {
	Email       string `json:"email"`
	DisplayName string `json:"display_name"`
	Password    string `json:"password"`
}

type LoginUserRequest struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

type UserUpdateRequest struct {
	Location Location `json:"location"`
}

type UserResponse struct {
	ID          uint     `json:"id"`
	Email       string   `json:"email"`
	DisplayName string   `json:"display_name"`
	Token       string   `json:"token"`
	Location    Location `json:"location,omitempty"`
}
