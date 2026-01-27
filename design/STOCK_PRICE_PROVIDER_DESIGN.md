# Stock Price Provider Implementation Design

## 1. Overview
The goal is to implement a stock price fetching service for ezbookkeeping, enabling users to track the real-time valuation of their stock/securities holdings. This service will follow the existing "Provider/DataSource" architecture used in the cryptocurrency and exchange rates modules.

## 2. Architecture
The implementation will be located in a new package `pkg/stocks`. It consists of the following components:

- **DataProvider Interface**: Defines the high-level method to fetch prices.
- **Common HTTP Provider**: Handles generic HTTP request execution, error handling, and proxy/timeout settings.
- **DataSource Interface**: Defines how to build specific API requests and parse their responses.
- **Container**: Manages the initialization and retrieval of the configured data provider.

### 2.1 Proposed Directory Structure
```text
pkg/stocks/
├── stock_price_data_provider.go           # Interface definition
├── common_http_stock_price_data_provider.go # Base HTTP implementation
├── stock_price_data_provider_container.go  # Registry and factory
└── yahoo_finance_datasource.go             # Yahoo Finance implementation
```

## 3. Data Models
New models will be added to `pkg/models/stock_price.go` to standardize the response format.

```go
type LatestStockPriceResponse struct {
    DataSource   string               `json:"dataSource"`
    ReferenceUrl string               `json:"referenceUrl"`
    UpdateTime   int64                `json:"updateTime"`
    BaseCurrency string               `json:"baseCurrency"`
    Prices       LatestStockPriceSlice `json:"prices"`
}

type LatestStockPrice struct {
    Symbol string `json:"symbol"`
    Price  string `json:"price"`
}
```

## 4. Configuration
A new `[stocks]` section will be added to the system configuration.

| Item | Description | Default |
|------|-------------|---------|
| `data_source` | Source of stock data (e.g., `yahoo_finance`) | - |
| `stocks` | List of stock symbols (e.g., `AAPL,TSLA,0700.HK`) | - |
| `request_timeout` | API request timeout in milliseconds | `10000` |
| `proxy` | Proxy server for requests | `system` |
| `api_key` | Optional API key for specific sources | - |

## 5. Performance and Optimization

To improve performance and ensure service availability, the backend implements several optimization strategies:

1. **In-Memory Caching**: The stock price container caches the latest successful fetch results for 5 minutes.
2. **Request Coalescing (Singleflight)**: Uses `singleflight` to prevent redundant concurrent requests to the same data source.
3. **Stale Cache Fallback**: If a remote API request fails, the system returns the last successful result from the cache.
4. **Auto-Updating**: A background cron job (`UpdateStockPricesJob`) runs every 5 minutes to refresh the cache if `enable_auto_update_stock_prices` is enabled in the `[cron]` section.

## 6. Planned Data Sources
1.  **Yahoo Finance (`yahoo_finance`)**:
    - Primary source due to wide coverage of global markets.
    - Supports symbols like `AAPL`, `0700.HK` (Hong Kong), and `600519.SS` (Shanghai).
2.  **Alpha Vantage (`alphavantage`)**:
    - Secondary/Fallback source (requires API Key).

## 6. Implementation Phases
1.  **Phase 1**: Define data models in `pkg/models`.
2.  **Phase 2**: Extend `pkg/settings` to support stock-related configuration items.
3.  **Phase 3**: Implement the core interfaces and `CommonHttpStockPriceDataProvider` in `pkg/stocks`.
4.  **Phase 4**: Implement `YahooFinanceDataSource` with robust parsing logic.
5.  **Phase 5**: Integrate the provider into the container and expose it for future business logic (valuation calculation).
