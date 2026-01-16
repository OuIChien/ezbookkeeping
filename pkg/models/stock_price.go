package models

// LatestStockPriceResponse represents the response of latest stock price
type LatestStockPriceResponse struct {
	DataSource   string               `json:"dataSource"`
	ReferenceUrl string               `json:"referenceUrl"`
	UpdateTime   int64                `json:"updateTime"`
	BaseCurrency string               `json:"baseCurrency"`
	Prices       LatestStockPriceSlice `json:"prices"`
}

// LatestStockPrice represents the latest stock price
type LatestStockPrice struct {
	Symbol   string `json:"symbol"`
	Price    string `json:"price"`
	Currency string `json:"currency"`
}

// LatestStockPriceSlice represents the slice of latest stock price
type LatestStockPriceSlice []*LatestStockPrice

// Len returns the length of the slice
func (s LatestStockPriceSlice) Len() int {
	return len(s)
}

// Swap swaps the elements with indexes i and j
func (s LatestStockPriceSlice) Swap(i, j int) {
	s[i], s[j] = s[j], s[i]
}

// Less returns true if the element with index i should be sorted before the element with index j
func (s LatestStockPriceSlice) Less(i, j int) bool {
	return s[i].Symbol < s[j].Symbol
}
