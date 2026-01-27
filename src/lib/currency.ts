import { CurrencyDisplaySymbol, CurrencyDisplayLocation, type CurrencyPrependAndAppendText, CurrencyDisplayType } from '@/core/currency.ts';
import { AccountAssetType } from '@/core/account.ts';
import { ALL_CURRENCIES, ALL_CRYPTOCURRENCIES, DEFAULT_CURRENCY_SYMBOL } from '@/consts/currency.ts';
import { isNumber } from './common.ts';

export function getCurrencyFraction(currencyCode?: string): number | undefined {
    if (!currencyCode) {
        return undefined;
    }

    const currencyInfo = ALL_CURRENCIES[currencyCode] || ALL_CRYPTOCURRENCIES[currencyCode];
    if (!currencyInfo?.fraction) {
        return undefined;
    }

    // For cryptocurrencies with fraction > 8, limit to 8 to avoid int64 overflow
    return Math.min(currencyInfo.fraction, 8);
}

export function inferAssetTypeFromCurrencyCode(currencyCode: string): number | undefined {
    if (!currencyCode) {
        return undefined;
    }

    // Check if it's a cryptocurrency first
    if (ALL_CRYPTOCURRENCIES[currencyCode]) {
        return AccountAssetType.Crypto.type;
    }

    // Check if it's a fiat currency
    if (ALL_CURRENCIES[currencyCode]) {
        return AccountAssetType.Fiat.type;
    }

    // If not found in either, it might be a stock symbol or unknown
    // Return Stock as fallback if it looks like a symbol (usually uppercase/numbers)
    // For now, return undefined to let the caller decide or return Stock
    return AccountAssetType.Stock.type;
}

export function getExchangedAmount(
    amount: number,
    fromCurrency: string,
    toCurrency: string,
    exchangeRatesStore: any,
    cryptocurrencyPricesStore: any,
    stockPricesStore: any
): number | null {
    if (fromCurrency === toCurrency) {
        return amount;
    }

    const fromAssetType = inferAssetTypeFromCurrencyCode(fromCurrency);
    const toAssetType = inferAssetTypeFromCurrencyCode(toCurrency);

    if (fromAssetType === AccountAssetType.Fiat.type && toAssetType === AccountAssetType.Fiat.type) {
        return exchangeRatesStore.getExchangedAmount(amount, fromCurrency, toCurrency);
    } else if (fromAssetType === AccountAssetType.Crypto.type && toAssetType === AccountAssetType.Fiat.type) {
        const price = cryptocurrencyPricesStore.getCryptocurrencyPriceInFiat(fromCurrency, toCurrency);
        if (isNumber(price)) {
            const fromFraction = getCurrencyFraction(fromCurrency) || 0;
            const toFraction = getCurrencyFraction(toCurrency) || 0;
            return amount * price * Math.pow(10, toFraction - fromFraction);
        }
    } else if (fromAssetType === AccountAssetType.Stock.type && toAssetType === AccountAssetType.Fiat.type) {
        const price = stockPricesStore.getStockPriceInFiat(fromCurrency, toCurrency, exchangeRatesStore);
        if (isNumber(price)) {
            const fromFraction = getCurrencyFraction(fromCurrency) || 0;
            const toFraction = getCurrencyFraction(toCurrency) || 0;
            return amount * price * Math.pow(10, toFraction - fromFraction);
        }
    } else if (fromAssetType === AccountAssetType.Fiat.type && toAssetType === AccountAssetType.Crypto.type) {
        const price = cryptocurrencyPricesStore.getCryptocurrencyPriceInFiat(toCurrency, fromCurrency);
        if (isNumber(price) && price > 0) {
            const fromFraction = getCurrencyFraction(fromCurrency) || 0;
            const toFraction = getCurrencyFraction(toCurrency) || 0;
            return (amount / price) * Math.pow(10, toFraction - fromFraction);
        }
    } else if (fromAssetType === AccountAssetType.Fiat.type && toAssetType === AccountAssetType.Stock.type) {
        const price = stockPricesStore.getStockPriceInFiat(toCurrency, fromCurrency, exchangeRatesStore);
        if (isNumber(price) && price > 0) {
            const fromFraction = getCurrencyFraction(fromCurrency) || 0;
            const toFraction = getCurrencyFraction(toCurrency) || 0;
            return (amount / price) * Math.pow(10, toFraction - fromFraction);
        }
    }

    return null;
}

export function appendCurrencySymbol(value: string, currencyDisplayType: CurrencyDisplayType, currencyCode: string, currencyUnit: string, currencyName: string, isPlural: boolean): string {
    const symbol = getAmountPrependAndAppendCurrencySymbol(currencyDisplayType, currencyCode, currencyUnit, currencyName, isPlural);

    if (!symbol) {
        return value;
    }

    const separator = currencyDisplayType.separator || '';
    let ret = value;

    if (symbol.prependText) {
        ret = symbol.prependText + separator + ret;
    }

    if (symbol.appendText) {
        ret = ret + separator + symbol.appendText;
    }

    return ret;
}

export function getAmountPrependAndAppendCurrencySymbol(currencyDisplayType: CurrencyDisplayType, currencyCode: string, currencyUnit: string, currencyName: string, isPlural: boolean): CurrencyPrependAndAppendText | null {
    if (!currencyDisplayType) {
        return null;
    }

    let symbol = '';

    if (currencyDisplayType.symbol === CurrencyDisplaySymbol.Symbol) {
        const currencyInfo = ALL_CURRENCIES[currencyCode] || ALL_CRYPTOCURRENCIES[currencyCode];

        if (currencyInfo && currencyInfo.symbol && currencyInfo.symbol.normal) {
            symbol = currencyInfo.symbol.normal;

            if (isPlural && currencyInfo.symbol.plural) {
                symbol = currencyInfo.symbol.plural;
            }
        }

        if (!symbol) {
            symbol = DEFAULT_CURRENCY_SYMBOL;
        }
    } else if (currencyDisplayType.symbol === CurrencyDisplaySymbol.Code) {
        symbol = currencyCode;
    } else if (currencyDisplayType.symbol === CurrencyDisplaySymbol.Unit) {
        symbol = currencyUnit;
    } else if (currencyDisplayType.symbol === CurrencyDisplaySymbol.Name) {
        symbol = currencyName;
    }

    if (currencyDisplayType.location === CurrencyDisplayLocation.BeforeAmount) {
        return {
            prependText: symbol
        };
    } else if (currencyDisplayType.location === CurrencyDisplayLocation.AfterAmount) {
        return {
            appendText: symbol
        };
    } else {
        return null;
    }
}
