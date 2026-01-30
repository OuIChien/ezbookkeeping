package api

import (
	"github.com/mayswind/ezbookkeeping/pkg/core"
	"github.com/mayswind/ezbookkeeping/pkg/errs"
	"github.com/mayswind/ezbookkeeping/pkg/models"
	"github.com/mayswind/ezbookkeeping/pkg/services"
	"github.com/mayswind/ezbookkeeping/pkg/settings"
	"github.com/mayswind/ezbookkeeping/pkg/stocks"
)

// StockApi represents stock api
type StockApi struct {
	ApiUsingConfig
	externalDataSourceConfigs *services.ExternalDataSourceConfigService
	stocks                    *services.StockService
}

// Initialize a stock api singleton instance
var (
	Stocks = &StockApi{
		ApiUsingConfig: ApiUsingConfig{
			container: settings.Container,
		},
		externalDataSourceConfigs: services.ExternalDataSourceConfigs,
		stocks:                    services.Stocks,
	}
)

// LatestStockPriceHandler returns latest stock price data
func (a *StockApi) LatestStockPriceHandler(c *core.WebContext) (any, *errs.Error) {
	config, err := a.externalDataSourceConfigs.GetConfig(c, models.EXTERNAL_DATA_SOURCE_TYPE_STOCK)
	
	if err != nil {
		return nil, errs.Or(err, errs.ErrOperationFailed)
	}

	stockList, err := a.stocks.GetAllVisibleStocks(c)
	
	if err != nil {
		return nil, errs.Or(err, errs.ErrOperationFailed)
	}

	symbols := make([]string, len(stockList))
	for i, stock := range stockList {
		symbols[i] = stock.Symbol
	}

	stockPriceResponse, err := stocks.Container.GetLatestStockPrices(c, c.GetCurrentUid(), config, symbols)

	if err != nil {
		return nil, errs.Or(err, errs.ErrOperationFailed)
	}

	return stockPriceResponse, nil
}

// StockListHandler returns stock list
func (a *StockApi) StockListHandler(c *core.WebContext) (any, *errs.Error) {
	stockList, err := a.stocks.GetAllStocks(c)
	if err != nil {
		return nil, errs.Or(err, errs.ErrOperationFailed)
	}

	responses := make([]*models.StockInfoResponse, len(stockList))
	for i, stock := range stockList {
		responses[i] = stock.ToStockInfoResponse()
	}

	return responses, nil
}

// StockAddHandler adds a new stock
func (a *StockApi) StockAddHandler(c *core.WebContext) (any, *errs.Error) {
	var req models.StockCreateRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		return nil, errs.NewIncompleteOrIncorrectSubmissionError(err)
	}

	stock := &models.Stock{
		Symbol:       req.Symbol,
		Name:         req.Name,
		Market:       req.Market,
		DisplayOrder: req.DisplayOrder,
		IsHidden:     false,
	}

	if err := a.stocks.CreateStock(c, stock); err != nil {
		return nil, errs.Or(err, errs.ErrOperationFailed)
	}

	return stock.ToStockInfoResponse(), nil
}

// StockModifyHandler modifies a stock
func (a *StockApi) StockModifyHandler(c *core.WebContext) (any, *errs.Error) {
	var req models.StockModifyRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		return nil, errs.NewIncompleteOrIncorrectSubmissionError(err)
	}

	stock := &models.Stock{
		Symbol:       req.Symbol,
		Name:         req.Name,
		Market:       req.Market,
		IsHidden:     req.IsHidden,
		DisplayOrder: req.DisplayOrder,
	}

	if err := a.stocks.UpdateStock(c, stock); err != nil {
		return nil, errs.Or(err, errs.ErrOperationFailed)
	}

	return stock.ToStockInfoResponse(), nil
}

// StockHideHandler hides or unhides a stock
func (a *StockApi) StockHideHandler(c *core.WebContext) (any, *errs.Error) {
	var req models.StockHideRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		return nil, errs.NewIncompleteOrIncorrectSubmissionError(err)
	}

	if err := a.stocks.HideStock(c, req.Symbol, req.Hidden); err != nil {
		return nil, errs.Or(err, errs.ErrOperationFailed)
	}

	return true, nil
}

// StockDeleteHandler deletes a stock
func (a *StockApi) StockDeleteHandler(c *core.WebContext) (any, *errs.Error) {
	var req models.StockDeleteRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		return nil, errs.NewIncompleteOrIncorrectSubmissionError(err)
	}

	if err := a.stocks.DeleteStock(c, req.Symbol); err != nil {
		return nil, errs.Or(err, errs.ErrOperationFailed)
	}

	return true, nil
}

// StockConfigGetHandler returns stock config
func (a *StockApi) StockConfigGetHandler(c *core.WebContext) (any, *errs.Error) {
	config, err := a.externalDataSourceConfigs.GetConfig(c, models.EXTERNAL_DATA_SOURCE_TYPE_STOCK)
	if err != nil {
		return nil, errs.Or(err, errs.ErrOperationFailed)
	}

	if config == nil {
		return nil, nil // Or default
	}

	return config.ToExternalDataSourceConfigResponse(), nil
}

// StockConfigSaveHandler saves stock config
func (a *StockApi) StockConfigSaveHandler(c *core.WebContext) (any, *errs.Error) {
	var req models.ExternalDataSourceConfigSaveRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		return nil, errs.NewIncompleteOrIncorrectSubmissionError(err)
	}

	if req.Type != models.EXTERNAL_DATA_SOURCE_TYPE_STOCK {
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
