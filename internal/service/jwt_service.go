package service

import (
	"context"
	"time"

	"github.com/golang-jwt/jwt/v5"

	"github.com/danblok/auth/pkg/types"
)

type JWTClaim struct {
	Payload string `json:"payload"`
	jwt.RegisteredClaims
}

// TokenService implementation.
type jwtTokenService struct {
	key []byte
}

// Constructs TokenService implementation.
func NewJWTService(key []byte) types.TokenService {
	return &jwtTokenService{
		key: key,
	}
}

// Validates given token.
func (s jwtTokenService) Validate(_ context.Context, token []byte) error {
	_, err := jwt.Parse(string(token), func(t *jwt.Token) (interface{}, error) {
		return []byte(s.key), nil
	}, jwt.WithValidMethods([]string{"HS256"}))

	return err
}

// Issues new token with given body.
func (s jwtTokenService) Sign(_ context.Context, body []byte) ([]byte, error) {
	claims := &JWTClaim{
		RegisteredClaims: jwt.RegisteredClaims{
			IssuedAt:  &jwt.NumericDate{Time: time.Now()},
			ExpiresAt: &jwt.NumericDate{Time: time.Now().AddDate(0, 0, 1)},
		},
		Payload: string(body),
	}

	tkn := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	ss, err := tkn.SignedString(s.key)
	if err != nil {
		return nil, err
	}

	return []byte(ss), nil
}
