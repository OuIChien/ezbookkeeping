<template>
    <v-row class="match-height">
        <v-col cols="12">
            <v-card>
                <v-layout>
                    <v-navigation-drawer :permanent="alwaysShowNav" v-model="showNav">
                        <div class="mx-6 my-4">
                            <span class="text-subtitle-2">{{ tt('Data source') }}</span>
                            <p class="text-body-1 mt-1 mb-3">
                                <a tabindex="-1" target="_blank" :href="stockPricesData?.referenceUrl" v-if="!loading && stockPricesData && stockPricesData.referenceUrl">{{ stockPricesData.dataSource }}</a>
                                <span v-else-if="!loading && stockPricesData && !stockPricesData.referenceUrl">{{ stockPricesData.dataSource }}</span>
                                <span v-else-if="!loading && !stockPricesData">{{ tt('None') }}</span>
                                <span v-else-if="loading">
                                    <v-skeleton-loader class="skeleton-no-margin mt-3 mb-4" type="text" :loading="true"></v-skeleton-loader>
                                </span>
                            </p>
                            <span class="text-subtitle-2" v-if="stockPricesDataUpdateTime || loading">{{ tt('Last Updated') }}</span>
                            <p class="text-body-1 mt-1" v-if="stockPricesDataUpdateTime || loading">
                                <span v-if="!loading">{{ stockPricesDataUpdateTime }}</span>
                                <span v-if="loading">
                                    <v-skeleton-loader class="skeleton-no-margin mt-3 mb-4" type="text" :loading="true"></v-skeleton-loader>
                                </span>
                            </p>
                        </div>
                    </v-navigation-drawer>

                    <v-main>
                        <v-window class="d-flex flex-grow-1 disable-tab-transition w-100-window-container" v-model="activeTab">
                            <v-window-item value="stockPricesPage">
                                <v-card variant="flat" min-height="680">
                                    <template #title>
                                        <div class="title-and-toolbar d-flex align-center">
                                            <v-btn class="me-3 d-md-none" density="compact" color="default" variant="plain"
                                                   :ripple="false" :icon="true" @click="showNav = !showNav">
                                                <v-icon :icon="mdiMenu" size="24" />
                                            </v-btn>
                                            <span>{{ tt('Stock Prices Data') }}</span>
                                            <v-btn density="compact" color="default" variant="text" size="24"
                                                   class="ms-2" :icon="true" :loading="loading" @click="refreshStockPrices()">
                                                <template #loader>
                                                    <v-progress-circular indeterminate size="20"/>
                                                </template>
                                                <v-icon :icon="mdiRefresh" size="24" />
                                                <v-tooltip activator="parent">{{ tt('Refresh') }}</v-tooltip>
                                            </v-btn>
                                        </div>
                                    </template>

                                    <v-table class="exchange-rates-table table-striped" :hover="!loading">
                                        <thead>
                                        <tr>
                                            <th>
                                                <div class="d-flex align-center">
                                                    <span>{{ tt('Symbol') }}</span>
                                                    <v-spacer/>
                                                    <span>{{ tt('Price') }}</span>
                                                </div>
                                            </th>
                                        </tr>
                                        </thead>

                                        <tbody>
                                        <tr :key="itemIdx"
                                            v-for="itemIdx in (loading && (!stockPricesData || !stockPricesData.prices || stockPricesData.prices.length < 1) ? [ 1, 2, 3, 4, 5, 6, 7, 8, 9, 10, 11, 12 ] : [])">
                                            <td class="px-0">
                                                <v-skeleton-loader type="text" :loading="true"></v-skeleton-loader>
                                            </td>
                                        </tr>

                                        <tr v-if="!loading && (!stockPricesData || !stockPricesData.prices || !stockPricesData.prices.length)">
                                            <td>{{ tt('No stock prices available') }}</td>
                                        </tr>

                                        <tr class="exchange-rates-table-row-data" :key="price.symbol"
                                            v-for="price in (stockPricesData?.prices || [])">
                                            <td>
                                                <div class="d-flex align-center">
                                                    <span class="text-sm">{{ price.symbol }}</span>

                                                    <v-spacer/>

                                                    <span class="ms-3 text-mono">{{ formatNumberToWesternArabicNumerals(parseFloat(price.price), 2) }} {{ price.currency }}</span>
                                                </div>
                                            </td>
                                        </tr>
                                        </tbody>
                                    </v-table>
                                </v-card>
                            </v-window-item>
                        </v-window>
                    </v-main>
                </v-layout>
            </v-card>
        </v-col>
    </v-row>
</template>

<script setup lang="ts">
import { ref, computed, onMounted, watch } from 'vue';
import { useDisplay } from 'vuetify';
import { useI18n } from '@/locales/helpers.ts';

import type { LatestStockPriceResponse } from '@/models/stock_price.ts';

import { useStockPricesStore } from '@/stores/stockPrices.ts';

import { parseDateTimeFromUnixTimeWithBrowserTimezone } from '@/lib/datetime.ts';
import { getTimeZone } from '@/lib/settings.ts';

import {
    mdiRefresh,
    mdiMenu
} from '@mdi/js';

const { mdAndUp } = useDisplay();

const { tt, formatNumberToWesternArabicNumerals, formatDateTimeToLongDateTime } = useI18n();
const stockPricesStore = useStockPricesStore();

const activeTab = ref<string>('stockPricesPage');
const loading = ref(false);
const alwaysShowNav = ref<boolean>(mdAndUp.value);
const showNav = ref<boolean>(mdAndUp.value);

const stockPricesData = computed<LatestStockPriceResponse | undefined>(() => {
    return stockPricesStore.latestStockPrices?.data;
});

const stockPricesDataUpdateTime = computed<string>(() => {
    if (!stockPricesStore.latestStockPrices?.time) {
        return '';
    }

    const timezone = getTimeZone();
    let updateTime;
    
    if (timezone && timezone.trim().length > 0) {
        updateTime = parseDateTimeFromUnixTimeWithBrowserTimezone(stockPricesStore.latestStockPrices.time).setTimezoneByIANATimeZoneName(timezone);
    } else {
        updateTime = parseDateTimeFromUnixTimeWithBrowserTimezone(stockPricesStore.latestStockPrices.time);
    }
    
    return formatDateTimeToLongDateTime(updateTime);
});

async function refreshStockPrices(): Promise<void> {
    loading.value = true;
    try {
        await stockPricesStore.getLatestStockPrices({ silent: false, force: true });
    } finally {
        loading.value = false;
    }
}

watch(mdAndUp, (newValue) => {
    alwaysShowNav.value = newValue;

    if (!showNav.value) {
        showNav.value = newValue;
    }
});

onMounted(() => {
    if (!stockPricesStore.latestStockPrices?.data) {
        refreshStockPrices();
    }
});
</script>

<style>
.exchange-rates-table tr.exchange-rates-table-row-data .hover-display {
    display: none;
}

.exchange-rates-table tr.exchange-rates-table-row-data:hover .hover-display {
    display: grid;
}

.text-mono {
    font-family: 'JetBrains Mono', 'Fira Code', 'Source Code Pro', monospace;
}
</style>
