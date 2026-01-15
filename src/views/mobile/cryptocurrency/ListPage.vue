<template>
    <f7-page ptr @ptr:refresh="reload">
        <f7-navbar>
            <f7-nav-left :back-link="tt('Back')"></f7-nav-left>
            <f7-nav-title :title="tt('Cryptocurrency Prices Data')"></f7-nav-title>
            <f7-nav-right>
                <f7-link icon-f7="ellipsis" @click="showMoreActionSheet = true"></f7-link>
            </f7-nav-right>
        </f7-navbar>

        <f7-block class="no-margin-vertical">
            <div class="data-source-info" v-if="cryptocurrencyPricesData">
                <div class="info-row">
                    <span class="label">{{ tt('Data source') }}:</span>
                    <span class="value">
                        <a :href="cryptocurrencyPricesData.referenceUrl" target="_blank" v-if="cryptocurrencyPricesData.referenceUrl">{{ cryptocurrencyPricesData.dataSource }}</a>
                        <span v-else>{{ cryptocurrencyPricesData.dataSource }}</span>
                    </span>
                </div>
                <div class="info-row" v-if="cryptocurrencyPricesDataUpdateTime">
                    <span class="label">{{ tt('Last Updated') }}:</span>
                    <span class="value">{{ cryptocurrencyPricesDataUpdateTime }}</span>
                </div>
            </div>
        </f7-block>

        <f7-block class="no-margin-vertical">
            <div class="base-amount-section">
                <div class="amount-input-container">
                    <span class="amount-label">{{ tt('Base Amount') }} (USD):</span>
                    <amount-input
                        class="amount-input"
                        :currency="'USD'"
                        :disabled="loading || !cryptocurrencyPricesData || !cryptocurrencyPricesData.prices || !cryptocurrencyPricesData.prices.length"
                        v-model="baseAmount"/>
                </div>
            </div>
        </f7-block>

        <f7-list v-if="!loading && cryptocurrencyPricesData && cryptocurrencyPricesData.prices && cryptocurrencyPricesData.prices.length" class="no-margin-top">
            <f7-list-item
                v-for="price in cryptocurrencyPricesData.prices"
                :key="price.symbol"
                :title="price.symbol"
                :after="formatNumberToWesternArabicNumerals(parseFloat(price.price), 2)"
                class="cryptocurrency-price-item">
                <template #media>
                    <div class="crypto-icon">
                        <f7-icon f7="bitcoin" size="20"></f7-icon>
                    </div>
                </template>
                <template #subtitle v-if="baseAmount && baseAmount > 0">
                    {{ tt('Equivalent') }}: {{ formatNumberToWesternArabicNumerals(parseFloat((baseAmount / parseFloat(price.price)).toFixed(8)), 8) }}
                </template>
            </f7-list-item>
        </f7-list>

        <f7-block v-else-if="!loading" class="text-align-center no-margin-vertical">
            <div class="empty-state">
                <f7-icon f7="bitcoin" size="64" class="empty-icon"></f7-icon>
                <div class="empty-title">{{ tt('No cryptocurrency prices available') }}</div>
                <div class="empty-subtitle">{{ tt('Cryptocurrency price data is not configured or failed to load.') }}</div>
            </div>
        </f7-block>

        <f7-block v-if="loading" class="text-align-center no-margin-vertical">
            <f7-preloader></f7-preloader>
            <div class="loading-text">{{ tt('Loading cryptocurrency prices...') }}</div>
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

import type { LatestCryptocurrencyPriceResponse } from '@/models/cryptocurrency_price.ts';

import { useCryptocurrencyPricesStore } from '@/stores/cryptocurrencyPrices.ts';

import { parseDateTimeFromUnixTimeWithBrowserTimezone } from '@/lib/datetime.ts';
import { getTimeZone } from '@/lib/settings.ts';

const { tt, formatNumberToWesternArabicNumerals, formatDateTimeToLongDateTime } = useI18n();
const { showToast } = useI18nUIComponents();
const cryptocurrencyPricesStore = useCryptocurrencyPricesStore();

const loading = ref(false);
const baseAmount = ref(1);
const showMoreActionSheet = ref<boolean>(false);

const cryptocurrencyPricesData = computed<LatestCryptocurrencyPriceResponse | undefined>(() => {
    return cryptocurrencyPricesStore.latestCryptocurrencyPrices?.data;
});

const cryptocurrencyPricesDataUpdateTime = computed<string>(() => {
    if (!cryptocurrencyPricesStore.latestCryptocurrencyPrices?.time) {
        return '';
    }

    const timezone = getTimeZone();
    let updateTime;
    
    if (timezone && timezone.trim().length > 0) {
        updateTime = parseDateTimeFromUnixTimeWithBrowserTimezone(cryptocurrencyPricesStore.latestCryptocurrencyPrices.time).setTimezoneByIANATimeZoneName(timezone);
    } else {
        updateTime = parseDateTimeFromUnixTimeWithBrowserTimezone(cryptocurrencyPricesStore.latestCryptocurrencyPrices.time);
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

    cryptocurrencyPricesStore.getLatestCryptocurrencyPrices({
        silent: false,
        force: true
    }).then(() => {
        done?.();

        loading.value = false;
        hideLoading();

        showToast('Cryptocurrency prices data has been updated');
    }).catch(error => {
        done?.();

        loading.value = false;
        hideLoading();

        if (!error.processed) {
            showToast(error.message || error);
        }
    });
}

// Load data on page mount if not already loaded
if (!cryptocurrencyPricesStore.latestCryptocurrencyPrices?.data) {
    cryptocurrencyPricesStore.getLatestCryptocurrencyPrices({
        silent: true,
        force: false
    });
}
</script>

<style scoped>
.data-source-info {
    background: var(--f7-block-bg-color);
    padding: 16px;
    border-radius: 8px;
    margin: 16px;
}

.info-row {
    display: flex;
    justify-content: space-between;
    align-items: center;
    margin-bottom: 8px;
}

.info-row:last-child {
    margin-bottom: 0;
}

.label {
    font-weight: 500;
    color: var(--f7-text-color);
}

.value {
    color: var(--f7-text-color-secondary);
}

.value a {
    color: var(--f7-link-color);
    text-decoration: none;
}

.base-amount-section {
    background: var(--f7-block-bg-color);
    padding: 16px;
    border-radius: 8px;
    margin: 16px;
}

.amount-input-container {
    display: flex;
    flex-direction: column;
    gap: 8px;
}

.amount-label {
    font-weight: 500;
    color: var(--f7-text-color);
}

.amount-input {
    flex: 1;
}

.cryptocurrency-price-item {
    --f7-list-item-media-margin: 12px;
}

.crypto-icon {
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