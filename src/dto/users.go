package dto

type CreateUserRequest struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

type LoginUserRequest = CreateUserRequest

type UserUpdateRequest struct {
	Location Location `json:"location"`
}

type UserResponse struct {
	ID       uint     `json:"id"`
	Username string   `json:"username"`
	Token    string   `json:"token"`
	Location Location `json:"location,omitempty"`
}
