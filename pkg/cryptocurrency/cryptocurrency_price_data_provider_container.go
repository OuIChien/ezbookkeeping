package cryptocurrency

import (
	"github.com/mayswind/ezbookkeeping/pkg/core"
	"github.com/mayswind/ezbookkeeping/pkg/errs"
	"github.com/mayswind/ezbookkeeping/pkg/models"
	"github.com/mayswind/ezbookkeeping/pkg/settings"
)

// CryptocurrencyPriceDataProviderContainer contains the current cryptocurrency price data provider
type CryptocurrencyPriceDataProviderContainer struct {
	current CryptocurrencyPriceDataProvider
}

// Initialize a cryptocurrency price data provider container singleton instance
var (
	Container = &CryptocurrencyPriceDataProviderContainer{}
)

// InitializeCryptocurrencyDataSource initializes the current cryptocurrency price data source according to the config
func InitializeCryptocurrencyDataSource(config *settings.Config) error {
	if config.CryptocurrencyDataSource == "" {
		// Cryptocurrency feature is optional, return nil if not configured
		Container.current = nil
		return nil
	}

	if config.CryptocurrencyDataSource == settings.CoinGeckoDataSource {
		Container.current = newCommonHttpCryptocurrencyPriceDataProvider(config, &CoinGeckoDataSource{})
		return nil
	} else if config.CryptocurrencyDataSource == settings.CoinMarketCapDataSource {
		// TODO: Implement CoinMarketCap data source
		return errs.ErrInvalidCryptocurrencyDataSource
	} else if config.CryptocurrencyDataSource == settings.BinanceDataSource {
		// TODO: Implement Binance data source
		return errs.ErrInvalidCryptocurrencyDataSource
	}

	return errs.ErrInvalidCryptocurrencyDataSource
}

// GetLatestCryptocurrencyPrices returns the latest cryptocurrency prices data from the current data source
func (c *CryptocurrencyPriceDataProviderContainer) GetLatestCryptocurrencyPrices(core core.Context, uid int64, currentConfig *settings.Config) (*models.LatestCryptocurrencyPriceResponse, error) {
	if Container.current == nil {
		return nil, errs.ErrInvalidCryptocurrencyDataSource
	}

	return Container.current.GetLatestCryptocurrencyPrices(core, uid, currentConfig)
}

