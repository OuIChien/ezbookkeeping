import { ref, computed } from 'vue';
import { defineStore } from 'pinia';

import { type BeforeResolveFunction, itemAndIndex } from '@/core/base.ts';

import type {
    UserCustomExchangeRateUpdateResponse,
    LatestExchangeRate,
    LatestExchangeRateResponse
} from '@/models/exchange_rate.ts';

import { isEquals } from '@/lib/common.ts';
import {
    isUnixTimeYearMonthDayEquals,
    isUnixTimeYearMonthDayHourEquals,
    getCurrentUnixTime
} from '@/lib/datetime.ts';
import { getExchangedAmountByRate } from '@/lib/numeral.ts';

import logger from '@/lib/logger.ts';
import services from '@/lib/services.ts';
import { useCryptocurrencyPricesStore } from '@/stores/cryptocurrencyPrices.ts';
import { getCurrencyType } from '@/consts/currency.ts';
import { CurrencyType } from '@/core/currency.ts';

const exchangeRatesLocalStorageKey = 'ebk_app_exchange_rates';
const userDataSourceType = 'user_custom';

function isCryptocurrency(currencyCode: string): boolean {
    return getCurrencyType(currencyCode) === CurrencyType.Cryptocurrency;
}

interface LatestExchangeRates {
    readonly time?: number;
    readonly data?: LatestExchangeRateResponse;
}

function getExchangeRatesFromLocalStorage(): LatestExchangeRates {
    const storageData = localStorage.getItem(exchangeRatesLocalStorageKey) || '{}';
    return JSON.parse(storageData) as LatestExchangeRates;
}

function setExchangeRatesToLocalStorage(value: LatestExchangeRates): void {
    const storageData = JSON.stringify(value);
    localStorage.setItem(exchangeRatesLocalStorageKey, storageData);
}

function clearExchangeRatesFromLocalStorage(): void {
    localStorage.removeItem(exchangeRatesLocalStorageKey);
}

export const useExchangeRatesStore = defineStore('exchangeRates', () => {
    const latestExchangeRates = ref<LatestExchangeRates>(getExchangeRatesFromLocalStorage());

    const isUserCustomExchangeRates = computed((): boolean => {
        if (!latestExchangeRates.value || !latestExchangeRates.value.data) {
            return false;
        }

        return latestExchangeRates.value.data.dataSource === userDataSourceType;
    });

    const exchangeRatesLastUpdateTime = computed<number | null>(() => {
        const exchangeRates = latestExchangeRates.value || {};
        return exchangeRates && exchangeRates.data ? exchangeRates.data.updateTime : null;
    });

    const latestExchangeRateMap = computed<Record<string, LatestExchangeRate>>(() => {
        const exchangeRateMap: Record<string, LatestExchangeRate> = {};

        if (!latestExchangeRates.value || !latestExchangeRates.value.data || !latestExchangeRates.value.data.exchangeRates) {
            return exchangeRateMap;
        }

        for (const exchangeRate of latestExchangeRates.value.data.exchangeRates) {
            exchangeRateMap[exchangeRate.currency] = exchangeRate;
        }

        return exchangeRateMap;
    });

    function updateExchangeRateToLatestExchangeRateList(latestExchangeRate: LatestExchangeRate, updateTime: number): void {
        if (!latestExchangeRates.value || !latestExchangeRates.value.data || !latestExchangeRates.value.data.exchangeRates) {
            return;
        }

        const exchangeRates = latestExchangeRates.value.data.exchangeRates;
        let changed = false;

        for (const [exchangeRate, index] of itemAndIndex(exchangeRates)) {
            if (exchangeRate.currency === latestExchangeRate.currency) {
                exchangeRates.splice(index, 1, latestExchangeRate);
                changed = true;
                break;
            }
        }

        if (!changed) {
            exchangeRates.push(latestExchangeRate);
            changed = true;
        }

        latestExchangeRates.value.data.updateTime = updateTime;

        if (changed) {
            setExchangeRatesToLocalStorage(latestExchangeRates.value);
        }
    }

    function removeExchangeRateFromLatestExchangeRateList(currency: string): void {
        if (!latestExchangeRates.value || !latestExchangeRates.value.data || !latestExchangeRates.value.data.exchangeRates) {
            return;
        }

        const exchangeRates = latestExchangeRates.value.data.exchangeRates;
        let changed = false;

        for (const [exchangeRate, index] of itemAndIndex(exchangeRates)) {
            if (exchangeRate.currency === currency) {
                exchangeRates.splice(index, 1);
                changed = true;
                break;
            }
        }

        if (changed) {
            setExchangeRatesToLocalStorage(latestExchangeRates.value);
        }
    }

    function resetLatestExchangeRates(): void {
        latestExchangeRates.value = {};
        clearExchangeRatesFromLocalStorage();
    }

    function getLatestExchangeRates({ silent, force }: { silent: boolean, force: boolean }): Promise<LatestExchangeRateResponse> {
        const currentExchangeRateData = latestExchangeRates.value;
        const now = getCurrentUnixTime();

        if (!force) {
            if (currentExchangeRateData && currentExchangeRateData.time && currentExchangeRateData.data && isUnixTimeYearMonthDayEquals(currentExchangeRateData.data.updateTime, now)) {
                return Promise.resolve(currentExchangeRateData.data);
            }

            if (currentExchangeRateData && currentExchangeRateData.time && currentExchangeRateData.data && isUnixTimeYearMonthDayHourEquals(currentExchangeRateData.time, now)) {
                return Promise.resolve(currentExchangeRateData.data);
            }
        }

        return new Promise((resolve, reject) => {
            services.getLatestExchangeRates({
                ignoreError: silent
            }).then(response => {
                const data = response.data;

                if (!data || !data.success || !data.result) {
                    reject({ message: 'Unable to retrieve exchange rates data' });
                    return;
                }

                const currentData = getExchangeRatesFromLocalStorage();

                if (force && currentData && currentData.data && isEquals(currentData.data, data.result)) {
                    reject({ message: 'Exchange rates data is up to date', isUpToDate: true });
                    return;
                }

                latestExchangeRates.value = {
                    time: now,
                    data: data.result
                };
                setExchangeRatesToLocalStorage(latestExchangeRates.value);

                resolve(data.result);
            }).catch(error => {
                logger.error('failed to retrieve latest exchange rates data', error);

                if (error && error.processed) {
                    reject(error);
                } else if (error.response && error.response.data && error.response.data.errorMessage) {
                    reject({ error: error.response.data });
                } else {
                    reject({ message: 'Unable to retrieve exchange rates data' });
                }
            });
        });
    }

    function updateUserCustomExchangeRate({ currency, rate }: { currency: string, rate: number }): Promise<UserCustomExchangeRateUpdateResponse> {
        return new Promise((resolve, reject) => {
            services.updateUserCustomExchangeRate({
                currency: currency,
                rate: rate.toString()
            }).then(response => {
                const data = response.data;

                if (!data || !data.success || !data.result) {
                    reject({ message: 'Unable to update user custom exchange rate' });
                    return;
                }

                const exchangeRate: LatestExchangeRate = {
                    currency: data.result.currency,
                    rate: data.result.rate
                };

                updateExchangeRateToLatestExchangeRateList(exchangeRate, data.result.updateTime);

                resolve(data.result);
            }).catch(error => {
                logger.error('failed to update user custom exchange rate', error);

                if (error.response && error.response.data && error.response.data.errorMessage) {
                    reject({ error: error.response.data });
                } else if (!error.processed) {
                    reject({ message: 'Unable to update user custom exchange rate' });
                } else {
                    reject(error);
                }
            });
        });
    }

    function deleteUserCustomExchangeRate({ currency, beforeResolve }: { currency: string, beforeResolve?: BeforeResolveFunction }): Promise<boolean> {
        return new Promise((resolve, reject) => {
            services.deleteUserCustomExchangeRate({
                currency: currency
            }).then(response => {
                const data = response.data;

                if (!data || !data.success || !data.result) {
                    reject({ message: 'Unable to delete this user custom exchange rate' });
                    return;
                }

                if (beforeResolve) {
                    beforeResolve(() => {
                        removeExchangeRateFromLatestExchangeRateList(currency);
                    });
                } else {
                    removeExchangeRateFromLatestExchangeRateList(currency);
                }

                resolve(data.result);
            }).catch(error => {
                logger.error('failed to delete user custom exchange rate', error);

                if (error.response && error.response.data && error.response.data.errorMessage) {
                    reject({ error: error.response.data });
                } else if (!error.processed) {
                    reject({ message: 'Unable to delete this user custom exchange rate' });
                } else {
                    reject(error);
                }
            });
        });
    }

    // Internal helper function to convert between fiat currencies using exchange rates
    // This avoids recursion when converting cryptocurrencies through USDT
    function getFiatExchangedAmount(amount: number, fromCurrency: string, toCurrency: string): number | null {
        if (amount === 0) {
            return 0;
        }

        // If both currencies are the same, return the amount directly
        if (fromCurrency === toCurrency) {
            return amount;
        }

        // Special handling for USDT: if USDT is not in exchange rates, assume USDT = USD = 1
        // This is a reasonable assumption since USDT is a stablecoin pegged to USD
        if (fromCurrency === 'USDT' || toCurrency === 'USDT') {
            if (!latestExchangeRates.value || !latestExchangeRates.value.data || !latestExchangeRates.value.data.exchangeRates) {
                // If no exchange rates available, assume USDT = USD = 1
                if (fromCurrency === 'USDT' && toCurrency === 'USD') {
                    return amount;
                }
                if (fromCurrency === 'USD' && toCurrency === 'USDT') {
                    return amount;
                }
                return null;
            }

            const exchangeRates = latestExchangeRates.value.data.exchangeRates;
            const exchangeRateMap: Record<string, LatestExchangeRate> = {};

            for (const exchangeRate of exchangeRates) {
                exchangeRateMap[exchangeRate.currency] = exchangeRate;
            }

            // If USDT is not in exchange rates, treat it as USD
            if (fromCurrency === 'USDT' && !exchangeRateMap['USDT']) {
                if (toCurrency === 'USD') {
                    return amount;
                }
                // Convert USDT -> USD -> target currency
                const usdToTarget = getFiatExchangedAmount(amount, 'USD', toCurrency);
                return usdToTarget;
            }

            if (toCurrency === 'USDT' && !exchangeRateMap['USDT']) {
                if (fromCurrency === 'USD') {
                    return amount;
                }
                // Convert source currency -> USD -> USDT
                const sourceToUsd = getFiatExchangedAmount(amount, fromCurrency, 'USD');
                return sourceToUsd;
            }
        }

        if (!latestExchangeRates.value || !latestExchangeRates.value.data || !latestExchangeRates.value.data.exchangeRates) {
            return null;
        }

        const exchangeRates = latestExchangeRates.value.data.exchangeRates;
        const exchangeRateMap: Record<string, LatestExchangeRate> = {};

        for (const exchangeRate of exchangeRates) {
            exchangeRateMap[exchangeRate.currency] = exchangeRate;
        }

        const fromCurrencyExchangeRate = exchangeRateMap[fromCurrency];
        const toCurrencyExchangeRate = exchangeRateMap[toCurrency];

        if (!fromCurrencyExchangeRate || !toCurrencyExchangeRate) {
            return null;
        }

        return getExchangedAmountByRate(amount, fromCurrencyExchangeRate.rate, toCurrencyExchangeRate.rate);
    }

    function getExchangedAmount(amount: number, fromCurrency: string, toCurrency: string): number | null {
        if (amount === 0) {
            return 0;
        }

        // If both currencies are the same, return the amount directly
        if (fromCurrency === toCurrency) {
            return amount;
        }

        const fromIsCrypto = isCryptocurrency(fromCurrency);
        const toIsCrypto = isCryptocurrency(toCurrency);

        // Handle cryptocurrency conversions
        if (fromIsCrypto || toIsCrypto) {
            const cryptocurrencyPricesStore = useCryptocurrencyPricesStore();

            // Case 1: Both are cryptocurrencies - convert through USD
            if (fromIsCrypto && toIsCrypto) {
                // Convert from crypto to USD
                const priceInUSD = cryptocurrencyPricesStore.getCryptocurrencyPriceInUSD(fromCurrency);
                if (!priceInUSD) {
                    return null;
                }

                const priceInUSDNum = parseFloat(priceInUSD);
                if (isNaN(priceInUSDNum) || priceInUSDNum === 0) {
                    return null;
                }

                // Convert amount to USD
                const amountInUSD = amount * priceInUSDNum;

                // Convert from USD to target crypto
                const targetPriceInUSD = cryptocurrencyPricesStore.getCryptocurrencyPriceInUSD(toCurrency);
                if (!targetPriceInUSD) {
                    return null;
                }

                const targetPriceInUSDNum = parseFloat(targetPriceInUSD);
                if (isNaN(targetPriceInUSDNum) || targetPriceInUSDNum === 0) {
                    return null;
                }

                // Convert USD amount to target crypto
                return amountInUSD / targetPriceInUSDNum;
            }

            // Case 2: From cryptocurrency to fiat currency - convert crypto -> USD -> fiat
            if (fromIsCrypto && !toIsCrypto) {
                const priceInUSD = cryptocurrencyPricesStore.getCryptocurrencyPriceInUSD(fromCurrency);
                if (!priceInUSD) {
                    return null;
                }

                const priceInUSDNum = parseFloat(priceInUSD);
                if (isNaN(priceInUSDNum) || priceInUSDNum === 0) {
                    return null;
                }

                // Convert amount to USD
                const amountInUSD = amount * priceInUSDNum;

                // Convert from USD to fiat using exchange rates (use helper to avoid recursion)
                return getFiatExchangedAmount(amountInUSD, 'USD', toCurrency);
            }

            // Case 3: From fiat currency to cryptocurrency - convert fiat -> USD -> crypto
            if (!fromIsCrypto && toIsCrypto) {
                // First convert fiat to USD (use helper to avoid recursion)
                const amountInUSD = getFiatExchangedAmount(amount, fromCurrency, 'USD');
                if (amountInUSD === null) {
                    return null;
                }

                // Then convert USD to crypto
                const targetPriceInUSD = cryptocurrencyPricesStore.getCryptocurrencyPriceInUSD(toCurrency);
                if (!targetPriceInUSD) {
                    return null;
                }

                const targetPriceInUSDNum = parseFloat(targetPriceInUSD);
                if (isNaN(targetPriceInUSDNum) || targetPriceInUSDNum === 0) {
                    return null;
                }

                // Convert USD amount to target crypto
                return amountInUSD / targetPriceInUSDNum;
            }
        }

        // Handle traditional fiat currency conversions
        return getFiatExchangedAmount(amount, fromCurrency, toCurrency);
    }

    return {
        // states
        latestExchangeRates,
        // computed states
        isUserCustomExchangeRates,
        exchangeRatesLastUpdateTime,
        latestExchangeRateMap,
        // functions
        resetLatestExchangeRates,
        getLatestExchangeRates,
        updateUserCustomExchangeRate,
        deleteUserCustomExchangeRate,
        getExchangedAmount
    };
});
