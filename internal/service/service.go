package service

import (
	"context"

	"github.com/danblok/auth/pkg/types"
)

// TokenService implementation with
// custom sign and validation functions.
type tokenService struct {
	sign     signFunc
	validate validationFunc
}

// Custom sign func type.
type signFunc func(context.Context, []byte) ([]byte, error)

// Custom validation func type.
type validationFunc func(context.Context, []byte) error

// Constructor for a bare TokenService with ability to create your own sign and validation functions.
func NewTokenService(sign signFunc, validate validationFunc) types.TokenService {
	return &tokenService{
		sign:     sign,
		validate: validate,
	}
}

// Validation func for TokenService.
func (s *tokenService) Validate(ctx context.Context, token []byte) error {
	return s.validate(ctx, token)
}

// Token function for TokenService.
func (s *tokenService) Token(ctx context.Context, payload []byte) ([]byte, error) {
	return s.sign(ctx, payload)
}
