package models

// OptionsChainResponse represents the response from Massive API
type OptionsChainResponse struct {
	Status    string           `json:"status"`
	RequestID string           `json:"request_id"`
	Results   []OptionContract `json:"results"`
	NextURL   *string          `json:"next_url,omitempty"`
}

// OptionContract represents a single options contract with all market data
type OptionContract struct {
	Details          *ContractDetails `json:"details,omitempty"`
	Greeks           *Greeks          `json:"greeks,omitempty"`
	ImpliedVol       *float64         `json:"implied_volatility,omitempty"`
	OpenInterest     *int64           `json:"open_interest,omitempty"`
	LastQuote        *LastQuote       `json:"last_quote,omitempty"`
	LastTrade        *LastTrade       `json:"last_trade,omitempty"`
	Day              *DayBar          `json:"day,omitempty"`
	UnderlyingAsset  *UnderlyingAsset `json:"underlying_asset,omitempty"`
}

// ContractDetails contains the contract specifications
type ContractDetails struct {
	Ticker            *string  `json:"ticker,omitempty"`
	ContractType      *string  `json:"contract_type,omitempty"` // "call" or "put"
	StrikePrice       *float64 `json:"strike_price,omitempty"`
	ExpirationDate    *string  `json:"expiration_date,omitempty"`
	ExerciseStyle     *string  `json:"exercise_style,omitempty"` // "american", "european", "bermudan"
	SharesPerContract *int     `json:"shares_per_contract,omitempty"`
}

// Greeks contains the option Greeks
type Greeks struct {
	Delta *float64 `json:"delta,omitempty"`
	Gamma *float64 `json:"gamma,omitempty"`
	Theta *float64 `json:"theta,omitempty"`
	Vega  *float64 `json:"vega,omitempty"`
	Rho   *float64 `json:"rho,omitempty"`
}

// LastQuote contains the most recent bid/ask data
type LastQuote struct {
	Bid     *float64 `json:"bid,omitempty"`
	Ask     *float64 `json:"ask,omitempty"`
	BidSize *int64   `json:"bid_size,omitempty"`
	AskSize *int64   `json:"ask_size,omitempty"`
}

// MidPrice calculates the mid-price between bid and ask
func (q *LastQuote) MidPrice() *float64 {
	if q.Bid != nil && q.Ask != nil {
		mid := (*q.Bid + *q.Ask) / 2.0
		return &mid
	}
	return nil
}

// Spread calculates the bid-ask spread
func (q *LastQuote) Spread() *float64 {
	if q.Bid != nil && q.Ask != nil {
		spread := *q.Ask - *q.Bid
		return &spread
	}
	return nil
}

// LastTrade contains the most recent trade data
type LastTrade struct {
	Price *float64 `json:"price,omitempty"`
	Size  *int64   `json:"size,omitempty"`
}

// DayBar contains OHLCV data for the trading day
type DayBar struct {
	Open   *float64 `json:"open,omitempty"`
	High   *float64 `json:"high,omitempty"`
	Low    *float64 `json:"low,omitempty"`
	Close  *float64 `json:"close,omitempty"`
	Volume *int64   `json:"volume,omitempty"`
}

// PercentChange calculates the % change from open to current price
func (d *DayBar) PercentChange(currentPrice float64) *float64 {
	if d.Open != nil && *d.Open > 0 {
		change := ((currentPrice - *d.Open) / *d.Open) * 100.0
		return &change
	}
	return nil
}

// UnderlyingAsset contains information about the underlying stock
type UnderlyingAsset struct {
	Ticker *string  `json:"ticker,omitempty"`
	Price  *float64 `json:"price,omitempty"`
}
