import { ref, computed } from 'vue';

import { useI18n } from '@/locales/helpers.ts';

import { useUserStore } from '@/stores/user.ts';
import { useCryptocurrencyPricesStore } from '@/stores/cryptocurrencyPrices.ts';
import { useExchangeRatesStore } from '@/stores/exchangeRates.ts';

import type {
    LatestCryptocurrencyPrice,
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

    const userStore = useUserStore();
    const cryptocurrencyPricesStore = useCryptocurrencyPricesStore();
    const exchangeRatesStore = useExchangeRatesStore();

    const baseSymbol = ref<string>('BTC');
    const baseAmount = ref<number>(1);

    const defaultCurrency = computed<string>(() => userStore.currentUserDefaultCurrency);
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

    function getConvertedAmount(baseAmount: number | '', fromPrice?: LatestCryptocurrencyPrice | LocalizedLatestCryptocurrencyPrice, toPrice?: LatestCryptocurrencyPrice | LocalizedLatestCryptocurrencyPrice): number | '' | null {
        if (!fromPrice || !toPrice) {
            return '';
        }

        if (baseAmount === '') {
            return 0;
        }

        const fromPriceNum = parseFloat(fromPrice.price);
        const toPriceNum = parseFloat(toPrice.price);

        if (isNaN(fromPriceNum) || isNaN(toPriceNum) || fromPriceNum === 0) {
            return null;
        }

        // Convert: baseAmount * (toPrice / fromPrice)
        return (baseAmount as number) * (toPriceNum / fromPriceNum);
    }

    function getConvertedAmountInFiat(baseAmount: number | '', fromPrice?: LatestCryptocurrencyPrice | LocalizedLatestCryptocurrencyPrice, fiatCurrency?: string): number | '' | null {
        if (!fromPrice || !fiatCurrency) {
            return '';
        }

        if (baseAmount === '') {
            return 0;
        }

        const fromPriceNum = parseFloat(fromPrice.price);
        if (isNaN(fromPriceNum) || fromPriceNum === 0) {
            return null;
        }

        // Convert crypto to USDT first
        const amountInUSDT = (baseAmount as number) * fromPriceNum;

        // Then convert USDT to fiat currency
        const fiatAmount = exchangeRatesStore.getExchangedAmount(amountInUSDT, 'USDT', fiatCurrency);
        return fiatAmount;
    }

    function setAsBaseline(symbol: string, amount: string): void {
        baseSymbol.value = symbol;
        const amountNum = parseFloat(amount);
        if (!isNaN(amountNum)) {
            baseAmount.value = amountNum;
        }
    }

    return {
        // states
        baseSymbol,
        baseAmount,
        // computed states
        defaultCurrency,
        cryptocurrencyPricesData,
        cryptocurrencyPricesDataUpdateTime,
        availableCryptocurrencyPrices,
        // functions
        getConvertedAmount,
        getConvertedAmountInFiat,
        setAsBaseline
    };
}

