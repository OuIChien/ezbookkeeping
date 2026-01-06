package cryptocurrency

import (
	"context"
	"io"
	"net"
	"net/http"
	"net/url"
	"sort"
	"strings"

	"github.com/mayswind/ezbookkeeping/pkg/core"
	"github.com/mayswind/ezbookkeeping/pkg/errs"
	"github.com/mayswind/ezbookkeeping/pkg/log"
	"github.com/mayswind/ezbookkeeping/pkg/models"
	"github.com/mayswind/ezbookkeeping/pkg/settings"
	"github.com/mayswind/ezbookkeeping/pkg/utils"
)

// HttpCryptocurrencyPriceDataSource defines the structure of http cryptocurrency price data source
type HttpCryptocurrencyPriceDataSource interface {
	// BuildRequests returns the http requests
	BuildRequests(symbols []string, apiKey string) ([]*http.Request, error)

	// Parse returns the common response entity according to the data source raw response
	Parse(c core.Context, content []byte) (*models.LatestCryptocurrencyPriceResponse, error)
}

// CommonHttpCryptocurrencyPriceDataProvider defines the structure of common http cryptocurrency price data provider
type CommonHttpCryptocurrencyPriceDataProvider struct {
	CryptocurrencyPriceDataProvider
	dataSource HttpCryptocurrencyPriceDataSource
	httpClient *http.Client
	config     *settings.Config
}

func (c *CommonHttpCryptocurrencyPriceDataProvider) GetLatestCryptocurrencyPrices(core core.Context, uid int64, currentConfig *settings.Config) (*models.LatestCryptocurrencyPriceResponse, error) {
	if len(currentConfig.CryptocurrencySymbols) == 0 {
		log.Warnf(core, "[common_http_cryptocurrency_price_data_provider.GetLatestCryptocurrencyPrices] no cryptocurrency symbols configured for user \"uid:%d\"", uid)
		return &models.LatestCryptocurrencyPriceResponse{
			BaseCurrency: "USDT",
			Prices:       make(models.LatestCryptocurrencyPriceSlice, 0),
		}, nil
	}

	requests, err := c.dataSource.BuildRequests(currentConfig.CryptocurrencySymbols, currentConfig.CryptocurrencyAPIKey)

	if err != nil {
		log.Errorf(core, "[common_http_cryptocurrency_price_data_provider.GetLatestCryptocurrencyPrices] failed to build requests for user \"uid:%d\", because %s", uid, err.Error())
		return nil, errs.ErrFailedToRequestRemoteApi
	}

	priceResps := make([]*models.LatestCryptocurrencyPriceResponse, 0, len(requests))

	for i := 0; i < len(requests); i++ {
		req := requests[i]
		resp, err := c.httpClient.Do(req)

		if err != nil {
			log.Errorf(core, "[common_http_cryptocurrency_price_data_provider.GetLatestCryptocurrencyPrices] failed to request latest cryptocurrency price data for user \"uid:%d\", because %s", uid, err.Error())
			
			// Check if it's a timeout or network error
			if isTimeoutOrNetworkError(err) {
				return nil, errs.ErrOperationFailed
			}
			
			return nil, errs.ErrFailedToRequestRemoteApi
		}

		defer resp.Body.Close()
		
		var body []byte
		
		// Check Content-Length to ensure we read the complete response
		if resp.ContentLength > 0 {
			body = make([]byte, resp.ContentLength)
			_, err = io.ReadFull(resp.Body, body)
		} else {
			body, err = io.ReadAll(resp.Body)
		}

		if err != nil {
			log.Errorf(core, "[common_http_cryptocurrency_price_data_provider.GetLatestCryptocurrencyPrices] failed to read response body for user \"uid:%d\", because %s", uid, err.Error())
			return nil, errs.ErrFailedToRequestRemoteApi
		}

		log.Debugf(core, "[common_http_cryptocurrency_price_data_provider.GetLatestCryptocurrencyPrices] response#%d is %s", i, body)

		if resp.StatusCode != 200 {
			log.Errorf(core, "[common_http_cryptocurrency_price_data_provider.GetLatestCryptocurrencyPrices] failed to get latest cryptocurrency price data response for user \"uid:%d\", because response code is %d", uid, resp.StatusCode)
			return nil, errs.ErrFailedToRequestRemoteApi
		}

		priceResp, err := c.dataSource.Parse(core, body)

		if err != nil {
			log.Errorf(core, "[common_http_cryptocurrency_price_data_provider.GetLatestCryptocurrencyPrices] failed to parse response for user \"uid:%d\", because %s", uid, err.Error())
			return nil, errs.Or(err, errs.ErrFailedToRequestRemoteApi)
		}

		priceResps = append(priceResps, priceResp)
	}

	lastPriceResponse := priceResps[len(priceResps)-1]
	allPricesMap := make(map[string]string)

	for i := 0; i < len(priceResps); i++ {
		priceResp := priceResps[i]

		for j := 0; j < len(priceResp.Prices); j++ {
			price := priceResp.Prices[j]
			allPricesMap[price.Symbol] = price.Price
		}
	}

	// Add USDT as base currency with price "1"
	allPricesMap["USDT"] = "1"
	allPrices := make(models.LatestCryptocurrencyPriceSlice, 0, len(allPricesMap))

	for symbol, price := range allPricesMap {
		allPrices = append(allPrices, &models.LatestCryptocurrencyPrice{
			Symbol: symbol,
			Price:  price,
		})
	}

	sort.Sort(allPrices)

	finalPriceResponse := &models.LatestCryptocurrencyPriceResponse{
		DataSource:   lastPriceResponse.DataSource,
		ReferenceUrl: lastPriceResponse.ReferenceUrl,
		UpdateTime:   lastPriceResponse.UpdateTime,
		BaseCurrency: "USDT",
		Prices:       allPrices,
	}

	return finalPriceResponse, nil
}

func newCommonHttpCryptocurrencyPriceDataProvider(config *settings.Config, dataSource HttpCryptocurrencyPriceDataSource) *CommonHttpCryptocurrencyPriceDataProvider {
	return &CommonHttpCryptocurrencyPriceDataProvider{
		dataSource: dataSource,
		httpClient: utils.NewHttpClient(config.CryptocurrencyRequestTimeout, config.CryptocurrencyProxy, config.CryptocurrencySkipTLSVerify, settings.GetUserAgent()),
		config:     config,
	}
}

// isTimeoutOrNetworkError checks if the error is a timeout or network-related error
func isTimeoutOrNetworkError(err error) bool {
	if err == nil {
		return false
	}

	errStr := err.Error()
	
	// Check for timeout errors
	if strings.Contains(errStr, "timeout") || 
	   strings.Contains(errStr, "deadline exceeded") ||
	   strings.Contains(errStr, "context deadline exceeded") ||
	   strings.Contains(errStr, "Client.Timeout exceeded") {
		return true
	}

	// Check for network errors
	if urlErr, ok := err.(*url.Error); ok {
		if urlErr.Timeout() {
			return true
		}
		if urlErr.Err != nil {
			if netErr, ok := urlErr.Err.(net.Error); ok {
				if netErr.Timeout() {
					return true
				}
			}
		}
	}

	// Check for context errors
	if err == context.DeadlineExceeded {
		return true
	}

	return false
}

