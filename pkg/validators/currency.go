package validators

import (
	"github.com/go-playground/validator/v10"

	"github.com/mayswind/ezbookkeeping/pkg/models"
)

// ParentAccountCurrencyPlaceholder represents the currency field of parent account stored in database
const ParentAccountCurrencyPlaceholder = "---"

// AllCurrencyNames re-exports models.AllCurrencyNames for backward compatibility
var AllCurrencyNames = models.AllCurrencyNames

// AllCryptocurrencySymbols re-exports models.AllCryptocurrencySymbols for backward compatibility
var AllCryptocurrencySymbols = models.AllCryptocurrencySymbols


// ValidCurrency returns whether the given currency is valid
func ValidCurrency(fl validator.FieldLevel) bool {
	if value, ok := fl.Field().Interface().(string); ok {
		if value == ParentAccountCurrencyPlaceholder {
			return true
		}

		if _, ok := AllCurrencyNames[value]; ok {
			return true
		}

		if _, ok := AllCryptocurrencySymbols[value]; ok {
			return true
		}

		if len(value) < 1 || len(value) > 10 {
			return false
		}

		for _, r := range value {
			if !((r >= 'a' && r <= 'z') || (r >= 'A' && r <= 'Z') || (r >= '0' && r <= '9')) {
				return false
			}
		}

		return true
	}

	return false
}
