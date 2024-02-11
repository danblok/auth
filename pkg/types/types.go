package types

import "context"

type TokenService interface {
	Validate(context.Context, []byte) error
	Sign(context.Context, []byte) ([]byte, error)
}
