package stocks

import (
	"io"
	"net/http"
	"strings"

	"github.com/mayswind/ezbookkeeping/pkg/core"
	"github.com/mayswind/ezbookkeeping/pkg/errs"
	"github.com/mayswind/ezbookkeeping/pkg/log"
	"github.com/mayswind/ezbookkeeping/pkg/models"
	"github.com/mayswind/ezbookkeeping/pkg/settings"
	"github.com/mayswind/ezbookkeeping/pkg/utils"
)

// HttpStockPriceDataSource defines the interface for http stock price data source
type HttpStockPriceDataSource interface {
	// BuildRequests builds the http requests
	BuildRequests(symbols []string, apiKey string) ([]*http.Request, error)

	// Parse parses the response content
	Parse(c core.Context, content []byte) (*models.LatestStockPriceResponse, error)
}

// CommonHttpStockPriceDataProvider represents the common http stock price data provider
type CommonHttpStockPriceDataProvider struct {
	dataSource HttpStockPriceDataSource
}

// NewCommonHttpStockPriceDataProvider returns a new common http stock price data provider
func NewCommonHttpStockPriceDataProvider(dataSource HttpStockPriceDataSource) *CommonHttpStockPriceDataProvider {
	return &CommonHttpStockPriceDataProvider{
		dataSource: dataSource,
	}
}

// GetLatestStockPrices returns the latest stock prices
func (p *CommonHttpStockPriceDataProvider) GetLatestStockPrices(c core.Context, uid int64, config *models.ExternalDataSourceConfig, symbols []string) (*models.LatestStockPriceResponse, error) {
	if len(symbols) == 0 {
		return nil, nil
	}

	requests, err := p.dataSource.BuildRequests(symbols, config.ApiKey)

	if err != nil {
		return nil, err
	}

	client := utils.NewHttpClient(uint32(config.RequestTimeout), config.Proxy, false, settings.GetUserAgent())

	// Currently we only support single request for most data sources,
	// but some (like Alpha Vantage) may require multiple requests.
	if len(requests) == 1 {
		return p.executeRequest(c, client, requests[0])
	}

	if len(requests) > 1 {
		finalResult := &models.LatestStockPriceResponse{
			Prices: make(models.LatestStockPriceSlice, 0),
		}

		for _, req := range requests {
			result, err := p.executeRequest(c, client, req)

			if err != nil {
				log.Warnf(c, "[stocks.CommonHttpStockPriceDataProvider] failed to request stock price data for %s, because %s", req.URL.String(), err.Error())
				continue
			}

			if result != nil {
				finalResult.DataSource = result.DataSource
				finalResult.ReferenceUrl = result.ReferenceUrl
				finalResult.UpdateTime = result.UpdateTime
				finalResult.BaseCurrency = result.BaseCurrency
				finalResult.Prices = append(finalResult.Prices, result.Prices...)
			}
		}

		if len(finalResult.Prices) == 0 {
			return nil, errs.ErrFailedToRequestRemoteApi
		}

		return finalResult, nil
	}

	if len(requests) > 0 {
		return p.executeRequest(c, client, requests[0])
	}

	return nil, errs.ErrSystemError
}

func (p *CommonHttpStockPriceDataProvider) executeRequest(c core.Context, client *http.Client, req *http.Request) (*models.LatestStockPriceResponse, error) {
	log.Debugf(c, "[stocks.CommonHttpStockPriceDataProvider] requesting %s", req.URL.String())

	resp, err := client.Do(req)

	if err != nil {
		log.Errorf(c, "[stocks.CommonHttpStockPriceDataProvider] failed to request stock price data for URL %s, because %s", req.URL.String(), err.Error())
		return nil, errs.ErrFailedToRequestRemoteApi
	}

	defer resp.Body.Close()

	content, err := io.ReadAll(resp.Body)

	if err != nil {
		log.Errorf(c, "[stocks.CommonHttpStockPriceDataProvider] failed to read response body for URL %s, because %s", req.URL.String(), err.Error())
		return nil, errs.ErrFailedToRequestRemoteApi
	}

	if resp.StatusCode != 200 {
		log.Errorf(c, "[stocks.CommonHttpStockPriceDataProvider] response status code is %d (expected 200) for URL %s, response content is %s", resp.StatusCode, req.URL.String(), string(content))
		if resp.StatusCode == 401 && strings.Contains(req.URL.Host, "yahoo") {
			log.Warnf(c, "[stocks.CommonHttpStockPriceDataProvider] Yahoo Finance public API has been restricted (401). Please switch to Alpha Vantage in Settings -> Stock Prices -> Data Source, and set a free API key from https://www.alphavantage.co/support/#api-key")
		}
		return nil, errs.ErrFailedToRequestRemoteApi
	}

	result, err := p.dataSource.Parse(c, content)

	if err != nil {
		log.Errorf(c, "[stocks.CommonHttpStockPriceDataProvider] failed to parse response, because %s", err.Error())
		return nil, errs.ErrFailedToRequestRemoteApi
	}

	return result, nil
}
