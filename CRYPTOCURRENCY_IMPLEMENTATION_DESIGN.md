# Cryptocurrency Price Implementation Design

## Overview

This document provides the design for the cryptocurrency price fetching feature in ezBookkeeping. The system will fetch cryptocurrency prices (in USDT) from remote sources, similar to how exchange rates are currently handled. This feature allows users to track cryptocurrency account balances and convert them to their default currency.

## 1. Architecture Overview

The cryptocurrency price system will follow the same **Strategy Pattern** with **Container Pattern** as the exchange rates system:

```
┌─────────────────────────────────────────────────────────┐
│              Frontend (Vue 3 + Pinia)                    │
│  ┌──────────────────────────────────────────────────┐   │
│  │  CryptocurrencyPricesStore                       │   │
│  │  - Manages cryptocurrency prices state          │   │
│  │  - LocalStorage caching                         │   │
│  │  - Price conversion calculations                 │   │
│  └──────────────────────────────────────────────────┘   │
└─────────────────────────────────────────────────────────┘
                        │ HTTP API
                        │ GET /api/v1/cryptocurrency/latest.json
                        ▼
┌─────────────────────────────────────────────────────────┐
│              Backend (Go + Gin)                         │
│  ┌──────────────────────────────────────────────────┐   │
│  │  API Layer (pkg/api/cryptocurrency.go)          │   │
│  │  - LatestCryptocurrencyPriceHandler              │   │
│  └──────────────────────────────────────────────────┘   │
│  ┌──────────────────────────────────────────────────┐   │
│  │  CryptocurrencyPriceDataProviderContainer        │   │
│  │  - Manages data provider instances               │   │
│  │  - Strategy pattern for data sources             │   │
│  └──────────────────────────────────────────────────┘   │
│  ┌──────────────────────────────────────────────────┐   │
│  │  Data Providers                                   │   │
│  │  - CommonHttpCryptocurrencyPriceDataProvider      │   │
│  └──────────────────────────────────────────────────┘   │
│  ┌──────────────────────────────────────────────────┐   │
│  │  Data Sources (HttpCryptocurrencyPriceDataSource)│   │
│  │  - CoinGeckoDataSource                           │   │
│  │  - CoinMarketCapDataSource                       │   │
│  │  - BinanceDataSource                             │   │
│  │  - ... (other data sources)                       │   │
│  └──────────────────────────────────────────────────┘   │
└─────────────────────────────────────────────────────────┘
                        │ HTTP Request
                        ▼
┌─────────────────────────────────────────────────────────┐
│         Remote Cryptocurrency Price APIs                │
│  - CoinGecko API                                        │
│  - CoinMarketCap API                                    │
│  - Binance API                                          │
│  - ...                                                   │
└─────────────────────────────────────────────────────────┘
```

## 2. Design Principles

### 2.1 Consistency with Exchange Rates System

- Follow the same architectural patterns as the exchange rates system
- Use similar naming conventions and code structure
- Reuse existing infrastructure (HTTP client, configuration, error handling)

### 2.2 Configuration-Based Cryptocurrency Selection

Similar to how exchange rates data source is configured, cryptocurrency selection will be:
- **Configuration-based**: Specified in `conf/ezbookkeeping.ini`
- **List-based**: Support multiple cryptocurrencies in configuration
- **Flexible**: Easy to add/remove cryptocurrencies without code changes

### 2.3 USDT as Base Currency

- All cryptocurrency prices will be fetched in USDT
- USDT will be treated as the base currency (rate = "1")
- This allows easy conversion to fiat currencies via existing exchange rates

## 3. Configuration Design

### 3.1 Configuration Section

Add a new section in `conf/ezbookkeeping.ini`:

```ini
[cryptocurrency]
# Cryptocurrency price data source, supports:
# "coingecko": CoinGecko API (free tier available)
# "coinmarketcap": CoinMarketCap API (requires API key)
# "binance": Binance Public API
data_source = coingecko

# Comma-separated list of cryptocurrency symbols to fetch
# Examples: BTC,ETH,BNB,SOL,ADA
# All prices will be in USDT
cryptocurrencies = BTC,ETH,BNB

# Request timeout (0 - 4294967295 milliseconds)
# Default is 10000 (10 seconds)
request_timeout = 10000

# Proxy setting (same as exchange_rates)
proxy = system

# Skip TLS verification
skip_tls_verify = false

# API key (optional, required for some data sources like CoinMarketCap)
api_key = 
```

### 3.2 Configuration Loading

- Add configuration loading in `pkg/settings/setting.go`
- Similar to `loadExchangeRatesConfiguration()`
- Function: `loadCryptocurrencyConfiguration()`
- Store in `Config` struct:
  - `CryptocurrencyDataSource`
  - `CryptocurrencySymbols` (slice of strings)
  - `CryptocurrencyRequestTimeout`
  - `CryptocurrencyProxy`
  - `CryptocurrencySkipTLSVerify`
  - `CryptocurrencyAPIKey`

## 4. Backend Implementation Design

### 4.1 Package Structure

Create new package: `pkg/cryptocurrency/`

**Core Files**:
- `cryptocurrency_price_data_provider.go`: Interface definition
- `cryptocurrency_price_data_provider_container.go`: Container and initialization
- `common_http_cryptocurrency_price_data_provider.go`: Common HTTP provider

**Data Source Files**:
- `coingecko_datasource.go`: CoinGecko API implementation
- `coinmarketcap_datasource.go`: CoinMarketCap API implementation
- `binance_datasource.go`: Binance API implementation
- (Additional data sources as needed)

### 4.2 Data Provider Interface

```go
type CryptocurrencyPriceDataProvider interface {
    GetLatestCryptocurrencyPrices(
        c core.Context, 
        uid int64, 
        currentConfig *settings.Config
    ) (*models.LatestCryptocurrencyPriceResponse, error)
}
```

### 4.3 HTTP Data Source Interface

```go
type HttpCryptocurrencyPriceDataSource interface {
    BuildRequests(symbols []string, apiKey string) ([]*http.Request, error)
    Parse(c core.Context, content []byte) (*models.LatestCryptocurrencyPriceResponse, error)
}
```

### 4.4 Data Models

**New Model**: `pkg/models/cryptocurrency_price.go`

```go
type LatestCryptocurrencyPriceResponse struct {
    DataSource    string
    ReferenceUrl  string
    UpdateTime    int64
    BaseCurrency  string  // Always "USDT"
    Prices        LatestCryptocurrencyPriceSlice
}

type LatestCryptocurrencyPrice struct {
    Symbol string  // e.g., "BTC", "ETH"
    Price  string  // Price in USD
}
```

### 4.5 Data Source Implementation Strategy

**CoinGecko (Recommended for Free Tier)**:
- API: `https://api.coingecko.com/api/v3/simple/price`
- Parameters: `ids=bitcoin,ethereum&vs_currencies=usdt`
- Free tier: No API key required, rate limits apply
- Response: JSON format

**CoinMarketCap**:
- API: `https://pro-api.coinmarketcap.com/v1/cryptocurrency/quotes/latest`
- Requires API key in header
- Response: JSON format

**Binance**:
- API: `https://api.binance.com/api/v3/ticker/price`
- No API key required
- Returns prices in USDT by default
- Response: JSON format

### 4.6 Initialization Flow

1. System startup: `InitializeCryptocurrencyDataSource(config)`
2. Read configuration: data source type and cryptocurrency list
3. Create appropriate data provider instance
4. Store in container singleton

### 4.7 Price Fetching Flow

1. API handler receives request
2. Container calls `GetLatestCryptocurrencyPrices()`
3. Data provider:
   - Gets cryptocurrency symbols from config
   - Builds HTTP requests using data source
   - Executes requests with timeout/proxy settings
   - Parses responses
   - Normalizes all prices to USDT base
   - Returns unified response

## 5. API Design

### 5.1 Endpoint

**Route**: `GET /api/v1/cryptocurrency/latest.json`

**Handler**: `CryptocurrencyApi.LatestCryptocurrencyPriceHandler`

**Response Format**:
```json
{
  "success": true,
  "result": {
    "dataSource": "CoinGecko",
    "referenceUrl": "https://www.coingecko.com",
    "updateTime": 1234567890,
    "baseCurrency": "USDT",
    "prices": [
      {
        "symbol": "BTC",
        "price": "45000.50"
      },
      {
        "symbol": "ETH",
        "price": "3000.25"
      }
    ]
  }
}
```

### 5.2 API File

**File**: `pkg/api/cryptocurrency.go`

- Similar structure to `pkg/api/exchange_rates.go`
- Single handler for fetching latest prices
- Error handling consistent with existing patterns

## 6. Frontend Implementation Design

### 6.1 Store Design

**File**: `src/stores/cryptocurrencyPrices.ts`

**State**:
- `latestCryptocurrencyPrices`: Current prices with timestamp

**Computed Properties**:
- `cryptocurrencyPricesLastUpdateTime`: Last update timestamp
- `latestCryptocurrencyPriceMap`: Map of symbol → price

**Key Methods**:
- `getLatestCryptocurrencyPrices({ silent, force })`: Fetch prices with caching
- `getCryptocurrencyPriceInUSDT(symbol)`: Get price for specific symbol
- `getCryptocurrencyPriceInFiat(symbol, fiatCurrency)`: Convert to fiat via exchange rates

### 6.2 LocalStorage Caching

- Key: `ebk_app_cryptocurrency_prices`
- Cache validity: Same as exchange rates (same day or same hour)
- Structure: `{ time: number, data: LatestCryptocurrencyPriceResponse }`

### 6.3 Integration with Exchange Rates

To convert cryptocurrency to fiat currency:
1. Get cryptocurrency Price in USD from cryptocurrency store
2. Get USDT to fiat exchange rate from exchange rates store
3. Calculate: `cryptoPrice * usdtToFiatRate`

### 6.4 Service Layer

**File**: `src/lib/services.ts`

Add method:
```typescript
getLatestCryptocurrencyPrices(options?: RequestOptions): Promise<ApiResponse<LatestCryptocurrencyPriceResponse>>
```

## 7. Currency Code Handling

### 7.1 Cryptocurrency Symbols

- Use standard cryptocurrency symbols (BTC, ETH, etc.)
- These are different from ISO 4217 currency codes
- Need to extend currency validation or create separate validation

### 7.2 Account Currency Field

- Current system uses ISO 4217 codes in `Account.Currency`
- Options:
  1. **Extend currency validator**: Allow cryptocurrency symbols in addition to ISO codes
  2. **Separate field**: Add `CryptocurrencySymbol` field (more complex)
  3. **Prefix convention**: Use prefix like "CRYPTO_BTC" (hacky)

**Recommended**: Option 1 - Extend validator to accept cryptocurrency symbols
- Add `AllCryptocurrencySymbols` map in `pkg/validators/currency.go`
- Update `ValidCurrency` to check both ISO codes and crypto symbols
- Add cryptocurrency symbols to frontend `ALL_CURRENCIES` constant

### 7.3 Supported Cryptocurrencies

Initial support for major cryptocurrencies:
- BTC (Bitcoin)
- ETH (Ethereum)
- BNB (Binance Coin)
- SOL (Solana)
- ADA (Cardano)
- XRP (Ripple)
- DOT (Polkadot)
- DOGE (Dogecoin)
- MATIC (Polygon)
- USDT (Tether) - base currency

Can be extended via configuration.

## 8. Data Flow

### 8.1 Initial Fetch Flow

```
1. User logs in → DesktopApp.vue / MobileApp.vue
2. Check autoUpdateCryptocurrencyPrices setting (if added)
3. Call cryptocurrencyPricesStore.getLatestCryptocurrencyPrices()
4. Check localStorage cache validity
5. If invalid/forced:
   a. Call services.getLatestCryptocurrencyPrices() (API)
   b. Backend: Cryptocurrency.LatestCryptocurrencyPriceHandler()
   c. Backend: cryptocurrency.Container.GetLatestCryptocurrencyPrices()
   d. Data Provider: GetLatestCryptocurrencyPrices()
   e. Data Source: BuildRequests() → HTTP request with symbols
   f. Remote API: Returns JSON
   g. Data Source: Parse() → LatestCryptocurrencyPriceResponse
   h. Data Provider: Normalize to USDT base
   i. API: Returns JSON response
   j. Frontend: Updates store and localStorage
6. Return cached or fresh data
```

### 8.2 Price Conversion Flow

```
1. User views account with cryptocurrency currency (e.g., BTC)
2. Component calls cryptocurrencyPricesStore.getCryptocurrencyPriceInFiat("BTC", "USD")
3. Store:
   a. Gets BTC Price in USD from cryptocurrencyPricesStore
   b. Gets USDT to USD rate from exchangeRatesStore
   c. Calculates: btcPriceInUSDT * usdtToUsdRate
4. Display converted amount
```

## 9. Error Handling

### 9.1 Backend Errors

- **Network errors**: Returns `ErrFailedToRequestRemoteApi`
- **Parse errors**: Returns `ErrFailedToRequestRemoteApi` with details
- **Invalid data source**: Returns `ErrInvalidCryptocurrencyDataSource`
- **Missing symbols**: Logs warning, continues with available prices
- **API key errors**: Returns appropriate error for data sources requiring keys

### 9.2 Frontend Errors

- **API failures**: 
  - If `silent: true`: Logs error, doesn't notify user
  - If `silent: false`: Shows error notification
- **Missing prices**: 
  - `getCryptocurrencyPriceInUSDT()` returns `null`
  - UI shows incomplete amount indicator
- **Cache failures**: Falls back to API request

## 10. Performance Considerations

### 10.1 Caching Strategy

- **LocalStorage**: Cache prices with timestamp
- **Cache validity**: Same day or same hour (configurable)
- **Force refresh**: Bypasses cache when `force: true`

### 10.2 Request Optimization

- **Batch requests**: Fetch all configured cryptocurrencies in single API call when possible
- **Rate limiting**: Respect API rate limits (especially for free tiers)
- **Timeout protection**: Configurable request timeout

### 10.3 Update Frequency

- Cryptocurrency prices change frequently (unlike daily exchange rates)
- Consider shorter cache validity (e.g., same hour or 15 minutes)
- Allow manual refresh option

## 11. Extension Points

### 11.1 Adding New Data Source

1. Create new file: `pkg/cryptocurrency/new_source_datasource.go`
2. Implement `HttpCryptocurrencyPriceDataSource` interface
3. Add data source constant in `pkg/settings/setting.go`
4. Add initialization case in `InitializeCryptocurrencyDataSource()`

### 11.2 Adding New Cryptocurrency

1. Add symbol to configuration: `cryptocurrencies = BTC,ETH,NEW_COIN`
2. Add symbol to validator: `pkg/validators/currency.go`
3. Add symbol to frontend: `src/consts/currency.ts`
4. No code changes needed if data source supports the symbol

## 12. Configuration Examples

### 12.1 Using CoinGecko (Free)

```ini
[cryptocurrency]
data_source = coingecko
cryptocurrencies = BTC,ETH,BNB,SOL,ADA
request_timeout = 10000
proxy = system
skip_tls_verify = false
```

### 12.2 Using CoinMarketCap (Requires API Key)

```ini
[cryptocurrency]
data_source = coinmarketcap
cryptocurrencies = BTC,ETH,BNB
api_key = your_api_key_here
request_timeout = 10000
proxy = system
skip_tls_verify = false
```

### 12.3 Using Binance (Free)

```ini
[cryptocurrency]
data_source = binance
cryptocurrencies = BTC,ETH,BNB,SOL
request_timeout = 10000
proxy = system
skip_tls_verify = false
```

## 13. Key Files Reference

### Backend
- `pkg/cryptocurrency/cryptocurrency_price_data_provider_container.go`: Container and initialization
- `pkg/cryptocurrency/common_http_cryptocurrency_price_data_provider.go`: HTTP provider implementation
- `pkg/cryptocurrency/coingecko_datasource.go`: CoinGecko data source
- `pkg/api/cryptocurrency.go`: API handlers
- `pkg/models/cryptocurrency_price.go`: Data models
- `pkg/settings/setting.go`: Configuration loading
- `pkg/validators/currency.go`: Extended currency validation

### Frontend
- `src/stores/cryptocurrencyPrices.ts`: Cryptocurrency prices store
- `src/lib/services.ts`: API service methods
- `src/consts/currency.ts`: Extended currency definitions
- `cmd/webserver.go`: API route registration

## 14. Implementation Phases

### Phase 1: Backend Core
1. Create package structure
2. Implement data provider interfaces
3. Implement CoinGecko data source (free, no API key)
4. Add configuration loading
5. Create API endpoint
6. Add data models

### Phase 2: Frontend Integration
1. Create cryptocurrency prices store
2. Add service methods
3. Implement caching
4. Add price conversion utilities

### Phase 3: Currency Support
1. Extend currency validator
2. Add cryptocurrency symbols to frontend constants
3. Update account currency handling

### Phase 4: Additional Data Sources
1. Implement Binance data source
2. Implement CoinMarketCap data source (optional, requires API key)
3. Add more data sources as needed

### Phase 5: Testing & Optimization
1. Test with various cryptocurrencies
2. Test error handling
3. Optimize caching strategy
4. Performance testing

## 15. Summary

The cryptocurrency price system will:

1. **Follow Existing Patterns**: Reuse exchange rates system architecture
2. **Configuration-Driven**: Easy to configure which cryptocurrencies to track
3. **Multiple Data Sources**: Support CoinGecko, Binance, CoinMarketCap, etc.
4. **USDT Base**: All prices in USDT for easy conversion to fiat
5. **Efficient Caching**: LocalStorage caching with smart invalidation
6. **Seamless Integration**: Works with existing exchange rates for fiat conversion
7. **Extensible**: Easy to add new data sources and cryptocurrencies

The implementation maintains consistency with the existing codebase while providing a flexible foundation for cryptocurrency price tracking.
