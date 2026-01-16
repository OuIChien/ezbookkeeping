package api

import (
	"github.com/mayswind/ezbookkeeping/pkg/core"
	"github.com/mayswind/ezbookkeeping/pkg/errs"
	"github.com/mayswind/ezbookkeeping/pkg/settings"
	"github.com/mayswind/ezbookkeeping/pkg/stocks"
)

// StockApi represents stock api
type StockApi struct {
	ApiUsingConfig
}

// Initialize a stock api singleton instance
var (
	Stocks = &StockApi{
		ApiUsingConfig: ApiUsingConfig{
			container: settings.Container,
		},
	}
)

// LatestStockPriceHandler returns latest stock price data
func (a *StockApi) LatestStockPriceHandler(c *core.WebContext) (any, *errs.Error) {
	stockPriceResponse, err := stocks.Container.GetLatestStockPrices(c, c.GetCurrentUid(), a.CurrentConfig())

	if err != nil {
		return nil, errs.Or(err, errs.ErrOperationFailed)
	}

	return stockPriceResponse, nil
}
