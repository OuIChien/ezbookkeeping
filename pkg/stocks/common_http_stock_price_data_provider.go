package stocks

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
func (p *CommonHttpStockPriceDataProvider) GetLatestStockPrices(c core.Context, uid int64, currentConfig *settings.Config) (*models.LatestStockPriceResponse, error) {
	if len(currentConfig.StockSymbols) == 0 {
		return nil, nil
	}

	requests, err := p.dataSource.BuildRequests(currentConfig.StockSymbols, currentConfig.StockAPIKey)

	if err != nil {
		return nil, err
	}

	client := utils.NewHttpClient(currentConfig.StockRequestTimeout, currentConfig.StockProxy, currentConfig.StockSkipTLSVerify, settings.GetUserAgent())

	// Currently we only support single request for most data sources
	if len(requests) == 1 {
		return p.executeRequest(c, client, requests[0])
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
		log.Errorf(c, "[stocks.CommonHttpStockPriceDataProvider] failed to request stock price data, because %s", err.Error())
		return nil, errs.ErrFailedToRequestRemoteApi
	}

	defer resp.Body.Close()

	content, err := io.ReadAll(resp.Body)

	if err != nil {
		log.Errorf(c, "[stocks.CommonHttpStockPriceDataProvider] failed to read response body, because %s", err.Error())
		return nil, errs.ErrFailedToRequestRemoteApi
	}

	if resp.StatusCode != 200 {
		log.Errorf(c, "[stocks.CommonHttpStockPriceDataProvider] response status code is not 200, content is %s", string(content))
		return nil, errs.ErrFailedToRequestRemoteApi
	}

	result, err := p.dataSource.Parse(c, content)

	if err != nil {
		log.Errorf(c, "[stocks.CommonHttpStockPriceDataProvider] failed to parse response, because %s", err.Error())
		return nil, errs.ErrFailedToRequestRemoteApi
	}

	return result, nil
}
