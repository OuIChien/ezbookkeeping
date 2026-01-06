<template>
    <f7-page ptr @ptr:refresh="reload">
        <f7-navbar>
            <f7-nav-left :back-link="tt('Back')"></f7-nav-left>
            <f7-nav-title :title="tt('Cryptocurrency Prices Data')"></f7-nav-title>
            <f7-nav-right>
                <f7-link icon-f7="ellipsis" @click="showMoreActionSheet = true"></f7-link>
            </f7-nav-right>
        </f7-navbar>


        <f7-list strong inset dividers class="margin-vertical" v-if="!cryptocurrencyPricesData || !cryptocurrencyPricesData.prices || !cryptocurrencyPricesData.prices.length">
            <f7-list-item :title="tt('No cryptocurrency prices data')"></f7-list-item>
        </f7-list>

        <f7-list strong inset dividers class="margin-vertical" v-if="cryptocurrencyPricesData && cryptocurrencyPricesData.prices && cryptocurrencyPricesData.prices.length">
            <f7-list-item
                          :id="getCryptocurrencyPriceDomId(price)"
                          :key="price.symbol" v-for="price in availableCryptocurrencyPrices">
                <template #title>
                    <div class="no-padding no-margin">
                        <span style="margin-inline-end: 5px">{{ price.symbolDisplayName }}</span>
                        <small class="smaller">{{ price.symbol }}</small>
                    </div>
                </template>
                <template #after>
                    <span class="text-sm">{{ formatPrice(price.price) }}</span>
                </template>
            </f7-list-item>
        </f7-list>

        <f7-list strong inset dividers class="margin-vertical" v-if="cryptocurrencyPricesData && cryptocurrencyPricesData.prices && cryptocurrencyPricesData.prices.length">
            <f7-list-item v-if="cryptocurrencyPricesDataUpdateTime">
                <small>{{ tt('Last Updated') }}</small>
                <small>{{ cryptocurrencyPricesDataUpdateTime }}</small>
            </f7-list-item>
            <f7-list-item>
                <small>{{ tt('Data source') }}</small>
                <small>
                    <f7-link @click="openExternalUrl(cryptocurrencyPricesData.referenceUrl)" v-if="cryptocurrencyPricesData.referenceUrl">{{ cryptocurrencyPricesData.dataSource }}</f7-link>
                    <span v-else>{{ cryptocurrencyPricesData.dataSource }}</span>
                </small>
            </f7-list-item>
            <f7-list-item>
                <small>{{ tt('Base Currency') }}</small>
                <small>{{ cryptocurrencyPricesData.baseCurrency }}</small>
            </f7-list-item>
        </f7-list>

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
import type { Router } from 'framework7/types';

import { useI18n } from '@/locales/helpers.ts';
import { useI18nUIComponents, showLoading, hideLoading } from '@/lib/ui/mobile.ts';
import { useCryptocurrencyPricesPageBase } from '@/views/base/CryptocurrencyPricesPageBase.ts';

import { useCryptocurrencyPricesStore } from '@/stores/cryptocurrencyPrices.ts';

import { NumeralSystem } from '@/core/numeral.ts';

import type { LocalizedLatestCryptocurrencyPrice } from '@/views/base/CryptocurrencyPricesPageBase.ts';

const props = defineProps<{
    f7router: Router.Router;
}>();

const {
    tt,
    getCurrentNumeralSystemType,
    formatExchangeRateAmountToWesternArabicNumerals
} = useI18n();

const { showToast, openExternalUrl } = useI18nUIComponents();

const {
    cryptocurrencyPricesData,
    cryptocurrencyPricesDataUpdateTime,
    availableCryptocurrencyPrices
} = useCryptocurrencyPricesPageBase();

const cryptocurrencyPricesStore = useCryptocurrencyPricesStore();

const loading = ref<boolean>(false);
const showMoreActionSheet = ref<boolean>(false);

const numeralSystem = computed<NumeralSystem>(() => getCurrentNumeralSystemType());

function getCryptocurrencyPriceDomId(price: LocalizedLatestCryptocurrencyPrice): string {
    return 'cryptocurrencyPrice_' + price.symbol;
}

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

function formatPrice(price: string): string {
    const priceNum = parseFloat(price);
    if (isNaN(priceNum)) {
        return '0';
    }

    let ret = formatExchangeRateAmountToWesternArabicNumerals(priceNum);
    ret = numeralSystem.value.replaceWesternArabicDigitsToLocalizedDigits(ret);
    return ret;
}

cryptocurrencyPricesStore.getLatestCryptocurrencyPrices({
    silent: true,
    force: false
});
</script>

<style>
</style>

