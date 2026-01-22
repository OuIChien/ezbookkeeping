package models

// LatestCryptocurrencyPriceResponse represents the response of latest cryptocurrency price
type LatestCryptocurrencyPriceResponse struct {
	DataSource   string                        `json:"dataSource"`
	ReferenceUrl string                        `json:"referenceUrl"`
	UpdateTime   int64                         `json:"updateTime"`
	BaseCurrency string                        `json:"baseCurrency"`
	Prices       LatestCryptocurrencyPriceSlice `json:"prices"`
}

// LatestCryptocurrencyPrice represents the latest cryptocurrency price
type LatestCryptocurrencyPrice struct {
	Symbol string `json:"symbol"`
	Price  string `json:"price"`
}

// LatestCryptocurrencyPriceSlice represents the slice of latest cryptocurrency price
type LatestCryptocurrencyPriceSlice []*LatestCryptocurrencyPrice

// Len returns the length of the slice
func (s LatestCryptocurrencyPriceSlice) Len() int {
	return len(s)
}

// Swap swaps the elements with indexes i and j
func (s LatestCryptocurrencyPriceSlice) Swap(i, j int) {
	s[i], s[j] = s[j], s[i]
}

// Less returns true if the element with index i should be sorted before the element with index j
func (s LatestCryptocurrencyPriceSlice) Less(i, j int) bool {
	return s[i].Symbol < s[j].Symbol
}
