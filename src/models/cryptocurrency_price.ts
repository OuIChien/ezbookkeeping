export interface LatestCryptocurrencyPriceResponse {
    dataSource: string;
    referenceUrl: string;
    updateTime: number;
    baseCurrency: string;
    prices: LatestCryptocurrencyPrice[];
}

export interface LatestCryptocurrencyPrice {
    symbol: string;
    price: string;
}
