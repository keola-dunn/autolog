package auth

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"

	"github.com/keola-dunn/autolog/internal/jwt"
)

type Client struct {
	http.Client

	urlHost *url.URL
}

func NewClient(host string) (*Client, error) {
	u, err := url.Parse(host)
	if err != nil {
		return nil, fmt.Errorf("failed to parse URL host: %w", err)
	}

	return &Client{
		Client:  *http.DefaultClient,
		urlHost: u,
	}, nil
}

func (c *Client) GetWellKnownJWKS(ctx context.Context) (*jwt.JWKS, error) {
	req, err := http.NewRequestWithContext(ctx, http.MethodGet,
		fmt.Sprintf("%s/.well-known/jwks.json", c.urlHost.String()), nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}

	resp, err := c.Client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("failed to do request: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("received unexpected status code: %d", resp.StatusCode)
	}

	respBody, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read response body: %w", err)
	}

	var jwks jwt.JWKS
	if err := json.Unmarshal(respBody, &jwks); err != nil {
		return nil, fmt.Errorf("failed to unmarshal response body: %w", err)
	}

	return &jwks, nil

}
