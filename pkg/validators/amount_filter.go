package validators

import (
	"strings"

	"github.com/go-playground/validator/v10"

	"github.com/mayswind/ezbookkeeping/pkg/utils"
)

// ValidAmountFilter returns whether the given amount filter is valid
func ValidAmountFilter(fl validator.FieldLevel) bool {
	if value, ok := fl.Field().Interface().(string); ok {
		if value == "" {
			return true
		}

		amountFilterItems := strings.Split(value, ":")

		if len(amountFilterItems) < 2 {
			return false
		}

		amount1, err := utils.StringToInt64(amountFilterItems[1])

		if err != nil {
			return false
		}

		switch amountFilterItems[0] {
		case "gt", "lt", "eq", "ne":
			if len(amountFilterItems) != 2 {
				return false
			}
		case "bt", "nb":
			if len(amountFilterItems) != 3 {
				return false
			}

			amount2, err := utils.StringToInt64(amountFilterItems[2])

			if err != nil {
				return false
			}

			if amount2 < amount1 {
				return false
			}
		}

		return true
	}

	return false
}
