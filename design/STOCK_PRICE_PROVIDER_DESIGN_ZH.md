# 股票价格获取服务 (Stock Price Provider) 实现设计方案

## 1. 概述
本方案旨在为 ezbookkeeping 实现股票价格获取服务，使用户能够实时跟踪其股票/证券资产的估值。该服务将遵循项目中现有的“提供商/数据源 (Provider/DataSource)”架构模式，这与加密货币和汇率模块的设计保持一致。

## 2. 架构设计
实现将位于新的包 `pkg/stocks` 中，包含以下核心组件：

- **DataProvider 接口**: 定义获取价格的高层方法。
- **通用 HTTP 提供商 (Common HTTP Provider)**: 处理通用的 HTTP 请求执行、错误处理以及代理/超时设置。
- **DataSource 接口**: 定义如何构建特定的 API 请求并解析其响应。
- **容器 (Container)**: 管理已配置的提供商的初始化和获取。

### 2.1 建议目录结构
```text
pkg/stocks/
├── stock_price_data_provider.go           # 接口定义
├── common_http_stock_price_data_provider.go # 基础 HTTP 实现
├── stock_price_data_provider_container.go  # 注册与工厂类
└── yahoo_finance_datasource.go             # Yahoo Finance 具体实现
```

## 3. 数据模型
将在 `pkg/models/stock_price.go` 中添加新模型，以标准化响应格式。

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

## 4. 配置项
系统配置中将增加一个新的 `[stocks]` 模块。

| 配置项 | 描述 | 默认值 |
|------|-------------|---------|
| `data_source` | 股票数据源 (例如：`yahoo_finance`) | - |
| `stocks` | 股票代码列表 (例如：`AAPL,TSLA,0700.HK`) | - |
| `request_timeout` | API 请求超时时间（毫秒） | `10000` |
| `proxy` | 请求使用的代理服务器 | `system` |
| `api_key` | 特定数据源的可选 API 密钥 | - |

## 5. 计划集成的数据源
1.  **Yahoo Finance (`yahoo_finance`)**:
    - 首选数据源，因其覆盖全球市场范围广。
    - 支持如 `AAPL`, `0700.HK` (港股), `600519.SS` (A股) 等代码。
2.  **Alpha Vantage (`alphavantage`)**:
    - 备选/后续扩展数据源（需要 API Key）。

## 6. 实施阶段
1.  **第一阶段**: 在 `pkg/models` 中定义数据模型。
2.  **第二阶段**: 扩展 `pkg/settings` 以支持股票相关的配置项。
3.  **第三阶段**: 在 `pkg/stocks` 中实现核心接口和 `CommonHttpStockPriceDataProvider`。
4.  **第四阶段**: 实现 `YahooFinanceDataSource` 及其解析逻辑。
5.  **第五阶段**: 将提供商集成到容器中，并为后续的资产估值业务逻辑提供调用接口。
