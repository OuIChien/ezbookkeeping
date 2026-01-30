package models

// Stock represents stock data stored in database
type Stock struct {
	StockId         int64  `xorm:"PK AUTOINCR"`
	Symbol          string `xorm:"VARCHAR(20) NOT NULL UNIQUE"`
	Name            string `xorm:"VARCHAR(100) NOT NULL"`
	Market          string `xorm:"VARCHAR(20)"`
	IsHidden        bool   `xorm:"NOT NULL DEFAULT 0"`
	DisplayOrder    int    `xorm:"NOT NULL DEFAULT 0"`
	CreatedUnixTime int64
	UpdatedUnixTime int64
}

// StockInfoResponse represents a view-object of stock
type StockInfoResponse struct {
	Symbol       string `json:"symbol"`
	Name         string `json:"name"`
	Market       string `json:"market"`
	IsHidden     bool   `json:"isHidden"`
	DisplayOrder int    `json:"displayOrder"`
}

// ToStockInfoResponse returns a view-object according to database model
func (s *Stock) ToStockInfoResponse() *StockInfoResponse {
	return &StockInfoResponse{
		Symbol:       s.Symbol,
		Name:         s.Name,
		Market:       s.Market,
		IsHidden:     s.IsHidden,
		DisplayOrder: s.DisplayOrder,
	}
}

// StockCreateRequest represents all parameters of stock creation request
type StockCreateRequest struct {
	Symbol       string `json:"symbol" binding:"required,notBlank,max=20"`
	Name         string `json:"name" binding:"required,notBlank,max=100"`
	Market       string `json:"market" binding:"max=20"`
	DisplayOrder int    `json:"displayOrder"`
}

// StockModifyRequest represents all parameters of stock modification request
type StockModifyRequest struct {
	Symbol       string `json:"symbol" binding:"required,notBlank,max=20"`
	Name         string `json:"name" binding:"required,notBlank,max=100"`
	Market       string `json:"market" binding:"max=20"`
	IsHidden     bool   `json:"isHidden"`
	DisplayOrder int    `json:"displayOrder"`
}

// StockHideRequest represents all parameters of stock hiding request
type StockHideRequest struct {
	Symbol string `json:"symbol" binding:"required,notBlank,max=20"`
	Hidden bool   `json:"hidden"`
}

// StockDeleteRequest represents all parameters of stock deleting request
type StockDeleteRequest struct {
	Symbol string `json:"symbol" binding:"required,notBlank,max=20"`
}
