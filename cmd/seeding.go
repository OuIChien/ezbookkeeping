package cmd

import (
	"time"

	"github.com/mayswind/ezbookkeeping/pkg/core"
	"github.com/mayswind/ezbookkeeping/pkg/datastore"
	"github.com/mayswind/ezbookkeeping/pkg/log"
	"github.com/mayswind/ezbookkeeping/pkg/models"
	"github.com/mayswind/ezbookkeeping/pkg/settings"
)

func seedDefaultData(c *core.CliContext) error {
	err := seedCryptocurrencies(c)
	if err != nil {
		return err
	}

	err = seedStocks(c)
	if err != nil {
		return err
	}

	err = seedExternalDataSourceConfigs(c)
	if err != nil {
		return err
	}

	return nil
}

func seedCryptocurrencies(c *core.CliContext) error {
	db := datastore.Container.UserDataStore.Choose(0)
	sess := db.NewSession(c)
	defer sess.Close()

	count, err := sess.Count(new(models.Cryptocurrency))
	if err != nil {
		return err
	}

	if count > 0 {
		return nil
	}

	log.BootInfof(c, "[seeding.seedCryptocurrencies] seeding default cryptocurrencies")

	now := time.Now().Unix()
	defaultCryptos := []models.Cryptocurrency{
		{Symbol: "BTC", Name: "Bitcoin", DisplayOrder: 1, CreatedUnixTime: now, UpdatedUnixTime: now},
		{Symbol: "ETH", Name: "Ethereum", DisplayOrder: 2, CreatedUnixTime: now, UpdatedUnixTime: now},
		{Symbol: "ATOM", Name: "Cosmos", DisplayOrder: 3, CreatedUnixTime: now, UpdatedUnixTime: now},
		{Symbol: "SOL", Name: "Solana", DisplayOrder: 4, CreatedUnixTime: now, UpdatedUnixTime: now},
		{Symbol: "ADA", Name: "Cardano", DisplayOrder: 5, CreatedUnixTime: now, UpdatedUnixTime: now},
	}

	for _, crypto := range defaultCryptos {
		if _, err := sess.Insert(&crypto); err != nil {
			return err
		}
	}

	return nil
}

func seedStocks(c *core.CliContext) error {
	db := datastore.Container.UserDataStore.Choose(0)
	sess := db.NewSession(c)
	defer sess.Close()

	count, err := sess.Count(new(models.Stock))
	if err != nil {
		return err
	}

	if count > 0 {
		return nil
	}

	log.BootInfof(c, "[seeding.seedStocks] seeding default stocks")

	now := time.Now().Unix()
	defaultStocks := []models.Stock{
		{Symbol: "VOO", Name: "Vanguard S&P 500 ETF", Market: "US", DisplayOrder: 1, CreatedUnixTime: now, UpdatedUnixTime: now},
		{Symbol: "TSLA", Name: "Tesla, Inc.", Market: "US", DisplayOrder: 2, CreatedUnixTime: now, UpdatedUnixTime: now},
		{Symbol: "TSLL", Name: "Direxion Daily TSLA Bull 1.5X Shares", Market: "US", DisplayOrder: 3, CreatedUnixTime: now, UpdatedUnixTime: now},
		{Symbol: "GOOG", Name: "Alphabet Inc.", Market: "US", DisplayOrder: 4, CreatedUnixTime: now, UpdatedUnixTime: now},
		{Symbol: "MSFT", Name: "Microsoft Corporation", Market: "US", DisplayOrder: 5, CreatedUnixTime: now, UpdatedUnixTime: now},
		{Symbol: "AAPL", Name: "Apple Inc.", Market: "US", DisplayOrder: 6, CreatedUnixTime: now, UpdatedUnixTime: now},
		{Symbol: "AMD", Name: "Advanced Micro Devices, Inc.", Market: "US", DisplayOrder: 7, CreatedUnixTime: now, UpdatedUnixTime: now},
		{Symbol: "NVDA", Name: "NVIDIA Corporation", Market: "US", DisplayOrder: 8, CreatedUnixTime: now, UpdatedUnixTime: now},
	}

	for _, stock := range defaultStocks {
		if _, err := sess.Insert(&stock); err != nil {
			return err
		}
	}

	return nil
}

func seedExternalDataSourceConfigs(c *core.CliContext) error {
	db := datastore.Container.UserDataStore.Choose(0)
	sess := db.NewSession(c)
	defer sess.Close()

	count, err := sess.Count(new(models.ExternalDataSourceConfig))
	if err != nil {
		return err
	}

	if count > 0 {
		return nil
	}

	log.BootInfof(c, "[seeding.seedExternalDataSourceConfigs] seeding default external data source configs")

	now := time.Now().Unix()
	defaultConfigs := []models.ExternalDataSourceConfig{
		{
			Type:            models.EXTERNAL_DATA_SOURCE_TYPE_CRYPTOCURRENCY,
			DataSource:      "coingecko", // Default as per design
			BaseCurrency:    "USD",
			RequestTimeout:  10000,
			CreatedUnixTime: now,
			UpdatedUnixTime: now,
		},
		{
			Type:            models.EXTERNAL_DATA_SOURCE_TYPE_STOCK,
			DataSource:      settings.FinancialModelingPrepDataSource, // One request = all symbols; free 250 req/day. API key: https://site.financialmodelingprep.com/developer/docs
			RequestTimeout:  10000,
			CreatedUnixTime: now,
			UpdatedUnixTime: now,
		},
	}

	for _, config := range defaultConfigs {
		if _, err := sess.Insert(&config); err != nil {
			return err
		}
	}

	return nil
}
