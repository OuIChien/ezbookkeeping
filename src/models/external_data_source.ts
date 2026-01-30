export enum ExternalDataSourceType {
    Cryptocurrency = 1,
    Stock = 2,
    ExchangeRate = 3
}

export interface ExternalDataSourceConfigResponse {
    type: ExternalDataSourceType;
    dataSource: string;
    baseCurrency: string;
    apiKey: string;
    requestTimeout: number;
    proxy: string;
    updateFrequency: string;
}

export interface ExternalDataSourceConfigSaveRequest {
    type: ExternalDataSourceType;
    dataSource: string;
    baseCurrency: string;
    apiKey: string;
    requestTimeout: number;
    proxy: string;
    updateFrequency: string;
}
