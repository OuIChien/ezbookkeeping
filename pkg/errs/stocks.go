package errs

import (
	"net/http"
)

// Stock error subcategories use NormalSubcategoryStocks = 20

// Error codes related to stocks
var (
	ErrInvalidStockDataSource = NewSystemError(SystemSubcategorySetting, 27, http.StatusInternalServerError, "invalid stock data source")
	ErrStockServiceNotEnabled = NewNormalError(NormalSubcategoryStocks, 0, http.StatusBadRequest, "stock service not enabled")
	ErrInvalidStockSymbol     = NewNormalError(NormalSubcategoryStocks, 1, http.StatusBadRequest, "invalid stock symbol")
)
