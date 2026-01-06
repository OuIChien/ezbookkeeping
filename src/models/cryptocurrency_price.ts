export interface LatestCryptocurrencyPriceResponse {
    readonly dataSource: string;
    readonly referenceUrl: string;
    readonly updateTime: number;
    readonly baseCurrency: string;
    readonly prices: readonly LatestCryptocurrencyPrice[];
}

export interface LatestCryptocurrencyPrice {
    readonly symbol: string;
    readonly price: string;
}

