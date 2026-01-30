# Cryptocurrency Price Implementation Design

## Overview

This document provides the design for the cryptocurrency price fetching feature in ezBookkeeping. The system will fetch cryptocurrency prices from remote sources, similar to how exchange rates are currently handled. This feature allows users to track cryptocurrency account balances and convert them to their default currency.

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
│  ┌──────────────────────────────────────────────────┐   │
│  │  Settings / Management UI                        │   │
│  │  - Manage supported cryptocurrencies             │   │
│  │  - Configure data source                         │   │
│  └──────────────────────────────────────────────────┘   │
└─────────────────────────────────────────────────────────┘
                        │ HTTP API
                        │ GET /api/v1/cryptocurrencies (List/Config)
                        │ POST /api/v1/cryptocurrencies (Update)
                        │ GET /api/v1/cryptocurrency/latest.json (Prices)
                        ▼
┌─────────────────────────────────────────────────────────┐
│              Backend (Go + Gin)                         │
│  ┌──────────────────────────────────────────────────┐   │
│  │  API Layer (pkg/api/cryptocurrency.go)          │   │
│  │  - LatestCryptocurrencyPriceHandler              │   │
│  │  - GetCryptocurrencyConfigsHandler               │   │
│  │  - UpdateCryptocurrencyConfigsHandler            │   │
│  └──────────────────────────────────────────────────┘   │
│  ┌──────────────────────────────────────────────────┐   │
│  │  Service/Data Layer                              │   │
│  │  - Load configs from Database                    │   │
│  │  - Persist configs to Database                   │   │
│  └──────────────────────────────────────────────────┘   │
│  ┌──────────────────────────────────────────────────┐   │
│  │  CryptocurrencyPriceDataProviderContainer        │   │
│  │  - Manages data provider instances               │   │
│  │  - Strategy pattern for data sources             │   │
│  └──────────────────────────────────────────────────┘   │
│  ┌──────────────────────────────────────────────────┐   │
│  │  Data Sources (HttpCryptocurrencyPriceDataSource)│   │
│  │  - CoinGeckoDataSource                           │   │
│  │  - CoinMarketCapDataSource                       │   │
│  │  - BinanceDataSource                             │   │
│  └──────────────────────────────────────────────────┘   │
└─────────────────────────────────────────────────────────┘
                        │ HTTP Request
                        ▼
┌─────────────────────────────────────────────────────────┐
│         Remote Cryptocurrency Price APIs                │
│  - CoinGecko API                                        │
│  - CoinMarketCap API                                    │
│  - Binance API                                          │
└─────────────────────────────────────────────────────────┘
```

## 2. Design Principles

### 2.1 Independence from Exchange Rates System

- The cryptocurrency price system must operate independently from the fiat exchange rate system.
- Data sources, providers, and stores for cryptocurrency prices should be distinct from those used for fiat exchange rates.

### 2.2 Database-Driven Configuration

- **Dynamic**: Cryptocurrency selection and data source configuration are stored in the database, not in `conf/app.ini`.
- **UI-Managed**: Users can add/remove cryptocurrencies and change data sources via the "Application Settings" page.
- **Persisted**: Configurations persist across restarts and are synchronized across devices (if using a central DB).

### 2.3 Flexible Base Currency Support

- Cryptocurrency prices can be fetched in various base currencies (e.g., USD, CNY, EUR) as supported by the data source.
- This allows direct valuation in the user's primary currency without relying on internal fiat exchange rate conversion.

## 3. Configuration & Database Design

### 3.1 Database Schema

We will introduce a new table to store the supported cryptocurrencies and system-wide settings for the module.

**Table: `cryptocurrencies`** (Stores the list of coins to track)

| Column | Type | Constraints | Description |
| :--- | :--- | :--- | :--- |
| `cryptocurrency_id` | BIGINT | PK, Auto Increment | Internal ID |
| `symbol` | VARCHAR(20) | NOT NULL, UNIQUE | e.g., "BTC", "ETH" |
| `name` | VARCHAR(100) | NOT NULL | e.g., "Bitcoin" |
| `is_hidden` | BOOLEAN | Default FALSE | If true, do not fetch price |
| `display_order` | INT | Default 0 | Sorting order in UI |

**Global Settings** (Stored in existing `settings` table or a new `external_data_sources` table)

Instead of a specific table for global config, we can use the existing `settings` table key-value structure (if available) or add a dedicated single-row table `cryptocurrency_exchange_rate_configs`:

| Column | Type | Constraints | Description |
| :--- | :--- | :--- | :--- |
| `config_id` | BIGINT | PK | Single row ID |
| `data_source` | VARCHAR(50) | Default 'coingecko' | "coingecko", "coinmarketcap", etc. |
| `base_currency` | VARCHAR(10) | Default 'USD' | Fiat currency to fetch prices in |
| `api_key` | VARCHAR(255) | Nullable | API Key for source |
| `request_timeout` | INT | Default 10000 | Timeout in ms |
| `proxy_type` | VARCHAR(20) | | Proxy setting |
| `proxy_address` | VARCHAR(255)| | Proxy address |

*Note: For simplicity, we might store these global configs as JSON in a general settings table or as individual rows in a key-value settings table.*

### 3.2 Configuration Loading

- `pkg/settings` currently loads from `.ini`.
- New service `CryptocurrencyService` will load these configurations from the database on startup and on demand.
- The `CryptocurrencyPriceDataProviderContainer` needs to be re-initialized if the data source changes in the settings.

## 4. Backend Implementation Design

### 4.1 Package Structure

Create new package: `pkg/cryptocurrency/`

**Core Files**:
- `service.go`: CRUD for cryptocurrencies and configs.
- `cryptocurrency_price_data_provider.go`: Interface definition.
- `cryptocurrency_price_data_provider_container.go`: Container.
- `common_http_cryptocurrency_price_data_provider.go`: Common HTTP provider.

### 4.2 Data Models (Backend)

**Config Model**:
```go
type CryptocurrencyConfig struct {
    DataSource   string `json:"dataSource"`
    BaseCurrency string `json:"baseCurrency"`
    ApiKey       string `json:"apiKey"`
    // ...
}

type Cryptocurrency struct {
    Symbol       string `json:"symbol"`
    Name         string `json:"name"`
    IsHidden     bool   `json:"isHidden"`
    DisplayOrder int    `json:"displayOrder"`
}
```

### 4.3 Initialization Flow

1. System startup: `InitializeCryptocurrencyService()`
2. Load config from DB.
3. Initialize `CryptocurrencyPriceDataProviderContainer` with DB config.
4. **Default Seeding**: If the `cryptocurrencies` table is empty, the system will automatically seed it with the following default cryptocurrencies:
   - **BTC** (Bitcoin)
   - **ETH** (Ethereum)
   - **ATOM** (Cosmos)
   - **SOL** (Solana)
   - **ADA** (Cardano)

## 5. API Design

### 5.1 Configuration Management Endpoints

**Get Configuration & List**: `GET /api/v1/cryptocurrencies`
- Returns global settings (data source) and list of cryptocurrencies.

**Update Configuration**: `POST /api/v1/cryptocurrencies/config`
- Updates data source, API key, base currency.

**Add/Update Cryptocurrency**: `POST /api/v1/cryptocurrencies`
- Adds a new symbol or updates existing (hidden/visible).

**Delete Cryptocurrency**: `DELETE /api/v1/cryptocurrencies/:symbol`
- Removes a symbol from tracking.

### 5.2 Price Endpoint

**Route**: `GET /api/v1/cryptocurrency/latest.json`
- Logic remains similar but uses the list of *visible* cryptocurrencies from the DB.

## 6. Frontend Implementation Design

### 6.1 Settings UI

**Desktop**:
Add a new tab in **Application Settings**.
**Location**: `src/views/desktop/app/settings/tabs/AppCryptocurrencySettingTab.vue`

**Mobile**:
Add a new settings page.
**Location**: `src/views/mobile/settings/CryptocurrencySettingsPage.vue`
**Entry Point**: Add a link in the main "Settings" list on mobile.

**Features (Both Platforms)**:
1. **Data Source Configuration**:
   - Dropdown for "Data Source" (CoinGecko, CoinMarketCap, Binance).
   - Input for "API Key" (shown if source requires it).
   - Input for "Base Currency" (USD, CNY, etc.).
2. **Cryptocurrency Management**:
   - List of tracked cryptocurrencies.
   - "Add" button to add a new symbol (e.g., input "SOL", "Solana").
   - Toggle switch for "Enable/Disable" (hide).
   - Delete button.
   - Drag-and-drop reordering (optional).

### 6.2 Store Design

**Store**: `src/stores/cryptocurrencyPrices.ts`
- Actions:
  - `loadConfig()`: Fetch settings from API.
  - `saveConfig(config)`: Save settings to API.
  - `addCryptocurrency(symbol, name)`
  - `deleteCryptocurrency(symbol)`

## 7. Migration Plan

Since this is a new feature, no complex data migration is needed for existing user data. However:
1. **Schema Migration**: Create the new tables.
2. **Default Seeding**: On first run (or migration), populate `cryptocurrencies` table with the default set (BTC, ETH, ATOM, SOL, ADA) and default config (CoinGecko).

## 8. Summary of Changes

| Feature | Old Design (File-based) | New Design (DB-based) |
| :--- | :--- | :--- |
| **Config Storage** | `conf/app.ini` | Database Tables |
| **Coin List** | Static list in `.ini` | Dynamic list in DB |
| **Management** | Edit file & Restart | UI Settings Page |
| **Extensibility** | Manual code/config change | User can add any coin supported by source |

## 9. Auto-Updating Mechanism

The system includes a background cron job to keep prices up-to-date even when no users are actively requesting them:

1.  **Cron Job**: `UpdateCryptocurrencyPricesJob` runs periodically (default every 5 minutes).
2.  **Logic**:
    *   Reads the current configuration and list of `is_hidden=false` cryptocurrencies from the database.
    *   Calls the data provider to fetch prices for these symbols.
    *   Updates the in-memory cache with the new prices.
3.  **Configuration**: Can be enabled/disabled via the "Auto-update Cryptocurrency Prices" setting in the UI (persisted to DB).
