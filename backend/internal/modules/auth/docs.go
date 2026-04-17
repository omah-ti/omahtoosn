package auth

type UserSuccessResponse struct {
	Success   bool         `json:"success" example:"true"`
	Message   string       `json:"message" example:"profile fetched"`
	RequestID string       `json:"request_id" example:"12cd4404-a385-47fb-bc1a-eea247de4516"`
	Data      UserResponse `json:"data"`
}

type LoginSuccessResponse struct {
	Success   bool          `json:"success" example:"true"`
	Message   string        `json:"message" example:"login successful"`
	RequestID string        `json:"request_id" example:"12cd4404-a385-47fb-bc1a-eea247de4516"`
	Data      LoginResponse `json:"data"`
}
