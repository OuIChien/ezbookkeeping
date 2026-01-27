import { ref, computed } from 'vue';
import { defineStore } from 'pinia';

import type {
    LatestCryptocurrencyPriceResponse
} from '@/models/cryptocurrency_price.ts';

import { isEquals } from '@/lib/common.ts';
import {
    isUnixTimeYearMonthDayHourEquals,
    getCurrentUnixTime
} from '@/lib/datetime.ts';

import logger from '@/lib/logger.ts';
import services from '@/lib/services.ts';
import { useExchangeRatesStore } from '@/stores/exchangeRates.ts';

const cryptocurrencyPricesLocalStorageKey = 'ebk_app_cryptocurrency_prices';

interface LatestCryptocurrencyPrices {
    readonly time?: number;
    readonly data?: LatestCryptocurrencyPriceResponse;
}

function getCryptocurrencyPricesFromLocalStorage(): LatestCryptocurrencyPrices {
    const storageData = localStorage.getItem(cryptocurrencyPricesLocalStorageKey) || '{}';
    return JSON.parse(storageData) as LatestCryptocurrencyPrices;
}

function setCryptocurrencyPricesToLocalStorage(value: LatestCryptocurrencyPrices): void {
    const storageData = JSON.stringify(value);
    localStorage.setItem(cryptocurrencyPricesLocalStorageKey, storageData);
}

function clearCryptocurrencyPricesFromLocalStorage(): void {
    localStorage.removeItem(cryptocurrencyPricesLocalStorageKey);
}

export const useCryptocurrencyPricesStore = defineStore('cryptocurrencyPrices', () => {
    const latestCryptocurrencyPrices = ref<LatestCryptocurrencyPrices>(getCryptocurrencyPricesFromLocalStorage());
    const exchangeRatesStore = useExchangeRatesStore();

    const cryptocurrencyPricesLastUpdateTime = computed<number | null>(() => {
        const prices = latestCryptocurrencyPrices.value || {};
        return prices && prices.data ? prices.data.updateTime : null;
    });

    const latestCryptocurrencyPriceMap = computed<Record<string, string>>(() => {
        const priceMap: Record<string, string> = {};

        if (!latestCryptocurrencyPrices.value || !latestCryptocurrencyPrices.value.data || !latestCryptocurrencyPrices.value.data.prices) {
            return priceMap;
        }

        for (const price of latestCryptocurrencyPrices.value.data.prices) {
            priceMap[price.symbol] = price.price;
        }

        return priceMap;
    });

    function resetLatestCryptocurrencyPrices(): void {
        latestCryptocurrencyPrices.value = {};
        clearCryptocurrencyPricesFromLocalStorage();
    }

    function getLatestCryptocurrencyPrices({ silent, force }: { silent: boolean, force: boolean }): Promise<LatestCryptocurrencyPriceResponse> {
        const currentPriceData = latestCryptocurrencyPrices.value;
        const now = getCurrentUnixTime();

        if (!force) {
            // Check if data is fresh (same day or same hour) - Logic from design doc
            // "unlike daily exchange rates ... Consider shorter cache validity (e.g., same hour or 15 minutes)"
            // For now, using same hour validity as per basic requirement or similar to exchange rates logic for simplicity in first pass
            // But design doc suggests shorter. Let's use hour check.
            if (currentPriceData && currentPriceData.time && currentPriceData.data && isUnixTimeYearMonthDayHourEquals(currentPriceData.time, now)) {
                return Promise.resolve(currentPriceData.data);
            }
        }

        return new Promise((resolve, reject) => {
            services.getLatestCryptocurrencyPrices({
                ignoreError: silent
            }).then(response => {
                const data = response.data;

                if (!data || !data.success || !data.result) {
                    reject({ message: 'Unable to retrieve cryptocurrency prices data' });
                    return;
                }

                const currentData = getCryptocurrencyPricesFromLocalStorage();

                if (force && currentData && currentData.data && isEquals(currentData.data, data.result)) {
                    reject({ message: 'Data is up to date', isUpToDate: true });
                    return;
                }

                latestCryptocurrencyPrices.value = {
                    time: now,
                    data: data.result
                };
                setCryptocurrencyPricesToLocalStorage(latestCryptocurrencyPrices.value);

                resolve(data.result);
            }).catch(error => {
                logger.error('failed to retrieve latest cryptocurrency prices data', error);

                if (error && error.processed) {
                    reject(error);
                } else if (error.response && error.response.data && error.response.data.errorMessage) {
                    reject({ error: error.response.data });
                } else {
                    reject({ message: 'Unable to retrieve cryptocurrency prices data' });
                }
            });
        });
    }

    function getCryptocurrencyPriceInUSDT(symbol: string): number | null {
        const priceStr = latestCryptocurrencyPriceMap.value[symbol];
        if (!priceStr) {
            return null;
        }
        return parseFloat(priceStr);
    }

    function getCryptocurrencyPriceInFiat(symbol: string, fiatCurrency: string): number | null {
        const priceInUSDT = getCryptocurrencyPriceInUSDT(symbol);
        if (priceInUSDT === null) {
            return null;
        }

        // USDT is treated as USD in exchange rates usually, or we can look up USDT->Fiat rate.
        // Design doc says: "All cryptocurrency prices will be fetched in USDT" ... "USDT will be treated as the base currency (rate = "1")"
        // "This allows easy conversion to fiat currencies via existing exchange rates"
        // And "Get USDT to fiat exchange rate from exchange rates store"

        // Assuming USDT is available in exchange rates or equivalent to USD.
        // If "USDT" is in exchange rates, use it. If not, maybe use "USD"?
        // The design says "USDT will be treated as the base currency".
        // But exchange rates usually have a base currency (e.g. USD or EUR).
        // If our exchange rates base is USD, and USDT ~= USD.
        
        // Let's assume we can convert USDT -> Fiat using getExchangedAmount from exchangeRatesStore.
        // We act as if we have `priceInUSDT` amount of USDT, and we want to convert it to `fiatCurrency`.
        
        // However, exchangeRatesStore.getExchangedAmount needs `fromCurrency` and `toCurrency`.
        // We will pass "USDT" as fromCurrency. If exchange rates data doesn't have USDT, this might fail unless we assume USDT=USD.
        // For safety, and as per design doc usually implying standard crypto symbols, let's try "USDT".
        // If the user's exchange rate source provides USDT, it works.
        // If not, and if we want to support it, we might need a fallback or mapping USDT->USD.
        // Design doc 7.3 says USDT is supported cryptocurrency.
        
        // Let's try direct conversion first.
        const exchangedAmount = exchangeRatesStore.getExchangedAmount(priceInUSDT, 'USDT', fiatCurrency);
        
        if (exchangedAmount !== null) {
            return exchangedAmount;
        }
        
        // Fallback: If USDT not found in exchange rates, try USD if available (1 USDT ~= 1 USD assumption for simple display)
        // This is a common practical assumption for personal finance if precise USDT rate isn't available.
        const exchangedAmountUSD = exchangeRatesStore.getExchangedAmount(priceInUSDT, 'USD', fiatCurrency);
        return exchangedAmountUSD;
    }

    return {
        // states
        latestCryptocurrencyPrices,
        // computed states
        cryptocurrencyPricesLastUpdateTime,
        latestCryptocurrencyPriceMap,
        // functions
        resetLatestCryptocurrencyPrices,
        getLatestCryptocurrencyPrices,
        getCryptocurrencyPriceInUSDT,
        getCryptocurrencyPriceInFiat
    };
});
