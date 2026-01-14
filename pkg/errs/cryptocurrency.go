package errs

import (
	"net/http"
)

// Cryptocurrency error subcategories use NormalSubcategoryCryptocurrency = 19

// Error codes related to cryptocurrency
var (
	ErrInvalidCryptocurrencyDataSource = NewSystemError(SystemSubcategorySetting, 26, http.StatusInternalServerError, "invalid cryptocurrency data source")
	ErrCryptocurrencyServiceNotEnabled = NewNormalError(NormalSubcategoryCryptocurrency, 0, http.StatusBadRequest, "cryptocurrency service not enabled")
	ErrInvalidCryptocurrencySymbol     = NewNormalError(NormalSubcategoryCryptocurrency, 1, http.StatusBadRequest, "invalid cryptocurrency symbol")
)
