# 加密货币价格功能实现设计

## 概述

本文档提供了 ezBookkeeping 中加密货币价格获取功能的设计方案。系统将从远程源获取加密货币价格，类似于目前处理汇率的方式。此功能允许用户追踪加密货币账户余额并将其转换为本位币。

## 1. 架构概览

加密货币价格系统将遵循与汇率系统相同的 **策略模式 (Strategy Pattern)** 和 **容器模式 (Container Pattern)**：

```
┌─────────────────────────────────────────────────────────┐
│              前端 (Vue 3 + Pinia)                        │
│  ┌──────────────────────────────────────────────────┐   │
│  │  CryptocurrencyPricesStore                       │   │
│  │  - 管理加密货币价格状态                             │   │
│  │  - LocalStorage 缓存                              │   │
│  │  - 价格转换计算                                   │   │
│  └──────────────────────────────────────────────────┘   │
└─────────────────────────────────────────────────────────┘
                        │ HTTP API
                        │ GET /api/v1/cryptocurrency/latest.json
                        ▼
┌─────────────────────────────────────────────────────────┐
│              后端 (Go + Gin)                             │
│  ┌──────────────────────────────────────────────────┐   │
│  │  API 层 (pkg/api/cryptocurrency.go)               │   │
│  │  - LatestCryptocurrencyPriceHandler              │   │
│  └──────────────────────────────────────────────────┘   │
│  ┌──────────────────────────────────────────────────┐   │
│  │  CryptocurrencyPriceDataProviderContainer        │   │
│  │  - 管理数据提供者实例                               │   │
│  │  - 数据源策略模式                                 │   │
│  └──────────────────────────────────────────────────┘   │
│  ┌──────────────────────────────────────────────────┐   │
│  │  数据提供者 (Data Providers)                      │   │
│  │  - CommonHttpCryptocurrencyPriceDataProvider      │   │
│  └──────────────────────────────────────────────────┘   │
│  ┌──────────────────────────────────────────────────┐   │
│  │  数据源 (HttpCryptocurrencyPriceDataSource)        │   │
│  │  - CoinGeckoDataSource                           │   │
│  │  - CoinMarketCapDataSource                       │   │
│  │  - BinanceDataSource                             │   │
│  │  - ... (其他数据源)                                │   │
│  └──────────────────────────────────────────────────┘   │
└─────────────────────────────────────────────────────────┘
                        │ HTTP 请求
                        ▼
┌─────────────────────────────────────────────────────────┐
│              远程加密货币价格 API                        │
│  - CoinGecko API                                        │
│  - CoinMarketCap API                                    │
│  - Binance API                                          │
│  - ...                                                   │
└─────────────────────────────────────────────────────────┘
```

## 2. 设计原则

### 2.1 与汇率系统独立

- 加密货币价格系统必须独立于法币汇率系统运行。
- 加密货币价格的数据源、提供者和存储应与法币汇率使用的区分开。
- 尽可能避免将法币汇率作为加密货币估值的强制中间环节。

### 2.2 基于配置的加密货币选择

与汇率数据源的配置方式类似，加密货币的选择将是：
- **基于配置**：在 `conf/ezbookkeeping.ini` 中指定。
- **基于列表**：支持在配置中列出多个加密货币。
- **灵活**：无需修改代码即可轻松添加/删除加密货币。

### 2.3 灵活的基准货币支持

- 加密货币价格可以根据数据源的支持，以各种基准法币（如 USD, CNY, EUR）获取。
- 这样可以直接以用户的主币种进行估值，而无需依赖内部法币汇率转换。
- USDT 仍可用作通用参考，但不强制作为唯一的基准货币。

## 3. 配置设计

### 3.1 配置节

在 `conf/ezbookkeeping.ini` 中添加新节：

```ini
[cryptocurrency]
# 加密货币价格数据源，支持：
# "coingecko": CoinGecko API (有免费档)
# "coinmarketcap": CoinMarketCap API (需要 API key)
# "binance": Binance 公开 API
data_source = coingecko

# 待获取的加密货币符号列表，以逗号分隔
# 例如：BTC,ETH,BNB,SOL,ADA
cryptocurrencies = BTC,ETH,BNB

# 加密货币价格的基准法币
# 如果数据源支持，价格将以此货币获取。
# 默认为 USD。
base_currency = USD

# 请求超时时间 (0 - 4294967295 毫秒)
# 默认为 10000 (10 秒)
request_timeout = 10000

# 代理设置
proxy = system

# 跳过 TLS 验证
skip_tls_verify = false

# API key (可选，某些数据源如 CoinMarketCap 需要)
api_key = 
```

### 3.2 配置加载

- 在 `pkg/settings/setting.go` 中添加配置加载逻辑。
- 类似于 `loadExchangeRatesConfiguration()`。
- 函数：`loadCryptocurrencyConfiguration()`。
- 存储在 `Config` 结构体中：
  - `CryptocurrencyDataSource`
  - `CryptocurrencySymbols` (字符串切片)
  - `CryptocurrencyBaseCurrency`
  - `CryptocurrencyRequestTimeout`
  - `CryptocurrencyProxy`
  - `CryptocurrencySkipTLSVerify`
  - `CryptocurrencyAPIKey`

## 4. 后端实现设计

### 4.1 包结构

创建新包：`pkg/cryptocurrency/`

**核心文件**：
- `cryptocurrency_price_data_provider.go`：接口定义
- `cryptocurrency_price_data_provider_container.go`：容器和初始化
- `common_http_cryptocurrency_price_data_provider.go`：通用 HTTP 提供者实现

**数据源文件**：
- `coingecko_datasource.go`：CoinGecko API 实现
- `coinmarketcap_datasource.go`：CoinMarketCap API 实现
- `binance_datasource.go`：Binance API 实现

### 4.2 数据提供者接口

```go
type CryptocurrencyPriceDataProvider interface {
    GetLatestCryptocurrencyPrices(
        c core.Context, 
        uid int64, 
        currentConfig *settings.Config
    ) (*models.LatestCryptocurrencyPriceResponse, error)
}
```

### 4.3 HTTP 数据源接口

```go
type HttpCryptocurrencyPriceDataSource interface {
    BuildRequests(symbols []string, baseCurrency string, apiKey string) ([]*http.Request, error)
    Parse(c core.Context, content []byte) (*models.LatestCryptocurrencyPriceResponse, error)
}
```

### 4.4 数据模型

**新模型**：`pkg/models/cryptocurrency_price.go`

```go
type LatestCryptocurrencyPriceResponse struct {
    DataSource    string                        `json:"dataSource"`
    ReferenceUrl  string                        `json:"referenceUrl"`
    UpdateTime    int64                         `json:"updateTime"`
    BaseCurrency  string                        `json:"baseCurrency"` // 例如 "USD", "CNY" 或 "USDT"
    Prices        LatestCryptocurrencyPriceSlice `json:"prices"`
}

type LatestCryptocurrencyPrice struct {
    Symbol string `json:"symbol"` // 例如 "BTC", "ETH"
    Price  string `json:"price"`  // 以基准货币计的价格
}
```

### 4.5 数据源实现策略

**CoinGecko (推荐免费使用)**：
- API: `https://api.coingecko.com/api/v3/simple/price`
- 参数: `ids=bitcoin,ethereum&vs_currencies=usd,cny`
- 免费档: 无需 API key，有频率限制。
- 响应: JSON 格式。

**CoinMarketCap**：
- API: `https://pro-api.coinmarketcap.com/v1/cryptocurrency/quotes/latest`
- 需要在 Header 中提供 API key。
- 支持通过 `convert` 参数转换多种法币。
- 响应: JSON 格式。

**Binance**：
- API: `https://api.binance.com/api/v3/ticker/price`
- 无需 API key。
- 默认返回以 USDT 计价的价格，但也支持其他交易对。
- 响应: JSON 格式。

### 4.6 初始化流程

1. 系统启动：`InitializeCryptocurrencyDataSource(config)`。
2. 读取配置：数据源类型、加密货币列表和基准货币。
3. 创建相应的数据提供者实例。
4. 存储在容器单例中。

### 4.7 价格获取流程

1. API 处理器接收请求。
2. 容器调用 `GetLatestCryptocurrencyPrices()`。
3. 数据提供者：
   - 从配置获取加密货币符号和基准货币。
   - 使用数据源构建 HTTP 请求。
   - 使用超时/代理设置执行请求。
   - 解析响应。
   - 将所有价格归一化到配置的基准货币。
   - 返回统一响应。

## 5. API 设计

### 5.1 端点

**路由**: `GET /api/v1/cryptocurrency/latest.json`

**处理器**: `CryptocurrencyApi.LatestCryptocurrencyPriceHandler`

**响应格式**:
```json
{
  "success": true,
  "result": {
    "dataSource": "CoinGecko",
    "referenceUrl": "https://www.coingecko.com",
    "updateTime": 1234567890,
    "baseCurrency": "USD",
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

### 5.2 API 文件

**文件**: `pkg/api/cryptocurrency.go`

- 结构类似于 `pkg/api/exchange_rates.go`。
- 用于获取最新价格的单个处理器。
- 错误处理与现有模式保持一致。

## 6. 前端实现设计

### 6.1 Store 设计

**文件**: `src/stores/cryptocurrencyPrices.ts`

**状态 (State)**:
- `latestCryptocurrencyPrices`: 带有时间戳的当前价格。

**计算属性 (Computed Properties)**:
- `cryptocurrencyPricesLastUpdateTime`: 最后更新时间戳。
- `latestCryptocurrencyPriceMap`: 符号 → 价格的映射。

**关键方法**:
- `getLatestCryptocurrencyPrices({ silent, force })`: 获取价格并处理缓存。
- `getCryptocurrencyPrice(symbol)`: 获取特定符号的价格。
- `getCryptocurrencyPriceInFiat(symbol, fiatCurrency)`: 转换为特定法币。

### 6.2 LocalStorage 缓存

- 键: `ebk_app_cryptocurrency_prices`
- 缓存有效期: 与汇率相同（当天或当前小时）。
- 结构: `{ time: number, data: LatestCryptocurrencyPriceResponse }`

### 6.3 独立性与转换

转换为法币的逻辑：
1. **直接获取 (推荐)**：通过在配置中设置 `base_currency`（如 CNY），直接获取所需法币的价格。这是最准确的估值方式，不会与内部汇率数据混淆。
2. **实时转换 (可选)**：如果目标法币与加密货币基准币种不一致，前端可以进行转换。但两个系统保持独立：
   - `CryptocurrencyPricesStore` 处理 加密货币 -> 基准法币。
   - `ExchangeRatesStore` 处理 法币 -> 法币。
3. UI 应清楚区分直接从市场源获取的价值和通过内部汇率转换得到的价值。

## 7. 货币代码处理

### 7.1 加密货币符号

- 使用标准的加密货币符号（BTC, ETH 等）。
- 这些不同于 ISO 4217 货币代码。
- 需要扩展货币验证或创建单独的验证逻辑。

### 7.2 账户币种字段

- 目前系统在 `Account.Currency` 中使用 ISO 4217 代码。
- 方案：
  1. **独立验证**：为 ISO 代码和加密货币符号维护单独的列表。更新验证逻辑以同时检查两者，但在代码中保持概念上的独立。
  2. **独立字段**：增加 `AssetSymbol` 或 `CryptocurrencySymbol` 字段（较复杂，但分离度最好）。
  3. **元数据**：在 `AccountExtend` 中存储资产类型（法币 vs 加密货币）。

**推荐**：方案 1，并保持明确的概念分离。
- `pkg/validators/currency.go` 将拥有 `AllCurrencyNames` 和 `AllCryptocurrencySymbols` 两个独立的 Map。
- `ValidCurrency` 将检查两者，但系统应感知其处理的类型。
- 前端 `ALL_CURRENCIES` 常量可以由这两组独立的数据组成。

## 8. 数据流

### 8.1 初始获取流程

```
1. 用户登录 → DesktopApp.vue / MobileApp.vue
2. 检查 autoUpdateCryptocurrencyPrices 设置
3. 调用 cryptocurrencyPricesStore.getLatestCryptocurrencyPrices()
4. 检查 localStorage 缓存有效性
5. 如果无效或强制刷新:
   a. 调用 services.getLatestCryptocurrencyPrices() (API)
   b. 后端: Cryptocurrency.LatestCryptocurrencyPriceHandler()
   c. 后端: cryptocurrency.Container.GetLatestCryptocurrencyPrices()
   d. 提供者: GetLatestCryptocurrencyPrices()
   e. 数据源: BuildRequests() → 带符号的 HTTP 请求
   f. 远程 API: 返回 JSON
   g. 数据源: Parse() → LatestCryptocurrencyPriceResponse
   h. 提供者: 归一化到配置的基准货币
   i. API: 返回 JSON 响应
   j. 前端: 更新 Store 和 localStorage
6. 返回缓存或新鲜数据
```

### 8.2 价格转换流程

```
1. 用户查看加密货币账户 (如 BTC)
2. 组件调用 cryptocurrencyPricesStore.getCryptocurrencyPriceInFiat("BTC", "CNY")
3. Store:
   a. 如果 BTC 以 CNY 计价的价格已获取 (因为 base_currency=CNY):
      - 直接返回该价格。
   b. 如果 BTC 以 CNY 计价的价格不可用 (例如 base_currency=USD):
      - 从 cryptocurrencyPricesStore 获取 BTC 的 USD 价格。
      - 从 exchangeRatesStore 获取 USD 到 CNY 的汇率。
      - 计算: btcPriceInUSD * usdToCnyRate。
4. 显示转换后的金额 (可选指示这是直接转换还是间接转换)。
```

## 15. 总结

加密货币价格系统将：

1. **系统独立**：加密货币价格和法币汇率由不同的模块处理。
2. **配置驱动**：易于配置追踪哪些加密货币以及使用哪种基准货币。
3. **多数据源**：支持 CoinGecko, Binance, CoinMarketCap 等。
4. **灵活基准币**：可直接以 USD, CNY 或其他受支持的法币获取价格。
5. **高效缓存**：具有智能失效机制的 LocalStorage 缓存。
6. **解耦集成**：仅在必要时配合汇率系统进行二次转换。
7. **可扩展性**：易于添加新的数据源和加密货币符号。
