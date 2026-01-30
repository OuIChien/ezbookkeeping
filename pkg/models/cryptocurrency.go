package models

// Cryptocurrency represents cryptocurrency data stored in database
type Cryptocurrency struct {
	CryptocurrencyId int64  `xorm:"PK AUTOINCR"`
	Symbol           string `xorm:"VARCHAR(20) NOT NULL UNIQUE"`
	Name             string `xorm:"VARCHAR(100) NOT NULL"`
	IsHidden         bool   `xorm:"NOT NULL DEFAULT 0"`
	DisplayOrder     int    `xorm:"NOT NULL DEFAULT 0"`
	CreatedUnixTime  int64
	UpdatedUnixTime  int64
}

// CryptocurrencyInfoResponse represents a view-object of cryptocurrency
type CryptocurrencyInfoResponse struct {
	Symbol       string `json:"symbol"`
	Name         string `json:"name"`
	IsHidden     bool   `json:"isHidden"`
	DisplayOrder int    `json:"displayOrder"`
}

// ToCryptocurrencyInfoResponse returns a view-object according to database model
func (c *Cryptocurrency) ToCryptocurrencyInfoResponse() *CryptocurrencyInfoResponse {
	return &CryptocurrencyInfoResponse{
		Symbol:       c.Symbol,
		Name:         c.Name,
		IsHidden:     c.IsHidden,
		DisplayOrder: c.DisplayOrder,
	}
}

// CryptocurrencyCreateRequest represents all parameters of cryptocurrency creation request
type CryptocurrencyCreateRequest struct {
	Symbol       string `json:"symbol" binding:"required,notBlank,max=20"`
	Name         string `json:"name" binding:"required,notBlank,max=100"`
	DisplayOrder int    `json:"displayOrder"`
}

// CryptocurrencyModifyRequest represents all parameters of cryptocurrency modification request
type CryptocurrencyModifyRequest struct {
	Symbol       string `json:"symbol" binding:"required,notBlank,max=20"`
	Name         string `json:"name" binding:"required,notBlank,max=100"`
	IsHidden     bool   `json:"isHidden"`
	DisplayOrder int    `json:"displayOrder"`
}

// CryptocurrencyHideRequest represents all parameters of cryptocurrency hiding request
type CryptocurrencyHideRequest struct {
	Symbol string `json:"symbol" binding:"required,notBlank,max=20"`
	Hidden bool   `json:"hidden"`
}

// CryptocurrencyDeleteRequest represents all parameters of cryptocurrency deleting request
type CryptocurrencyDeleteRequest struct {
	Symbol string `json:"symbol" binding:"required,notBlank,max=20"`
}
