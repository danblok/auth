package api

import (
	"bytes"
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/danblok/auth/internal/service"
	"github.com/danblok/auth/pkg/types"
)

func TestHandleTokenReceive(t *testing.T) {
	tests := map[string]struct {
		payload  any
		wantCode int
	}{
		"with payload": {
			payload:  Body{Payload: "some payload"},
			wantCode: http.StatusCreated,
		},
		"empty payload": {
			payload:  Body{Payload: ""},
			wantCode: http.StatusBadRequest,
		},
		"incorrect body payload name": {
			payload:  map[string]any{"random-name": "some text"},
			wantCode: http.StatusBadRequest,
		},
		"incorrect payload type number": {
			payload:  map[string]any{"payload": 5},
			wantCode: http.StatusBadRequest,
		},
		"incorrect payload type bool": {
			payload:  map[string]any{"payload": true},
			wantCode: http.StatusBadRequest,
		},
	}

	svc := service.NewJWTService([]byte("secret-key"))
	srv := NewHTTPServer(svc, "localhost:3000")

	for name, tt := range tests {
		t.Run(name, func(t *testing.T) {
			body, _ := json.Marshal(tt.payload)
			r := httptest.NewRequest("POST", "/token", bytes.NewReader(body))
			w := httptest.NewRecorder()
			h := makeHTTPHandler(srv.handleTokenReceive)
			h(w, r)

			resp := w.Result()
			if resp.StatusCode != tt.wantCode {
				t.Errorf("status code is not the same: want=%d, got=%d", tt.wantCode, resp.StatusCode)
			}
		})
	}
}

func TestHandleTokenValidate(t *testing.T) {
	svc := service.NewJWTService([]byte("secret-key"))
	srv := NewHTTPServer(svc, "localhost:3000")
	tkn, _ := svc.Token(context.Background(), []byte("some payload"))

	tests := map[string]struct {
		token    string
		wantCode int
		valid    bool
	}{
		"valid token": {
			token:    string(tkn),
			wantCode: http.StatusOK,
			valid:    true,
		},
		"invalid token": {
			token:    "invalid-token",
			wantCode: http.StatusOK,
			valid:    false,
		},
		"empty token": {
			token:    "",
			wantCode: http.StatusBadRequest,
			valid:    false,
		},
	}

	for name, tt := range tests {
		t.Run(name, func(t *testing.T) {
			r := httptest.NewRequest("POST", "/validate", nil)
			q := r.URL.Query()
			q.Add("token", tt.token)
			r.URL.RawQuery = q.Encode()
			w := httptest.NewRecorder()
			h := makeHTTPHandler(srv.handleTokenValidation)
			h(w, r)

			resp := w.Result()
			if resp.StatusCode != tt.wantCode {
				t.Errorf("status code is not the same: want=%d, got=%d", tt.wantCode, resp.StatusCode)
			}

			var got types.TokenValidationResponse
			_ = json.NewDecoder(resp.Body).Decode(&got)
			if got.Valid != tt.valid {
				t.Errorf("valid is not the same: want=%d, got=%d", tt.wantCode, resp.StatusCode)
			}
		})
	}
}
