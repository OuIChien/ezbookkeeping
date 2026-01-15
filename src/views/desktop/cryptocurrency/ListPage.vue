<template>
    <v-row class="match-height">
        <v-col cols="12">
            <v-card>
                <v-layout>
                    <v-navigation-drawer :permanent="alwaysShowNav" v-model="showNav">
                        <div class="mx-6 my-4">
                            <span class="text-subtitle-2">{{ tt('Data source') }}</span>
                            <p class="text-body-1 mt-1 mb-3">
                                <a tabindex="-1" target="_blank" :href="cryptocurrencyPricesData?.referenceUrl" v-if="!loading && cryptocurrencyPricesData && cryptocurrencyPricesData.referenceUrl">{{ cryptocurrencyPricesData.dataSource }}</a>
                                <span v-else-if="!loading && cryptocurrencyPricesData && !cryptocurrencyPricesData.referenceUrl">{{ cryptocurrencyPricesData.dataSource }}</span>
                                <span v-else-if="!loading && !cryptocurrencyPricesData">{{ tt('None') }}</span>
                                <span v-else-if="loading">
                                    <v-skeleton-loader class="skeleton-no-margin mt-3 mb-4" type="text" :loading="true"></v-skeleton-loader>
                                </span>
                            </p>
                            <span class="text-subtitle-2" v-if="cryptocurrencyPricesDataUpdateTime || loading">{{ tt('Last Updated') }}</span>
                            <p class="text-body-1 mt-1" v-if="cryptocurrencyPricesDataUpdateTime || loading">
                                <span v-if="!loading">{{ cryptocurrencyPricesDataUpdateTime }}</span>
                                <span v-if="loading">
                                    <v-skeleton-loader class="skeleton-no-margin mt-3 mb-4" type="text" :loading="true"></v-skeleton-loader>
                                </span>
                            </p>
                        </div>
                        <v-divider />
                        <div class="mx-6 mt-4">
                            <span class="text-subtitle-2">{{ tt('Base Currency') }}</span>
                            <p class="text-body-1 mt-1 mb-3">USD ({{ tt('US Dollar') }})</p>
                        </div>
                    </v-navigation-drawer>

                    <v-main>
                        <div class="px-6 py-4">
                            <div class="d-flex align-center mb-4">
                                <h4 class="text-h4 me-4">{{ tt('Cryptocurrency Prices') }}</h4>
                                <v-spacer />
                                <v-btn :icon="mdiRefresh" variant="text" size="small"
                                       :loading="loading"
                                       :disabled="loading"
                                       @click="refreshCryptocurrencyPrices()">
                                    <v-tooltip activator="parent">{{ tt('Refresh') }}</v-tooltip>
                                </v-btn>
                            </div>

                            <div class="cryptocurrency-prices-table">
                                <v-data-table
                                    :headers="headers"
                                    :items="displayItems"
                                    :loading="loading"
                                    :items-per-page="-1"
                                    density="compact"
                                    item-key="symbol"
                                    class="elevation-1"
                                    hide-default-footer>
                                    <template #item.symbol="{ item }">
                                        <div class="d-flex align-center">
                                            <span class="font-weight-medium">{{ item.symbol }}</span>
                                        </div>
                                    </template>

                                    <template #item.price="{ item }">
                                        <span class="text-mono">{{ formatNumberToWesternArabicNumerals(parseFloat(item.price), 2) }}</span>
                                    </template>

                                    <template #no-data>
                                        <div class="text-center py-4">
                                            <v-icon size="64" class="mb-4 text-disabled">mdi-bitcoin</v-icon>
                                            <div class="text-h6 mb-2">{{ tt('No cryptocurrency prices available') }}</div>
                                            <div class="text-body-2 text-disabled">{{ tt('Cryptocurrency price data is not configured or failed to load.') }}</div>
                                        </div>
                                    </template>
                                </v-data-table>
                            </div>
                        </div>
                    </v-main>
                </v-layout>
            </v-card>
        </v-col>
    </v-row>
</template>

<script setup lang="ts">
import { ref, computed, onMounted } from 'vue';
import { useI18n } from '@/locales/helpers.ts';

import type { LatestCryptocurrencyPriceResponse } from '@/models/cryptocurrency_price.ts';

import { useCryptocurrencyPricesStore } from '@/stores/cryptocurrencyPrices.ts';

import { mdiRefresh } from '@mdi/js';

const { tt, formatNumberToWesternArabicNumerals } = useI18n();
const cryptocurrencyPricesStore = useCryptocurrencyPricesStore();

const loading = ref(false);
const showNav = ref(true);

const alwaysShowNav = computed(() => true);

const cryptocurrencyPricesData = computed<LatestCryptocurrencyPriceResponse | undefined>(() => {
    return cryptocurrencyPricesStore.latestCryptocurrencyPrices?.data;
});

const cryptocurrencyPricesDataUpdateTime = computed<string>(() => {
    if (!cryptocurrencyPricesStore.latestCryptocurrencyPrices?.time) {
        return '';
    }

    return formatDateTime(cryptocurrencyPricesStore.latestCryptocurrencyPrices.time);
});

const headers = computed(() => [
    {
        title: tt('Symbol'),
        key: 'symbol',
        sortable: false
    },
    {
        title: tt('Price (USD)'),
        key: 'price',
        sortable: false
    }
]);

const displayItems = computed(() => {
    if (!cryptocurrencyPricesData.value?.prices) {
        return [];
    }

    return cryptocurrencyPricesData.value.prices.map(price => ({
        symbol: price.symbol,
        price: price.price
    }));
});

function formatDateTime(timestamp: number): string {
    const date = new Date(timestamp * 1000);
    return date.toLocaleString();
}

async function refreshCryptocurrencyPrices(): Promise<void> {
    loading.value = true;
    try {
        await cryptocurrencyPricesStore.getLatestCryptocurrencyPrices({ silent: false, force: true });
    } finally {
        loading.value = false;
    }
}

onMounted(() => {
    if (!cryptocurrencyPricesStore.latestCryptocurrencyPrices?.data) {
        refreshCryptocurrencyPrices();
    }
});
</script>

<style scoped>
.cryptocurrency-prices-table {
    border-radius: 8px;
    overflow: hidden;
}

.text-mono {
    font-family: 'JetBrains Mono', 'Fira Code', 'Source Code Pro', monospace;
}
</style>
