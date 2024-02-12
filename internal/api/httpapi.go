package api

import (
	"context"
	"encoding/json"
	"errors"
	"log"
	"net/http"
	"time"

	"github.com/google/uuid"

	"github.com/danblok/auth/pkg/types"
)

// HTTP API implementation for TokenService.
type HTTPServer struct {
	svc types.TokenService
	srv *http.Server
}

// HTTP helper handler func.
type HTTPHandlerFunc func(context.Context, http.ResponseWriter, *http.Request) error

// Constructs new HTTPServer that signs and validates tokens via HTTP.
func NewHTTPService(svc types.TokenService, addr string) *HTTPServer {
	return &HTTPServer{
		svc: svc,
		srv: &http.Server{
			Addr:        addr,
			ReadTimeout: 3 * time.Second,
			IdleTimeout: 3 * time.Second,
		},
	}
}

// Runs the HTTPServer
func (s *HTTPServer) Run() error {
	mux := http.NewServeMux()
	mux.Handle("POST /token", makeHTTPHandler(s.handleTokenSign))
	mux.Handle("GET /validate", makeHTTPHandler(s.handleTokenValidation))
	s.srv.Handler = mux

	return s.srv.ListenAndServe()
}

// Attaches request_id to the context and returns http.Handler.
func makeHTTPHandler(fn HTTPHandlerFunc) http.HandlerFunc {
	ctx := context.WithValue(context.Background(), types.RequestID("request_id"), uuid.NewString())

	return func(w http.ResponseWriter, r *http.Request) {
		if err := fn(ctx, w, r); err != nil {
			_ = writeJSON(w, http.StatusBadRequest, map[string]string{"error": err.Error()})
		}
	}
}

// Handles token validation
func (s *HTTPServer) handleTokenValidation(ctx context.Context, w http.ResponseWriter, r *http.Request) error {
	token := r.URL.Query().Get("token")
	if token == "" {
		return errors.New("token not provided")
	}

	err := s.svc.Validate(ctx, []byte(token))
	if err != nil {
		return writeJSON(w, http.StatusBadRequest, map[string]bool{"valid": false})
	}
	return writeJSON(w, http.StatusOK, map[string]bool{"valid": true})
}

// Handles token sign.
func (s *HTTPServer) handleTokenSign(ctx context.Context, w http.ResponseWriter, r *http.Request) error {
	type body struct {
		Msg string `json:"msg"`
	}

	var b body
	err := json.NewDecoder(r.Body).Decode(&b)
	if err != nil {
		return err
	}
	r.Body.Close()
	log.Println(b)

	token, err := s.svc.Sign(ctx, []byte(b.Msg))
	if err != nil {
		return err
	}

	return writeJSON(w, http.StatusCreated, map[string]string{"token": string(token)})
}

// Helper func for responding with JSON.
func writeJSON(w http.ResponseWriter, code int, body any) error {
	w.WriteHeader(code)
	return json.NewEncoder(w).Encode(body)
}
