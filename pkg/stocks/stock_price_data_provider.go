package stocks

import (
	"github.com/mayswind/ezbookkeeping/pkg/core"
	"github.com/mayswind/ezbookkeeping/pkg/models"
	"github.com/mayswind/ezbookkeeping/pkg/settings"
)

// StockPriceDataProvider defines the interface for stock price data provider
type StockPriceDataProvider interface {
	// GetLatestStockPrices returns the latest stock prices
	GetLatestStockPrices(c core.Context, uid int64, currentConfig *settings.Config) (*models.LatestStockPriceResponse, error)
}
