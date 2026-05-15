package apiresponse

type ErrorResponse struct {
	Error string `json:"error" example:"The server encountered a problem and could not process your request"`
}
