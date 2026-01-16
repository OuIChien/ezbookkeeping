package stocks

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"strconv"
	"strings"
	"time"

	"github.com/mayswind/ezbookkeeping/pkg/core"
	"github.com/mayswind/ezbookkeeping/pkg/models"
)

const (
	yahooFinanceQuoteApiUrl = "https://query1.finance.yahoo.com/v7/finance/quote"
)

// YahooFinanceDataSource defines the structure of Yahoo Finance data source
type YahooFinanceDataSource struct {
}

// YahooFinanceQuoteResponse represents the response from Yahoo Finance API
type YahooFinanceQuoteResponse struct {
	QuoteResponse struct {
		Result []struct {
			Symbol             string  `json:"symbol"`
			RegularMarketPrice float64 `json:"regularMarketPrice"`
			Currency           string  `json:"currency"`
		} `json:"result"`
		Error interface{} `json:"error"`
	} `json:"quoteResponse"`
}

// BuildRequests builds the http requests
func (s *YahooFinanceDataSource) BuildRequests(symbols []string, apiKey string) ([]*http.Request, error) {
	if len(symbols) == 0 {
		return nil, nil
	}

	u, err := url.Parse(yahooFinanceQuoteApiUrl)
	if err != nil {
		return nil, err
	}

	q := u.Query()
	q.Set("symbols", strings.Join(symbols, ","))
	q.Set("fields", "symbol,regularMarketPrice,currency")
	u.RawQuery = q.Encode()

	req, err := http.NewRequest("GET", u.String(), nil)
	if err != nil {
		return nil, err
	}

	return []*http.Request{req}, nil
}

// Parse parses the response content
func (s *YahooFinanceDataSource) Parse(c core.Context, content []byte) (*models.LatestStockPriceResponse, error) {
	var response YahooFinanceQuoteResponse
	err := json.Unmarshal(content, &response)
	if err != nil {
		return nil, err
	}

	if response.QuoteResponse.Error != nil {
		return nil, fmt.Errorf("yahoo finance api error: %v", response.QuoteResponse.Error)
	}

	prices := make(models.LatestStockPriceSlice, 0, len(response.QuoteResponse.Result))

	for _, result := range response.QuoteResponse.Result {
		prices = append(prices, &models.LatestStockPrice{
			Symbol:   strings.ToUpper(result.Symbol),
			Price:    strconv.FormatFloat(result.RegularMarketPrice, 'f', -1, 64),
			Currency: strings.ToUpper(result.Currency),
		})
	}

	return &models.LatestStockPriceResponse{
		DataSource:   "Yahoo Finance",
		ReferenceUrl: "https://finance.yahoo.com/",
		UpdateTime:   time.Now().Unix(),
		BaseCurrency: "USD", // This is a bit complex as Yahoo returns prices in various currencies based on exchange, but for now we standardise on response field if needed. The proposal says "Total Account Value = Held Quantity × Real-time Market Price", and "converts the valuation to user's Default Currency".
		Prices:       prices,
	}, nil
}
