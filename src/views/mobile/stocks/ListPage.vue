<template>
    <f7-page ptr @ptr:refresh="reload">
        <f7-navbar>
            <f7-nav-left :back-link="tt('Back')"></f7-nav-left>
            <f7-nav-title :title="tt('Stock Prices Data')"></f7-nav-title>
            <f7-nav-right>
                <f7-link icon-f7="bars" @click="openPanel"></f7-link>
                <f7-link icon-f7="ellipsis" @click="showMoreActionSheet = true"></f7-link>
            </f7-nav-right>
        </f7-navbar>

        <f7-panel left reveal :opened="panelOpened" @panel:close="panelOpened = false">
            <f7-page>
                <f7-navbar>
                    <f7-nav-title>{{ tt('Information') }}</f7-nav-title>
                    <f7-nav-right>
                        <f7-link icon-f7="xmark" @click="closePanel"></f7-link>
                    </f7-nav-right>
                </f7-navbar>
                <f7-block class="no-margin-vertical">
                    <div class="panel-info-section">
                        <div class="info-item">
                            <span class="info-label">{{ tt('Data source') }}</span>
                            <p class="info-value">
                                <a :href="stockPricesData?.referenceUrl" target="_blank" v-if="!loading && stockPricesData && stockPricesData.referenceUrl">{{ stockPricesData.dataSource }}</a>
                                <span v-else-if="!loading && stockPricesData && !stockPricesData.referenceUrl">{{ stockPricesData.dataSource }}</span>
                                <span v-else-if="!loading && !stockPricesData">{{ tt('None') }}</span>
                                <span v-else>{{ tt('Loading...') }}</span>
                            </p>
                        </div>
                        <div class="info-item" v-if="stockPricesDataUpdateTime || loading">
                            <span class="info-label">{{ tt('Last Updated') }}</span>
                            <p class="info-value">
                                <span v-if="!loading">{{ stockPricesDataUpdateTime }}</span>
                                <span v-else>{{ tt('Loading...') }}</span>
                            </p>
                        </div>
                    </div>
                </f7-block>
            </f7-page>
        </f7-panel>

        <f7-list v-if="!loading && stockPricesData && stockPricesData.prices && stockPricesData.prices.length" class="no-margin-top">
            <f7-list-item
                v-for="price in stockPricesData.prices"
                :key="price.symbol"
                :title="price.symbol"
                :after="formatNumberToWesternArabicNumerals(parseFloat(price.price), 2) + ' ' + price.currency"
                class="stock-price-item">
                <template #media>
                    <div class="stock-icon">
                        <f7-icon f7="chart_bar_fill" size="20"></f7-icon>
                    </div>
                </template>
            </f7-list-item>
        </f7-list>

        <f7-block v-else-if="!loading" class="text-align-center no-margin-vertical">
            <div class="empty-state">
                <f7-icon f7="chart_bar_fill" size="64" class="empty-icon"></f7-icon>
                <div class="empty-title">{{ tt('No stock prices available') }}</div>
                <div class="empty-subtitle">{{ tt('Stock price data is not configured or failed to load.') }}</div>
            </div>
        </f7-block>

        <f7-block v-if="loading" class="text-align-center no-margin-vertical">
            <f7-preloader></f7-preloader>
            <div class="loading-text">{{ tt('Loading stock prices...') }}</div>
        </f7-block>

        <f7-actions close-by-outside-click close-on-escape :opened="showMoreActionSheet" @actions:closed="showMoreActionSheet = false">
            <f7-actions-group>
                <f7-actions-button :class="{ 'disabled': loading }" @click="reload(undefined)">
                    <span>{{ tt('Refresh') }}</span>
                </f7-actions-button>
            </f7-actions-group>
            <f7-actions-group>
                <f7-actions-button bold close>{{ tt('Cancel') }}</f7-actions-button>
            </f7-actions-group>
        </f7-actions>
    </f7-page>
</template>

<script setup lang="ts">
import { ref, computed } from 'vue';
import { useI18n } from '@/locales/helpers.ts';
import { useI18nUIComponents, showLoading, hideLoading } from '@/lib/ui/mobile.ts';

import type { LatestStockPriceResponse } from '@/models/stock_price.ts';

import { useStockPricesStore } from '@/stores/stockPrices.ts';

import { parseDateTimeFromUnixTimeWithBrowserTimezone } from '@/lib/datetime.ts';
import { getTimeZone } from '@/lib/settings.ts';

const { tt, formatNumberToWesternArabicNumerals, formatDateTimeToLongDateTime } = useI18n();
const { showToast } = useI18nUIComponents();
const stockPricesStore = useStockPricesStore();

const loading = ref(false);
const showMoreActionSheet = ref<boolean>(false);
const panelOpened = ref<boolean>(false);

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

function reload(done?: () => void): void {
    if (loading.value) {
        done?.();
        return;
    }

    loading.value = true;

    if (!done) {
        showLoading();
    }

    stockPricesStore.getLatestStockPrices({
        silent: false,
        force: true
    }).then(() => {
        done?.();

        loading.value = false;
        hideLoading();

        showToast('Stock prices data has been updated');
    }).catch(error => {
        done?.();

        loading.value = false;
        hideLoading();

        if (!error.processed) {
            showToast(error.message || error);
        }
    });
}

function openPanel(): void {
    panelOpened.value = true;
}

function closePanel(): void {
    panelOpened.value = false;
}

// Load data on page mount if not already loaded
if (!stockPricesStore.latestStockPrices?.data) {
    stockPricesStore.getLatestStockPrices({
        silent: true,
        force: false
    });
}
</script>

<style scoped>
.panel-info-section {
    padding: 16px;
}

.info-item {
    margin-bottom: 24px;
}

.info-item:last-child {
    margin-bottom: 0;
}

.info-label {
    display: block;
    font-size: 14px;
    font-weight: 500;
    color: var(--f7-text-color);
    margin-bottom: 8px;
}

.info-value {
    font-size: 16px;
    color: var(--f7-text-color-secondary);
    margin: 0;
}

.info-value a {
    color: var(--f7-link-color);
    text-decoration: none;
}

.stock-price-item {
    --f7-list-item-media-margin: 12px;
}

.stock-icon {
    width: 32px;
    height: 32px;
    border-radius: 50%;
    background: var(--f7-color-primary);
    display: flex;
    align-items: center;
    justify-content: center;
    color: white;
}

.empty-state {
    padding: 32px 16px;
}

.empty-icon {
    color: var(--f7-text-color-secondary);
    margin-bottom: 16px;
}

.empty-title {
    font-size: 18px;
    font-weight: 500;
    color: var(--f7-text-color);
    margin-bottom: 8px;
}

.empty-subtitle {
    font-size: 14px;
    color: var(--f7-text-color-secondary);
}

.loading-text {
    margin-top: 16px;
    color: var(--f7-text-color-secondary);
}
</style>
