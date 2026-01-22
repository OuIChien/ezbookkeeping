package settings

import (
	"strings"

	"gopkg.in/ini.v1"

	"github.com/mayswind/ezbookkeeping/pkg/errs"
)

func loadStockConfiguration(config *Config, configFile *ini.File, sectionName string) error {
	dataSource := getConfigItemStringValue(configFile, sectionName, "data_source")

	if dataSource == "" {
		config.StockDataSource = ""
	} else if dataSource == YahooFinanceDataSource || dataSource == AlphaVantageDataSource {
		config.StockDataSource = dataSource
	} else {
		return errs.ErrInvalidStockDataSource
	}

	stocks := getConfigItemStringValue(configFile, sectionName, "stocks")

	if stocks != "" {
		config.StockSymbols = strings.Split(stocks, ",")
	} else {
		config.StockSymbols = nil
	}

	config.StockRequestTimeout = getConfigItemUint32Value(configFile, sectionName, "request_timeout", defaultExchangeRatesDataRequestTimeout) // Reuse exchange rates default timeout
	config.StockProxy = getConfigItemStringValue(configFile, sectionName, "proxy", "system")
	config.StockSkipTLSVerify = getConfigItemBoolValue(configFile, sectionName, "skip_tls_verify", false)
	config.StockAPIKey = getConfigItemStringValue(configFile, sectionName, "api_key")

	return nil
}
