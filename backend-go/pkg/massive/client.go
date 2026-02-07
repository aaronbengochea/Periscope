package massive

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"time"

	"github.com/aaronbengochea/periscope/backend-go/internal/models"
	"golang.org/x/time/rate"
)

// Client wraps the Massive API client
type Client struct {
	httpClient *http.Client
	baseURL    string
	apiKey     string
	limiter    *rate.Limiter
}

// OptionsChainParams contains optional query parameters for options chain requests
type OptionsChainParams struct {
	StrikePrice    *float64
	ExpirationDate *string
	ContractType   *string // "call" or "put"
	Limit          *int
}

// NewClient creates a new Massive API client
func NewClient(baseURL, apiKey string) *Client {
	return &Client{
		httpClient: &http.Client{
			Timeout: 10 * time.Second,
		},
		baseURL: baseURL,
		apiKey:  apiKey,
		limiter: rate.NewLimiter(rate.Limit(10), 1), // 10 requests per second, burst of 1
	}
}

// GetOptionsChain fetches the options chain for a given underlying ticker
func (c *Client) GetOptionsChain(ctx context.Context, underlyingTicker string, params *OptionsChainParams) (*models.OptionsChainResponse, error) {
	// Apply rate limiting
	if err := c.limiter.Wait(ctx); err != nil {
		return nil, fmt.Errorf("rate limit wait failed: %w", err)
	}

	// Build URL
	u, err := url.Parse(fmt.Sprintf("%s/snapshot/options/%s", c.baseURL, underlyingTicker))
	if err != nil {
		return nil, fmt.Errorf("failed to parse URL: %w", err)
	}

	// Add query parameters
	q := u.Query()
	q.Set("apiKey", c.apiKey)

	if params != nil {
		if params.Limit != nil {
			q.Set("limit", fmt.Sprintf("%d", *params.Limit))
		}
		if params.StrikePrice != nil {
			q.Set("strike_price", fmt.Sprintf("%f", *params.StrikePrice))
		}
		if params.ExpirationDate != nil {
			q.Set("expiration_date", *params.ExpirationDate)
		}
		if params.ContractType != nil {
			q.Set("contract_type", *params.ContractType)
		}
	}

	u.RawQuery = q.Encode()

	// Create request
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, u.String(), nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}

	// Execute request
	resp, err := c.httpClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("failed to execute request: %w", err)
	}
	defer resp.Body.Close()

	// Check status code
	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		return nil, fmt.Errorf("unexpected status code %d: %s", resp.StatusCode, string(body))
	}

	// Parse response
	var result models.OptionsChainResponse
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return nil, fmt.Errorf("failed to decode response: %w", err)
	}

	return &result, nil
}
