package handlers

import (
	"net/http"
	"strconv"

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

	if limitStr := c.Query("limit"); limitStr != "" {
		limit, err := strconv.Atoi(limitStr)
		if err != nil {
			appErr := errors.NewBadRequestError("invalid limit parameter", err)
			c.JSON(appErr.StatusCode, gin.H{"error": appErr.Message})
			return
		}
		params.Limit = &limit
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

	// Fetch options chain from Massive API
	response, err := h.massiveClient.GetOptionsChain(c.Request.Context(), ticker, params)
	if err != nil {
		appErr := errors.NewInternalError("failed to fetch options chain", err)
		c.JSON(appErr.StatusCode, gin.H{"error": appErr.Message})
		return
	}

	c.JSON(http.StatusOK, response)
}
