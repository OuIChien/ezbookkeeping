package exchangerates

import (
	"sync"
	"time"

	"golang.org/x/sync/singleflight"

	"github.com/mayswind/ezbookkeeping/pkg/core"
	"github.com/mayswind/ezbookkeeping/pkg/errs"
	"github.com/mayswind/ezbookkeeping/pkg/log"
	"github.com/mayswind/ezbookkeeping/pkg/models"
	"github.com/mayswind/ezbookkeeping/pkg/settings"
)

const (
	exchangeRatesCacheTimeout = 5 * time.Minute
)

// ExchangeRatesDataProviderContainer contains the current exchange rates data provider
type ExchangeRatesDataProviderContainer struct {
	current      ExchangeRatesDataProvider
	isCustom     bool
	lastResult   *models.LatestExchangeRateResponse
	lastTime     time.Time
	mu           sync.RWMutex
	requestGroup singleflight.Group
}

// Initialize a exchange rates data provider container singleton instance
var (
	Container = &ExchangeRatesDataProviderContainer{}
)

// InitializeExchangeRatesDataSource initializes the current exchange rates data source according to the config
func InitializeExchangeRatesDataSource(config *settings.Config) error {
	if config.ExchangeRatesDataSource == settings.ReserveBankOfAustraliaDataSource {
		Container.current = newCommonHttpExchangeRatesDataProvider(config, &ReserveBankOfAustraliaDataSource{})
		Container.isCustom = false
		return nil
	} else if config.ExchangeRatesDataSource == settings.BankOfCanadaDataSource {
		Container.current = newCommonHttpExchangeRatesDataProvider(config, &BankOfCanadaDataSource{})
		Container.isCustom = false
		return nil
	} else if config.ExchangeRatesDataSource == settings.CzechNationalBankDataSource {
		Container.current = newCommonHttpExchangeRatesDataProvider(config, &CzechNationalBankDataSource{})
		Container.isCustom = false
		return nil
	} else if config.ExchangeRatesDataSource == settings.DanmarksNationalbankDataSource {
		Container.current = newCommonHttpExchangeRatesDataProvider(config, &DanmarksNationalbankDataSource{})
		Container.isCustom = false
		return nil
	} else if config.ExchangeRatesDataSource == settings.EuroCentralBankDataSource {
		Container.current = newCommonHttpExchangeRatesDataProvider(config, &EuroCentralBankDataSource{})
		Container.isCustom = false
		return nil
	} else if config.ExchangeRatesDataSource == settings.NationalBankOfGeorgiaDataSource {
		Container.current = newCommonHttpExchangeRatesDataProvider(config, &NationalBankOfGeorgiaDataSource{})
		Container.isCustom = false
		return nil
	} else if config.ExchangeRatesDataSource == settings.CentralBankOfHungaryDataSource {
		Container.current = newCommonHttpExchangeRatesDataProvider(config, &CentralBankOfHungaryDataSource{})
		Container.isCustom = false
		return nil
	} else if config.ExchangeRatesDataSource == settings.BankOfIsraelDataSource {
		Container.current = newCommonHttpExchangeRatesDataProvider(config, &BankOfIsraelDataSource{})
		Container.isCustom = false
		return nil
	} else if config.ExchangeRatesDataSource == settings.CentralBankOfMyanmarDataSource {
		Container.current = newCommonHttpExchangeRatesDataProvider(config, &CentralBankOfMyanmarDataSource{})
		Container.isCustom = false
		return nil
	} else if config.ExchangeRatesDataSource == settings.NorgesBankDataSource {
		Container.current = newCommonHttpExchangeRatesDataProvider(config, &NorgesBankDataSource{})
		Container.isCustom = false
		return nil
	} else if config.ExchangeRatesDataSource == settings.NationalBankOfPolandDataSource {
		Container.current = newCommonHttpExchangeRatesDataProvider(config, &NationalBankOfPolandDataSource{})
		Container.isCustom = false
		return nil
	} else if config.ExchangeRatesDataSource == settings.NationalBankOfRomaniaDataSource {
		Container.current = newCommonHttpExchangeRatesDataProvider(config, &NationalBankOfRomaniaDataSource{})
		Container.isCustom = false
		return nil
	} else if config.ExchangeRatesDataSource == settings.BankOfRussiaDataSource {
		Container.current = newCommonHttpExchangeRatesDataProvider(config, &BankOfRussiaDataSource{})
		Container.isCustom = false
		return nil
	} else if config.ExchangeRatesDataSource == settings.SwissNationalBankDataSource {
		Container.current = newCommonHttpExchangeRatesDataProvider(config, &SwissNationalBankDataSource{})
		Container.isCustom = false
		return nil
	} else if config.ExchangeRatesDataSource == settings.NationalBankOfUkraineDataSource {
		Container.current = newCommonHttpExchangeRatesDataProvider(config, &NationalBankOfUkraineDataSource{})
		Container.isCustom = false
		return nil
	} else if config.ExchangeRatesDataSource == settings.CentralBankOfUzbekistanDataSource {
		Container.current = newCommonHttpExchangeRatesDataProvider(config, &CentralBankOfUzbekistanDataSource{})
		Container.isCustom = false
		return nil
	} else if config.ExchangeRatesDataSource == settings.InternationalMonetaryFundDataSource {
		Container.current = newCommonHttpExchangeRatesDataProvider(config, &InternationalMonetaryFundDataSource{})
		Container.isCustom = false
		return nil
	} else if config.ExchangeRatesDataSource == settings.UserCustomExchangeRatesDataSource {
		Container.current = newUserCustomExchangeRatesDataProvider()
		Container.isCustom = true
		return nil
	}

	return errs.ErrInvalidExchangeRatesDataSource
}

// GetLatestExchangeRates returns the latest exchange rates data from the current exchange rates data source
func (e *ExchangeRatesDataProviderContainer) GetLatestExchangeRates(c core.Context, uid int64, currentConfig *settings.Config) (*models.LatestExchangeRateResponse, error) {
	if Container.current == nil {
		return nil, errs.ErrInvalidExchangeRatesDataSource
	}

	if !e.isCustom {
		e.mu.RLock()
		if e.lastResult != nil && time.Since(e.lastTime) < exchangeRatesCacheTimeout {
			result := e.lastResult
			e.mu.RUnlock()
			return result, nil
		}
		e.mu.RUnlock()

		result, err, _ := e.requestGroup.Do("GetLatestExchangeRates", func() (interface{}, error) {
			res, fetchErr := e.current.GetLatestExchangeRates(c, uid, currentConfig)

			if fetchErr == nil {
				e.mu.Lock()
				e.lastResult = res
				e.lastTime = time.Now()
				e.mu.Unlock()
			}

			return res, fetchErr
		})

		if err == nil {
			return result.(*models.LatestExchangeRateResponse), nil
		}

		e.mu.RLock()
		defer e.mu.RUnlock()

		if e.lastResult != nil {
			log.Warnf(c, "[exchangerates.Container] failed to get latest rates, using stale cache from %s", e.lastTime.Format("2006-01-02 15:04:05"))
			return e.lastResult, nil
		}

		return nil, err
	}

	return e.current.GetLatestExchangeRates(c, uid, currentConfig)
}
