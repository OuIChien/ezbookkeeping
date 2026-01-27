package stocks

import (
	"sync"
	"time"

	"golang.org/x/sync/singleflight"

	"github.com/mayswind/ezbookkeeping/pkg/core"
	"github.com/mayswind/ezbookkeeping/pkg/errs"
	"github.com/mayswind/ezbookkeeping/pkg/log"
	"github.com/mayswind/ezbookkeeping/pkg/models"
	"github.com/mayswind/ezbookkeeping/pkg/settings"
)

const (
	stockPriceCacheTimeout = 5 * time.Minute
)

// StockPriceDataProviderContainer contains the stock price data provider
type StockPriceDataProviderContainer struct {
	Current      StockPriceDataProvider
	IsEnabled    bool
	lastResult   *models.LatestStockPriceResponse
	lastTime     time.Time
	mu           sync.RWMutex
	requestGroup singleflight.Group
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

	c.mu.RLock()
	if c.lastResult != nil && time.Since(c.lastTime) < stockPriceCacheTimeout {
		result := c.lastResult
		c.mu.RUnlock()
		return result, nil
	}
	c.mu.RUnlock()

	result, err, _ := c.requestGroup.Do("GetLatestStockPrices", func() (interface{}, error) {
		res, fetchErr := c.Current.GetLatestStockPrices(ctx, uid, currentConfig)

		if fetchErr == nil {
			c.mu.Lock()
			c.lastResult = res
			c.lastTime = time.Now()
			c.mu.Unlock()
		}

		return res, fetchErr
	})

	if err == nil {
		return result.(*models.LatestStockPriceResponse), nil
	}

	c.mu.RLock()
	defer c.mu.RUnlock()

	if c.lastResult != nil {
		log.Warnf(ctx, "[stocks.Container] failed to get latest prices, using stale cache from %s", c.lastTime.Format("2006-01-02 15:04:05"))
		return c.lastResult, nil
	}

	return nil, err
}
