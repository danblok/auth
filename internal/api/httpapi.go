package api

import (
	"context"
	"crypto/tls"
	"encoding/json"
	"errors"
	"net/http"
	"time"

	"github.com/google/uuid"
	"golang.org/x/net/http2"

	"github.com/danblok/auth/pkg/types"
)

// HTTPServer implementation for TokenService.
type HTTPServer struct {
	svc types.TokenService
	srv *http.Server
	tls bool
}

// HTTPHandlerFunc is a helper handler func.
type HTTPHandlerFunc func(context.Context, http.ResponseWriter, *http.Request) error

// Body represents the body
// of a request to receive a token.
type Body struct {
	Payload string `json:"payload"`
}

// NewHTTPServer constructs new HTTPServer that signs and validates tokens via HTTP.
func NewHTTPServer(svc types.TokenService, addr string) *HTTPServer {
	return &HTTPServer{
		svc: svc,
		srv: &http.Server{
			Addr:        addr,
			ReadTimeout: 3 * time.Second,
			IdleTimeout: 3 * time.Second,
		},
		tls: false,
	}
}

// NewHTTPServerTLS constructs new HTTPServer that signs and validates tokens via HTTP securely.
func NewHTTPServerTLS(svc types.TokenService, addr string, cert tls.Certificate) (*HTTPServer, error) {
	srv := &http.Server{
		Addr:        addr,
		ReadTimeout: 3 * time.Second,
		IdleTimeout: 3 * time.Second,
		TLSConfig:   &tls.Config{Certificates: []tls.Certificate{cert}},
	}

	err := http2.ConfigureServer(srv, nil)
	if err != nil {
		return nil, err
	}

	return &HTTPServer{
		svc: svc,
		tls: true,
		srv: srv,
	}, nil
}

// Run starts the HTTPServer
func (s *HTTPServer) Run() error {
	mux := http.NewServeMux()
	mux.Handle("POST /token", makeHTTPHandler(s.handleTokenReceive))
	mux.Handle("GET /validate", makeHTTPHandler(s.handleTokenValidation))
	s.srv.Handler = mux

	if s.tls {
		return s.srv.ListenAndServeTLS("", "")
	}

	return s.srv.ListenAndServe()
}

// Attaches request_id to the context and returns http.Handler.
func makeHTTPHandler(fn HTTPHandlerFunc) http.HandlerFunc {
	ctx := context.WithValue(context.Background(), types.RequestID("request_id"), uuid.NewString())

	return func(w http.ResponseWriter, r *http.Request) {
		if err := fn(ctx, w, r); err != nil {
			_ = writeJSON(w, http.StatusBadRequest, HTTPErrResponse{Error: err.Error()})
		}
	}
}

// Handles token validation.
func (s *HTTPServer) handleTokenValidation(ctx context.Context, w http.ResponseWriter, r *http.Request) error {
	token := r.URL.Query().Get("token")
	if token == "" {
		return errors.New("token not provided")
	}

	err := s.svc.Validate(ctx, []byte(token))
	if err != nil {
		return writeJSON(w, http.StatusOK, types.TokenValidationResponse{Valid: false})
	}
	return writeJSON(w, http.StatusOK, types.TokenValidationResponse{Valid: true})
}

// Handles token receive.
func (s *HTTPServer) handleTokenReceive(ctx context.Context, w http.ResponseWriter, r *http.Request) error {
	var b Body
	err := json.NewDecoder(r.Body).Decode(&b)
	if err != nil {
		return err
	}
	r.Body.Close()

	if b.Payload == "" {
		return errors.New("incorrect payload")
	}

	token, err := s.svc.Token(ctx, []byte(b.Payload))
	if err != nil {
		return err
	}

	return writeJSON(w, http.StatusCreated, types.TokenResponse{Token: string(token)})
}

// Helper func for responding with JSON.
func writeJSON(w http.ResponseWriter, code int, body any) error {
	w.WriteHeader(code)
	return json.NewEncoder(w).Encode(body)
}
