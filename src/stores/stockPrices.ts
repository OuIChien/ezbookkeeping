import { ref, computed } from 'vue';
import { defineStore } from 'pinia';

import { type BeforeResolveFunction, itemAndIndex } from '@/core/base.ts';

import type {
    LatestStockPriceResponse
} from '@/models/stock_price.ts';
import type {
    StockInfoResponse,
    StockCreateRequest,
    StockModifyRequest
} from '@/models/stock.ts';
import type {
    ExternalDataSourceConfigResponse,
    ExternalDataSourceConfigSaveRequest
} from '@/models/external_data_source.ts';

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

    const allStocks = ref<StockInfoResponse[]>([]);
    const stockConfig = ref<ExternalDataSourceConfigResponse | null>(null);

    const allVisibleStocks = computed<StockInfoResponse[]>(() => {
        return allStocks.value.filter(s => !s.isHidden);
    });

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

    function loadAllStocks({ force }: { force: boolean }): Promise<StockInfoResponse[]> {
        if (!force && allStocks.value.length > 0) {
            return Promise.resolve(allStocks.value);
        }

        return new Promise((resolve, reject) => {
            services.getAllStocks().then(response => {
                const data = response.data;

                if (!data || !data.success || !data.result) {
                    reject({ message: 'Unable to retrieve stocks list' });
                    return;
                }

                allStocks.value = data.result;
                resolve(data.result);
            }).catch(error => {
                logger.error('failed to retrieve stocks list', error);
                reject(error);
            });
        });
    }

    function addStock(req: StockCreateRequest): Promise<StockInfoResponse> {
        return new Promise((resolve, reject) => {
            services.addStock(req).then(response => {
                const data = response.data;

                if (!data || !data.success || !data.result) {
                    reject({ message: 'Unable to add stock' });
                    return;
                }

                allStocks.value.push(data.result);
                resolve(data.result);
            }).catch(error => {
                logger.error('failed to add stock', error);
                reject(error);
            });
        });
    }

    function modifyStock(req: StockModifyRequest): Promise<StockInfoResponse> {
        return new Promise((resolve, reject) => {
            services.modifyStock(req).then(response => {
                const data = response.data;

                if (!data || !data.success || !data.result) {
                    reject({ message: 'Unable to modify stock' });
                    return;
                }

                for (const [stock, index] of itemAndIndex(allStocks.value)) {
                    if (stock.symbol === data.result.symbol) {
                        allStocks.value.splice(index, 1, data.result);
                        break;
                    }
                }

                resolve(data.result);
            }).catch(error => {
                logger.error('failed to modify stock', error);
                reject(error);
            });
        });
    }

    function hideStock({ symbol, hidden }: { symbol: string, hidden: boolean }): Promise<boolean> {
        return new Promise((resolve, reject) => {
            services.hideStock({ symbol, hidden }).then(response => {
                const data = response.data;

                if (!data || !data.success || !data.result) {
                    reject({ message: 'Unable to hide/unhide stock' });
                    return;
                }

                for (const stock of allStocks.value) {
                    if (stock.symbol === symbol) {
                        stock.isHidden = hidden;
                        break;
                    }
                }

                resolve(data.result);
            }).catch(error => {
                logger.error('failed to hide/unhide stock', error);
                reject(error);
            });
        });
    }

    function deleteStock({ symbol, beforeResolve }: { symbol: string, beforeResolve?: BeforeResolveFunction }): Promise<boolean> {
        return new Promise((resolve, reject) => {
            services.deleteStock({ symbol }).then(response => {
                const data = response.data;

                if (!data || !data.success || !data.result) {
                    reject({ message: 'Unable to delete stock' });
                    return;
                }

                if (beforeResolve) {
                    beforeResolve(() => {
                        for (const [stock, index] of itemAndIndex(allStocks.value)) {
                            if (stock.symbol === symbol) {
                                allStocks.value.splice(index, 1);
                                break;
                            }
                        }
                    });
                } else {
                    for (const [stock, index] of itemAndIndex(allStocks.value)) {
                        if (stock.symbol === symbol) {
                            allStocks.value.splice(index, 1);
                            break;
                        }
                    }
                }

                resolve(data.result);
            }).catch(error => {
                logger.error('failed to delete stock', error);
                reject(error);
            });
        });
    }

    function loadStockConfig(): Promise<ExternalDataSourceConfigResponse> {
        return new Promise((resolve, reject) => {
            services.getStockConfig().then(response => {
                const data = response.data;

                if (!data || !data.success) {
                    reject({ message: 'Unable to retrieve stock config' });
                    return;
                }

                if (data.result) {
                    stockConfig.value = data.result;
                }
                
                resolve(data.result);
            }).catch(error => {
                logger.error('failed to retrieve stock config', error);
                reject(error);
            });
        });
    }

    function saveStockConfig(req: ExternalDataSourceConfigSaveRequest): Promise<ExternalDataSourceConfigResponse> {
        return new Promise((resolve, reject) => {
            services.saveStockConfig(req).then(response => {
                const data = response.data;

                if (!data || !data.success || !data.result) {
                    reject({ message: 'Unable to save stock config' });
                    return;
                }

                stockConfig.value = data.result;
                resolve(data.result);
            }).catch(error => {
                logger.error('failed to save stock config', error);
                reject(error);
            });
        });
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
                    reject({ message: 'Data is up to date', isUpToDate: true });
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

    function getStockPriceInFiat(symbol: string, fiatCurrency: string, exchangeRatesStore: any): number | null {
        const priceStr = latestStockPriceMap.value[symbol];
        if (!priceStr) {
            return null;
        }

        const price = parseFloat(priceStr);
        const priceData = latestStockPrices.value.data?.prices?.find(p => p.symbol === symbol);
        
        if (!priceData) {
            return null;
        }

        // 1 stock = price priceData.currency
        // We want to convert `price` of `priceData.currency` to `fiatCurrency`.
        return exchangeRatesStore.getExchangedAmount(price, priceData.currency, fiatCurrency);
    }

    return {
        // states
        latestStockPrices,
        allStocks,
        stockConfig,
        // computed states
        allVisibleStocks,
        stockPricesLastUpdateTime,
        latestStockPriceMap,
        // functions
        loadAllStocks,
        addStock,
        modifyStock,
        hideStock,
        deleteStock,
        loadStockConfig,
        saveStockConfig,
        resetLatestStockPrices,
        getLatestStockPrices,
        getStockPriceInFiat
    };
});
