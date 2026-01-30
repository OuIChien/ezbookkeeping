package cron

import (
	"time"

	"github.com/mayswind/ezbookkeeping/pkg/core"
	"github.com/mayswind/ezbookkeeping/pkg/cryptocurrency"
	"github.com/mayswind/ezbookkeeping/pkg/exchangerates"
	"github.com/mayswind/ezbookkeeping/pkg/models"
	"github.com/mayswind/ezbookkeeping/pkg/services"
	"github.com/mayswind/ezbookkeeping/pkg/settings"
	"github.com/mayswind/ezbookkeeping/pkg/stocks"
)

// RemoveExpiredTokensJob represents the cron job which periodically remove expired user tokens from the database
var RemoveExpiredTokensJob = &CronJob{
	Name:        "RemoveExpiredTokens",
	Description: "Periodically remove expired user tokens from the database.",
	Period: CronJobFixedHourPeriod{
		Hour: 0,
	},
	Run: func(c *core.CronContext) error {
		return services.Tokens.DeleteAllExpiredTokens(c)
	},
}

// CreateScheduledTransactionJob represents the cron job which periodically create transaction by scheduled transaction template
var CreateScheduledTransactionJob = &CronJob{
	Name:        "CreateScheduledTransaction",
	Description: "Periodically create transaction by scheduled transaction template.",
	Period: CronJobEvery15MinutesPeriod{
		Second: 0,
	},
	Run: func(c *core.CronContext) error {
		return services.Transactions.CreateScheduledTransactions(c, time.Now().Unix(), c.GetInterval())
	},
}

// UpdateCryptocurrencyPricesJob represents the cron job which periodically update cryptocurrency prices
var UpdateCryptocurrencyPricesJob = &CronJob{
	Name:        "UpdateCryptocurrencyPrices",
	Description: "Periodically update cryptocurrency prices.",
	Period: CronJobIntervalPeriod{
		Interval: 5 * time.Minute,
	},
	Run: func(c *core.CronContext) error {
		config, err := services.ExternalDataSourceConfigs.GetConfig(c, models.EXTERNAL_DATA_SOURCE_TYPE_CRYPTOCURRENCY)

		if err != nil {
			return err
		}

		cryptos, err := services.Cryptocurrencies.GetAllVisibleCryptocurrencies(c)

		if err != nil {
			return err
		}

		symbols := make([]string, len(cryptos))
		for i, crypto := range cryptos {
			symbols[i] = crypto.Symbol
		}

		_, err = cryptocurrency.Container.GetLatestCryptocurrencyPrices(c, 0, config, symbols)
		return err
	},
}

// UpdateStockPricesJob represents the cron job which periodically update stock prices
var UpdateStockPricesJob = &CronJob{
	Name:        "UpdateStockPrices",
	Description: "Periodically update stock prices.",
	Period: CronJobIntervalPeriod{
		Interval: 5 * time.Minute,
	},
	Run: func(c *core.CronContext) error {
		config, err := services.ExternalDataSourceConfigs.GetConfig(c, models.EXTERNAL_DATA_SOURCE_TYPE_STOCK)

		if err != nil {
			return err
		}

		stockList, err := services.Stocks.GetAllVisibleStocks(c)

		if err != nil {
			return err
		}

		symbols := make([]string, len(stockList))
		for i, stock := range stockList {
			symbols[i] = stock.Symbol
		}

		_, err = stocks.Container.GetLatestStockPrices(c, 0, config, symbols)
		return err
	},
}

// UpdateExchangeRatesJob represents the cron job which periodically update exchange rates
var UpdateExchangeRatesJob = &CronJob{
	Name:        "UpdateExchangeRates",
	Description: "Periodically update exchange rates.",
	Period: CronJobIntervalPeriod{
		Interval: 5 * time.Minute,
	},
	Run: func(c *core.CronContext) error {
		_, err := exchangerates.Container.GetLatestExchangeRates(c, 0, settings.Container.GetCurrentConfig())
		return err
	},
}
