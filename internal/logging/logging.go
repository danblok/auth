package logging

import (
	"context"
	"fmt"
	"log/slog"
	"time"

	"github.com/danblok/auth/pkg/types"
)

// Logging for TokenService
type loggingService struct {
	svc types.TokenService
	log *slog.Logger
}

func NewLoggingService(svc types.TokenService) types.TokenService {
	return &loggingService{
		svc: svc,
		log: slog.Default(),
	}
}

func (s *loggingService) Sign(ctx context.Context, body []byte) (token []byte, err error) {
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

	return s.svc.Sign(ctx, body)
}

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
