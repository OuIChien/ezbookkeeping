package models

// ExternalDataSourceType represents the type of external data source
type ExternalDataSourceType int

const (
	EXTERNAL_DATA_SOURCE_TYPE_CRYPTOCURRENCY ExternalDataSourceType = 1
	EXTERNAL_DATA_SOURCE_TYPE_STOCK          ExternalDataSourceType = 2
	EXTERNAL_DATA_SOURCE_TYPE_EXCHANGE_RATE  ExternalDataSourceType = 3
)

// ExternalDataSourceConfig represents external data source configuration stored in database
type ExternalDataSourceConfig struct {
	ConfigId        int64                  `xorm:"PK AUTOINCR"`
	Type            ExternalDataSourceType `xorm:"NOT NULL UNIQUE"`
	DataSource      string                 `xorm:"VARCHAR(50) NOT NULL"`
	BaseCurrency    string                 `xorm:"VARCHAR(10)"`
	ApiKey          string                 `xorm:"VARCHAR(255)"`
	RequestTimeout  int                    `xorm:"INT"`
	Proxy           string                 `xorm:"VARCHAR(255)"`
	UpdateFrequency string                 `xorm:"VARCHAR(20)"`
	CreatedUnixTime int64
	UpdatedUnixTime int64
}

// ExternalDataSourceConfigResponse represents a view-object of external data source configuration
type ExternalDataSourceConfigResponse struct {
	Type            ExternalDataSourceType `json:"type"`
	DataSource      string                 `json:"dataSource"`
	BaseCurrency    string                 `json:"baseCurrency"`
	ApiKey          string                 `json:"apiKey"`
	RequestTimeout  int                    `json:"requestTimeout"`
	Proxy           string                 `json:"proxy"`
	UpdateFrequency string                 `json:"updateFrequency"`
}

// ToExternalDataSourceConfigResponse returns a view-object according to database model
func (c *ExternalDataSourceConfig) ToExternalDataSourceConfigResponse() *ExternalDataSourceConfigResponse {
	return &ExternalDataSourceConfigResponse{
		Type:            c.Type,
		DataSource:      c.DataSource,
		BaseCurrency:    c.BaseCurrency,
		ApiKey:          c.ApiKey,
		RequestTimeout:  c.RequestTimeout,
		Proxy:           c.Proxy,
		UpdateFrequency: c.UpdateFrequency,
	}
}

// ExternalDataSourceConfigSaveRequest represents all parameters of external data source config saving request
type ExternalDataSourceConfigSaveRequest struct {
	Type            ExternalDataSourceType `json:"type" binding:"required"`
	DataSource      string                 `json:"dataSource" binding:"required,notBlank,max=50"`
	BaseCurrency    string                 `json:"baseCurrency" binding:"max=10"`
	ApiKey          string                 `json:"apiKey" binding:"max=255"`
	RequestTimeout  int                    `json:"requestTimeout" binding:"min=0"`
	Proxy           string                 `json:"proxy" binding:"max=255"`
	UpdateFrequency string                 `json:"updateFrequency" binding:"max=20"`
}

// ExternalDataSourceConfigGetRequest represents all parameters of external data source config getting request
type ExternalDataSourceConfigGetRequest struct {
	Type ExternalDataSourceType `form:"type" binding:"required"`
}
