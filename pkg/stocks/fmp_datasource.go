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
	// FMP v3 quote: one request for multiple symbols (comma-separated in path). Free tier: 250 requests/day.
	fmpQuoteApiUrl = "https://financialmodelingprep.com/api/v3/quote/"
)

// FMPDataSource defines the structure of Financial Modeling Prep data source
type FMPDataSource struct {
}

// FMPQuoteItem represents one quote in FMP v3 quote API response (API may use symbol/price or Symbol/Price)
type FMPQuoteItem struct {
	Symbol   string  `json:"symbol"`
	Price    float64 `json:"price"`
	Currency string  `json:"currency,omitempty"`
	// FMP sometimes returns PascalCase in stable endpoints
	SymbolAlt   string  `json:"Symbol,omitempty"`
	PriceAlt    float64 `json:"Price,omitempty"`
	CurrencyAlt string  `json:"Currency,omitempty"`
}

// BuildRequests builds the http requests (single request for all symbols)
func (s *FMPDataSource) BuildRequests(symbols []string, apiKey string) ([]*http.Request, error) {
	if len(symbols) == 0 {
		return nil, nil
	}

	// One request: /api/v3/quote/AAPL,MSFT,GOOG?apikey=xxx
	path := fmpQuoteApiUrl + strings.Join(symbols, ",")
	u, err := url.Parse(path)
	if err != nil {
		return nil, err
	}
	q := u.Query()
	q.Set("apikey", apiKey)
	u.RawQuery = q.Encode()

	req, err := http.NewRequest("GET", u.String(), nil)
	if err != nil {
		return nil, err
	}
	return []*http.Request{req}, nil
}

// Parse parses the response content (array of quote objects)
func (s *FMPDataSource) Parse(c core.Context, content []byte) (*models.LatestStockPriceResponse, error) {
	var items []FMPQuoteItem
	if err := json.Unmarshal(content, &items); err != nil {
		// FMP may return error as object, e.g. {"Error Message": "..."}
		var errObj map[string]interface{}
		if json.Unmarshal(content, &errObj) == nil {
			if msg, ok := errObj["Error Message"]; ok {
				return nil, fmt.Errorf("fmp api error: %v", msg)
			}
			if msg, ok := errObj["message"]; ok {
				return nil, fmt.Errorf("fmp api error: %v", msg)
			}
		}
		return nil, err
	}

	if len(items) == 0 {
		return nil, fmt.Errorf("fmp api returned empty quote list")
	}

	prices := make(models.LatestStockPriceSlice, 0, len(items))
	for _, item := range items {
		symbol := item.Symbol
		if symbol == "" {
			symbol = item.SymbolAlt
		}
		if symbol == "" {
			continue
		}
		price := item.Price
		if price == 0 {
			price = item.PriceAlt
		}
		currency := strings.ToUpper(item.Currency)
		if currency == "" {
			currency = strings.ToUpper(item.CurrencyAlt)
		}
		if currency == "" {
			currency = "USD"
		}
		if strings.HasSuffix(strings.ToUpper(symbol), ".HK") {
			currency = "HKD"
		} else if strings.HasSuffix(strings.ToUpper(symbol), ".SS") || strings.HasSuffix(strings.ToUpper(symbol), ".SZ") {
			currency = "CNY"
		}
		priceStr := strconv.FormatFloat(price, 'f', -1, 64)
		prices = append(prices, &models.LatestStockPrice{
			Symbol:   strings.ToUpper(symbol),
			Price:    priceStr,
			Currency: currency,
		})
	}

	if len(prices) == 0 {
		return nil, fmt.Errorf("fmp api returned no valid quotes")
	}

	return &models.LatestStockPriceResponse{
		DataSource:   "Financial Modeling Prep",
		ReferenceUrl: "https://site.financialmodelingprep.com/",
		UpdateTime:   time.Now().Unix(),
		BaseCurrency: "USD",
		Prices:       prices,
	}, nil
}
