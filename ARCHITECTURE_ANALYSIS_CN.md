# ezBookkeeping 架构分析

## 1. 项目概述

**ezBookkeeping** 是一个轻量级、自托管的个人财务管理应用，具有友好的用户界面和强大的记账功能。它被设计为资源高效且高度可扩展，可以在小到树莓派的设备上运行，也可以扩展到大型集群环境。

### 核心特性
- 开源且自托管
- 多平台支持（桌面、移动、PWA）
- 多语言和多货币支持
- AI 驱动的收据图像识别
- 支持多种数据导入/导出格式
- 双因素认证（2FA）
- OAuth2 认证
- MCP（模型上下文协议）支持

## 2. 技术栈

### 后端
- **语言**: Go 1.25+
- **Web 框架**: Gin（HTTP Web 框架）
- **ORM**: XORM（支持 SQLite、MySQL、PostgreSQL）
- **CLI 框架**: urfave/cli/v3
- **认证**: JWT (golang-jwt/jwt/v5)、OAuth2、2FA (pquerna/otp)
- **任务调度**: go-co-op/gocron/v2
- **日志**: logrus
- **配置**: gopkg.in/ini.v1

### 前端
- **框架**: Vue 3（组合式 API）
- **语言**: TypeScript
- **构建工具**: Vite 7.x
- **桌面 UI**: Vuetify 3
- **移动 UI**: Framework7 9.x
- **状态管理**: Pinia
- **路由**: Vue Router 4
- **图表**: ECharts (vue-echarts)
- **地图**: Leaflet
- **日期处理**: Moment.js + moment-timezone
- **PWA**: vite-plugin-pwa

### 数据库
- **支持**: SQLite3、MySQL、PostgreSQL
- **ORM**: XORM，带连接池

### 存储
- **本地存储**: 文件系统
- **对象存储**: MinIO（S3 兼容）
- **WebDAV**: 支持

## 3. 整体架构

### 架构模式
项目采用**分层架构**，职责清晰分离：

```
┌─────────────────────────────────────────┐
│         前端 (Vue 3)                     │
│  ┌──────────┐  ┌──────────┐            │
│  │  桌面端   │  │  移动端   │            │
│  │ (Vuetify)│  │(Framework7)│           │
│  └──────────┘  └──────────┘            │
└─────────────────────────────────────────┘
                  │ HTTP/REST API
                  │ JSON
┌─────────────────────────────────────────┐
│      后端 (Go + Gin)                     │
│  ┌──────────────────────────────────┐  │
│  │  API 层 (pkg/api)                │  │
│  │  - RESTful 端点                  │  │
│  │  - JSON-RPC (MCP)                │  │
│  └──────────────────────────────────┘  │
│  ┌──────────────────────────────────┐  │
│  │  服务层 (pkg/services)            │  │
│  │  - 业务逻辑                       │  │
│  └──────────────────────────────────┘  │
│  ┌──────────────────────────────────┐  │
│  │  数据层 (pkg/datastore)           │  │
│  │  - XORM ORM                       │  │
│  │  - 数据库抽象                      │  │
│  └──────────────────────────────────┘  │
└─────────────────────────────────────────┘
                  │
┌─────────────────────────────────────────┐
│   数据库 (SQLite/MySQL/PostgreSQL)      │
└─────────────────────────────────────────┘
```

### 请求流程
1. **客户端请求** → 前端（Vue Router）
2. **API 调用** → Axios → 后端 API 端点
3. **中间件** → 认证、请求 ID、日志记录
4. **API 处理器** → 业务逻辑验证
5. **服务层** → 核心业务操作
6. **数据存储** → 通过 XORM 进行数据库操作
7. **响应** → JSON 响应返回前端

## 4. 目录结构

### 根目录
```
ezbookkeeping/
├── cmd/              # CLI 命令（server、database、user_data 等）
├── pkg/              # 后端包（Go）
├── src/              # 前端源代码（Vue/TypeScript）
├── public/           # 静态资源
├── conf/             # 配置文件
├── templates/        # 邮件和提示模板
├── docker/           # Docker 构建脚本
└── testdata/         # 测试数据文件
```

### 后端结构 (`pkg/`)

#### 核心模块
- **`api/`**: HTTP API 处理器（RESTful 端点）
  - `accounts.go`: 账户管理
  - `transactions.go`: 交易 CRUD 操作
  - `authorizations.go`: 认证端点
  - `users.go`: 用户管理
  - `exchange_rates.go`: 货币汇率
  - `large_language_models.go`: AI/LLM 集成
  - `model_context_protocols.go`: MCP 协议处理器

- **`services/`**: 业务逻辑层
  - 包含各种业务操作的服务实现
  - 作为 API 处理器和数据存储之间的中介

- **`datastore/`**: 数据库抽象层
  - `database.go`: 数据库连接和事务管理
  - `datastore.go`: 数据存储容器，支持分片
  - 使用 XORM 进行 ORM 操作

- **`core/`**: 核心类型和工具
  - 上下文包装器（CLI、Web、Cron）
  - 通用类型和处理器
  - 日历、货币、日期时间工具

- **`models/`**: 数据库模型（XORM 的结构体定义）

- **`middlewares/`**: HTTP 中间件
  - JWT 认证
  - 请求日志
  - 错误恢复
  - 请求 ID 生成

#### 支持模块
- **`auth/oauth2/`**: OAuth2 认证提供者
- **`avatars/`**: 头像提供者实现（Gravatar、内部、空）
- **`converters/`**: 数据导入转换器（CSV、OFX、QIF、Excel 等）
- **`cron/`**: 定时任务管理
- **`errs/`**: 错误定义
- **`exchangerates/`**: 汇率数据源
- **`llm/`**: 大语言模型提供者
- **`locales/`**: 国际化（22 种语言）
- **`log/`**: 日志工具
- **`mail/`**: 邮件发送
- **`mcp/`**: 模型上下文协议实现
- **`settings/`**: 配置管理
- **`storage/`**: 对象存储（本地、MinIO、WebDAV）
- **`utils/`**: 工具函数
- **`uuid/`**: UUID 生成
- **`validators/`**: 输入验证

### 前端结构 (`src/`)

#### 入口点
- **`index-main.ts`**: 检测设备类型并重定向的入口点
- **`desktop-main.ts`**: 桌面应用入口
- **`mobile-main.ts`**: 移动应用入口

#### 核心目录
- **`components/`**: Vue 组件
  - `base/`: 基础组件（共享逻辑）
  - `common/`: 通用组件（DateTimePicker、MapView 等）
  - `desktop/`: 桌面专用组件
  - `mobile/`: 移动专用组件

- **`views/`**: 页面视图
  - `base/`: 基础视图类
  - `desktop/`: 桌面页面
  - `mobile/`: 移动页面

- **`stores/`**: Pinia 状态管理
  - `account.ts`: 账户状态
  - `transaction.ts`: 交易状态
  - `user.ts`: 用户状态
  - `token.ts`: 认证令牌
  - `exchangeRates.ts`: 汇率
  - 等等

- **`router/`**: Vue Router 配置
  - `desktop.ts`: 桌面路由
  - `mobile.ts`: 移动路由

- **`lib/`**: 工具库
  - `api.ts`: API 客户端（Axios 包装器）
  - `services.ts`: 服务层
  - `map/`: 地图提供者（Leaflet、Google Maps、百度、高德）
  - `ui/`: UI 工具

- **`core/`**: 核心工具
  - `account.ts`、`transaction.ts`、`category.ts`: 领域模型
  - `datetime.ts`、`currency.ts`: 格式化工具
  - `statistics.ts`: 统计计算

- **`consts/`**: 常量
- **`models/`**: TypeScript 类型定义
- **`locales/`**: 前端国际化文件（JSON）
- **`styles/`**: SCSS 样式表

## 5. 核心组件

### 5.1 后端组件

#### 命令结构 (`cmd/`)
- **`webserver.go`**: Web 服务器初始化和路由
- **`database.go`**: 数据库管理命令
- **`user_data.go`**: 用户数据管理 CLI
- **`cron_jobs.go`**: 定时任务管理
- **`security.go`**: 安全工具
- **`initializer.go`**: 系统初始化（配置、数据库、存储等）

#### API 层 (`pkg/api/`)
每个 API 模块遵循相似的模式：
- 接收 `core.Context` 的处理器函数
- 输入验证
- 服务层调用
- 响应格式化

示例结构：
```go
var Accounts = &accountsApi{}

type accountsApi struct{}

func (a *accountsApi) AccountListHandler(c core.Context) (interface{}, error) {
    // 验证、服务调用、响应
}
```

#### 数据存储 (`pkg/datastore/`)
- **容器模式**: `DataStoreContainer` 持有多个存储
  - `UserStore`: 用户相关数据
  - `TokenStore`: 认证令牌
  - `UserDataStore`: 用户特定数据（账户、交易等）
- **分片支持**: 为水平扩展而设计（当前使用单个数据库）
- **事务管理**: `DoTransaction` 方法用于原子操作

#### 服务层 (`pkg/services/`)
- 业务逻辑实现
- 数据验证和转换
- 调用数据存储层

### 5.2 前端组件

#### 应用结构
- **桌面应用**: `DesktopApp.vue` - 主桌面应用组件
- **移动应用**: `MobileApp.vue` - 主移动应用组件
- 两者都使用 Vue 3 组合式 API

#### 状态管理（Pinia）
- 按领域组织存储（账户、交易、用户等）
- 存储处理 API 调用和状态更新
- 响应式状态用于 UI 更新

#### API 客户端 (`lib/api.ts`)
- 基于 Axios 的 HTTP 客户端
- 请求/响应拦截器
- 错误处理
- 令牌管理

#### 路由
- 桌面和移动端独立的路由器
- 路由守卫用于认证
- 支持懒加载

## 6. 关键设计模式

### 6.1 依赖注入
- 服务的容器模式（如 `DataStoreContainer`、`AvatarProviderContainer`）
- 共享资源的单例实例

### 6.2 中间件模式
- HTTP 中间件的责任链
- 认证、日志、错误恢复

### 6.3 仓库模式
- 数据存储抽象层
- 数据访问与业务逻辑分离

### 6.4 工厂模式
- 提供者工厂（头像、存储、LLM、汇率）
- 数据导入的转换器工厂

### 6.5 策略模式
- 同一接口的多种实现
  - 头像提供者（Gravatar、内部、空）
  - 存储提供者（本地、MinIO、WebDAV）
  - 汇率数据源
  - LLM 提供者

## 7. 数据流

### 7.1 用户认证流程
1. 用户提交凭据 → `/api/authorize.json`
2. 后端验证 → JWT 令牌生成
3. 令牌存储在前端（Pinia 存储）
4. 令牌包含在后续 API 请求中（通过 Axios 拦截器）
5. 中间件在每个请求上验证令牌

### 7.2 交易创建流程
1. 用户在前端填写表单
2. 前端验证输入
3. API 调用: `POST /api/v1/transactions/add.json`
4. 中间件: JWT 认证
5. API 处理器: 输入验证
6. 服务层: 业务逻辑
7. 数据存储: 数据库插入
8. 响应: 成功/错误 JSON
9. 前端: 更新 UI 状态

### 7.3 数据导入流程
1. 用户上传文件（CSV、OFX、Excel 等）
2. 前端: 文件解析
3. API: `POST /api/v1/transactions/parse_import.json`
4. 后端: 根据文件类型选择转换器
5. 转换器: 解析文件 → 交易数据
6. API: `POST /api/v1/transactions/import.json`
7. 后端: 批量插入并验证
8. 响应: 导入结果

## 8. 配置

### 配置文件 (`conf/ezbookkeeping.ini`)
- INI 格式配置
- 部分: `[global]`、`[server]`、`[database]`、`[mcp]` 等
- 支持环境变量
- 默认路径: `conf/ezbookkeeping.ini`

### 关键配置区域
- **服务器**: 协议、地址、端口、静态文件
- **数据库**: 类型、连接字符串、连接池
- **认证**: 内部认证、OAuth2、2FA
- **存储**: 本地路径、MinIO、WebDAV
- **LLM**: AI 提供者配置
- **汇率**: 数据源配置
- **邮件**: SMTP 配置

## 9. 安全特性

### 认证
- 基于 JWT 的认证
- 令牌刷新机制
- 双因素认证（TOTP）
- 恢复码
- OAuth2 支持（Google、GitHub 等）

### 授权
- 基于中间件的路由保护
- 用户特定数据隔离
- API 令牌管理

### 数据保护
- 密码哈希（bcrypt）
- 加密密钥
- HTTPS 支持
- 请求速率限制

## 10. 国际化

### 后端 (`pkg/locales/`)
- 22 种语言文件（Go 结构体）
- 特定于语言环境的格式化（日期、数字、货币）
- 日历支持（公历、农历、波斯历）

### 前端 (`src/locales/`)
- 基于 JSON 的翻译文件
- Vue i18n 集成
- 移动端 RTL（从右到左）支持

## 11. 构建和部署

### 构建过程
- **后端**: Go build → 单个二进制文件
- **前端**: Vite build → 静态文件
- **Docker**: 多阶段构建（后端 + 前端）

### 构建脚本
- `build.sh`: Linux/macOS 构建脚本
- `build.bat` / `build.ps1`: Windows 构建脚本
- `docker-bake.hcl`: Docker 构建配置

### 部署
- 单个二进制文件部署
- Docker 容器部署
- 静态文件由 Go 服务器提供
- 支持反向代理（Nginx 等）

## 12. 测试

### 后端测试
- Go test 框架
- 测试文件: `*_test.go`
- 示例: `gravatar_provider_test.go`、`cron_container_test.go`

### 前端测试
- Jest 配置 (`jest.config.ts`)
- TypeScript 测试文件
- 示例: `fiscal_year.ts` 测试

## 13. 扩展点

### 添加新功能
1. **新 API 端点**:
   - 在 `pkg/api/` 中添加处理器
   - 在 `cmd/webserver.go` 中添加路由
   - 在 `pkg/services/` 中添加服务（如需要）

2. **新数据导入格式**:
   - 在 `pkg/converters/` 中创建转换器
   - 在 `transaction_data_converters.go` 中注册

3. **新前端页面**:
   - 在 `src/views/desktop/` 或 `src/views/mobile/` 中创建视图
   - 在 `src/router/` 中添加路由
   - 如需要，在 `src/stores/` 中创建存储

4. **新存储提供者**:
   - 在 `pkg/storage/` 中实现接口
   - 在存储容器中注册

## 14. 需要理解的关键文件

### 后端入口点
- `ezbookkeeping.go`: 主入口点
- `cmd/webserver.go`: Web 服务器设置
- `cmd/initializer.go`: 系统初始化

### 前端入口点
- `src/index-main.ts`: 设备检测和路由
- `src/desktop-main.ts`: 桌面应用入口
- `src/mobile-main.ts`: 移动应用入口

### 核心配置
- `conf/ezbookkeeping.ini`: 主配置文件
- `pkg/settings/`: 配置加载和管理

### API 结构
- `pkg/api/`: 所有 API 处理器
- `pkg/services/`: 业务逻辑
- `pkg/datastore/`: 数据访问

### 前端结构
- `src/stores/`: 状态管理
- `src/lib/api.ts`: API 客户端
- `src/router/`: 路由配置

## 15. 开发工作流

### 本地开发
1. **后端**: `go run ezbookkeeping.go server run`
2. **前端**: `npm run serve`（Vite 开发服务器，端口 8081）
3. **数据库**: 默认 SQLite（或配置 MySQL/PostgreSQL）

### 代码组织原则
- **关注点分离**: 清晰的层边界
- **单一职责**: 每个模块都有明确的职责
- **依赖注入**: 共享资源的容器
- **基于接口的设计**: 易于交换实现

## 16. 深入理解的下一步

要深入了解特定领域：

1. **认证流程**: 研究 `pkg/api/authorizations.go` 和 `pkg/middlewares/`
2. **交易管理**: 研究 `pkg/api/transactions.go` 和 `pkg/services/`
3. **数据导入**: 研究 `pkg/converters/` 和导入处理器
4. **前端状态**: 研究 `src/stores/` 和组件使用
5. **数据库模型**: 研究 `pkg/models/` 了解数据结构
6. **API 设计**: 研究 `pkg/api/base.go` 了解通用模式

---

**注意**: 本文档提供高级概述。有关特定功能的详细实现，请参考源代码和内联注释。

