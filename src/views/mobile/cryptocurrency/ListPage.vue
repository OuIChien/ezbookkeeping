<template>
    <f7-page ptr @ptr:refresh="reload">
        <f7-navbar>
            <f7-nav-left :back-link="tt('Back')"></f7-nav-left>
            <f7-nav-title :title="tt('Cryptocurrency Prices Data')"></f7-nav-title>
            <f7-nav-right>
                <f7-link icon-f7="ellipsis" @click="showMoreActionSheet = true"></f7-link>
            </f7-nav-right>
        </f7-navbar>

        <f7-list strong inset dividers class="margin-vertical" v-if="cryptocurrencyPricesData && cryptocurrencyPricesData.prices && cryptocurrencyPricesData.prices.length">
            <f7-list-item
                class="list-item-with-header-and-title list-item-no-item-after"
                link="#"
                :header="tt('Base Cryptocurrency')"
                @click="showBaseSymbolPopup = true"
            >
                <template #title>
                    <div class="no-padding no-margin">
                        <span>{{ getCurrencyName(baseSymbol) }}&nbsp;</span>
                        <small class="smaller">{{ baseSymbol }}</small>
                    </div>
                </template>
                <list-item-selection-popup value-type="item"
                                           key-field="symbol" value-field="symbol"
                                           title-field="symbolDisplayName" after-field="symbol"
                                           :title="tt('Base Cryptocurrency')"
                                           :enable-filter="true"
                                           :filter-placeholder="tt('Cryptocurrency')"
                                           :filter-no-items-text="tt('No results')"
                                           :items="availableCryptocurrencyPrices"
                                           v-model:show="showBaseSymbolPopup"
                                           v-model="baseSymbol">
                </list-item-selection-popup>
            </f7-list-item>
            <f7-list-item
                class="cryptocurrency-base-amount"
                link="#" no-chevron
                :class="baseAmountFontSizeClass"
                :header="tt('Base Amount')"
                :title="displayBaseAmount"
                @click="showBaseAmountSheet = true"
            >
                <number-pad-sheet :min-value="TRANSACTION_MIN_AMOUNT"
                                  :max-value="TRANSACTION_MAX_AMOUNT"
                                  :currency="baseSymbol"
                                  v-model:show="showBaseAmountSheet"
                                  v-model="baseAmount"
                ></number-pad-sheet>
            </f7-list-item>
        </f7-list>

        <f7-list strong inset dividers class="margin-vertical" v-if="!cryptocurrencyPricesData || !cryptocurrencyPricesData.prices || !cryptocurrencyPricesData.prices.length">
            <f7-list-item :title="tt('No cryptocurrency prices data')"></f7-list-item>
        </f7-list>

        <f7-list strong inset dividers class="margin-vertical" v-if="cryptocurrencyPricesData && cryptocurrencyPricesData.prices && cryptocurrencyPricesData.prices.length">
            <f7-list-item swipeout
                          :id="getCryptocurrencyPriceDomId(price)"
                          :after="getFinalConvertedAmount(price, true)"
                          :key="baseSymbolChangedTime + '_' + price.symbol" v-for="price in availableCryptocurrencyPrices"
                          @swipeout:closed="onCryptocurrencyPriceSwipeoutClosed()">
                <template #title>
                    <div class="no-padding no-margin">
                        <span style="margin-inline-end: 5px">{{ price.symbolDisplayName }}</span>
                        <small class="smaller">{{ price.symbol }}</small>
                    </div>
                </template>
                <template #after>
                    <div class="d-flex align-center">
                        <span class="me-2 text-sm">{{ formatPrice(price.price) }}</span>
                        <span>{{ getFinalConvertedAmount(price, true) }}</span>
                    </div>
                </template>
                <f7-swipeout-actions :left="textDirection === TextDirection.RTL"
                                     :right="textDirection === TextDirection.LTR"
                                     v-if="price.symbol !== baseSymbol">
                    <f7-swipeout-button color="primary" close
                                        :text="tt('Set as Base')"
                                        :class="{ 'disabled': price.symbol === baseSymbol }"
                                        @click="setAsBaseline(price.symbol, getFinalConvertedAmount(price, false)); settingBaseLine = true"
                                        v-if="settingBaseLine || price.symbol !== baseSymbol"></f7-swipeout-button>
                </f7-swipeout-actions>
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

import { TextDirection } from '@/core/text.ts';
import { NumeralSystem } from '@/core/numeral.ts';
import { TRANSACTION_MIN_AMOUNT, TRANSACTION_MAX_AMOUNT } from '@/consts/transaction.ts';

import type { LocalizedLatestCryptocurrencyPrice } from '@/views/base/CryptocurrencyPricesPageBase.ts';

import {
    getCurrentUnixTime
} from '@/lib/datetime.ts';

const props = defineProps<{
    f7router: Router.Router;
}>();

const {
    tt,
    getCurrentLanguageTextDirection,
    getCurrentNumeralSystemType,
    getCurrencyName,
    formatAmountToLocalizedNumerals,
    formatExchangeRateAmountToWesternArabicNumerals
} = useI18n();

const { showToast, openExternalUrl } = useI18nUIComponents();

const {
    baseSymbol,
    baseAmount,
    cryptocurrencyPricesData,
    cryptocurrencyPricesDataUpdateTime,
    availableCryptocurrencyPrices,
    getConvertedAmount,
    setAsBaseline
} = useCryptocurrencyPricesPageBase();

const cryptocurrencyPricesStore = useCryptocurrencyPricesStore();

const loading = ref<boolean>(false);
const baseSymbolChangedTime = ref<number>(getCurrentUnixTime());
const settingBaseLine = ref<boolean>(false);
const showMoreActionSheet = ref<boolean>(false);
const showBaseSymbolPopup = ref<boolean>(false);
const showBaseAmountSheet = ref<boolean>(false);

const textDirection = computed<TextDirection>(() => getCurrentLanguageTextDirection());
const numeralSystem = computed<NumeralSystem>(() => getCurrentNumeralSystemType());
const displayBaseAmount = computed<string>(() => formatAmountToLocalizedNumerals(baseAmount.value, baseSymbol.value));
const baseAmountFontSizeClass = computed<string>(() => {
    if (baseAmount.value >= 100000000 || baseAmount.value <= -100000000) {
        return 'ebk-small-amount';
    } else if (baseAmount.value >= 1000000 || baseAmount.value <= -1000000) {
        return 'ebk-normal-amount';
    } else {
        return 'ebk-large-amount';
    }
});

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

        if (cryptocurrencyPricesData.value && cryptocurrencyPricesData.value.prices) {
            const prices = cryptocurrencyPricesData.value.prices;
            let foundBaseSymbol = false;

            for (const price of prices) {
                if (price.symbol === baseSymbol.value) {
                    foundBaseSymbol = true;
                    break;
                }
            }

            if (!foundBaseSymbol && prices.length > 0) {
                const firstPrice = prices[0];
                if (firstPrice) {
                    baseSymbol.value = firstPrice.symbol;
                }
            }
        }

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

function getFinalConvertedAmount(toPrice: LocalizedLatestCryptocurrencyPrice, displayLocalizedDigits: boolean): string {
    const fromPrice = cryptocurrencyPricesStore.latestCryptocurrencyPriceMap[baseSymbol.value];
    const convertedAmount = getConvertedAmount(baseAmount.value, fromPrice, toPrice);

    if (!convertedAmount) {
        if (displayLocalizedDigits) {
            return numeralSystem.value.digitZero;
        } else {
            return NumeralSystem.WesternArabicNumerals.digitZero;
        }
    }

    let ret = formatExchangeRateAmountToWesternArabicNumerals(convertedAmount);

    if (displayLocalizedDigits) {
        ret = numeralSystem.value.replaceWesternArabicDigitsToLocalizedDigits(ret);
    }

    return ret;
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

function onCryptocurrencyPriceSwipeoutClosed(): void {
    baseSymbolChangedTime.value = getCurrentUnixTime();
    settingBaseLine.value = false;
}

cryptocurrencyPricesStore.getLatestCryptocurrencyPrices({
    silent: true,
    force: false
}).then(() => {
    if (cryptocurrencyPricesData.value && cryptocurrencyPricesData.value.prices) {
        const prices = cryptocurrencyPricesData.value.prices;
        let hasBaseSymbol = false;

        for (const price of prices) {
            if (price.symbol === baseSymbol.value) {
                hasBaseSymbol = true;
                break;
            }
        }

        if (!hasBaseSymbol && prices.length > 0) {
            const firstPrice = prices[0];
            if (firstPrice) {
                baseSymbol.value = firstPrice.symbol;
            }
        }
    }
});
</script>

<style>
.cryptocurrency-base-amount {
    line-height: 53px;
}

.cryptocurrency-base-amount .item-header {
    padding-top: calc(var(--f7-typography-padding) / 2);
}
</style>

