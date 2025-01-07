package response

type ErrorResponse struct {
	RequestID string            `json:"request_id"`
	Status    string            `json:"status"`
	Code      int               `json:"code"`
	Message   string            `json:"message"`
	Details   map[string]string `json:"details"`
}
