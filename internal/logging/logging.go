package logging

import (
	"context"
	"fmt"
	"log/slog"
	"time"

	"github.com/danblok/auth/pkg/types"
)

// Logging for TokenService.
type loggingService struct {
	svc types.TokenService
	log *slog.Logger
}

// NewLoggingService creates logging for TokenService.
func NewLoggingService(svc types.TokenService) types.TokenService {
	return &loggingService{
		svc: svc,
		log: slog.Default(),
	}
}

// Token passes call to Token to the next TokenService implmentator
// and logs time since start, request_id, err and a new token to stdout.
func (s *loggingService) Token(ctx context.Context, payload []byte) (token []byte, err error) {
	defer func(t time.Time) {
		s.log.InfoContext(
			ctx,
			fmt.Sprintf(
				"time=%+v, request_id=%+v, err=%+v, token=%+v",
				time.Since(t),
				ctx.Value(types.RequestID("request_id")),
				err,
				string(token),
			),
		)
	}(time.Now())

	return s.svc.Token(ctx, payload)
}

// Validate passes call to Validate to the next TokenService implmentator
// and logs time since start, request_id, err if validation failed.
func (s *loggingService) Validate(ctx context.Context, token []byte) (err error) {
	defer func(t time.Time) {
		s.log.InfoContext(
			ctx,
			fmt.Sprintf(
				"time=%+v, request_id=%+v, err=%+v, token=%+v",
				time.Since(t),
				ctx.Value(types.RequestID("request_id")),
				err,
				string(token),
			),
		)
	}(time.Now())

	return s.svc.Validate(ctx, token)
}
