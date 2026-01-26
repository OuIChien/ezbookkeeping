package stocks

import (
	"sync"
	"time"

	"github.com/mayswind/ezbookkeeping/pkg/core"
	"github.com/mayswind/ezbookkeeping/pkg/errs"
	"github.com/mayswind/ezbookkeeping/pkg/log"
	"github.com/mayswind/ezbookkeeping/pkg/models"
	"github.com/mayswind/ezbookkeeping/pkg/settings"
)

// StockPriceDataProviderContainer contains the stock price data provider
type StockPriceDataProviderContainer struct {
	Current    StockPriceDataProvider
	IsEnabled  bool
	lastResult *models.LatestStockPriceResponse
	lastTime   time.Time
	mu         sync.RWMutex
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
	case settings.AlphaVantageDataSource:
		provider = NewCommonHttpStockPriceDataProvider(&AlphaVantageDataSource{})
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

	result, err := c.Current.GetLatestStockPrices(ctx, uid, currentConfig)

	c.mu.Lock()
	defer c.mu.Unlock()

	if err == nil {
		c.lastResult = result
		c.lastTime = time.Now()
		return result, nil
	}

	if c.lastResult != nil {
		log.Warnf(ctx, "[stocks.Container] failed to get latest prices, using stale cache from %s", c.lastTime.Format("2006-01-02 15:04:05"))
		return c.lastResult, nil
	}

	return nil, err
}
