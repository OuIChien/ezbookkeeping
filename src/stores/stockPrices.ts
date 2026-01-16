import { ref, computed } from 'vue';
import { defineStore } from 'pinia';

import type {
    LatestStockPriceResponse
} from '@/models/stock_price.ts';

import { isEquals } from '@/lib/common.ts';
import {
    isUnixTimeYearMonthDayHourEquals,
    getCurrentUnixTime
} from '@/lib/datetime.ts';

import logger from '@/lib/logger.ts';
import services from '@/lib/services.ts';

const stockPricesLocalStorageKey = 'ebk_app_stock_prices';

interface LatestStockPrices {
    readonly time?: number;
    readonly data?: LatestStockPriceResponse;
}

function getStockPricesFromLocalStorage(): LatestStockPrices {
    const storageData = localStorage.getItem(stockPricesLocalStorageKey) || '{}';
    return JSON.parse(storageData) as LatestStockPrices;
}

function setStockPricesToLocalStorage(value: LatestStockPrices): void {
    const storageData = JSON.stringify(value);
    localStorage.setItem(stockPricesLocalStorageKey, storageData);
}

function clearStockPricesFromLocalStorage(): void {
    localStorage.removeItem(stockPricesLocalStorageKey);
}

export const useStockPricesStore = defineStore('stockPrices', () => {
    const latestStockPrices = ref<LatestStockPrices>(getStockPricesFromLocalStorage());

    const stockPricesLastUpdateTime = computed<number | null>(() => {
        const prices = latestStockPrices.value || {};
        return prices && prices.data ? prices.data.updateTime : null;
    });

    const latestStockPriceMap = computed<Record<string, string>>(() => {
        const priceMap: Record<string, string> = {};

        if (!latestStockPrices.value || !latestStockPrices.value.data || !latestStockPrices.value.data.prices) {
            return priceMap;
        }

        for (const price of latestStockPrices.value.data.prices) {
            priceMap[price.symbol] = price.price;
        }

        return priceMap;
    });

    function resetLatestStockPrices(): void {
        latestStockPrices.value = {};
        clearStockPricesFromLocalStorage();
    }

    function getLatestStockPrices({ silent, force }: { silent: boolean, force: boolean }): Promise<LatestStockPriceResponse> {
        const currentPriceData = latestStockPrices.value;
        const now = getCurrentUnixTime();

        if (!force) {
            if (currentPriceData && currentPriceData.time && currentPriceData.data && isUnixTimeYearMonthDayHourEquals(currentPriceData.time, now)) {
                return Promise.resolve(currentPriceData.data);
            }
        }

        return new Promise((resolve, reject) => {
            services.getLatestStockPrices({
                ignoreError: silent
            }).then(response => {
                const data = response.data;

                if (!data || !data.success || !data.result) {
                    reject({ message: 'Unable to retrieve stock prices data' });
                    return;
                }

                const currentData = getStockPricesFromLocalStorage();

                if (force && currentData && currentData.data && isEquals(currentData.data, data.result)) {
                    reject({ message: 'Stock prices data is up to date', isUpToDate: true });
                    return;
                }

                latestStockPrices.value = {
                    time: now,
                    data: data.result
                };
                setStockPricesToLocalStorage(latestStockPrices.value);

                resolve(data.result);
            }).catch(error => {
                logger.error('failed to retrieve latest stock prices data', error);

                if (error && error.processed) {
                    reject(error);
                } else if (error.response && error.response.data && error.response.data.errorMessage) {
                    reject({ error: error.response.data });
                } else {
                    reject({ message: 'Unable to retrieve stock prices data' });
                }
            });
        });
    }

    return {
        // states
        latestStockPrices,
        // computed states
        stockPricesLastUpdateTime,
        latestStockPriceMap,
        // functions
        resetLatestStockPrices,
        getLatestStockPrices
    };
});
