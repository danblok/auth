package client

import (
	"context"
	"testing"
	"time"

	"github.com/danblok/auth/internal/api"
	"github.com/danblok/auth/internal/service"
)

func TestToken(t *testing.T) {
	payload := []byte("some payload")

	ctx := context.Background()
	svc := service.NewJWTService([]byte("secret"))

	go func() {
		s := api.NewHTTPService(svc, ":42069")
		_ = s.Run()
	}()

	// 100 ms should be fine for the server to startup
	time.Sleep(100 * time.Millisecond)

	c := NewHTPPClient("localhost:42069")
	got, err := c.Token(ctx, payload)
	if got == nil {
		t.Error("got shouldn't be nil")
	}

	if err != nil {
		t.Error("error should be nil")
	}

	c = NewHTPPClient("localhost:80")
	got, err = c.Token(ctx, payload)
	if got != nil {
		t.Error("got should be nil")
	}

	if err == nil {
		t.Error("error shouldn't be nil")
	}
}

func TestValidate(t *testing.T) {
	payload := []byte("some payload")
	ctx := context.Background()
	svc := service.NewJWTService([]byte("secret"))
	tkn, _ := svc.Token(ctx, payload)

	go func() {
		s := api.NewHTTPService(svc, ":42069")
		_ = s.Run()
	}()

	// 100 ms should be fine for the server to startup
	time.Sleep(100 * time.Millisecond)

	c := NewHTPPClient("localhost:42069")
	got, err := c.Validate(ctx, tkn)
	if got.Valid != true {
		t.Error("token should be valid")
	}

	if err != nil {
		t.Error("error should be nil")
	}

	c = NewHTPPClient("localhost:42069")
	got, err = c.Validate(ctx, []byte("some-random-text"))
	if got.Valid != false {
		t.Error("token shouldn't be valid")
	}

	if err != nil {
		t.Error("error should be nil")
	}

	c = NewHTPPClient("localhost:80")
	got, err = c.Validate(ctx, tkn)
	if got != nil {
		t.Error("got should be nil")
	}

	if err == nil {
		t.Error("error shouldn't be nil")
	}
}
