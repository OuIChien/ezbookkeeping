package cryptocurrency

import (
	"encoding/json"
	"net/http"
	"net/url"
	"strconv"
	"strings"
	"time"

	"github.com/mayswind/ezbookkeeping/pkg/core"
	"github.com/mayswind/ezbookkeeping/pkg/errs"
	"github.com/mayswind/ezbookkeeping/pkg/models"
)

const (
	coinGeckoPriceApiUrl = "https://api.coingecko.com/api/v3/simple/price"
)

// CoinGeckoDataSource defines the structure of CoinGecko data source
type CoinGeckoDataSource struct {
}

// CoinGeckoPriceResponse represents the response from CoinGecko API
type CoinGeckoPriceResponse map[string]map[string]float64

// ToLatestCryptocurrencyPriceResponse converts CoinGecko response to LatestCryptocurrencyPriceResponse
func (r CoinGeckoPriceResponse) ToLatestCryptocurrencyPriceResponse(symbolMap map[string]string) *models.LatestCryptocurrencyPriceResponse {
	prices := make(models.LatestCryptocurrencyPriceSlice, 0, len(r))

	for coinId, priceData := range r {
		symbol := ""
		
		// Find symbol by coinId (reverse lookup)
		for s, id := range symbolMap {
			if id == coinId {
				symbol = s
				break
			}
		}

		if symbol == "" {
			continue
		}

		if priceInUsd, exists := priceData["usd"]; exists {
			prices = append(prices, &models.LatestCryptocurrencyPrice{
				Symbol: strings.ToUpper(symbol),
				Price:  strconv.FormatFloat(priceInUsd, 'f', -1, 64),
			})
		}
	}

	return &models.LatestCryptocurrencyPriceResponse{
		DataSource:   "CoinGecko",
		ReferenceUrl: "https://www.coingecko.com/",
		UpdateTime:   time.Now().Unix(),
		BaseCurrency: "USDT", // Treating USD as USDT for simplicity as per design
		Prices:       prices,
	}
}

// BuildRequests builds the http requests
func (s *CoinGeckoDataSource) BuildRequests(symbols []string, apiKey string) ([]*http.Request, error) {
	if len(symbols) == 0 {
		return nil, nil
	}

	// Map symbols to CoinGecko IDs
	// In a real implementation, we might need a more robust mapping or allow user to configure mapping
	// For now, we use a simple mapping for common coins
	coinIds := make([]string, 0, len(symbols))
	for _, symbol := range symbols {
		if id, ok := coinGeckoSymbolMap[strings.ToUpper(symbol)]; ok {
			coinIds = append(coinIds, id)
		}
	}

	if len(coinIds) == 0 {
		return nil, errs.ErrInvalidCryptocurrencySymbol
	}

	u, err := url.Parse(coinGeckoPriceApiUrl)
	if err != nil {
		return nil, err
	}

	q := u.Query()
	q.Set("ids", strings.Join(coinIds, ","))
	q.Set("vs_currencies", "usd")
	u.RawQuery = q.Encode()

	req, err := http.NewRequest("GET", u.String(), nil)
	if err != nil {
		return nil, err
	}

	return []*http.Request{req}, nil
}

// Parse parses the response content
func (s *CoinGeckoDataSource) Parse(c core.Context, content []byte) (*models.LatestCryptocurrencyPriceResponse, error) {
	var response CoinGeckoPriceResponse
	err := json.Unmarshal(content, &response)
	if err != nil {
		return nil, err
	}

	return response.ToLatestCryptocurrencyPriceResponse(coinGeckoSymbolMap), nil
}

// coinGeckoSymbolMap maps common cryptocurrency symbols to CoinGecko IDs
// This is a simplified list. In production, this might need to be more comprehensive or configurable.
var coinGeckoSymbolMap = map[string]string{
	"BTC":   "bitcoin",
	"ETH":   "ethereum",
	"BNB":   "binancecoin",
	"SOL":   "solana",
	"ADA":   "cardano",
	"XRP":   "ripple",
	"DOT":   "polkadot",
	"DOGE":  "dogecoin",
	"MATIC": "matic-network",
	"USDT":  "tether",
	"USDC":  "usd-coin",
	"DAI":   "dai",
	"LTC":   "litecoin",
	"BCH":   "bitcoin-cash",
	"LINK":  "chainlink",
	"XLM":   "stellar",
	"UNI":   "uniswap",
	"ATOM":  "cosmos",
	"XMR":   "monero",
	"ETC":   "ethereum-classic",
}
