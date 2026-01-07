package cryptocurrency

import (
	"encoding/json"
	"net/http"
	"net/url"
	"strings"
	"time"

	"github.com/mayswind/ezbookkeeping/pkg/core"
	"github.com/mayswind/ezbookkeeping/pkg/errs"
	"github.com/mayswind/ezbookkeeping/pkg/log"
	"github.com/mayswind/ezbookkeeping/pkg/models"
	"github.com/mayswind/ezbookkeeping/pkg/utils"
	"github.com/mayswind/ezbookkeeping/pkg/validators"
)

const coinGeckoAPIUrl = "https://api.coingecko.com/api/v3/simple/price"
const coinGeckoReferenceUrl = "https://www.coingecko.com"
const coinGeckoDataSource = "CoinGecko"
const coinGeckoBaseCurrency = "USD"

// CoinGecko symbol to ID mapping
var coinGeckoSymbolToID = map[string]string{
	"BTC":  "bitcoin",
	"ETH":  "ethereum",
	"BNB":  "binancecoin",
	"SOL":  "solana",
	"ADA":  "cardano",
	"XRP":  "ripple",
	"DOT":  "polkadot",
	"DOGE": "dogecoin",
	"MATIC": "matic-network",
	"USDT": "tether",
}

// CoinGeckoDataSource defines the structure of cryptocurrency price data source of CoinGecko
type CoinGeckoDataSource struct {
	HttpCryptocurrencyPriceDataSource
}

// CoinGeckoPriceResponse represents the response from CoinGecko API
type CoinGeckoPriceResponse map[string]map[string]float64

// BuildRequests returns the CoinGecko cryptocurrency prices http requests
func (c *CoinGeckoDataSource) BuildRequests(symbols []string, apiKey string) ([]*http.Request, error) {
	if len(symbols) == 0 {
		return nil, errs.ErrFailedToRequestRemoteApi
	}

	// Convert symbols to CoinGecko IDs
	ids := make([]string, 0, len(symbols))
	for _, symbol := range symbols {
		if id, ok := coinGeckoSymbolToID[strings.ToUpper(symbol)]; ok {
			ids = append(ids, id)
		}
	}

	if len(ids) == 0 {
		return nil, errs.ErrFailedToRequestRemoteApi
	}

	// Build URL with parameters
	u, err := url.Parse(coinGeckoAPIUrl)
	if err != nil {
		return nil, err
	}

	q := u.Query()
	q.Set("ids", strings.Join(ids, ","))
	// CoinGecko API uses "usd" as the base currency
	q.Set("vs_currencies", "usd")
	u.RawQuery = q.Encode()

	req, err := http.NewRequest("GET", u.String(), nil)
	if err != nil {
		return nil, err
	}

	return []*http.Request{req}, nil
}

// Parse returns the common response entity according to the CoinGecko data source raw response
func (c *CoinGeckoDataSource) Parse(core core.Context, content []byte) (*models.LatestCryptocurrencyPriceResponse, error) {
	var coinGeckoData CoinGeckoPriceResponse
	err := json.Unmarshal(content, &coinGeckoData)

	if err != nil {
		log.Errorf(core, "[coingecko_datasource.Parse] failed to parse json data, content is %s, because %s", string(content), err.Error())
		return nil, errs.ErrFailedToRequestRemoteApi
	}

	// Create reverse mapping from ID to symbol
	idToSymbol := make(map[string]string)
	for symbol, id := range coinGeckoSymbolToID {
		idToSymbol[id] = symbol
	}

	prices := make(models.LatestCryptocurrencyPriceSlice, 0, len(coinGeckoData))

	for id, priceData := range coinGeckoData {
		symbol, ok := idToSymbol[id]
		if !ok {
			continue
		}

		// Validate symbol
		if _, exists := validators.AllCryptocurrencySymbols[symbol]; !exists {
			continue
		}

		// CoinGecko API returns prices in USD
		usdPrice, ok := priceData["usd"]
		if !ok {
			continue
		}

		priceStr := utils.Float64ToString(usdPrice)
		if _, err := utils.StringToFloat64(priceStr); err != nil {
			continue
		}

		prices = append(prices, &models.LatestCryptocurrencyPrice{
			Symbol: symbol,
			Price:  priceStr,
		})
	}

	latestPriceResponse := &models.LatestCryptocurrencyPriceResponse{
		DataSource:   coinGeckoDataSource,
		ReferenceUrl: coinGeckoReferenceUrl,
		UpdateTime:   time.Now().Unix(),
		BaseCurrency:  coinGeckoBaseCurrency,
		Prices:       prices,
	}

	return latestPriceResponse, nil
}

