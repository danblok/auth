package client

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/danblok/auth/internal/api"
	"github.com/danblok/auth/pkg/types"
)

// HTTP client of TokenService.
type HTTPClient struct {
	host string
}

// Constructs a new HTTPClient with given host of the Token service server.
func NewHTPPClient(host string) *HTTPClient {
	return &HTTPClient{
		host: host,
	}
}

// Token fetches a new token and returns it.
func (c *HTTPClient) Token(ctx context.Context, payload []byte) (*types.TokenResponse, error) {
	url := fmt.Sprintf("http://%s/token", c.host)
	body, err := json.Marshal(api.Body{Payload: string(payload)})
	if err != nil {
		return nil, err
	}
	req, err := http.NewRequestWithContext(ctx, "POST", url, bytes.NewReader(body))
	if err != nil {
		return nil, err
	}

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		fmt.Println("aboba")
		return nil, err
	}

	if resp.StatusCode != http.StatusCreated {
		var httpErr api.HTTPErrResponse
		err := json.NewDecoder(resp.Body).Decode(&httpErr)
		if err != nil {
			fmt.Println("aboba")
			return nil, err
		}

		return nil, fmt.Errorf("server responded with non Created status: %v", httpErr.Error)
	}

	token := new(types.TokenResponse)
	if err := json.NewDecoder(resp.Body).Decode(token); err != nil {
		return nil, err
	}

	return token, nil
}

// Validate sends the given token to the server to validate it and returns validation result.
func (c *HTTPClient) Validate(ctx context.Context, token []byte) (*types.TokenValidationResponse, error) {
	url := fmt.Sprintf("http://%s/validate?token=%s", c.host, string(token))
	req, err := http.NewRequestWithContext(ctx, "GET", url, nil)
	if err != nil {
		return nil, err
	}

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, err
	}

	if resp.StatusCode != http.StatusOK {
		var httpErr api.HTTPErrResponse
		if err := json.NewDecoder(resp.Body).Decode(&httpErr); err != nil {
			return nil, err
		}
		defer resp.Body.Close()

		return nil, fmt.Errorf("server responded with non OK status: %v", httpErr.Error)
	}

	valid := new(types.TokenValidationResponse)
	if err := json.NewDecoder(resp.Body).Decode(&valid); err != nil {
		return nil, err
	}

	return valid, nil
}
