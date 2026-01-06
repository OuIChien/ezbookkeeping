package models

import (
	"strings"
)

// LatestCryptocurrencyPriceResponse returns a view-object which contains latest cryptocurrency prices
type LatestCryptocurrencyPriceResponse struct {
	DataSource   string                        `json:"dataSource"`
	ReferenceUrl string                        `json:"referenceUrl"`
	UpdateTime   int64                         `json:"updateTime"`
	BaseCurrency  string                       `json:"baseCurrency"`
	Prices       LatestCryptocurrencyPriceSlice `json:"prices"`
}

// LatestCryptocurrencyPrice represents a data pair of cryptocurrency symbol and price
type LatestCryptocurrencyPrice struct {
	Symbol string `json:"symbol"`
	Price  string `json:"price"`
}

// LatestCryptocurrencyPriceSlice represents the slice data structure of LatestCryptocurrencyPrice
type LatestCryptocurrencyPriceSlice []*LatestCryptocurrencyPrice

// Len returns the count of items
func (s LatestCryptocurrencyPriceSlice) Len() int {
	return len(s)
}

// Swap swaps two items
func (s LatestCryptocurrencyPriceSlice) Swap(i, j int) {
	s[i], s[j] = s[j], s[i]
}

// Less reports whether the first item is less than the second one
func (s LatestCryptocurrencyPriceSlice) Less(i, j int) bool {
	return strings.Compare(s[i].Symbol, s[j].Symbol) < 0
}

