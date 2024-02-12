package types

import "context"

type TokenService interface {
	Validate(context.Context, []byte) error
	Token(context.Context, []byte) ([]byte, error)
}

type RequestID string
