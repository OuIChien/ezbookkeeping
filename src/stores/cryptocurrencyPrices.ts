import { ref, computed } from 'vue';
import { defineStore } from 'pinia';

import { type BeforeResolveFunction, itemAndIndex } from '@/core/base.ts';

import type {
    LatestCryptocurrencyPrice,
    LatestCryptocurrencyPriceResponse
} from '@/models/cryptocurrency_price.ts';

import { isEquals } from '@/lib/common.ts';
import {
    isUnixTimeYearMonthDayEquals,
    isUnixTimeYearMonthDayHourEquals,
    getCurrentUnixTime
} from '@/lib/datetime.ts';
import { getExchangedAmountByRate } from '@/lib/numeral.ts';

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

    const cryptocurrencyPricesLastUpdateTime = computed<number | null>(() => {
        const prices = latestCryptocurrencyPrices.value || {};
        return prices && prices.data ? prices.data.updateTime : null;
    });

    const latestCryptocurrencyPriceMap = computed<Record<string, LatestCryptocurrencyPrice>>(() => {
        const priceMap: Record<string, LatestCryptocurrencyPrice> = {};

        if (!latestCryptocurrencyPrices.value || !latestCryptocurrencyPrices.value.data || !latestCryptocurrencyPrices.value.data.prices) {
            return priceMap;
        }

        for (const price of latestCryptocurrencyPrices.value.data.prices) {
            priceMap[price.symbol] = price;
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
            if (currentPriceData && currentPriceData.time && currentPriceData.data && isUnixTimeYearMonthDayEquals(currentPriceData.data.updateTime, now)) {
                return Promise.resolve(currentPriceData.data);
            }

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
                    reject({ message: 'Cryptocurrency prices data is up to date', isUpToDate: true });
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

    function getCryptocurrencyPriceInUSDT(symbol: string): string | null {
        const priceMap = latestCryptocurrencyPriceMap.value;
        const price = priceMap[symbol];

        if (!price) {
            return null;
        }

        return price.price;
    }

    function getCryptocurrencyPriceInFiat(symbol: string, fiatCurrency: string): number | null {
        const exchangeRatesStore = useExchangeRatesStore();
        const priceInUSDT = getCryptocurrencyPriceInUSDT(symbol);

        if (!priceInUSDT) {
            return null;
        }

        // Get USDT to fiat exchange rate
        const usdtToFiatRate = exchangeRatesStore.getExchangedAmount(1, 'USDT', fiatCurrency);

        if (usdtToFiatRate === null) {
            return null;
        }

        // Calculate: cryptoPriceInUSDT * usdtToFiatRate
        const priceInUSDTNum = parseFloat(priceInUSDT);
        if (isNaN(priceInUSDTNum)) {
            return null;
        }

        return priceInUSDTNum * usdtToFiatRate;
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

