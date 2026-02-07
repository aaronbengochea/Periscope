package handlers

import (
	"log"
	"net/http"
	"strconv"

	"github.com/aaronbengochea/periscope/backend-go/internal/models"
	"github.com/aaronbengochea/periscope/backend-go/pkg/errors"
	"github.com/aaronbengochea/periscope/backend-go/pkg/massive"
	"github.com/gin-gonic/gin"
)

// OptionsHandler handles options-related requests
type OptionsHandler struct {
	massiveClient *massive.Client
}

// NewOptionsHandler creates a new options handler
func NewOptionsHandler(massiveClient *massive.Client) *OptionsHandler {
	return &OptionsHandler{
		massiveClient: massiveClient,
	}
}

// GetOptionsChain handles GET /api/v1/options/:ticker
func (h *OptionsHandler) GetOptionsChain(c *gin.Context) {
	ticker := c.Param("ticker")
	if ticker == "" {
		err := errors.NewBadRequestError("ticker is required", nil)
		c.JSON(err.StatusCode, gin.H{"error": err.Message})
		return
	}

	// Parse optional query parameters
	params := &massive.OptionsChainParams{}

	// Set default limit to 250 (Massive API maximum per page)
	// Pagination will automatically fetch all additional pages
	defaultLimit := 250
	if limitStr := c.Query("limit"); limitStr != "" {
		limit, err := strconv.Atoi(limitStr)
		if err != nil {
			appErr := errors.NewBadRequestError("invalid limit parameter", err)
			c.JSON(appErr.StatusCode, gin.H{"error": appErr.Message})
			return
		}
		params.Limit = &limit
	} else {
		params.Limit = &defaultLimit
	}

	if strikeStr := c.Query("strike_price"); strikeStr != "" {
		strike, err := strconv.ParseFloat(strikeStr, 64)
		if err != nil {
			appErr := errors.NewBadRequestError("invalid strike_price parameter", err)
			c.JSON(appErr.StatusCode, gin.H{"error": appErr.Message})
			return
		}
		params.StrikePrice = &strike
	}

	if expiration := c.Query("expiration_date"); expiration != "" {
		params.ExpirationDate = &expiration
	}

	if contractType := c.Query("contract_type"); contractType != "" {
		params.ContractType = &contractType
	}

	log.Printf("[Handler] Fetching options chain for ticker: %s", ticker)

	// Fetch options chain from Massive API
	response, err := h.massiveClient.GetOptionsChain(c.Request.Context(), ticker, params)
	if err != nil {
		log.Printf("[Handler] ✗ Failed to fetch options chain: %v", err)
		appErr := errors.NewInternalError("failed to fetch options chain", err)
		c.JSON(appErr.StatusCode, gin.H{"error": appErr.Message})
		return
	}

	log.Printf("[Handler] ✓ Received %d option contracts", len(response.Results))

	// Fetch underlying stock price separately (required if user doesn't have stocks subscription)
	stockPrice, err := h.massiveClient.GetStockPrice(c.Request.Context(), ticker)
	if err != nil {
		// Log warning but don't fail the request - underlying price might be in options response
		log.Printf("[Handler] ⚠ Stock price fetch failed: %v", err)
		c.Writer.Header().Set("X-Stock-Price-Fetch-Failed", "true")
	} else if stockPrice != nil {
		log.Printf("[Handler] ✓ Stock price fetched: $%.2f", *stockPrice)
		log.Printf("[Handler] Injecting stock price into %d contracts", len(response.Results))

		// Count how many contracts needed price injection
		injected := 0
		for i := range response.Results {
			if response.Results[i].UnderlyingAsset == nil {
				response.Results[i].UnderlyingAsset = &models.UnderlyingAsset{}
			}
			if response.Results[i].UnderlyingAsset.Price == nil {
				response.Results[i].UnderlyingAsset.Price = stockPrice
				injected++
			}
			if response.Results[i].UnderlyingAsset.Ticker == nil {
				response.Results[i].UnderlyingAsset.Ticker = &ticker
			}
		}
		log.Printf("[Handler] ✓ Injected price into %d contracts (already had price: %d)", injected, len(response.Results)-injected)
		c.Writer.Header().Set("X-Stock-Price-Injected", "true")
	}

	log.Printf("[Handler] Sending response with %d contracts to client", len(response.Results))
	c.JSON(http.StatusOK, response)
}

// GetContractDetailsRequest represents the request body for fetching contract details
type GetContractDetailsRequest struct {
	ContractTickers []string `json:"contract_tickers" binding:"required,min=1,max=250"`
}

// GetContractDetails handles POST /api/v1/options/details
// Fetches detailed contract data using the unified snapshot endpoint
func (h *OptionsHandler) GetContractDetails(c *gin.Context) {
	var req GetContractDetailsRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		appErr := errors.NewBadRequestError("invalid request body", err)
		c.JSON(appErr.StatusCode, gin.H{"error": appErr.Message})
		return
	}

	log.Printf("[Handler] Fetching contract details for %d contracts", len(req.ContractTickers))

	// Fetch contract details from Massive API unified snapshot
	contracts, err := h.massiveClient.GetContractDetails(c.Request.Context(), req.ContractTickers)
	if err != nil {
		log.Printf("[Handler] ✗ Failed to fetch contract details: %v", err)
		appErr := errors.NewInternalError("failed to fetch contract details", err)
		c.JSON(appErr.StatusCode, gin.H{"error": appErr.Message})
		return
	}

	log.Printf("[Handler] ✓ Received %d contract details", len(contracts))

	// Return the contracts in the same format as options chain
	response := models.OptionsChainResponse{
		Status:    "OK",
		RequestID: c.GetString("request_id"),
		Results:   contracts,
	}

	c.JSON(http.StatusOK, response)
}
