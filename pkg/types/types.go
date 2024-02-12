package types

import "context"

// TokenService interfaces out
// implementations details of
// services of different levels.
type TokenService interface {
	Validate(context.Context, []byte) error
	Token(context.Context, []byte) ([]byte, error)
}

// RequestID type is used by a context
// in services to attach and receive
// the request id of each request.
type RequestID string

// TokenResponse is used in HTTP server and
// HTTP client for responses from server.
type TokenResponse struct {
	Token string `json:"token"`
}

// TokenValidationResponse is used in HTTP server and
// HTTP client for responses from server.
type TokenValidationResponse struct {
	Valid bool `json:"valid"`
}
