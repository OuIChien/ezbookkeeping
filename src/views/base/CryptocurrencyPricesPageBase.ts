import { computed } from 'vue';

import { useI18n } from '@/locales/helpers.ts';

import { useCryptocurrencyPricesStore } from '@/stores/cryptocurrencyPrices.ts';
import { useUserStore } from '@/stores/user.ts';

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
    const { formatDateTimeToLongDate, formatAmountToLocalizedNumeralsWithCurrency, getCurrencyName } = useI18n();

    const cryptocurrencyPricesStore = useCryptocurrencyPricesStore();
    const userStore = useUserStore();

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

    function formatCryptocurrencyPrice(symbol: string): string {
        const defaultCurrency = userStore.currentUserDefaultCurrency;

        // Convert float to amount (amount is stored as integer in cents: float * 100)
        function floatToAmount(floatValue: number): number {
            return Math.round(floatValue * 100);
        }

        // Get price in USD (prices are already in USD from the API)
        const priceInUSD = cryptocurrencyPricesStore.getCryptocurrencyPriceInFiat(symbol, 'USD');
        if (priceInUSD === null) {
            return '';
        }
        
        // Get price in default currency
        const priceInDefaultCurrency = cryptocurrencyPricesStore.getCryptocurrencyPriceInFiat(symbol, defaultCurrency);

        // If USD and default currency are the same, only show USD
        if (defaultCurrency === 'USD') {
            return formatAmountToLocalizedNumeralsWithCurrency(floatToAmount(priceInUSD), 'USD');
        }

        // Show default currency on the left, USD price in parentheses on the right
        const defaultCurrencyPrice = priceInDefaultCurrency !== null 
            ? formatAmountToLocalizedNumeralsWithCurrency(floatToAmount(priceInDefaultCurrency), defaultCurrency)
            : '';
        const usdPrice = formatAmountToLocalizedNumeralsWithCurrency(floatToAmount(priceInUSD), 'USD');

        if (defaultCurrencyPrice) {
            return `${defaultCurrencyPrice} (${usdPrice})`;
        }

        // Fallback: if default currency price is not available, show USD only
        return usdPrice;
    }

    return {
        // computed states
        cryptocurrencyPricesData,
        cryptocurrencyPricesDataUpdateTime,
        availableCryptocurrencyPrices,
        // functions
        formatCryptocurrencyPrice
    };
}

