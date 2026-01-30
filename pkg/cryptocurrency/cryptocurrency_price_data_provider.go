package cryptocurrency

import (
	"github.com/mayswind/ezbookkeeping/pkg/core"
	"github.com/mayswind/ezbookkeeping/pkg/models"
)

// CryptocurrencyPriceDataProvider defines the interface for cryptocurrency price data provider
type CryptocurrencyPriceDataProvider interface {
	// GetLatestCryptocurrencyPrices returns the latest cryptocurrency prices
	GetLatestCryptocurrencyPrices(c core.Context, uid int64, config *models.ExternalDataSourceConfig, symbols []string) (*models.LatestCryptocurrencyPriceResponse, error)
}
