package api

import (
	"github.com/mayswind/ezbookkeeping/pkg/core"
	"github.com/mayswind/ezbookkeeping/pkg/cryptocurrency"
	"github.com/mayswind/ezbookkeeping/pkg/errs"
	"github.com/mayswind/ezbookkeeping/pkg/models"
	"github.com/mayswind/ezbookkeeping/pkg/services"
	"github.com/mayswind/ezbookkeeping/pkg/settings"
)

// CryptocurrencyApi represents cryptocurrency api
type CryptocurrencyApi struct {
	ApiUsingConfig
	externalDataSourceConfigs *services.ExternalDataSourceConfigService
	cryptocurrencies          *services.CryptocurrencyService
}

// Initialize a cryptocurrency api singleton instance
var (
	Cryptocurrencies = &CryptocurrencyApi{
		ApiUsingConfig: ApiUsingConfig{
			container: settings.Container,
		},
		externalDataSourceConfigs: services.ExternalDataSourceConfigs,
		cryptocurrencies:          services.Cryptocurrencies,
	}
)

// LatestCryptocurrencyPriceHandler returns latest cryptocurrency price data
func (a *CryptocurrencyApi) LatestCryptocurrencyPriceHandler(c *core.WebContext) (any, *errs.Error) {
	config, err := a.externalDataSourceConfigs.GetConfig(c, models.EXTERNAL_DATA_SOURCE_TYPE_CRYPTOCURRENCY)
	
	if err != nil {
		return nil, errs.Or(err, errs.ErrOperationFailed)
	}

	cryptos, err := a.cryptocurrencies.GetAllVisibleCryptocurrencies(c)
	
	if err != nil {
		return nil, errs.Or(err, errs.ErrOperationFailed)
	}

	symbols := make([]string, len(cryptos))
	for i, crypto := range cryptos {
		symbols[i] = crypto.Symbol
	}

	cryptocurrencyPriceResponse, err := cryptocurrency.Container.GetLatestCryptocurrencyPrices(c, c.GetCurrentUid(), config, symbols)

	if err != nil {
		return nil, errs.Or(err, errs.ErrOperationFailed)
	}

	return cryptocurrencyPriceResponse, nil
}

// CryptocurrencyListHandler returns cryptocurrency list
func (a *CryptocurrencyApi) CryptocurrencyListHandler(c *core.WebContext) (any, *errs.Error) {
	cryptos, err := a.cryptocurrencies.GetAllCryptocurrencies(c)
	if err != nil {
		return nil, errs.Or(err, errs.ErrOperationFailed)
	}

	responses := make([]*models.CryptocurrencyInfoResponse, len(cryptos))
	for i, crypto := range cryptos {
		responses[i] = crypto.ToCryptocurrencyInfoResponse()
	}

	return responses, nil
}

// CryptocurrencyAddHandler adds a new cryptocurrency
func (a *CryptocurrencyApi) CryptocurrencyAddHandler(c *core.WebContext) (any, *errs.Error) {
	var req models.CryptocurrencyCreateRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		return nil, errs.NewIncompleteOrIncorrectSubmissionError(err)
	}

	crypto := &models.Cryptocurrency{
		Symbol:       req.Symbol,
		Name:         req.Name,
		DisplayOrder: req.DisplayOrder,
		IsHidden:     false,
	}

	if err := a.cryptocurrencies.CreateCryptocurrency(c, crypto); err != nil {
		return nil, errs.Or(err, errs.ErrOperationFailed)
	}

	return crypto.ToCryptocurrencyInfoResponse(), nil
}

// CryptocurrencyModifyHandler modifies a cryptocurrency
func (a *CryptocurrencyApi) CryptocurrencyModifyHandler(c *core.WebContext) (any, *errs.Error) {
	var req models.CryptocurrencyModifyRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		return nil, errs.NewIncompleteOrIncorrectSubmissionError(err)
	}

	crypto := &models.Cryptocurrency{
		Symbol:       req.Symbol,
		Name:         req.Name,
		IsHidden:     req.IsHidden,
		DisplayOrder: req.DisplayOrder,
	}

	if err := a.cryptocurrencies.UpdateCryptocurrency(c, crypto); err != nil {
		return nil, errs.Or(err, errs.ErrOperationFailed)
	}

	return crypto.ToCryptocurrencyInfoResponse(), nil
}

// CryptocurrencyHideHandler hides or unhides a cryptocurrency
func (a *CryptocurrencyApi) CryptocurrencyHideHandler(c *core.WebContext) (any, *errs.Error) {
	var req models.CryptocurrencyHideRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		return nil, errs.NewIncompleteOrIncorrectSubmissionError(err)
	}

	if err := a.cryptocurrencies.HideCryptocurrency(c, req.Symbol, req.Hidden); err != nil {
		return nil, errs.Or(err, errs.ErrOperationFailed)
	}

	return true, nil
}

// CryptocurrencyDeleteHandler deletes a cryptocurrency
func (a *CryptocurrencyApi) CryptocurrencyDeleteHandler(c *core.WebContext) (any, *errs.Error) {
	var req models.CryptocurrencyDeleteRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		return nil, errs.NewIncompleteOrIncorrectSubmissionError(err)
	}

	if err := a.cryptocurrencies.DeleteCryptocurrency(c, req.Symbol); err != nil {
		return nil, errs.Or(err, errs.ErrOperationFailed)
	}

	return true, nil
}

// CryptocurrencyConfigGetHandler returns cryptocurrency config
func (a *CryptocurrencyApi) CryptocurrencyConfigGetHandler(c *core.WebContext) (any, *errs.Error) {
	config, err := a.externalDataSourceConfigs.GetConfig(c, models.EXTERNAL_DATA_SOURCE_TYPE_CRYPTOCURRENCY)
	if err != nil {
		return nil, errs.Or(err, errs.ErrOperationFailed)
	}

	if config == nil {
		return nil, nil // Or default
	}

	return config.ToExternalDataSourceConfigResponse(), nil
}

// CryptocurrencyConfigSaveHandler saves cryptocurrency config
func (a *CryptocurrencyApi) CryptocurrencyConfigSaveHandler(c *core.WebContext) (any, *errs.Error) {
	var req models.ExternalDataSourceConfigSaveRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		return nil, errs.NewIncompleteOrIncorrectSubmissionError(err)
	}

	if req.Type != models.EXTERNAL_DATA_SOURCE_TYPE_CRYPTOCURRENCY {
		return nil, errs.ErrIncompleteOrIncorrectSubmission
	}

	config := &models.ExternalDataSourceConfig{
		Type:            req.Type,
		DataSource:      req.DataSource,
		BaseCurrency:    req.BaseCurrency,
		ApiKey:          req.ApiKey,
		RequestTimeout:  req.RequestTimeout,
		Proxy:           req.Proxy,
		UpdateFrequency: req.UpdateFrequency,
	}

	if err := a.externalDataSourceConfigs.SaveConfig(c, config); err != nil {
		return nil, errs.Or(err, errs.ErrOperationFailed)
	}

	return config.ToExternalDataSourceConfigResponse(), nil
}
