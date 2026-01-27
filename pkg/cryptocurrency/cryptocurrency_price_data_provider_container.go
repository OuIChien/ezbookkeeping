package cryptocurrency

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
	cryptocurrencyPriceCacheTimeout = 5 * time.Minute
)

// CryptocurrencyPriceDataProviderContainer contains the cryptocurrency price data provider
type CryptocurrencyPriceDataProviderContainer struct {
	Current      CryptocurrencyPriceDataProvider
	IsEnabled    bool
	lastResult   *models.LatestCryptocurrencyPriceResponse
	lastTime     time.Time
	mu           sync.RWMutex
	requestGroup singleflight.Group
}

// Initialize a cryptocurrency price data provider container singleton instance
var (
	Container = &CryptocurrencyPriceDataProviderContainer{}
)

// InitializeCryptocurrencyPriceDataProvider initializes the cryptocurrency price data provider
func InitializeCryptocurrencyPriceDataProvider(config *settings.Config) error {
	if config.CryptocurrencyDataSource == "" {
		Container.IsEnabled = false
		return nil
	}

	var provider CryptocurrencyPriceDataProvider

	switch config.CryptocurrencyDataSource {
	case settings.CoinGeckoDataSource:
		provider = NewCommonHttpCryptocurrencyPriceDataProvider(&CoinGeckoDataSource{})
	default:
		return errs.ErrInvalidCryptocurrencyDataSource
	}

	Container.Current = provider
	Container.IsEnabled = true

	return nil
}

// GetLatestCryptocurrencyPrices returns the latest cryptocurrency prices
func (c *CryptocurrencyPriceDataProviderContainer) GetLatestCryptocurrencyPrices(ctx core.Context, uid int64, currentConfig *settings.Config) (*models.LatestCryptocurrencyPriceResponse, error) {
	if !c.IsEnabled {
		return nil, errs.ErrCryptocurrencyServiceNotEnabled
	}

	c.mu.RLock()
	if c.lastResult != nil && time.Since(c.lastTime) < cryptocurrencyPriceCacheTimeout {
		result := c.lastResult
		c.mu.RUnlock()
		return result, nil
	}
	c.mu.RUnlock()

	result, err, _ := c.requestGroup.Do("GetLatestCryptocurrencyPrices", func() (interface{}, error) {
		res, fetchErr := c.Current.GetLatestCryptocurrencyPrices(ctx, uid, currentConfig)

		if fetchErr == nil {
			c.mu.Lock()
			c.lastResult = res
			c.lastTime = time.Now()
			c.mu.Unlock()
		}

		return res, fetchErr
	})

	if err == nil {
		return result.(*models.LatestCryptocurrencyPriceResponse), nil
	}

	c.mu.RLock()
	defer c.mu.RUnlock()

	if c.lastResult != nil {
		log.Warnf(ctx, "[cryptocurrency.Container] failed to get latest prices, using stale cache from %s", c.lastTime.Format("2006-01-02 15:04:05"))
		return c.lastResult, nil
	}

	return nil, err
}
