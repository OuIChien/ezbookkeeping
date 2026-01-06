package cryptocurrency

import (
	"github.com/mayswind/ezbookkeeping/pkg/core"
	"github.com/mayswind/ezbookkeeping/pkg/models"
	"github.com/mayswind/ezbookkeeping/pkg/settings"
)

// CryptocurrencyPriceDataProvider defines the structure of cryptocurrency price data provider
type CryptocurrencyPriceDataProvider interface {
	// GetLatestCryptocurrencyPrices returns the common response entities
	GetLatestCryptocurrencyPrices(c core.Context, uid int64, currentConfig *settings.Config) (*models.LatestCryptocurrencyPriceResponse, error)
}

