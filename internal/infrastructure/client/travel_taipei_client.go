package client

import (
	"context"
	"fmt"
	"io"
	"net/http"
	"time"
)

// TravelTaipeiClient is a pure HTTP client.
// Responsibilities: build request, send, read body.
// Contains no business logic (language validation and page defaults are not here).
type TravelTaipeiClient struct {
	baseURL    string
	httpClient *http.Client
}

func NewTravelTaipeiClient(baseURL string, timeout time.Duration) *TravelTaipeiClient {
	return &TravelTaipeiClient{
		baseURL: baseURL,
		httpClient: &http.Client{
			Timeout: timeout,
		},
	}
}

// FetchAudio calls the upstream API and returns the raw response body bytes.
// The repository impl is responsible for parsing the bytes into a domain entity.
func (c *TravelTaipeiClient) FetchAudio(ctx context.Context, lang string, page int) ([]byte, error) {
	url := fmt.Sprintf("%s/%s/Media/Audio?page=%d", c.baseURL, lang, page)

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, url, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to build request: %w", err)
	}
	req.Header.Set("Accept", "application/json")

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("HTTP request failed: %w", err)
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read response body: %w", err)
	}

	return body, nil
}
