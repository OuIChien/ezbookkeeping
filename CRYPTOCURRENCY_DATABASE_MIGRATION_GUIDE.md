# Cryptocurrency Support Database Migration Guide

## Overview

This guide explains how to migrate the database schema to support cryptocurrency symbols in addition to traditional ISO 4217 currency codes. The migration extends the `Currency` field length from `VARCHAR(3)` to `VARCHAR(10)` to accommodate cryptocurrency symbols like `BTC`, `ETH`, `DOGE`, `MATIC`, `USDT`, etc.

## Affected Database Tables

The following tables require schema updates:

1. **`account`** table - `currency` column
2. **`user`** table - `default_currency` column  
3. **`user_custom_exchange_rate`** table - `currency` column

## Migration Methods

### Method 1: Automatic Migration (Recommended)

If your system has `auto_update_database = true` in the configuration file (default setting), the database schema will be automatically updated when you start the server.

**Steps:**

1. Ensure `auto_update_database = true` is set in `conf/ezbookkeeping.ini`:
   ```ini
   [database]
   auto_update_database = true
   ```

2. Start the server normally:
   ```bash
   ./ezbookkeeping server run
   ```

3. Check the boot logs for confirmation:
   ```
   [database.updateAllDatabaseTablesStructure] account table maintained successfully
   [database.updateAllDatabaseTablesStructure] user table maintained successfully
   [database.updateAllDatabaseTablesStructure] user custom exchange rate table maintained successfully
   ```

**Note:** Automatic migration works well for SQLite and PostgreSQL. For MySQL, you may need to use Method 2 if automatic migration fails.

### Method 2: Manual Database Update Command

Use the built-in database update command to manually trigger schema synchronization.

**Steps:**

1. Run the database update command:
   ```bash
   ./ezbookkeeping database update
   ```

2. Verify the output shows successful table updates:
   ```
   [database.updateAllDatabaseTablesStructure] account table maintained successfully
   [database.updateAllDatabaseTablesStructure] user table maintained successfully
   [database.updateAllDatabaseTablesStructure] user custom exchange rate table maintained successfully
   ```

### Method 3: Manual SQL Migration (If Automatic Migration Fails)

If automatic migration fails (especially with MySQL), you can manually execute SQL statements to update the schema.

#### For MySQL/MariaDB

```sql
-- Update account table
ALTER TABLE account MODIFY COLUMN currency VARCHAR(10) NOT NULL;

-- Update user table
ALTER TABLE user MODIFY COLUMN default_currency VARCHAR(10) NOT NULL;

-- Update user_custom_exchange_rate table
ALTER TABLE user_custom_exchange_rate MODIFY COLUMN currency VARCHAR(10) NOT NULL;
```

#### For PostgreSQL

```sql
-- Update account table
ALTER TABLE account ALTER COLUMN currency TYPE VARCHAR(10);

-- Update user table
ALTER TABLE user ALTER COLUMN default_currency TYPE VARCHAR(10);

-- Update user_custom_exchange_rate table
ALTER TABLE user_custom_exchange_rate ALTER COLUMN currency TYPE VARCHAR(10);
```

#### For SQLite

SQLite does not support `ALTER COLUMN` directly. You need to recreate the table:

```sql
-- 1. Create new account table
CREATE TABLE account_new (
    account_id INTEGER PRIMARY KEY,
    uid INTEGER NOT NULL,
    deleted INTEGER NOT NULL,
    category INTEGER NOT NULL,
    type INTEGER NOT NULL,
    parent_account_id INTEGER NOT NULL,
    name VARCHAR(64) NOT NULL,
    display_order INTEGER NOT NULL,
    icon INTEGER NOT NULL,
    color VARCHAR(6) NOT NULL,
    currency VARCHAR(10) NOT NULL,
    balance INTEGER NOT NULL,
    comment VARCHAR(255) NOT NULL,
    extend BLOB,
    hidden INTEGER NOT NULL,
    created_unix_time INTEGER,
    updated_unix_time INTEGER,
    deleted_unix_time INTEGER
);

-- 2. Copy data from old table
INSERT INTO account_new SELECT * FROM account;

-- 3. Drop old table
DROP TABLE account;

-- 4. Rename new table
ALTER TABLE account_new RENAME TO account;

-- 5. Recreate indexes (adjust based on your actual indexes)
CREATE INDEX IDX_account_uid_deleted_parent_account_id_order ON account(uid, deleted, parent_account_id, display_order);

-- Repeat similar steps for user and user_custom_exchange_rate tables
```

**Important:** For SQLite, it's recommended to use Method 1 or Method 2, as xorm will handle the table recreation automatically.

## Verification Steps

After migration, verify the changes:

### 1. Check Column Definitions

#### MySQL/MariaDB
```sql
SHOW COLUMNS FROM account LIKE 'currency';
SHOW COLUMNS FROM user LIKE 'default_currency';
SHOW COLUMNS FROM user_custom_exchange_rate LIKE 'currency';
```

Expected output should show `VARCHAR(10)` instead of `VARCHAR(3)`.

#### PostgreSQL
```sql
SELECT column_name, data_type, character_maximum_length 
FROM information_schema.columns 
WHERE table_name = 'account' AND column_name = 'currency';

SELECT column_name, data_type, character_maximum_length 
FROM information_schema.columns 
WHERE table_name = 'user' AND column_name = 'default_currency';

SELECT column_name, data_type, character_maximum_length 
FROM information_schema.columns 
WHERE table_name = 'user_custom_exchange_rate' AND column_name = 'currency';
```

#### SQLite
```sql
PRAGMA table_info(account);
PRAGMA table_info(user);
PRAGMA table_info(user_custom_exchange_rate);
```

### 2. Test Cryptocurrency Symbol Storage

Create a test account with a cryptocurrency symbol to verify it works:

```sql
-- This should work after migration (example for testing)
-- Note: Use actual API or UI to create accounts, don't insert directly
```

### 3. Check Application Logs

After starting the application, check for any errors related to currency validation or database operations.

## Rollback Plan

If you need to rollback the migration (not recommended after using cryptocurrency features):

### MySQL/MariaDB
```sql
ALTER TABLE account MODIFY COLUMN currency VARCHAR(3) NOT NULL;
ALTER TABLE user MODIFY COLUMN default_currency VARCHAR(3) NOT NULL;
ALTER TABLE user_custom_exchange_rate MODIFY COLUMN currency VARCHAR(3) NOT NULL;
```

**Warning:** This will truncate any cryptocurrency symbols longer than 3 characters. Ensure you don't have any cryptocurrency data before rolling back.

### PostgreSQL
```sql
ALTER TABLE account ALTER COLUMN currency TYPE VARCHAR(3);
ALTER TABLE user ALTER COLUMN default_currency TYPE VARCHAR(3);
ALTER TABLE user_custom_exchange_rate ALTER COLUMN currency TYPE VARCHAR(3);
```

**Warning:** Same as above - ensure no cryptocurrency data exists.

## Pre-Migration Checklist

- [ ] Backup your database
- [ ] Review current currency data to ensure no conflicts
- [ ] Check database type (MySQL/PostgreSQL/SQLite)
- [ ] Verify `auto_update_database` setting in configuration
- [ ] Ensure sufficient disk space for migration
- [ ] Plan for minimal downtime if using production database

## Post-Migration Checklist

- [ ] Verify column definitions are updated correctly
- [ ] Test creating accounts with cryptocurrency symbols (BTC, ETH, etc.)
- [ ] Test updating existing accounts with cryptocurrency symbols
- [ ] Verify exchange rate functionality works with cryptocurrencies
- [ ] Check application logs for any errors
- [ ] Test cryptocurrency price fetching (if configured)

## Troubleshooting

### Issue: Automatic Migration Fails

**Symptoms:** Server fails to start or shows database errors.

**Solutions:**
1. Check database permissions - ensure the database user has `ALTER TABLE` privileges
2. Check for table locks - ensure no other processes are using the database
3. Use Method 3 (Manual SQL Migration) instead
4. Check database logs for specific error messages

### Issue: VARCHAR Length Not Updated (MySQL)

**Symptoms:** Column still shows `VARCHAR(3)` after automatic migration.

**Solutions:**
1. MySQL sometimes requires explicit `ALTER TABLE` statements
2. Use Method 3 (Manual SQL Migration) with MySQL-specific syntax
3. Verify MySQL version supports the operation

### Issue: Data Truncation Warnings

**Symptoms:** Warnings about data being truncated.

**Solutions:**
1. This should not happen if migration is done before using cryptocurrency features
2. If you have existing data, check for any values longer than 3 characters
3. Ensure all currency codes are valid before migration

### Issue: Index or Constraint Errors

**Symptoms:** Errors related to indexes or foreign keys during migration.

**Solutions:**
1. Some databases may need indexes to be dropped and recreated
2. Check for foreign key constraints that reference the currency column
3. Temporarily disable foreign key checks if necessary (MySQL: `SET FOREIGN_KEY_CHECKS=0;`)

## Database-Specific Notes

### MySQL/MariaDB

- May require explicit `ALTER TABLE` statements
- Check `sql_mode` settings - strict mode may cause issues
- Ensure `innodb_file_format` supports the operation
- Consider using `pt-online-schema-change` for large tables in production

### PostgreSQL

- Generally handles automatic migration well
- May require `VACUUM` after migration for large tables
- Check for any custom constraints or triggers

### SQLite

- Automatic migration works best
- Table recreation may be required (handled automatically by xorm)
- Ensure sufficient disk space for table recreation
- Backup is especially important for SQLite

## Configuration

After successful migration, you can configure cryptocurrency support in `conf/ezbookkeeping.ini`:

```ini
[cryptocurrency]
# Cryptocurrency price data source
# Options: "coingecko", "coinmarketcap", "binance"
data_source = coingecko

# Comma-separated list of cryptocurrency symbols to fetch
# Examples: BTC,ETH,BNB,SOL,ADA,XRP,DOT,DOGE,MATIC,USDT
cryptocurrencies = BTC,ETH,BNB,SOL,ADA

# Request timeout in milliseconds (default: 10000)
request_timeout = 10000

# Proxy setting (same as exchange_rates)
proxy = system

# Skip TLS verification
skip_tls_verify = false

# API key (optional, required for some data sources like CoinMarketCap)
api_key = 
```

## Support

If you encounter issues during migration:

1. Check the application logs for detailed error messages
2. Review database-specific documentation for your database type
3. Ensure you have proper database backups before attempting migration
4. Test the migration on a development/staging environment first

## Summary

The migration from `VARCHAR(3)` to `VARCHAR(10)` is a straightforward schema change that:

- Extends currency field length to support cryptocurrency symbols
- Maintains backward compatibility with existing ISO 4217 currency codes
- Can be performed automatically in most cases
- Requires manual SQL only if automatic migration fails (primarily MySQL)

Always backup your database before performing any migration, and test in a non-production environment first.

