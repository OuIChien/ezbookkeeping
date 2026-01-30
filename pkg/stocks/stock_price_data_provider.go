package stocks

import (
	"github.com/mayswind/ezbookkeeping/pkg/core"
	"github.com/mayswind/ezbookkeeping/pkg/models"
)

// StockPriceDataProvider defines the interface for stock price data provider
type StockPriceDataProvider interface {
	// GetLatestStockPrices returns the latest stock prices
	GetLatestStockPrices(c core.Context, uid int64, config *models.ExternalDataSourceConfig, symbols []string) (*models.LatestStockPriceResponse, error)
}
