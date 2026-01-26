package validators

import (
	"strconv"
	"strings"

	"github.com/go-playground/validator/v10"
)

// transactionNoTagFilterValue is the value for "no tag" filter, must match models.TransactionNoTagFilterValue
const transactionNoTagFilterValue = "none"

// ValidTagFilter returns whether the given tag filter is valid
func ValidTagFilter(fl validator.FieldLevel) bool {
	if value, ok := fl.Field().Interface().(string); ok {
		return isTagFilterFormatValid(value)
	}
	return false
}

// isTagFilterFormatValid checks tag filter format without depending on models.
// Format: "" or "none" = valid; else "type:id1,id2;type2:id3" where type is 0-3, ids are comma-separated int64.
func isTagFilterFormatValid(s string) bool {
	if s == "" || s == transactionNoTagFilterValue {
		return true
	}
	for _, filter := range strings.Split(s, ";") {
		parts := strings.Split(filter, ":")
		if len(parts) != 2 {
			return false
		}
		typ, err := strconv.Atoi(parts[0])
		if err != nil || typ < 0 || typ > 3 {
			return false
		}
		if parts[1] == "" {
			return false
		}
		for _, idStr := range strings.Split(parts[1], ",") {
			_, err := strconv.ParseInt(idStr, 10, 64)
			if err != nil {
				return false
			}
		}
	}
	return true
}
