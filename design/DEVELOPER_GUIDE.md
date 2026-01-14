# ezBookkeeping Architecture Analysis

## 1. Project Overview

**ezBookkeeping** is a lightweight, self-hosted personal finance application with a user-friendly interface and powerful bookkeeping features. It's designed to be resource-efficient and highly scalable, capable of running on devices as small as a Raspberry Pi or scaling up to large cluster environments.

### Key Features
- Open source and self-hosted
- Multi-platform support (Desktop, Mobile, PWA)
- Multi-language and multi-currency support
- AI-powered receipt image recognition
- Support for various data import/export formats
- Two-factor authentication (2FA)
- OAuth2 authentication
- MCP (Model Context Protocol) support

## 2. Technology Stack

### Backend
- **Language**: Go 1.25+
- **Web Framework**: Gin (HTTP web framework)
- **ORM**: XORM (supports SQLite, MySQL, PostgreSQL)
- **CLI Framework**: urfave/cli/v3
- **Authentication**: JWT (golang-jwt/jwt/v5), OAuth2, 2FA (pquerna/otp)
- **Task Scheduling**: go-co-op/gocron/v2
- **Logging**: logrus
- **Configuration**: gopkg.in/ini.v1

### Frontend
- **Framework**: Vue 3 (Composition API)
- **Language**: TypeScript
- **Build Tool**: Vite 7.x
- **Desktop UI**: Vuetify 3
- **Mobile UI**: Framework7 9.x
- **State Management**: Pinia
- **Routing**: Vue Router 4
- **Charts**: ECharts (vue-echarts)
- **Maps**: Leaflet
- **Date Handling**: Moment.js + moment-timezone
- **PWA**: vite-plugin-pwa

### Database
- **Supported**: SQLite3, MySQL, PostgreSQL
- **ORM**: XORM with connection pooling

### Storage
- **Local Storage**: File system
- **Object Storage**: MinIO (S3-compatible)
- **WebDAV**: Supported

## 3. Overall Architecture

### Architecture Pattern
The project follows a **layered architecture** with clear separation of concerns:

```
┌─────────────────────────────────────────┐
│         Frontend (Vue 3)                │
│  ┌──────────┐  ┌──────────┐             │
│  │ Desktop  │  │  Mobile  │             │
│  │ (Vuetify)│  │(Framework7)│            │
│  └──────────┘  └──────────┘             │
└─────────────────────────────────────────┘
                  │ HTTP/REST API
                  │ JSON
┌─────────────────────────────────────────┐
│      Backend (Go + Gin)                 │
│  ┌──────────────────────────────────┐  │
│  │  API Layer (pkg/api)             │  │
│  │  - RESTful endpoints             │  │
│  │  - JSON-RPC (MCP)                 │  │
│  └──────────────────────────────────┘  │
│  ┌──────────────────────────────────┐  │
│  │  Service Layer (pkg/services)    │  │
│  │  - Business logic                │  │
│  └──────────────────────────────────┘  │
│  ┌──────────────────────────────────┐  │
│  │  Data Layer (pkg/datastore)      │  │
│  │  - XORM ORM                      │  │
│  │  - Database abstraction           │  │
│  └──────────────────────────────────┘  │
└─────────────────────────────────────────┘
                  │
┌─────────────────────────────────────────┐
│      Database (SQLite/MySQL/PostgreSQL) │
└─────────────────────────────────────────┘
```

### Request Flow
1. **Client Request** → Frontend (Vue Router)
2. **API Call** → Axios → Backend API endpoint
3. **Middleware** → Authentication, Request ID, Logging
4. **API Handler** → Business logic validation
5. **Service Layer** → Core business operations
6. **Data Store** → Database operations via XORM
7. **Response** → JSON response back to frontend

## 4. Directory Structure

### Root Level
```
ezbookkeeping/
├── cmd/              # CLI commands (server, database, user_data, etc.)
├── pkg/              # Backend packages (Go)
├── src/              # Frontend source code (Vue/TypeScript)
├── public/           # Static assets
├── conf/             # Configuration files
├── templates/        # Email and prompt templates
├── docker/           # Docker build scripts
└── testdata/         # Test data files
```

### Backend Structure (`pkg/`)

#### Core Modules
- **`api/`**: HTTP API handlers (RESTful endpoints)
  - `accounts.go`: Account management
  - `transactions.go`: Transaction CRUD operations
  - `authorizations.go`: Authentication endpoints
  - `users.go`: User management
  - `exchange_rates.go`: Currency exchange rates
  - `large_language_models.go`: AI/LLM integration
  - `model_context_protocols.go`: MCP protocol handlers

- **`services/`**: Business logic layer
  - Contains service implementations for various business operations
  - Acts as an intermediary between API handlers and data store

- **`datastore/`**: Database abstraction layer
  - `database.go`: Database connection and transaction management
  - `datastore.go`: Data store container with sharding support
  - Uses XORM for ORM operations

- **`core/`**: Core types and utilities
  - Context wrappers (CLI, Web, Cron)
  - Common types and handlers
  - Calendar, currency, datetime utilities

- **`models/`**: Database models (struct definitions for XORM)

- **`middlewares/`**: HTTP middleware
  - JWT authentication
  - Request logging
  - Recovery
  - Request ID generation

#### Supporting Modules
- **`auth/oauth2/`**: OAuth2 authentication providers
- **`avatars/`**: Avatar provider implementations (Gravatar, Internal, Null)
- **`converters/`**: Data import converters (CSV, OFX, QIF, Excel, etc.)
- **`cron/`**: Scheduled task management
- **`errs/`**: Error definitions
- **`exchangerates/`**: Exchange rate data sources
- **`llm/`**: Large language model providers
- **`locales/`**: Internationalization (22 languages)
- **`log/`**: Logging utilities
- **`mail/`**: Email sending
- **`mcp/`**: Model Context Protocol implementation
- **`settings/`**: Configuration management
- **`storage/`**: Object storage (local, MinIO, WebDAV)
- **`utils/`**: Utility functions
- **`uuid/`**: UUID generation
- **`validators/`**: Input validation

### Frontend Structure (`src/`)

#### Entry Points
- **`index-main.ts`**: Entry point that detects device type and redirects
- **`desktop-main.ts`**: Desktop application entry
- **`mobile-main.ts`**: Mobile application entry

#### Core Directories
- **`components/`**: Vue components
  - `base/`: Base components (shared logic)
  - `common/`: Common components (DateTimePicker, MapView, etc.)
  - `desktop/`: Desktop-specific components
  - `mobile/`: Mobile-specific components

- **`views/`**: Page views
  - `base/`: Base view classes
  - `desktop/`: Desktop pages
  - `mobile/`: Mobile pages

- **`stores/`**: Pinia state management
  - `account.ts`: Account state
  - `transaction.ts`: Transaction state
  - `user.ts`: User state
  - `token.ts`: Authentication tokens
  - `exchangeRates.ts`: Exchange rates
  - etc.

- **`router/`**: Vue Router configuration
  - `desktop.ts`: Desktop routes
  - `mobile.ts`: Mobile routes

- **`lib/`**: Utility libraries
  - `api.ts`: API client (Axios wrapper)
  - `services.ts`: Service layer
  - `map/`: Map providers (Leaflet, Google Maps, Baidu, Amap)
  - `ui/`: UI utilities

- **`core/`**: Core utilities
  - `account.ts`, `transaction.ts`, `category.ts`: Domain models
  - `datetime.ts`, `currency.ts`: Formatting utilities
  - `statistics.ts`: Statistical calculations

- **`consts/`**: Constants
- **`models/`**: TypeScript type definitions
- **`locales/`**: Frontend i18n files (JSON)
- **`styles/`**: SCSS stylesheets

## 5. Core Components

### 5.1 Backend Components

#### Command Structure (`cmd/`)
- **`webserver.go`**: Main web server initialization and routing
- **`database.go`**: Database management commands
- **`user_data.go`**: User data management CLI
- **`cron_jobs.go`**: Cron job management
- **`security.go`**: Security utilities
- **`initializer.go`**: System initialization (config, database, storage, etc.)

#### API Layer (`pkg/api/`)
Each API module follows a similar pattern:
- Handler functions that receive `core.Context`
- Input validation
- Service layer calls
- Response formatting

Example structure:
```go
var Accounts = &accountsApi{}

type accountsApi struct{}

func (a *accountsApi) AccountListHandler(c core.Context) (interface{}, error) {
    // Validation, service call, response
}
```

#### Data Store (`pkg/datastore/`)
- **Container Pattern**: `DataStoreContainer` holds multiple stores
  - `UserStore`: User-related data
  - `TokenStore`: Authentication tokens
  - `UserDataStore`: User-specific data (accounts, transactions, etc.)
- **Sharding Support**: Designed for horizontal scaling (currently uses single database)
- **Transaction Management**: `DoTransaction` method for atomic operations

#### Service Layer (`pkg/services/`)
- Business logic implementation
- Data validation and transformation
- Calls to data store layer

### 5.2 Frontend Components

#### Application Structure
- **Desktop App**: `DesktopApp.vue` - Main desktop application component
- **Mobile App**: `MobileApp.vue` - Main mobile application component
- Both use Vue 3 Composition API

#### State Management (Pinia)
- Stores are organized by domain (account, transaction, user, etc.)
- Stores handle API calls and state updates
- Reactive state for UI updates

#### API Client (`lib/api.ts`)
- Axios-based HTTP client
- Request/response interceptors
- Error handling
- Token management

#### Routing
- Separate routers for desktop and mobile
- Route guards for authentication
- Lazy loading support

## 6. Key Design Patterns

### 6.1 Dependency Injection
- Container pattern for services (e.g., `DataStoreContainer`, `AvatarProviderContainer`)
- Singleton instances for shared resources

### 6.2 Middleware Pattern
- Chain of responsibility for HTTP middleware
- Authentication, logging, error recovery

### 6.3 Repository Pattern
- Data store abstraction layer
- Separation of data access from business logic

### 6.4 Factory Pattern
- Provider factories (avatar, storage, LLM, exchange rates)
- Converter factories for data import

### 6.5 Strategy Pattern
- Multiple implementations for same interface
  - Avatar providers (Gravatar, Internal, Null)
  - Storage providers (Local, MinIO, WebDAV)
  - Exchange rate sources
  - LLM providers

## 7. Data Flow

### 7.1 User Authentication Flow
1. User submits credentials → `/api/authorize.json`
2. Backend validates → JWT token generation
3. Token stored in frontend (Pinia store)
4. Token included in subsequent API requests (via Axios interceptor)
5. Middleware validates token on each request

### 7.2 Transaction Creation Flow
1. User fills form in frontend
2. Frontend validates input
3. API call: `POST /api/v1/transactions/add.json`
4. Middleware: JWT authentication
5. API handler: Input validation
6. Service layer: Business logic
7. Data store: Database insert
8. Response: Success/error JSON
9. Frontend: Update UI state

### 7.3 Data Import Flow
1. User uploads file (CSV, OFX, Excel, etc.)
2. Frontend: File parsing
3. API: `POST /api/v1/transactions/parse_import.json`
4. Backend: Converter selection based on file type
5. Converter: Parse file → Transaction data
6. API: `POST /api/v1/transactions/import.json`
7. Backend: Batch insert with validation
8. Response: Import results

## 8. Configuration

### Configuration File (`conf/ezbookkeeping.ini`)
- INI format configuration
- Sections: `[global]`, `[server]`, `[database]`, `[mcp]`, etc.
- Supports environment variables
- Default path: `conf/ezbookkeeping.ini`

### Key Configuration Areas
- **Server**: Protocol, address, port, static files
- **Database**: Type, connection string, pooling
- **Authentication**: Internal auth, OAuth2, 2FA
- **Storage**: Local path, MinIO, WebDAV
- **LLM**: AI provider configuration
- **Exchange Rates**: Data source configuration
- **Email**: SMTP configuration

## 9. Security Features

### Authentication
- JWT-based authentication
- Token refresh mechanism
- Two-factor authentication (TOTP)
- Recovery codes
- OAuth2 support (Google, GitHub, etc.)

### Authorization
- Middleware-based route protection
- User-specific data isolation
- API token management

### Data Protection
- Password hashing (bcrypt)
- Secret key for encryption
- HTTPS support
- Request rate limiting

## 10. Internationalization

### Backend (`pkg/locales/`)
- 22 language files (Go structs)
- Locale-specific formatting (dates, numbers, currencies)
- Calendar support (Gregorian, Chinese, Persian)

### Frontend (`src/locales/`)
- JSON-based translation files
- Vue i18n integration
- RTL (Right-to-Left) support for mobile

## 11. Build and Deployment

### Build Process
- **Backend**: Go build → Single binary
- **Frontend**: Vite build → Static files
- **Docker**: Multi-stage build (backend + frontend)

### Build Scripts
- `build.sh`: Linux/macOS build script
- `build.bat` / `build.ps1`: Windows build scripts
- `docker-bake.hcl`: Docker build configuration

### Deployment
- Single binary deployment
- Docker container deployment
- Static files served by Go server
- Supports reverse proxy (Nginx, etc.)

## 12. Testing

### Backend Testing
- Go test framework
- Test files: `*_test.go`
- Examples: `gravatar_provider_test.go`, `cron_container_test.go`

### Frontend Testing
- Jest configuration (`jest.config.ts`)
- TypeScript test files
- Example: `fiscal_year.ts` tests

## 13. Extension Points

### Adding New Features
1. **New API Endpoint**:
   - Add handler in `pkg/api/`
   - Add route in `cmd/webserver.go`
   - Add service in `pkg/services/` (if needed)

2. **New Data Import Format**:
   - Create converter in `pkg/converters/`
   - Register in `transaction_data_converters.go`

3. **New Frontend Page**:
   - Create view in `src/views/desktop/` or `src/views/mobile/`
   - Add route in `src/router/`
   - Create store if needed in `src/stores/`

4. **New Storage Provider**:
   - Implement interface in `pkg/storage/`
   - Register in storage container

## 14. Key Files to Understand

### Backend Entry Points
- `ezbookkeeping.go`: Main entry point
- `cmd/webserver.go`: Web server setup
- `cmd/initializer.go`: System initialization

### Frontend Entry Points
- `src/index-main.ts`: Device detection and routing
- `src/desktop-main.ts`: Desktop app entry
- `src/mobile-main.ts`: Mobile app entry

### Core Configuration
- `conf/ezbookkeeping.ini`: Main configuration file
- `pkg/settings/`: Configuration loading and management

### API Structure
- `pkg/api/`: All API handlers
- `pkg/services/`: Business logic
- `pkg/datastore/`: Data access

### Frontend Structure
- `src/stores/`: State management
- `src/lib/api.ts`: API client
- `src/router/`: Routing configuration

## 15. Development Workflow

### Local Development
1. **Backend**: `go run ezbookkeeping.go server run`
2. **Frontend**: `npm run serve` (Vite dev server on port 8081)
3. **Database**: SQLite by default (or configure MySQL/PostgreSQL)

### Code Organization Principles
- **Separation of Concerns**: Clear layer boundaries
- **Single Responsibility**: Each module has a focused purpose
- **Dependency Injection**: Containers for shared resources
- **Interface-based Design**: Easy to swap implementations

## 16. Next Steps for Understanding

To dive deeper into specific areas:

1. **Authentication Flow**: Study `pkg/api/authorizations.go` and `pkg/middlewares/`
2. **Transaction Management**: Study `pkg/api/transactions.go` and `pkg/services/`
3. **Data Import**: Study `pkg/converters/` and import handlers
4. **Frontend State**: Study `src/stores/` and component usage
5. **Database Models**: Study `pkg/models/` for data structure
6. **API Design**: Study `pkg/api/base.go` for common patterns

## 17. Coding Style and Conventions

### 17.1 Backend (Go) Coding Style

#### Naming Conventions
- **Packages**: Lowercase, single word (e.g., `api`, `utils`, `errs`)
- **Types**: PascalCase (e.g., `WebContext`, `Error`, `AccountService`)
- **Functions**: PascalCase for exported, camelCase for unexported (e.g., `AccountListHandler`, `parseFromUnixTime`)
- **Variables**: camelCase (e.g., `userAccount`, `transactionTime`)
- **Constants**: PascalCase or UPPER_SNAKE_CASE (e.g., `CATEGORY_SYSTEM`, `ErrApiNotFound`)

#### API Handler Pattern
All API handlers follow a consistent pattern:
```go
// 1. Define API struct with embedded base types
type AccountsApi struct {
    ApiUsingConfig
    ApiUsingDuplicateChecker
    accounts *services.AccountService
}

// 2. Initialize singleton instance
var Accounts = &AccountsApi{
    ApiUsingConfig: ApiUsingConfig{
        container: settings.Container,
    },
    // ... other initializations
}

// 3. Handler function signature
func (a *AccountsApi) AccountListHandler(c *core.WebContext) (any, *errs.Error) {
    // Input validation
    var req models.AccountListRequest
    err := c.ShouldBindQuery(&req)
    if err != nil {
        log.Warnf(c, "[accounts.AccountListHandler] parse request failed, because %s", err.Error())
        return nil, errs.NewIncompleteOrIncorrectSubmissionError(err)
    }
    
    // Business logic
    // ...
    
    // Return result or error
    return result, nil
}
```

#### Error Handling
- Always use `*errs.Error` for application errors
- Use `errs.NewIncompleteOrIncorrectSubmissionError()` for validation errors
- Use `errs.Or()` to convert unknown errors to known errors
- Log errors with context using `log.Errorf()`, `log.Warnf()`, etc.
- Include request context in log messages: `log.Errorf(c, "[module.function] message")`

#### Logging
- Use structured logging with `log.Debugf()`, `log.Infof()`, `log.Warnf()`, `log.Errorf()`
- Always include context: `log.Errorf(c, "[module.function] message")`
- Use appropriate log levels:
  - `Debugf`: Detailed debugging information
  - `Infof`: General informational messages
  - `Warnf`: Warning messages (non-critical issues)
  - `Errorf`: Error messages (critical issues)

#### Context Usage
- Always pass `core.Context` (or `*core.WebContext`) to functions that need logging or request context
- Use `c.GetCurrentUid()` to get current user ID
- Use `c.GetContextId()` to get request ID for logging

### 17.2 Frontend (TypeScript/Vue) Coding Style

#### Naming Conventions
- **Files**: kebab-case (e.g., `account-list.vue`, `user-profile-page-base.ts`)
- **Components**: PascalCase (e.g., `AccountList`, `UserProfilePage`)
- **Functions**: camelCase (e.g., `getAccountList`, `formatCurrency`)
- **Types/Interfaces**: PascalCase (e.g., `AccountInfo`, `TransactionRequest`)
- **Constants**: UPPER_SNAKE_CASE (e.g., `DEFAULT_CURRENCY`, `MAX_AMOUNT`)

#### Vue Component Pattern
- Use Composition API with `<script setup lang="ts">`
- Use TypeScript for type safety
- Extract reusable logic to composables (e.g., `useAccountListBase()`)
- Use Pinia stores for state management

#### Type Safety
- Always define types for function parameters and return values
- Use type guards (e.g., `isString()`, `isNumber()`) from `@/lib/common.ts`
- Avoid `any` type; use `unknown` when type is truly unknown

### 17.3 Code Organization
- **One file, one primary responsibility**: Each file should have a clear, single purpose
- **Group related functionality**: Keep related functions and types together
- **Use base classes/interfaces**: Share common functionality through base types
- **Avoid code duplication**: Extract common patterns into utility functions

## 18. Common Utility Classes

### 18.1 Backend Utilities (`pkg/utils/`)

#### String Utilities (`strings.go`)
- `SubString()`: Extract substring with rune support
- `ContainsAnyString()`: Check if string contains any substring
- `GetFirstLowerCharString()`: Convert first character to lowercase
- `GetRandomString()`: Generate random string
- `GetRandomNumberOrLetter()`: Generate random alphanumeric string
- `MD5Encode()`, `MD5EncodeToString()`: MD5 hashing
- `AESGCMEncrypt()`, `AESGCMDecrypt()`: AES-GCM encryption/decryption
- `EncodePassword()`: Password encoding with PBKDF2
- `EncryptSecret()`, `DecryptSecret()`: Secret encryption/decryption

#### Number Utilities (`converter.go`)
- `IntToString()`, `StringToInt()`: Integer conversion
- `Int64ToString()`, `StringToInt64()`: Int64 conversion
- `Float64ToString()`, `StringToFloat64()`: Float64 conversion
- `FormatAmount()`, `ParseAmount()`: Amount formatting (cent-based)
- `Int64ArrayToStringArray()`: Array conversion

#### DateTime Utilities (`datetimes.go`)
- `FormatUnixTimeToLongDate()`: Format unix time to date string
- `FormatUnixTimeToLongDateTime()`: Format unix time to datetime string
- `ParseFromLongDateTimeInFixedUtcOffset()`: Parse datetime string
- `GetTimezoneOffsetMinutes()`: Get timezone offset
- `GetTransactionTimeRangeByYearMonth()`: Get transaction time range
- `GetStartOfDay()`: Get start of day time

#### HTTP Utilities (`http.go`)
- `NewHttpClient()`: Create HTTP client with proxy and TLS settings
- `SetProxyUrl()`: Configure proxy for HTTP transport

#### API Utilities (`api.go`)
- `PrintJsonSuccessResult()`: Write JSON success response
- `PrintJsonErrorResult()`: Write JSON error response
- `PrintJSONRPCSuccessResult()`: Write JSON-RPC success response
- `GetDisplayErrorMessage()`: Get user-friendly error message
- `GetJsonErrorResult()`: Format error response

#### Validation Utilities (`validators.go`)
- `IsValidUsername()`: Validate username format
- `IsValidEmail()`: Validate email format
- `IsValidNickName()`: Validate nickname
- `IsValidHexRGBColor()`: Validate hex color
- `IsValidLongDateTimeFormat()`: Validate datetime format
- `IsValidLongDateFormat()`: Validate date format

#### Slice Utilities (`slices.go`)
- `Int64SliceEquals()`: Compare int64 slices
- `Int64SliceMinus()`: Subtract slices
- `ToUniqueInt64Slice()`: Remove duplicates
- `Int64Sort()`: Sort int64 slice
- `ToSet()`: Convert slice to map

#### I/O Utilities (`io.go`)
- `GetImageContentType()`: Get content type for image extension
- `ListFileNamesWithPrefixAndSuffix()`: List files matching pattern
- `IsExists()`: Check if file/directory exists
- `WriteFile()`: Write file content
- `GetFileNameWithoutExtension()`: Extract filename without extension
- `GetFileNameExtension()`: Get file extension

#### Object Utilities (`object.go`)
- `Clone()`: Deep clone object using gob encoding
- `PrintObjectFields()`: Print all fields of an object

### 18.2 Frontend Utilities (`src/lib/`)

#### Common Utilities (`common.ts`)
- `isDefined()`: Check if value is not null/undefined
- `isObject()`, `isArray()`, `isString()`, `isNumber()`, `isBoolean()`: Type guards
- `isEquals()`: Deep equality check
- `limitText()`: Limit text length with ellipsis
- `base64encode()`, `base64decode()`: Base64 encoding/decoding
- `getItemByKeyValue()`: Find item in array/object by key-value
- `arrayContainsFieldValue()`: Check if array contains value

#### DateTime Utilities (`datetime.ts`)
- `formatCurrentTime()`: Format current time
- `formatDateTime()`: Format datetime with timezone
- `parseDateTimeFromUnixTime()`: Parse unix time to datetime
- `getTimezoneOffset()`: Get timezone offset string
- `getBrowserTimezoneName()`: Get browser timezone
- `getFiscalYearTimeRangeFromUnixTime()`: Get fiscal year range

#### Currency Utilities (`currency.ts`)
- `getCurrencyFraction()`: Get currency decimal places
- `appendCurrencySymbol()`: Append currency symbol to amount
- `getAmountPrependAndAppendCurrencySymbol()`: Get currency symbol position

#### File Utilities (`file.ts`)
- `getFileExtension()`: Extract file extension
- `isFileExtensionSupported()`: Check if extension is supported
- `detectFileEncoding()`: Detect file encoding

#### Settings Utilities (`settings.ts`)
- `getApplicationSettings()`: Get application settings
- `updateApplicationSettingsValue()`: Update setting value
- `getTheme()`, `getTimeZone()`: Get specific settings

#### UI Utilities (`ui/`)
- `common.ts`: Common UI utilities (scroll, theme, clipboard, etc.)
- `desktop.ts`: Desktop-specific UI utilities
- `mobile.ts`: Mobile-specific UI utilities (Framework7)

## 19. Developer Guidelines and Best Practices

### 19.1 Error Handling Guidelines

#### Backend
1. **Always return `*errs.Error`**: Never return standard Go errors from API handlers
2. **Use appropriate error codes**: Use predefined errors from `pkg/errs/` when possible
3. **Wrap validation errors**: Use `errs.NewIncompleteOrIncorrectSubmissionError()` for input validation failures
4. **Log before returning**: Always log errors with context before returning
5. **Use `errs.Or()`**: Convert unknown errors to known errors when appropriate

Example:
```go
result, err := a.service.DoSomething(c, param)
if err != nil {
    log.Errorf(c, "[module.function] failed to do something, because %s", err.Error())
    return nil, errs.Or(err, errs.ErrOperationFailed)
}
```

#### Frontend
1. **Handle API errors**: Always check for errors in API responses
2. **Show user-friendly messages**: Display localized error messages to users
3. **Log errors for debugging**: Use console logging for debugging (only in development)

### 19.2 Logging Guidelines

1. **Include context**: Always include request context (`c`) in log calls
2. **Use appropriate levels**: 
   - Debug: Detailed debugging info (only in debug mode)
   - Info: Normal operation messages
   - Warn: Warning conditions (non-critical)
   - Error: Error conditions (critical)
3. **Format consistently**: Use format `[module.function] message` for log messages
4. **Include relevant data**: Log relevant parameters and error details

Example:
```go
log.Infof(c, "[accounts.AccountListHandler] getting accounts for user uid:%d", uid)
log.Warnf(c, "[accounts.AccountListHandler] account not found, accountId:%d", accountId)
log.Errorf(c, "[accounts.AccountListHandler] failed to get account, because %s", err.Error())
```

### 19.3 API Handler Guidelines

1. **Follow the pattern**: Use the standard API handler pattern (see section 17.1)
2. **Validate input**: Always validate input using struct tags and `ShouldBind*()` methods
3. **Check permissions**: Verify user has permission to access the resource
4. **Use service layer**: Don't access datastore directly from API handlers
5. **Return appropriate responses**: Use `utils.PrintJsonSuccessResult()` or `utils.PrintJsonErrorResult()`

### 19.4 Database Access Guidelines

1. **Use datastore layer**: Never access XORM directly; use datastore methods
2. **Handle transactions**: Use `DoTransaction()` for operations requiring atomicity
3. **Check user ownership**: Always verify user owns the resource before operations
4. **Use prepared statements**: XORM handles this automatically, but be aware of SQL injection risks

### 19.5 Frontend Development Guidelines

1. **Use TypeScript**: Always use TypeScript for type safety
2. **Use composables**: Extract reusable logic to composables
3. **Use Pinia stores**: Manage state through Pinia stores, not component data
4. **Handle loading states**: Always show loading indicators for async operations
5. **Handle errors gracefully**: Show user-friendly error messages
6. **Use i18n**: Always use i18n for user-facing text

### 19.6 Testing Guidelines

1. **Write tests**: Write unit tests for utility functions
2. **Test edge cases**: Test boundary conditions and error cases
3. **Use test data**: Use `testdata/` directory for test files
4. **Naming**: Test files should be named `*_test.go` (Go) or `*.test.ts` (TypeScript)

### 19.7 Code Review Checklist

- [ ] Follows naming conventions
- [ ] Includes proper error handling
- [ ] Includes logging with context
- [ ] Validates all inputs
- [ ] Checks user permissions
- [ ] Uses appropriate utility functions
- [ ] No code duplication
- [ ] Includes comments for complex logic
- [ ] TypeScript types are properly defined
- [ ] No hardcoded strings (use i18n)

### 19.8 Common Pitfalls to Avoid

1. **Don't ignore errors**: Always handle errors properly
2. **Don't log sensitive data**: Never log passwords, tokens, or sensitive user data
3. **Don't access datastore directly**: Always use service layer
4. **Don't use `any` type**: Use proper TypeScript types
5. **Don't hardcode strings**: Use i18n for all user-facing text
6. **Don't forget context**: Always pass context to functions that need it
7. **Don't skip validation**: Always validate user input
8. **Don't forget timezone**: Always consider timezone when working with dates

---

**Note**: This document provides a high-level overview. For detailed implementation of specific features, refer to the source code and inline comments.

