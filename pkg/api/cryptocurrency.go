package api

import (
	"github.com/mayswind/ezbookkeeping/pkg/core"
	"github.com/mayswind/ezbookkeeping/pkg/cryptocurrency"
	"github.com/mayswind/ezbookkeeping/pkg/errs"
	"github.com/mayswind/ezbookkeeping/pkg/settings"
)

// CryptocurrencyApi represents cryptocurrency price api
type CryptocurrencyApi struct {
	ApiUsingConfig
}

// Initialize a cryptocurrency price api singleton instance
var (
	Cryptocurrency = &CryptocurrencyApi{
		ApiUsingConfig: ApiUsingConfig{
			container: settings.Container,
		},
	}
)

// LatestCryptocurrencyPriceHandler returns latest cryptocurrency price data
func (a *CryptocurrencyApi) LatestCryptocurrencyPriceHandler(c *core.WebContext) (any, *errs.Error) {
	priceResponse, err := cryptocurrency.Container.GetLatestCryptocurrencyPrices(c, c.GetCurrentUid(), a.CurrentConfig())

	if err != nil {
		return nil, errs.Or(err, errs.ErrOperationFailed)
	}

	return priceResponse, nil
}

