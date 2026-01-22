export interface LatestStockPriceResponse {
    dataSource: string;
    referenceUrl: string;
    updateTime: number;
    baseCurrency: string;
    prices: LatestStockPrice[];
}

export interface LatestStockPrice {
    symbol: string;
    price: string;
    currency: string;
}
