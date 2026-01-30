# Stock Price Provider Implementation Design

## 1. Overview

The goal is to implement a stock price fetching service for ezbookkeeping, enabling users to track the real-time valuation of their stock/securities holdings. This service will follow the existing "Provider/DataSource" architecture but will be configured dynamically via the database and UI, rather than static configuration files.

## 2. Architecture

The implementation will be located in a new package `pkg/stocks`.

```
Frontend (Settings UI) <-> API <-> Database (Configs)
                                      ^
                                      |
Frontend (Display) <-> API <-> DataProviderContainer <-> Remote APIs (Yahoo, etc.)
```

## 3. Database Design

To support dynamic configuration, we will add tables to store stock symbols and provider settings.

### 3.1 Tables

**Table: `stock_symbols`**

| Column | Type | Constraints | Description |
| :--- | :--- | :--- | :--- |
| `stock_id` | BIGINT | PK, Auto Increment | Internal ID |
| `market` | VARCHAR(20) | | e.g., "US", "HK", "CN" (Optional, for grouping) |
| `symbol` | VARCHAR(20) | NOT NULL, UNIQUE | e.g., "AAPL", "0700.HK" |
| `name` | VARCHAR(100) | | e.g., "Apple Inc." |
| `is_hidden` | BOOLEAN | Default FALSE | |
| `display_order` | INT | Default 0 | |

**Table: `stock_global_configs`** (or integrated into general settings)

| Column | Type | Description |
| :--- | :--- | :--- |
| `data_source` | VARCHAR(50) | Default 'yahoo_finance' |
| `api_key` | VARCHAR(255) | For Alpha Vantage etc. |
| `request_timeout` | INT | |
| `update_frequency`| VARCHAR(20) | Cron schedule or interval |

## 4. API Design

### 4.1 Management Endpoints

*   `GET /api/v1/stocks/config`: Get global settings and list of tracked stocks.
*   `POST /api/v1/stocks/config`: Update global settings (data source, keys).
*   `POST /api/v1/stocks`: Add a new stock symbol.
*   `PUT /api/v1/stocks/:symbol`: Update stock details (hide/show).
*   `DELETE /api/v1/stocks/:symbol`: Remove stock.

### 4.2 Price Endpoint

*   `GET /api/v1/stocks/latest.json`: Returns latest prices for all visible stocks.

## 5. Frontend Implementation (Settings)

A new "Stock Prices" tab will be added to **Application Settings**.

### 5.1 Features

*   **Provider Settings**: Choose "Yahoo Finance" or "Alpha Vantage". Input API keys if needed.
*   **Stock Watchlist**:
    *   Table showing Symbol, Name, Market.
    *   "Add Stock" dialog.
    *   Action buttons: Delete, Hide/Show.

### 5.2 Store

*   `src/stores/stockPrices.ts` will handle both the price data (for display) and the configuration actions (for settings).

## 6. Backend Implementation Details

### 6.1 Package `pkg/stocks`

*   **Manager**: `StockService` handles DB operations for adding/removing stocks.
*   **Provider**: `StockPriceDataProvider` fetches prices for the list of symbols retrieved from `StockService`.

### 6.2 Data Sources

1.  **Yahoo Finance (`yahoo_finance`)**:
    *   Supports symbols like `AAPL`, `0700.HK`.
    *   No API key usually required for basic scraping/API.
2.  **Alpha Vantage (`alphavantage`)**:
    *   Requires API Key.

## 7. Performance and Optimization

*   **Database Caching**: Configuration is cached in memory and reloaded only on change (via signal or polling).
*   **Price Caching**: Prices are cached for ~5-15 minutes to avoid hitting rate limits.
*   **Batching**: Requests for multiple symbols should be batched if the provider supports it (e.g., `?symbols=AAPL,GOOG`).

## 8. Summary

By moving configuration to the database:
1.  **Flexibility**: Users can track any stock supported by the provider without editing server files.
2.  **Usability**: Configuration is done via a friendly GUI.
3.  **Persistence**: Settings are backed up with the database.
