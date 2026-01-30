package cryptocurrency

import (
	"io"
	"net/http"

	"github.com/mayswind/ezbookkeeping/pkg/core"
	"github.com/mayswind/ezbookkeeping/pkg/errs"
	"github.com/mayswind/ezbookkeeping/pkg/log"
	"github.com/mayswind/ezbookkeeping/pkg/models"
	"github.com/mayswind/ezbookkeeping/pkg/settings"
	"github.com/mayswind/ezbookkeeping/pkg/utils"
)

// HttpCryptocurrencyPriceDataSource defines the interface for http cryptocurrency price data source
type HttpCryptocurrencyPriceDataSource interface {
	// BuildRequests builds the http requests
	BuildRequests(symbols []string, apiKey string) ([]*http.Request, error)

	// Parse parses the response content
	Parse(c core.Context, content []byte) (*models.LatestCryptocurrencyPriceResponse, error)
}

// CommonHttpCryptocurrencyPriceDataProvider represents the common http cryptocurrency price data provider
type CommonHttpCryptocurrencyPriceDataProvider struct {
	dataSource HttpCryptocurrencyPriceDataSource
}

// NewCommonHttpCryptocurrencyPriceDataProvider returns a new common http cryptocurrency price data provider
func NewCommonHttpCryptocurrencyPriceDataProvider(dataSource HttpCryptocurrencyPriceDataSource) *CommonHttpCryptocurrencyPriceDataProvider {
	return &CommonHttpCryptocurrencyPriceDataProvider{
		dataSource: dataSource,
	}
}

// GetLatestCryptocurrencyPrices returns the latest cryptocurrency prices
func (p *CommonHttpCryptocurrencyPriceDataProvider) GetLatestCryptocurrencyPrices(c core.Context, uid int64, config *models.ExternalDataSourceConfig, symbols []string) (*models.LatestCryptocurrencyPriceResponse, error) {
	if len(symbols) == 0 {
		return nil, nil
	}

	requests, err := p.dataSource.BuildRequests(symbols, config.ApiKey)

	if err != nil {
		return nil, err
	}

	client := utils.NewHttpClient(uint32(config.RequestTimeout), config.Proxy, false, settings.GetUserAgent())
	
	// Currently we only support single request for most data sources
	// If we need to support multiple requests (e.g. batch limit), we need to merge results
	if len(requests) == 1 {
		return p.executeRequest(c, client, requests[0])
	}

	// For multiple requests, we need to handle them (implementation for future if needed)
	// For now, return error or handle first one
	if len(requests) > 0 {
		return p.executeRequest(c, client, requests[0])
	}

	return nil, errs.ErrSystemError
}

func (p *CommonHttpCryptocurrencyPriceDataProvider) executeRequest(c core.Context, client *http.Client, req *http.Request) (*models.LatestCryptocurrencyPriceResponse, error) {
	log.Debugf(c, "[cryptocurrency.CommonHttpCryptocurrencyPriceDataProvider] requesting %s", req.URL.String())

	resp, err := client.Do(req)

	if err != nil {
		log.Errorf(c, "[cryptocurrency.CommonHttpCryptocurrencyPriceDataProvider] failed to request cryptocurrency price data for URL %s, because %s", req.URL.String(), err.Error())
		return nil, errs.ErrFailedToRequestRemoteApi
	}

	defer resp.Body.Close()

	content, err := io.ReadAll(resp.Body)

	if err != nil {
		log.Errorf(c, "[cryptocurrency.CommonHttpCryptocurrencyPriceDataProvider] failed to read response body for URL %s, because %s", req.URL.String(), err.Error())
		return nil, errs.ErrFailedToRequestRemoteApi
	}

	if resp.StatusCode != 200 {
		log.Errorf(c, "[cryptocurrency.CommonHttpCryptocurrencyPriceDataProvider] response status code is %d (expected 200) for URL %s, response content is %s", resp.StatusCode, req.URL.String(), string(content))
		return nil, errs.ErrFailedToRequestRemoteApi
	}

	result, err := p.dataSource.Parse(c, content)

	if err != nil {
		log.Errorf(c, "[cryptocurrency.CommonHttpCryptocurrencyPriceDataProvider] failed to parse response, because %s", err.Error())
		return nil, errs.ErrFailedToRequestRemoteApi
	}

	return result, nil
}
