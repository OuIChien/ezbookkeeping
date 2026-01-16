package stocks

import (
	"github.com/mayswind/ezbookkeeping/pkg/core"
	"github.com/mayswind/ezbookkeeping/pkg/errs"
	"github.com/mayswind/ezbookkeeping/pkg/models"
	"github.com/mayswind/ezbookkeeping/pkg/settings"
)

// StockPriceDataProviderContainer contains the stock price data provider
type StockPriceDataProviderContainer struct {
	Current   StockPriceDataProvider
	IsEnabled bool
}

// Initialize a stock price data provider container singleton instance
var (
	Container = &StockPriceDataProviderContainer{}
)

// InitializeStockPriceDataProvider initializes the stock price data provider
func InitializeStockPriceDataProvider(config *settings.Config) error {
	if config.StockDataSource == "" {
		Container.IsEnabled = false
		return nil
	}

	var provider StockPriceDataProvider

	switch config.StockDataSource {
	case settings.YahooFinanceDataSource:
		provider = NewCommonHttpStockPriceDataProvider(&YahooFinanceDataSource{})
	default:
		return errs.ErrInvalidStockDataSource
	}

	Container.Current = provider
	Container.IsEnabled = true

	return nil
}

// GetLatestStockPrices returns the latest stock prices
func (c *StockPriceDataProviderContainer) GetLatestStockPrices(ctx core.Context, uid int64, currentConfig *settings.Config) (*models.LatestStockPriceResponse, error) {
	if !c.IsEnabled {
		return nil, errs.ErrStockServiceNotEnabled
	}

	return c.Current.GetLatestStockPrices(ctx, uid, currentConfig)
}
