package main

import (
	"bytes"
	"context"
	"errors"
	"flag"
	"log"
	"slices"

	"github.com/danblok/auth/internal/api"
	"github.com/danblok/auth/internal/logging"
	"github.com/danblok/auth/internal/service"
)

func main() {
	svc := service.NewTokenService(
		func(ctx context.Context, body []byte) ([]byte, error) {
			if len(string(body)) < 10 {
				return nil, errors.New("token sign failed")
			}
			return slices.Concat([]byte("signed "), body), nil
		},
		func(ctx context.Context, token []byte) error {
			if !bytes.Contains(token[:7], []byte("signed ")) {
				return errors.New("token validation failed")
			}
			return nil
		},
	)
	svc = logging.NewLoggingService(svc)
	addr := flag.String("addr", ":3000", "listen addr of the server")

	srv := api.NewHTTPService(svc, *addr)
	log.Printf("started server on port: %s", *addr)
	if err := srv.Run(); err != nil {
		log.Fatal("couldn't start the server")
	}
}
