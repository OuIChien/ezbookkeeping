package api

import (
	"github.com/mayswind/ezbookkeeping/pkg/core"
	"github.com/mayswind/ezbookkeeping/pkg/cryptocurrency"
	"github.com/mayswind/ezbookkeeping/pkg/errs"
	"github.com/mayswind/ezbookkeeping/pkg/exchangerates"
	"github.com/mayswind/ezbookkeeping/pkg/log"
	"github.com/mayswind/ezbookkeeping/pkg/models"
	"github.com/mayswind/ezbookkeeping/pkg/services"
	"github.com/mayswind/ezbookkeeping/pkg/settings"
	"github.com/mayswind/ezbookkeeping/pkg/stocks"
	"github.com/mayswind/ezbookkeeping/pkg/utils"
)

// ExchangeRatesApi represents exchange rate api
type ExchangeRatesApi struct {
	ApiUsingConfig
	users                   *services.UserService
	userCustomExchangeRates *services.UserCustomExchangeRatesService
	cryptocurrency          *cryptocurrency.CryptocurrencyPriceDataProviderContainer
	stocks                  *stocks.StockPriceDataProviderContainer
}

// Initialize a exchange rate api singleton instance
var (
	ExchangeRates = &ExchangeRatesApi{
		ApiUsingConfig: ApiUsingConfig{
			container: settings.Container,
		},
		users:                   services.Users,
		userCustomExchangeRates: services.UserCustomExchangeRates,
		cryptocurrency:          cryptocurrency.Container,
		stocks:                  stocks.Container,
	}
)

// LatestExchangeRateHandler returns latest exchange rate data
func (a *ExchangeRatesApi) LatestExchangeRateHandler(c *core.WebContext) (any, *errs.Error) {
	exchangeRateResponse, err := exchangerates.Container.GetLatestExchangeRates(c, c.GetCurrentUid(), a.CurrentConfig())

	if err != nil {
		return nil, errs.Or(err, errs.ErrOperationFailed)
	}

	exchangeRatesMap := make(map[string]float64)
	for _, rate := range exchangeRateResponse.ExchangeRates {
		val, _ := utils.StringToFloat64(rate.Rate)
		exchangeRatesMap[rate.Currency] = val
	}
	exchangeRatesMap[exchangeRateResponse.BaseCurrency] = 1.0

	// Fetch crypto prices
	cryptoPriceResponse, err := a.cryptocurrency.GetLatestCryptocurrencyPrices(c, c.GetCurrentUid(), a.CurrentConfig())
	if err == nil && cryptoPriceResponse != nil {
		rateToCryptoBase, ok := exchangeRatesMap[cryptoPriceResponse.BaseCurrency]
		if ok && rateToCryptoBase != 0 {
			for _, priceData := range cryptoPriceResponse.Prices {
				price, _ := utils.StringToFloat64(priceData.Price)
				if price > 0 {
					// 1 crypto = price cryptoBase
					// 1 base = rateToCryptoBase cryptoBase
					// 1 base = rateToCryptoBase / price crypto
					rate := rateToCryptoBase / price
					exchangeRateResponse.ExchangeRates = append(exchangeRateResponse.ExchangeRates, &models.LatestExchangeRate{
						Currency: priceData.Symbol,
						Rate:     utils.Float64ToString(rate),
					})
				}
			}
		}
	}

	// Fetch stock prices
	stockPriceResponse, err := a.stocks.GetLatestStockPrices(c, c.GetCurrentUid(), a.CurrentConfig())
	if err == nil && stockPriceResponse != nil {
		for _, stockData := range stockPriceResponse.Prices {
			price, _ := utils.StringToFloat64(stockData.Price)
			rateToStockBase, ok := exchangeRatesMap[stockData.Currency]
			if ok && rateToStockBase != 0 && price > 0 {
				// 1 stock = price stockBase
				// 1 base = rateToStockBase stockBase
				// 1 base = rateToStockBase / price stock
				rate := rateToStockBase / price
				exchangeRateResponse.ExchangeRates = append(exchangeRateResponse.ExchangeRates, &models.LatestExchangeRate{
					Currency: stockData.Symbol,
					Rate:     utils.Float64ToString(rate),
				})
			}
		}
	}

	return exchangeRateResponse, nil
}

// UserCustomExchangeRateUpdateHandler updates user custom exchange rates data by request parameters for current user
func (a *ExchangeRatesApi) UserCustomExchangeRateUpdateHandler(c *core.WebContext) (any, *errs.Error) {
	var customExchangeRateUpdateReq models.UserCustomExchangeRateUpdateRequest
	err := c.ShouldBindJSON(&customExchangeRateUpdateReq)

	if err != nil {
		log.Warnf(c, "[exchange_rates.UserCustomExchangeRateUpdateHandler] parse request failed, because %s", err.Error())
		return nil, errs.NewIncompleteOrIncorrectSubmissionError(err)
	}

	uid := c.GetCurrentUid()
	user, err := a.users.GetUserById(c, uid)

	if err != nil {
		log.Errorf(c, "[exchange_rates.UserCustomExchangeRateUpdateHandler] failed to get user \"uid:%d\", because %s", uid, err.Error())
		return nil, errs.Or(err, errs.ErrOperationFailed)
	}

	if customExchangeRateUpdateReq.Currency == user.DefaultCurrency {
		return nil, errs.ErrCannotUpdateExchangeRateForDefaultCurrency
	}

	newCustomExchangeRate, defaultCurrencyExchangeRate, err := a.userCustomExchangeRates.UpdateCustomExchangeRate(c, uid, customExchangeRateUpdateReq.Currency, customExchangeRateUpdateReq.Rate, user.DefaultCurrency)

	if err != nil {
		log.Errorf(c, "[exchange_rates.UserCustomExchangeRateUpdateHandler] failed to update user custom exchange rate \"currency:%s\" for user \"uid:%d\", because %s", customExchangeRateUpdateReq.Currency, uid, err.Error())
		return nil, errs.Or(err, errs.ErrOperationFailed)
	}

	log.Infof(c, "[exchange_rates.UserCustomExchangeRateUpdateHandler] user \"uid:%d\" has updated user custom exchange rate \"currency:%s\" successfully", uid, customExchangeRateUpdateReq.Currency)
	return newCustomExchangeRate.ToUserCustomExchangeRateUpdateResponse(defaultCurrencyExchangeRate.Rate), nil
}

// UserCustomExchangeRateDeleteHandler deletes an existed user custom exchange rates data by request parameters for current user
func (a *ExchangeRatesApi) UserCustomExchangeRateDeleteHandler(c *core.WebContext) (any, *errs.Error) {
	var customExchangeRateDeleteReq models.UserCustomExchangeRateDeleteRequest
	err := c.ShouldBindJSON(&customExchangeRateDeleteReq)

	if err != nil {
		log.Warnf(c, "[exchange_rates.UserCustomExchangeRateDeleteHandler] parse request failed, because %s", err.Error())
		return nil, errs.NewIncompleteOrIncorrectSubmissionError(err)
	}

	uid := c.GetCurrentUid()
	user, err := a.users.GetUserById(c, uid)

	if err != nil {
		log.Errorf(c, "[exchange_rates.UserCustomExchangeRateDeleteHandler] failed to get user \"uid:%d\", because %s", uid, err.Error())
		return nil, errs.Or(err, errs.ErrOperationFailed)
	}

	if customExchangeRateDeleteReq.Currency == user.DefaultCurrency {
		return nil, errs.ErrCannotDeleteExchangeRateForDefaultCurrency
	}

	err = a.userCustomExchangeRates.DeleteCustomExchangeRate(c, uid, customExchangeRateDeleteReq.Currency)

	if err != nil {
		log.Errorf(c, "[exchange_rates.UserCustomExchangeRateDeleteHandler] failed to delete user custom exchange rate \"currency:%s\" for user \"uid:%d\", because %s", customExchangeRateDeleteReq.Currency, uid, err.Error())
		return nil, errs.Or(err, errs.ErrOperationFailed)
	}

	log.Infof(c, "[exchange_rates.UserCustomExchangeRateDeleteHandler] user \"uid:%d\" has deleted user custom exchange rate \"currency:%s\"", uid, customExchangeRateDeleteReq.Currency)
	return true, nil
}
