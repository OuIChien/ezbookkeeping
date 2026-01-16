# ezbookkeeping Cryptocurrency and Stock Support Proposal

## 1. Goal
Enable ezbookkeeping to manage non-monetary assets (cryptocurrencies, stocks, funds, etc.) by tracking the held **Quantity**, recording the purchase price, and automatically calculating the total valuation based on real-time market data.

## 2. Proposed Account Structure Changes

### 2.1 Database Schema Extensions (`Account` Table)
- **Currency Field Length**: Currently, `Currency` is `VARCHAR(3)` (ISO 4217). It should be extended to `VARCHAR(10)` to support stock tickers (e.g., `AAPL`, `TSLA`) and longer crypto symbols (e.g., `SHIB`).
- **Balance Precision**:
    - The existing `Balance` (int64) represents sub-units (e.g., cents for fiat).
    - **Recommended**: For non-fiat assets, the `Balance` field represents the **Held Quantity**. To support high precision (e.g., 8 decimal places for BTC), the storage unit should be $10^{-8}$.
    - **Alternative**: Add a `Quantity` field in `AccountExtend` for asset quantity, while keeping the main `Balance` as the converted fiat value.

### 2.2 Account Category Refinement
- Maintain `ACCOUNT_CATEGORY_INVESTMENT`.
- Add an `AssetType` identifier in UI or `AccountExtend`:
    - `FIAT`: Standard currency accounts.
    - `CRYPTO`: Cryptocurrency accounts.
    - `STOCK`: Stock/Securities accounts.

---

## 3. Business Logic Enhancements

### 3.1 Valuation Logic
- **Real-time Pricing**:
    - **Cryptocurrency**: Directly integrate with the existing cryptocurrency price service in the project (located in `pkg/cryptocurrency`), which already supports fetching prices from sources like CoinGecko and Binance.
    - **Stocks/Securities**: Add a new stock price service (e.g., Yahoo Finance, Alpha Vantage).
- **Valuation Calculation**:
    - Total Account Value = Held Quantity × Real-time Market Price.
    - The system automatically converts the valuation to the user's "Default Currency" for total asset display.

### 3.2 Transaction Logic
- **Purchase**:
    - From: Fiat Account (decrease balance).
    - To: Asset Account (increase quantity).
    - Recording Price: Store the execution price in the transaction `Comment` or a new `Price` extension field.
- **Sale**:
    - From: Asset Account (decrease quantity).
    - To: Fiat Account (increase balance).
- **Dividends/Interest**:
    - Record as `INCOME` to either the asset account (e.g., stock splits) or a fiat account (e.g., cash dividends).

---

## 4. Validation and Constraints
- **Extended Currency Validation**: Update `pkg/validators/currency.go` to allow custom symbols beyond ISO 4217.
- **Deletion Constraints**: Strengthen constraints for accounts with active holdings to ensure audit trail integrity.

---

## 5. UI/UX Recommendations
- **Asset Overview**: On the account list page, display both the held quantity and its equivalent value in the user's default currency.
- **Visualizations**: Add "Asset Composition" charts to differentiate between Fiat, Crypto, and Securities.

---

## 6. Suggested Implementation Phases
1.  **Phase 1**: Modify DB schema for `Currency` field length and update validation logic.
2.  **Phase 2**: Integrate existing cryptocurrency price service and implement backend stock price fetching service (Stock Price Provider).
3.  **Phase 3**: Update the frontend account model to handle `AssetType` and display live valuations.
4.  **Phase 4**: Optimize the transaction entry flow to support "Buy/Sell by Unit Price".
