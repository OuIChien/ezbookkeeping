export interface StockInfoResponse {
    symbol: string;
    name: string;
    market: string;
    isHidden: boolean;
    displayOrder: number;
}

export interface StockCreateRequest {
    symbol: string;
    name: string;
    market: string;
    displayOrder: number;
}

export interface StockModifyRequest {
    symbol: string;
    name: string;
    market: string;
    isHidden: boolean;
    displayOrder: number;
}

export interface StockHideRequest {
    symbol: string;
    hidden: boolean;
}

export interface StockDeleteRequest {
    symbol: string;
}
