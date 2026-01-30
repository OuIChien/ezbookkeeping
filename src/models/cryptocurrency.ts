export interface CryptocurrencyInfoResponse {
    symbol: string;
    name: string;
    isHidden: boolean;
    displayOrder: number;
}

export interface CryptocurrencyCreateRequest {
    symbol: string;
    name: string;
    displayOrder: number;
}

export interface CryptocurrencyModifyRequest {
    symbol: string;
    name: string;
    isHidden: boolean;
    displayOrder: number;
}

export interface CryptocurrencyHideRequest {
    symbol: string;
    hidden: boolean;
}

export interface CryptocurrencyDeleteRequest {
    symbol: string;
}
