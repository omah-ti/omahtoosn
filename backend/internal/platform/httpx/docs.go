package httpx

type ErrorResponse struct {
	Success   bool              `json:"success" example:"false"`
	Message   string            `json:"message" example:"invalid request body"`
	RequestID string            `json:"request_id" example:"12cd4404-a385-47fb-bc1a-eea247de4516"`
	Fields    map[string]string `json:"fields,omitempty"`
}

type EmptySuccessResponse struct {
	Success   bool   `json:"success" example:"true"`
	Message   string `json:"message" example:"logout successful"`
	RequestID string `json:"request_id" example:"12cd4404-a385-47fb-bc1a-eea247de4516"`
}
