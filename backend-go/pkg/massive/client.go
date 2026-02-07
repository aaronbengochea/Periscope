package massive

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"net/url"
	"strings"
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
// Automatically follows pagination to get all available contracts
func (c *Client) GetOptionsChain(ctx context.Context, underlyingTicker string, params *OptionsChainParams) (*models.OptionsChainResponse, error) {
	allResults := []models.OptionContract{}
	var firstResponse *models.OptionsChainResponse
	nextURL := ""
	pageCount := 0
	maxPages := 20 // Safety limit to prevent infinite loops (20 pages = 5000 contracts)

	for {
		pageCount++
		if pageCount > maxPages {
			log.Printf("[Massive API] ⚠ Reached max page limit (%d), stopping pagination", maxPages)
			break
		}

		var response *models.OptionsChainResponse
		var err error

		if nextURL != "" {
			// Fetch next page using next_url
			response, err = c.fetchPage(ctx, nextURL)
		} else {
			// Fetch first page with params
			response, err = c.fetchFirstPage(ctx, underlyingTicker, params)
		}

		if err != nil {
			return nil, err
		}

		// Store first response for metadata
		if firstResponse == nil {
			firstResponse = response
		}

		// Accumulate results
		allResults = append(allResults, response.Results...)
		log.Printf("[Massive API] Page %d: fetched %d contracts (total: %d)", pageCount, len(response.Results), len(allResults))

		// Check if there are more pages
		if response.NextURL == nil || *response.NextURL == "" {
			break
		}

		nextURL = *response.NextURL
	}

	// Return combined response
	firstResponse.Results = allResults
	log.Printf("[Massive API] ✓ Total contracts fetched: %d across %d pages", len(allResults), pageCount)

	return firstResponse, nil
}

// fetchFirstPage fetches the first page of options chain
func (c *Client) fetchFirstPage(ctx context.Context, underlyingTicker string, params *OptionsChainParams) (*models.OptionsChainResponse, error) {
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

	return c.executeRequest(ctx, u.String())
}

// fetchPage fetches a page using the next_url from pagination
func (c *Client) fetchPage(ctx context.Context, nextURL string) (*models.OptionsChainResponse, error) {
	// Apply rate limiting
	if err := c.limiter.Wait(ctx); err != nil {
		return nil, fmt.Errorf("rate limit wait failed: %w", err)
	}

	// next_url doesn't include API key, so we need to append it
	u, err := url.Parse(nextURL)
	if err != nil {
		return nil, fmt.Errorf("failed to parse next_url: %w", err)
	}

	q := u.Query()
	q.Set("apiKey", c.apiKey)
	u.RawQuery = q.Encode()

	return c.executeRequest(ctx, u.String())
}

// executeRequest executes the HTTP request and parses the response
func (c *Client) executeRequest(ctx context.Context, urlStr string) (*models.OptionsChainResponse, error) {
	// Create request
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, urlStr, nil)
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

func getInt64PtrValue(ptr *int64) string {
	if ptr == nil {
		return "nil"
	}
	return fmt.Sprintf("%d", *ptr)
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

// GetContractDetails fetches detailed snapshot data for specific option contracts
// Uses the unified snapshot endpoint to get bid/ask, greeks, and session data
func (c *Client) GetContractDetails(ctx context.Context, contractTickers []string) ([]models.OptionContract, error) {
	if len(contractTickers) == 0 {
		return []models.OptionContract{}, nil
	}

	log.Printf("[Massive API] Fetching contract details for %d contracts", len(contractTickers))
	if len(contractTickers) > 0 {
		sampleSize := 3
		if len(contractTickers) < sampleSize {
			sampleSize = len(contractTickers)
		}
		log.Printf("[Massive API] First %d contract tickers: %v", sampleSize, contractTickers[:sampleSize])
	}

	// Unified snapshot endpoint supports up to 250 tickers per request
	// If we have more, we need to batch them
	const maxPerRequest = 250
	var allContracts []models.OptionContract

	for i := 0; i < len(contractTickers); i += maxPerRequest {
		end := i + maxPerRequest
		if end > len(contractTickers) {
			end = len(contractTickers)
		}
		batch := contractTickers[i:end]

		// Apply rate limiting
		if err := c.limiter.Wait(ctx); err != nil {
			return nil, fmt.Errorf("rate limit wait failed: %w", err)
		}

		// Build URL for unified snapshot (baseURL already includes /v3)
		u, err := url.Parse(fmt.Sprintf("%s/snapshot", c.baseURL))
		if err != nil {
			return nil, fmt.Errorf("failed to parse URL: %w", err)
		}

		// Add query parameters
		q := u.Query()
		q.Set("apiKey", c.apiKey)
		q.Set("ticker.any_of", strings.Join(batch, ","))

		u.RawQuery = q.Encode()

		// Log the URL (mask API key)
		logURL := u.String()
		if len(c.apiKey) > 3 {
			logURL = strings.Replace(logURL, c.apiKey, c.apiKey[:3]+"***", 1)
		}
		log.Printf("[Massive API] Unified snapshot URL: %s", logURL)
		log.Printf("[Massive API] Unified snapshot request for batch %d-%d", i, end)

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

		allContracts = append(allContracts, result.Results...)
		log.Printf("[Massive API] Batch %d-%d: fetched %d contracts", i, end, len(result.Results))
	}

	log.Printf("[Massive API] ✓ Total contract details fetched: %d", len(allContracts))

	// Log detailed data inspection for first few contracts
	if len(allContracts) > 0 {
		// Print full JSON of first contract for debugging
		contractJSON, _ := json.MarshalIndent(allContracts[0], "", "  ")
		log.Printf("[Massive API] First contract (full JSON):\n%s", string(contractJSON))

		log.Printf("[Massive API] Inspecting first contract for data availability:")
		contract := allContracts[0]

		log.Printf("  - Ticker: %v", getStringPtrValue(contract.Details.Ticker))

		// Check last_quote
		if contract.LastQuote != nil {
			log.Printf("  - LastQuote exists:")
			log.Printf("    - Bid: %v", getFloatPtrValue(contract.LastQuote.Bid))
			log.Printf("    - Ask: %v", getFloatPtrValue(contract.LastQuote.Ask))
			log.Printf("    - BidSize: %v", getInt64PtrValue(contract.LastQuote.BidSize))
			log.Printf("    - AskSize: %v", getInt64PtrValue(contract.LastQuote.AskSize))
		} else {
			log.Printf("  - LastQuote: nil")
		}

		// Check greeks
		if contract.Greeks != nil {
			log.Printf("  - Greeks exists:")
			log.Printf("    - Delta: %v", getFloatPtrValue(contract.Greeks.Delta))
			log.Printf("    - Gamma: %v", getFloatPtrValue(contract.Greeks.Gamma))
			log.Printf("    - Theta: %v", getFloatPtrValue(contract.Greeks.Theta))
			log.Printf("    - Vega: %v", getFloatPtrValue(contract.Greeks.Vega))
		} else {
			log.Printf("  - Greeks: nil")
		}

		// Check session
		if contract.Session != nil {
			log.Printf("  - Session exists:")
			log.Printf("    - Change ($): %v", getFloatPtrValue(contract.Session.Change))
			log.Printf("    - Change (%%): %v", getFloatPtrValue(contract.Session.ChangePercent))
			log.Printf("    - Close: %v", getFloatPtrValue(contract.Session.Close))
			log.Printf("    - Volume: %v", getInt64PtrValue(contract.Session.Volume))
		} else {
			log.Printf("  - Session: nil")
		}
	}

	return allContracts, nil
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
