package massive

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"log"
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

	// Log first contract in full for debugging
	if len(result.Results) > 0 {
		firstContractJSON, _ := json.MarshalIndent(result.Results[0], "", "  ")
		log.Printf("[Massive API] First Option Contract (full structure):\n%s", string(firstContractJSON))
	}

	// Log response details
	log.Printf("[Massive API] Options Chain Response for %s:", underlyingTicker)
	log.Printf("  - Status: %s", result.Status)
	log.Printf("  - Request ID: %s", result.RequestID)
	log.Printf("  - Total contracts: %d", len(result.Results))

	// Check first few contracts for underlying asset data
	if len(result.Results) > 0 {
		firstContract := result.Results[0]
		log.Printf("  - First contract ticker: %v", getStringValue(firstContract.Details, "ticker"))
		log.Printf("  - First contract has underlying_asset: %v", firstContract.UnderlyingAsset != nil)
		if firstContract.UnderlyingAsset != nil {
			log.Printf("    - underlying_asset.ticker: %v", getStringPtrValue(firstContract.UnderlyingAsset.Ticker))
			log.Printf("    - underlying_asset.price: %v", getFloatPtrValue(firstContract.UnderlyingAsset.Price))
		}
	}

	return &result, nil
}

// Helper functions for logging
func getStringValue(details *models.ContractDetails, field string) string {
	if details == nil || details.Ticker == nil {
		return "nil"
	}
	return *details.Ticker
}

func getStringPtrValue(ptr *string) string {
	if ptr == nil {
		return "nil"
	}
	return *ptr
}

func getFloatPtrValue(ptr *float64) string {
	if ptr == nil {
		return "nil"
	}
	return fmt.Sprintf("%.2f", *ptr)
}

// StockSnapshot represents a stock snapshot response
type StockSnapshot struct {
	Status    string        `json:"status"`
	RequestID string        `json:"request_id"`
	Results   []StockResult `json:"results"`
}

// StockResult represents a single stock in the snapshot
type StockResult struct {
	Ticker  string       `json:"ticker"`
	Name    string       `json:"name"`
	Type    string       `json:"type"`
	Session *StockSession `json:"session,omitempty"`
}

// StockSession contains the stock's session data including price
type StockSession struct {
	Close          *float64 `json:"close,omitempty"`
	High           *float64 `json:"high,omitempty"`
	Low            *float64 `json:"low,omitempty"`
	Open           *float64 `json:"open,omitempty"`
	PreviousClose  *float64 `json:"previous_close,omitempty"`
	Change         *float64 `json:"change,omitempty"`
	ChangePercent  *float64 `json:"change_percent,omitempty"`
}

// GetStockPrice fetches the current stock price for a given ticker
func (c *Client) GetStockPrice(ctx context.Context, ticker string) (*float64, error) {
	log.Printf("[Massive API] Fetching stock price for %s", ticker)

	// Apply rate limiting
	if err := c.limiter.Wait(ctx); err != nil {
		return nil, fmt.Errorf("rate limit wait failed: %w", err)
	}

	// Build URL for unified snapshot
	u, err := url.Parse(fmt.Sprintf("%s/snapshot", c.baseURL))
	if err != nil {
		return nil, fmt.Errorf("failed to parse URL: %w", err)
	}

	// Add query parameters
	// Note: Use just ticker parameter without type to get stock data
	q := u.Query()
	q.Set("apiKey", c.apiKey)
	q.Set("ticker", ticker)
	u.RawQuery = q.Encode()

	log.Printf("[Massive API] Stock snapshot URL: %s", u.String()[:len(u.String())-len(c.apiKey)]+"***")

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
	var result StockSnapshot
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return nil, fmt.Errorf("failed to decode response: %w", err)
	}

	// Log full response for debugging
	responseJSON, _ := json.MarshalIndent(result, "", "  ")
	log.Printf("[Massive API] Stock Snapshot Full Response for %s:\n%s", ticker, string(responseJSON))

	// Log response summary
	log.Printf("[Massive API] Stock Snapshot Response for %s:", ticker)
	log.Printf("  - Status: %s", result.Status)
	log.Printf("  - Request ID: %s", result.RequestID)
	log.Printf("  - Total results: %d", len(result.Results))

	if len(result.Results) > 0 {
		stock := result.Results[0]
		log.Printf("  - Ticker: %s", stock.Ticker)
		log.Printf("  - Name: %s", stock.Name)
		log.Printf("  - Type: %s", stock.Type)
		log.Printf("  - Has session data: %v", stock.Session != nil)

		if stock.Session != nil {
			log.Printf("  - Session.Close: %v", getFloatPtrValue(stock.Session.Close))
			log.Printf("  - Session.Open: %v", getFloatPtrValue(stock.Session.Open))
			log.Printf("  - Session.High: %v", getFloatPtrValue(stock.Session.High))
			log.Printf("  - Session.Low: %v", getFloatPtrValue(stock.Session.Low))
			log.Printf("  - Session.Change: %v", getFloatPtrValue(stock.Session.Change))
			log.Printf("  - Session.ChangePercent: %v", getFloatPtrValue(stock.Session.ChangePercent))
		}
	}

	// Extract price from session.close (most recent close price)
	if len(result.Results) > 0 && result.Results[0].Session != nil {
		if result.Results[0].Session.Close != nil {
			log.Printf("[Massive API] ✓ Extracted stock price: $%.2f", *result.Results[0].Session.Close)
			return result.Results[0].Session.Close, nil
		}
	}

	log.Printf("[Massive API] ✗ No price data available for ticker %s", ticker)
	return nil, fmt.Errorf("no price data available for ticker %s", ticker)
}
