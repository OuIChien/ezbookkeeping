package settings

import (
	"strings"

	"gopkg.in/ini.v1"

	"github.com/mayswind/ezbookkeeping/pkg/errs"
)

func loadCryptocurrencyConfiguration(config *Config, configFile *ini.File, sectionName string) error {
	dataSource := getConfigItemStringValue(configFile, sectionName, "data_source")

	if dataSource == "" {
		config.CryptocurrencyDataSource = ""
	} else if dataSource == CoinGeckoDataSource {
		config.CryptocurrencyDataSource = dataSource
	} else {
		return errs.ErrInvalidCryptocurrencyDataSource
	}

	cryptocurrencies := getConfigItemStringValue(configFile, sectionName, "cryptocurrencies")

	if cryptocurrencies != "" {
		config.CryptocurrencySymbols = strings.Split(cryptocurrencies, ",")
	} else {
		config.CryptocurrencySymbols = nil
	}

	config.CryptocurrencyRequestTimeout = getConfigItemUint32Value(configFile, sectionName, "request_timeout", defaultExchangeRatesDataRequestTimeout) // Reuse exchange rates default timeout
	config.CryptocurrencyProxy = getConfigItemStringValue(configFile, sectionName, "proxy", "system")
	config.CryptocurrencySkipTLSVerify = getConfigItemBoolValue(configFile, sectionName, "skip_tls_verify", false)
	config.CryptocurrencyAPIKey = getConfigItemStringValue(configFile, sectionName, "api_key")

	return nil
}
