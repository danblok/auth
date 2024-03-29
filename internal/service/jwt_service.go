package service

import (
	"context"
	"errors"
	"time"

	"github.com/golang-jwt/jwt/v5"

	"github.com/danblok/auth/pkg/types"
)

var errNotVaildToken = errors.New("token not valid")

// JWTClaim that supports payload.
type JWTClaim struct {
	Payload string `json:"payload"`
	jwt.RegisteredClaims
}

// TokenService implementation.
type jwtTokenService struct {
	key []byte
}

// NewJWTService creates a JWT TokenService implementation.
func NewJWTService(key []byte) types.TokenService {
	return &jwtTokenService{
		key: key,
	}
}

// Validates given token.
func (s jwtTokenService) Validate(_ context.Context, token []byte) error {
	tkn, err := jwt.Parse(string(token), func(_ *jwt.Token) (interface{}, error) {
		return []byte(s.key), nil
	}, jwt.WithValidMethods([]string{"HS256"}))
	if err != nil {
		return err
	}

	if !tkn.Valid {
		return errNotVaildToken
	}

	return nil
}

// Issues new token with given body.
func (s jwtTokenService) Token(_ context.Context, payload []byte) ([]byte, error) {
	claims := &JWTClaim{
		RegisteredClaims: jwt.RegisteredClaims{
			IssuedAt:  &jwt.NumericDate{Time: time.Now()},
			ExpiresAt: &jwt.NumericDate{Time: time.Now().AddDate(0, 0, 1)},
		},
		Payload: string(payload),
	}

	tkn := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	ss, err := tkn.SignedString(s.key)
	if err != nil {
		return nil, err
	}

	return []byte(ss), nil
}
