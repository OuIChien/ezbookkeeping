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
	CurrentType  string
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
	// Initialization is now dynamic based on DB config, so we don't set Current here from file config.
	// We just ensure container exists.
	return nil
}

// GetLatestCryptocurrencyPrices returns the latest cryptocurrency prices
func (c *CryptocurrencyPriceDataProviderContainer) GetLatestCryptocurrencyPrices(ctx core.Context, uid int64, config *models.ExternalDataSourceConfig, symbols []string) (*models.LatestCryptocurrencyPriceResponse, error) {
	if config == nil {
		return nil, errs.ErrCryptocurrencyServiceNotEnabled
	}

	c.mu.Lock()
	if c.CurrentType != config.DataSource {
		var provider CryptocurrencyPriceDataProvider

		switch config.DataSource {
		case "coingecko":
			provider = NewCommonHttpCryptocurrencyPriceDataProvider(&CoinGeckoDataSource{})
		default:
			c.mu.Unlock()
			return nil, errs.ErrInvalidCryptocurrencyDataSource
		}

		c.Current = provider
		c.CurrentType = config.DataSource
		c.IsEnabled = true
	}
	
	provider := c.Current
	c.mu.Unlock()

	c.mu.RLock()
	if c.lastResult != nil && time.Since(c.lastTime) < cryptocurrencyPriceCacheTimeout {
		result := c.lastResult
		c.mu.RUnlock()
		return result, nil
	}
	c.mu.RUnlock()

	result, err, _ := c.requestGroup.Do("GetLatestCryptocurrencyPrices", func() (interface{}, error) {
		res, fetchErr := provider.GetLatestCryptocurrencyPrices(ctx, uid, config, symbols)

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
