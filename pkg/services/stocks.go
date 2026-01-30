package services

import (
	"time"

	"github.com/mayswind/ezbookkeeping/pkg/core"
	"github.com/mayswind/ezbookkeeping/pkg/datastore"
	"github.com/mayswind/ezbookkeeping/pkg/errs"
	"github.com/mayswind/ezbookkeeping/pkg/models"
)

// StockService represents stock service
type StockService struct {
	ServiceUsingDB
}

// Initialize a stock service singleton instance
var (
	Stocks = &StockService{
		ServiceUsingDB: ServiceUsingDB{
			container: datastore.Container,
		},
	}
)

// GetAllStocks returns all stock models
func (s *StockService) GetAllStocks(c core.Context) ([]*models.Stock, error) {
	var stocks []*models.Stock
	err := s.UserDataDB(0).NewSession(c).OrderBy("display_order asc").Find(&stocks)
	return stocks, err
}

// GetAllVisibleStocks returns all visible stock models
func (s *StockService) GetAllVisibleStocks(c core.Context) ([]*models.Stock, error) {
	var stocks []*models.Stock
	err := s.UserDataDB(0).NewSession(c).Where("is_hidden=?", false).OrderBy("display_order asc").Find(&stocks)
	return stocks, err
}

// GetStockBySymbol returns stock model by symbol
func (s *StockService) GetStockBySymbol(c core.Context, symbol string) (*models.Stock, error) {
	stock := &models.Stock{}
	has, err := s.UserDataDB(0).NewSession(c).Where("symbol=?", symbol).Get(stock)
	if err != nil {
		return nil, err
	} else if !has {
		return nil, errs.ErrStockNotFound
	}
	return stock, nil
}

// CreateStock saves a new stock model to database
func (s *StockService) CreateStock(c core.Context, stock *models.Stock) error {
	now := time.Now().Unix()
	stock.CreatedUnixTime = now
	stock.UpdatedUnixTime = now

	_, err := s.UserDataDB(0).NewSession(c).Insert(stock)
	return err
}

// UpdateStock updates an existing stock model
func (s *StockService) UpdateStock(c core.Context, stock *models.Stock) error {
	now := time.Now().Unix()
	stock.UpdatedUnixTime = now

	updatedRows, err := s.UserDataDB(0).NewSession(c).Where("symbol=?", stock.Symbol).Update(stock)
	if err != nil {
		return err
	} else if updatedRows < 1 {
		return errs.ErrStockNotFound
	}
	return nil
}

// DeleteStock deletes a stock from database
func (s *StockService) DeleteStock(c core.Context, symbol string) error {
	_, err := s.UserDataDB(0).NewSession(c).Where("symbol=?", symbol).Delete(&models.Stock{})
	return err
}

// HideStock updates the hidden status of a stock
func (s *StockService) HideStock(c core.Context, symbol string, hidden bool) error {
	now := time.Now().Unix()
	updateModel := &models.Stock{
		IsHidden:        hidden,
		UpdatedUnixTime: now,
	}
	
	updatedRows, err := s.UserDataDB(0).NewSession(c).Where("symbol=?", symbol).Cols("is_hidden", "updated_unix_time").Update(updateModel)
	if err != nil {
		return err
	} else if updatedRows < 1 {
		return errs.ErrStockNotFound
	}
	return nil
}
