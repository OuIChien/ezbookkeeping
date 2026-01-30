package services

import (
	"time"

	"github.com/mayswind/ezbookkeeping/pkg/core"
	"github.com/mayswind/ezbookkeeping/pkg/datastore"
	"github.com/mayswind/ezbookkeeping/pkg/errs"
	"github.com/mayswind/ezbookkeeping/pkg/models"
)

// CryptocurrencyService represents cryptocurrency service
type CryptocurrencyService struct {
	ServiceUsingDB
}

// Initialize a cryptocurrency service singleton instance
var (
	Cryptocurrencies = &CryptocurrencyService{
		ServiceUsingDB: ServiceUsingDB{
			container: datastore.Container,
		},
	}
)

// GetAllCryptocurrencies returns all cryptocurrency models
func (s *CryptocurrencyService) GetAllCryptocurrencies(c core.Context) ([]*models.Cryptocurrency, error) {
	var cryptocurrencies []*models.Cryptocurrency
	err := s.UserDataDB(0).NewSession(c).OrderBy("display_order asc").Find(&cryptocurrencies)
	return cryptocurrencies, err
}

// GetAllVisibleCryptocurrencies returns all visible cryptocurrency models
func (s *CryptocurrencyService) GetAllVisibleCryptocurrencies(c core.Context) ([]*models.Cryptocurrency, error) {
	var cryptocurrencies []*models.Cryptocurrency
	err := s.UserDataDB(0).NewSession(c).Where("is_hidden=?", false).OrderBy("display_order asc").Find(&cryptocurrencies)
	return cryptocurrencies, err
}

// GetCryptocurrencyBySymbol returns cryptocurrency model by symbol
func (s *CryptocurrencyService) GetCryptocurrencyBySymbol(c core.Context, symbol string) (*models.Cryptocurrency, error) {
	crypto := &models.Cryptocurrency{}
	has, err := s.UserDataDB(0).NewSession(c).Where("symbol=?", symbol).Get(crypto)
	if err != nil {
		return nil, err
	} else if !has {
		return nil, errs.ErrCryptocurrencyNotFound
	}
	return crypto, nil
}

// CreateCryptocurrency saves a new cryptocurrency model to database
func (s *CryptocurrencyService) CreateCryptocurrency(c core.Context, crypto *models.Cryptocurrency) error {
	now := time.Now().Unix()
	crypto.CreatedUnixTime = now
	crypto.UpdatedUnixTime = now

	_, err := s.UserDataDB(0).NewSession(c).Insert(crypto)
	return err
}

// UpdateCryptocurrency updates an existing cryptocurrency model
func (s *CryptocurrencyService) UpdateCryptocurrency(c core.Context, crypto *models.Cryptocurrency) error {
	now := time.Now().Unix()
	crypto.UpdatedUnixTime = now

	updatedRows, err := s.UserDataDB(0).NewSession(c).Where("symbol=?", crypto.Symbol).Update(crypto)
	if err != nil {
		return err
	} else if updatedRows < 1 {
		return errs.ErrCryptocurrencyNotFound
	}
	return nil
}

// DeleteCryptocurrency deletes a cryptocurrency from database
func (s *CryptocurrencyService) DeleteCryptocurrency(c core.Context, symbol string) error {
	_, err := s.UserDataDB(0).NewSession(c).Where("symbol=?", symbol).Delete(&models.Cryptocurrency{})
	return err
}

// HideCryptocurrency updates the hidden status of a cryptocurrency
func (s *CryptocurrencyService) HideCryptocurrency(c core.Context, symbol string, hidden bool) error {
	now := time.Now().Unix()
	updateModel := &models.Cryptocurrency{
		IsHidden:        hidden,
		UpdatedUnixTime: now,
	}
	
	updatedRows, err := s.UserDataDB(0).NewSession(c).Where("symbol=?", symbol).Cols("is_hidden", "updated_unix_time").Update(updateModel)
	if err != nil {
		return err
	} else if updatedRows < 1 {
		return errs.ErrCryptocurrencyNotFound
	}
	return nil
}
