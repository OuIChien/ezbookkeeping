package cryptocurrency

import (
	"github.com/mayswind/ezbookkeeping/pkg/core"
	"github.com/mayswind/ezbookkeeping/pkg/errs"
	"github.com/mayswind/ezbookkeeping/pkg/models"
	"github.com/mayswind/ezbookkeeping/pkg/settings"
)

// CryptocurrencyPriceDataProviderContainer contains the cryptocurrency price data provider
type CryptocurrencyPriceDataProviderContainer struct {
	Current   CryptocurrencyPriceDataProvider
	IsEnabled bool
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

	return c.Current.GetLatestCryptocurrencyPrices(ctx, uid, currentConfig)
}
