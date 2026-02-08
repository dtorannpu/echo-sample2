package httpclient

import (
	"context"
	"net/http"
	"time"
)

type AuthHttpClient struct {
	wellKnownConfigUrl string
	httpClient         *http.Client
}

func NewAuthHttpClient(wellKnownConfigUrl string) *AuthHttpClient {
	return &AuthHttpClient{wellKnownConfigUrl: wellKnownConfigUrl, httpClient: &http.Client{Timeout: 10 * time.Second}}
}

func (g *AuthHttpClient) GetAuthWellKnownConfig(ctx context.Context) (*http.Response, error) {
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, g.wellKnownConfigUrl, nil)
	if err != nil {
		return nil, err
	}

	resp, err := g.httpClient.Do(req)
	if err != nil {
		return nil, err
	}

	return resp, nil
}
