# Exchange Rates Implementation Analysis

## Overview

This document provides a detailed analysis of the exchange rates feature in ezBookkeeping, which fetches currency exchange rates from remote sources and uses them to display account balances in the user's default currency.

## 1. Architecture Overview

The exchange rates system follows a **Strategy Pattern** with **Container Pattern** for dependency injection:

```
┌─────────────────────────────────────────────────────────┐
│              Frontend (Vue 3 + Pinia)                    │
│  ┌──────────────────────────────────────────────────┐   │
│  │  ExchangeRatesStore                              │   │
│  │  - Manages exchange rates state                  │   │
│  │  - LocalStorage caching                         │   │
│  │  - Currency conversion calculations              │   │
│  └──────────────────────────────────────────────────┘   │
└─────────────────────────────────────────────────────────┘
                        │ HTTP API
                        │ GET /api/v1/exchange_rates/latest.json
                        ▼
┌─────────────────────────────────────────────────────────┐
│              Backend (Go + Gin)                         │
│  ┌──────────────────────────────────────────────────┐   │
│  │  API Layer (pkg/api/exchange_rates.go)          │   │
│  │  - LatestExchangeRateHandler                     │   │
│  │  - UserCustomExchangeRateUpdateHandler           │   │
│  └──────────────────────────────────────────────────┘   │
│  ┌──────────────────────────────────────────────────┐   │
│  │  ExchangeRatesDataProviderContainer               │   │
│  │  - Manages data provider instances               │   │
│  │  - Strategy pattern for data sources             │   │
│  └──────────────────────────────────────────────────┘   │
│  ┌──────────────────────────────────────────────────┐   │
│  │  Data Providers                                   │   │
│  │  - CommonHttpExchangeRatesDataProvider           │   │
│  │  - UserCustomExchangeRatesDataProvider           │   │
│  └──────────────────────────────────────────────────┘   │
│  ┌──────────────────────────────────────────────────┐   │
│  │  Data Sources (HttpExchangeRatesDataSource)      │   │
│  │  - EuroCentralBankDataSource                     │   │
│  │  - BankOfCanadaDataSource                        │   │
│  │  - InternationalMonetaryFundDataSource          │   │
│  │  - ... (17+ data sources)                        │   │
│  └──────────────────────────────────────────────────┘   │
└─────────────────────────────────────────────────────────┘
                        │ HTTP Request
                        ▼
┌─────────────────────────────────────────────────────────┐
│         Remote Exchange Rate APIs                       │
│  - European Central Bank                                │
│  - National Banks (Poland, Canada, etc.)               │
│  - International Monetary Fund                           │
│  - ...                                                   │
└─────────────────────────────────────────────────────────┘
```

## 2. Supported Data Sources

The system supports **17+ exchange rate data sources**:

1. **Reserve Bank of Australia**
2. **Bank of Canada**
3. **Czech National Bank**
4. **Danmarks Nationalbank**
5. **European Central Bank** (EUR base)
6. **National Bank of Georgia**
7. **Central Bank of Hungary**
8. **Bank of Israel**
9. **Central Bank of Myanmar**
10. **Norges Bank** (Norway)
11. **National Bank of Poland**
12. **National Bank of Romania**
13. **Bank of Russia**
14. **Swiss National Bank**
15. **National Bank of Ukraine**
16. **Central Bank of Uzbekistan**
17. **International Monetary Fund**
18. **User Custom Exchange Rates** (manual entry)

Each data source implements the `HttpExchangeRatesDataSource` interface with:
- `BuildRequests()`: Creates HTTP requests to fetch data
- `Parse()`: Parses the response into standard format

## 3. Backend Implementation

### 3.1 Configuration

Exchange rates configuration is loaded from `conf/ezbookkeeping.ini`:

```ini
[exchange_rates]
data_source = euro_central_bank  # or other data source
proxy = system                   # HTTP proxy setting
request_timeout = 30            # Request timeout in seconds
skip_tls_verify = false         # Skip TLS certificate verification
```

Configuration is loaded in `pkg/settings/setting.go` via `loadExchangeRatesConfiguration()`.

### 3.2 Data Provider Container

**File**: `pkg/exchangerates/exchange_rates_data_provider_container.go`

The container manages the current data provider instance:

```go
type ExchangeRatesDataProviderContainer struct {
    current ExchangeRatesDataProvider
}

var Container = &ExchangeRatesDataProviderContainer{}
```

**Initialization** (`InitializeExchangeRatesDataSource`):
- Reads configuration to determine data source
- Creates appropriate provider instance
- For HTTP sources: wraps with `CommonHttpExchangeRatesDataProvider`
- For user custom: creates `UserCustomExchangeRatesDataProvider`

### 3.3 Common HTTP Data Provider

**File**: `pkg/exchangerates/common_http_exchange_rates_data_provider.go`

This provider handles HTTP-based data sources:

**Key Methods**:

1. **`GetLatestExchangeRates()`**:
   - Calls `dataSource.BuildRequests()` to get HTTP requests
   - Executes HTTP requests with configured timeout/proxy
   - Parses responses using `dataSource.Parse()`
   - Merges multiple responses if needed
   - Normalizes all rates to base currency
   - Returns unified `LatestExchangeRateResponse`

**Process Flow**:
```
1. Build HTTP requests from data source
2. Execute requests sequentially
3. Parse each response (XML/JSON)
4. Merge all exchange rates into single map
5. Add base currency with rate "1"
6. Sort by currency code
7. Return unified response
```

### 3.4 Data Source Implementation Example

**File**: `pkg/exchangerates/euro_central_bank_datasource.go`

Example implementation for European Central Bank:

**Data Source Details**:
- URL: `https://www.ecb.europa.eu/stats/eurofxref/eurofxref-daily.xml`
- Format: XML
- Base Currency: EUR
- Update Time: ~16:00 CET on working days

**Implementation**:
1. **`BuildRequests()`**: Creates GET request to ECB XML endpoint
2. **`Parse()`**: 
   - Decodes XML with charset handling
   - Extracts exchange rates from XML structure
   - Validates currency codes
   - Converts to `LatestExchangeRateResponse`
   - Sets update time from XML date

**Key Structures**:
```go
type EuroCentralBankExchangeRateData struct {
    XMLName          xml.Name
    AllExchangeRates []*EuroCentralBankExchangeRates
}

type EuroCentralBankExchangeRate struct {
    Currency string `xml:"currency,attr"`
    Rate     string `xml:"rate,attr"`
}
```

### 3.5 User Custom Exchange Rates

**File**: `pkg/exchangerates/user_custom_data_provider.go`

Allows users to manually set exchange rates:

**Features**:
- Stored in database (`UserCustomExchangeRate` model)
- Rates stored as integers (factor: 100000000) for precision
- Base currency rate must be set first
- Other rates calculated relative to base currency

**Process**:
1. Retrieves user's default currency
2. Fetches all custom exchange rates from database
3. Finds base currency rate
4. Converts all rates to relative format
5. Returns unified response

### 3.6 API Endpoints

**File**: `pkg/api/exchange_rates.go`

**Routes** (defined in `cmd/webserver.go`):
- `GET /api/v1/exchange_rates/latest.json`: Get latest exchange rates
- `POST /api/v1/exchange_rates/user_custom/update.json`: Update custom rate
- `POST /api/v1/exchange_rates/user_custom/delete.json`: Delete custom rate

**LatestExchangeRateHandler**:
```go
func (a *ExchangeRatesApi) LatestExchangeRateHandler(c *core.WebContext) (any, *errs.Error) {
    exchangeRateResponse, err := exchangerates.Container.GetLatestExchangeRates(
        c, c.GetCurrentUid(), a.CurrentConfig())
    if err != nil {
        return nil, errs.Or(err, errs.ErrOperationFailed)
    }
    return exchangeRateResponse, nil
}
```

## 4. Frontend Implementation

### 4.1 Exchange Rates Store

**File**: `src/stores/exchangeRates.ts`

Pinia store managing exchange rates state:

**State**:
- `latestExchangeRates`: Current exchange rates data with timestamp

**Computed Properties**:
- `isUserCustomExchangeRates`: Whether using custom rates
- `exchangeRatesLastUpdateTime`: Last update timestamp
- `latestExchangeRateMap`: Map of currency → exchange rate

**Key Methods**:

1. **`getLatestExchangeRates({ silent, force })`**:
   - Checks cache validity (same day or same hour)
   - If valid and not forced, returns cached data
   - Otherwise, calls API endpoint
   - Updates localStorage cache
   - Returns promise with latest rates

2. **`getExchangedAmount(amount, fromCurrency, toCurrency)`**:
   - Core conversion function
   - Looks up exchange rates for both currencies
   - Calculates: `amount * (toRate / fromRate)`
   - Returns `null` if rates unavailable

3. **`updateUserCustomExchangeRate({ currency, rate })`**:
   - Updates user's custom exchange rate
   - Updates local state after API success

**LocalStorage Caching**:
- Key: `ebk_app_exchange_rates`
- Stores: `{ time: number, data: LatestExchangeRateResponse }`
- Cache validity: Same day or same hour (configurable)

### 4.2 Currency Conversion Logic

**File**: `src/lib/numeral.ts`

**Function**: `getExchangedAmountByRate(amount, fromRate, toRate)`

```typescript
export function getExchangedAmountByRate(
    amount: number, 
    fromRate: string, 
    toRate: string
): number | null {
    const exchangeRate = parseFloat(toRate) / parseFloat(fromRate);
    if (!isNumber(exchangeRate)) {
        return null;
    }
    return amount * exchangeRate;
}
```

**Formula**: `convertedAmount = originalAmount × (targetRate / sourceRate)`

### 4.3 Automatic Refresh

**Files**: `src/DesktopApp.vue`, `src/MobileApp.vue`

On application startup (if user is logged in):
```typescript
if (settingsStore.appSettings.autoUpdateExchangeRatesData) {
    exchangeRatesStore.getLatestExchangeRates({ silent: true, force: false });
}
```

- `silent: true`: Errors won't show user notifications
- `force: false`: Uses cache if valid

## 5. Account Balance Calculation with Exchange Rates

### 5.1 Multi-Currency Account Balance

**File**: `src/stores/account.ts`

When calculating account balance with multiple sub-accounts in different currencies:

**Function**: `getAccountSubAccountBalance()`

**Process**:
1. Determines result currency (first sub-account currency or default)
2. Iterates through sub-accounts
3. For each sub-account:
   - If currency matches result currency: add directly
   - If different: convert using `exchangeRatesStore.getExchangedAmount()`
4. Sums all converted amounts
5. Returns total balance in result currency

**Code Flow**:
```typescript
if (subAccount.currency === resultCurrency) {
    totalBalance += subAccount.balance;
} else {
    const balance = exchangeRatesStore.getExchangedAmount(
        subAccount.balance, 
        subAccount.currency, 
        resultCurrency
    );
    if (isNumber(balance)) {
        totalBalance += Math.trunc(balance);
    } else {
        hasUnCalculatedAmount = true; // Mark incomplete
    }
}
```

### 5.2 Category Total Balance

**File**: `src/locales/helpers.ts`

When calculating total balance for account categories:

**Function**: `getCategorizedAccountsWithDisplayBalance()`

**Process**:
1. Gets all accounts in category
2. For each account:
   - If currency matches default: add directly
   - Otherwise: convert to default currency
3. Sums all amounts
4. Formats with currency symbol

### 5.3 Transaction Statistics

**File**: `src/stores/transaction.ts`

When calculating monthly transaction totals:

**Function**: `calculateMonthTotalAmount()`

**Process**:
1. Iterates through transactions
2. If transaction currency differs from default:
   - Converts amount using exchange rates
   - Marks as incomplete if rate unavailable
3. Accumulates income/expense totals

### 5.4 Statistics Charts

**File**: `src/stores/statistics.ts`

When displaying account total assets/liabilities:

**Function**: `accountTotalAmountAnalysisData()`

**Process**:
1. Iterates through all accounts
2. Converts each account balance to default currency
3. Aggregates by account category
4. Returns data for chart visualization

## 6. Data Flow

### 6.1 Initial Fetch Flow

```
1. User logs in → DesktopApp.vue / MobileApp.vue
2. Check autoUpdateExchangeRatesData setting
3. Call exchangeRatesStore.getLatestExchangeRates()
4. Check localStorage cache validity
5. If invalid/forced:
   a. Call services.getLatestExchangeRates() (API)
   b. Backend: ExchangeRates.LatestExchangeRateHandler()
   c. Backend: exchangerates.Container.GetLatestExchangeRates()
   d. Data Provider: GetLatestExchangeRates()
   e. Data Source: BuildRequests() → HTTP request
   f. Remote API: Returns XML/JSON
   g. Data Source: Parse() → LatestExchangeRateResponse
   h. Data Provider: Merge and normalize rates
   i. API: Returns JSON response
   j. Frontend: Updates store and localStorage
6. Return cached or fresh data
```

### 6.2 Balance Calculation Flow

```
1. User views account list
2. Component calls accountsStore.getAccountBalance()
3. For each account/sub-account:
   a. Check if currency matches default
   b. If different:
      - Call exchangeRatesStore.getExchangedAmount()
      - Lookup rates in latestExchangeRateMap
      - Calculate: amount × (toRate / fromRate)
      - Return converted amount or null
   c. Sum all amounts
4. Display formatted balance with currency
```

### 6.3 Custom Rate Update Flow

```
1. User updates custom exchange rate in UI
2. Frontend: exchangeRatesStore.updateUserCustomExchangeRate()
3. API: POST /api/v1/exchange_rates/user_custom/update.json
4. Backend: ExchangeRates.UserCustomExchangeRateUpdateHandler()
5. Service: UserCustomExchangeRatesService.UpdateCustomExchangeRate()
6. Database: Insert/Update UserCustomExchangeRate record
7. Response: Updated rate with timestamp
8. Frontend: Updates local store and localStorage
```

## 7. Error Handling

### 7.1 Backend Errors

- **Network errors**: Returns `ErrFailedToRequestRemoteApi`
- **Parse errors**: Returns `ErrFailedToRequestRemoteApi` with details
- **Invalid data source**: Returns `ErrInvalidExchangeRatesDataSource`
- **Missing rates**: Logs warning, continues with available rates

### 7.2 Frontend Errors

- **API failures**: 
  - If `silent: true`: Logs error, doesn't notify user
  - If `silent: false`: Shows error notification
- **Missing rates**: 
  - `getExchangedAmount()` returns `null`
  - UI shows incomplete amount indicator (`INCOMPLETE_AMOUNT_SUFFIX`)
- **Cache failures**: Falls back to API request

## 8. Performance Optimizations

### 8.1 Caching Strategy

- **LocalStorage**: Caches exchange rates with timestamp
- **Cache validity**: 
  - Same day: Uses cache without API call
  - Same hour: Uses cache for silent updates
  - Different day/hour: Fetches fresh data
- **Force refresh**: Bypasses cache when `force: true`

### 8.2 Request Optimization

- **Single request per data source**: Most sources return all rates in one request
- **Parallel requests**: If data source requires multiple requests, executed sequentially (could be optimized)
- **Timeout protection**: Configurable request timeout prevents hanging

### 8.3 Calculation Optimization

- **Rate map**: Pre-computed map of currency → rate for O(1) lookup
- **Memoization**: Computed properties cache conversion results
- **Truncation**: Uses `Math.trunc()` for display amounts to avoid floating-point precision issues

## 9. Configuration Examples

### 9.1 Using European Central Bank

```ini
[exchange_rates]
data_source = euro_central_bank
proxy = system
request_timeout = 30
skip_tls_verify = false
```

### 9.2 Using User Custom Rates

```ini
[exchange_rates]
data_source = user_custom
```

Then users can set rates via UI or API.

## 10. Extension Points

### 10.1 Adding New Data Source

1. Create new file: `pkg/exchangerates/new_bank_datasource.go`
2. Implement `HttpExchangeRatesDataSource` interface:
   - `BuildRequests()`: Return HTTP requests
   - `Parse()`: Parse response to `LatestExchangeRateResponse`
3. Add data source constant in `pkg/settings/setting.go`
4. Add initialization case in `InitializeExchangeRatesDataSource()`

### 10.2 Custom Rate Management

Users can:
- Set custom rates via UI
- Update rates via API
- Delete custom rates
- Mix custom rates with data source rates (if data source supports it)

## 11. Key Files Reference

### Backend
- `pkg/exchangerates/exchange_rates_data_provider_container.go`: Container and initialization
- `pkg/exchangerates/common_http_exchange_rates_data_provider.go`: HTTP provider implementation
- `pkg/exchangerates/euro_central_bank_datasource.go`: Example data source
- `pkg/api/exchange_rates.go`: API handlers
- `pkg/models/exchange_rate.go`: Data models
- `pkg/services/user_custom_exchange_rates.go`: Custom rates service

### Frontend
- `src/stores/exchangeRates.ts`: Exchange rates store
- `src/lib/numeral.ts`: Currency conversion functions
- `src/stores/account.ts`: Account balance calculations
- `src/stores/transaction.ts`: Transaction amount conversions
- `src/stores/statistics.ts`: Statistics calculations

## 12. Summary

The exchange rates system provides:

1. **Flexible Data Sources**: 17+ remote data sources plus user custom rates
2. **Automatic Updates**: Configurable auto-refresh on app startup
3. **Efficient Caching**: LocalStorage caching with smart invalidation
4. **Multi-Currency Support**: Seamless conversion across all currencies
5. **Error Resilience**: Graceful handling of missing rates with incomplete indicators
6. **User Control**: Custom rate management for manual overrides

The implementation follows clean architecture principles with clear separation between data fetching, processing, and presentation layers, making it easy to extend and maintain.

