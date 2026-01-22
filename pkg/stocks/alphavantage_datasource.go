package stocks

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"strings"
	"time"

	"github.com/mayswind/ezbookkeeping/pkg/core"
	"github.com/mayswind/ezbookkeeping/pkg/models"
)

const (
	alphaVantageApiUrl = "https://www.alphavantage.co/query"
)

// AlphaVantageDataSource defines the structure of Alpha Vantage data source
type AlphaVantageDataSource struct {
}

// AlphaVantageGlobalQuoteResponse represents the response from Alpha Vantage Global Quote API
type AlphaVantageGlobalQuoteResponse struct {
	GlobalQuote struct {
		Symbol           string `json:"01. symbol"`
		Price            string `json:"05. price"`
		LatestTradingDay string `json:"07. latest trading day"`
	} `json:"Global Quote"`
}

// BuildRequests builds the http requests
func (s *AlphaVantageDataSource) BuildRequests(symbols []string, apiKey string) ([]*http.Request, error) {
	if len(symbols) == 0 {
		return nil, nil
	}

	requests := make([]*http.Request, 0, len(symbols))

	for _, symbol := range symbols {
		u, err := url.Parse(alphaVantageApiUrl)
		if err != nil {
			return nil, err
		}

		q := u.Query()
		q.Set("function", "GLOBAL_QUOTE")
		q.Set("symbol", symbol)
		q.Set("apikey", apiKey)
		u.RawQuery = q.Encode()

		req, err := http.NewRequest("GET", u.String(), nil)
		if err != nil {
			return nil, err
		}

		requests = append(requests, req)
	}

	return requests, nil
}

// Parse parses the response content
func (s *AlphaVantageDataSource) Parse(c core.Context, content []byte) (*models.LatestStockPriceResponse, error) {
	var response AlphaVantageGlobalQuoteResponse
	err := json.Unmarshal(content, &response)
	if err != nil {
		return nil, err
	}

	if response.GlobalQuote.Symbol == "" {
		// Alpha Vantage sometimes returns an empty object or error message in a different field
		var errorResponse map[string]interface{}
		json.Unmarshal(content, &errorResponse)
		if msg, ok := errorResponse["Error Message"]; ok {
			return nil, fmt.Errorf("alpha vantage api error: %v", msg)
		}
		if msg, ok := errorResponse["Information"]; ok {
			return nil, fmt.Errorf("alpha vantage api info: %v", msg)
		}
		if msg, ok := errorResponse["Note"]; ok {
			return nil, fmt.Errorf("alpha vantage api rate limit: %v", msg)
		}
		return nil, fmt.Errorf("invalid response from alpha vantage")
	}

	price := response.GlobalQuote.Price
	if price == "" {
		price = "0"
	}

	// Alpha Vantage Global Quote does not return currency directly in the same call.
	// For simplicity, we default to USD or try to guess from symbol, but the better way is to use its Search API.
	// For now, let's just return what we have.
	currency := "USD"
	if strings.HasSuffix(strings.ToUpper(response.GlobalQuote.Symbol), ".HK") {
		currency = "HKD"
	} else if strings.HasSuffix(strings.ToUpper(response.GlobalQuote.Symbol), ".SS") || strings.HasSuffix(strings.ToUpper(response.GlobalQuote.Symbol), ".SZ") {
		currency = "CNY"
	}

	prices := []*models.LatestStockPrice{
		{
			Symbol:   strings.ToUpper(response.GlobalQuote.Symbol),
			Price:    price,
			Currency: currency,
		},
	}

	return &models.LatestStockPriceResponse{
		DataSource:   "Alpha Vantage",
		ReferenceUrl: "https://www.alphavantage.co/",
		UpdateTime:   time.Now().Unix(),
		BaseCurrency: currency,
		Prices:       prices,
	}, nil
}
