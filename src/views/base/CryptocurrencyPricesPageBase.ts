import { computed } from 'vue';

import { useI18n } from '@/locales/helpers.ts';

import { useCryptocurrencyPricesStore } from '@/stores/cryptocurrencyPrices.ts';

import type {
    LatestCryptocurrencyPriceResponse
} from '@/models/cryptocurrency_price.ts';

import { parseDateTimeFromUnixTime } from '@/lib/datetime.ts';

export interface LocalizedLatestCryptocurrencyPrice {
    readonly symbol: string;
    readonly symbolDisplayName: string;
    readonly price: string;
}

export function useCryptocurrencyPricesPageBase() {
    const { getCurrencyName, formatDateTimeToLongDate } = useI18n();

    const cryptocurrencyPricesStore = useCryptocurrencyPricesStore();

    const cryptocurrencyPricesData = computed<LatestCryptocurrencyPriceResponse | undefined>(() => cryptocurrencyPricesStore.latestCryptocurrencyPrices.data);

    const cryptocurrencyPricesDataUpdateTime = computed<string>(() => {
        if (!cryptocurrencyPricesStore.cryptocurrencyPricesLastUpdateTime) {
            return '';
        }

        const updateTime = parseDateTimeFromUnixTime(cryptocurrencyPricesStore.cryptocurrencyPricesLastUpdateTime);
        return formatDateTimeToLongDate(updateTime);
    });

    const availableCryptocurrencyPrices = computed<LocalizedLatestCryptocurrencyPrice[]>(() => {
        const availablePrices: LocalizedLatestCryptocurrencyPrice[] = [];

        if (!cryptocurrencyPricesData.value || !cryptocurrencyPricesData.value.prices) {
            return availablePrices;
        }

        for (const price of cryptocurrencyPricesData.value.prices) {
            availablePrices.push({
                symbol: price.symbol,
                symbolDisplayName: getCurrencyName(price.symbol),
                price: price.price
            });
        }

        // Sort by symbol
        availablePrices.sort(function (p1, p2) {
            return p1.symbol.localeCompare(p2.symbol);
        });

        return availablePrices;
    });

    return {
        // computed states
        cryptocurrencyPricesData,
        cryptocurrencyPricesDataUpdateTime,
        availableCryptocurrencyPrices
    };
}

