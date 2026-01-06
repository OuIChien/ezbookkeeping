<template>
    <v-row class="match-height">
        <v-col cols="12">
            <v-card>
                <v-layout>
                    <v-navigation-drawer :permanent="alwaysShowNav" v-model="showNav">
                        <div class="mx-6 my-4">
                            <span class="text-subtitle-2">{{ tt('Data source') }}</span>
                            <p class="text-body-1 mt-1 mb-3">
                                <a tabindex="-1" target="_blank" :href="cryptocurrencyPricesData.referenceUrl" v-if="!loading && cryptocurrencyPricesData && cryptocurrencyPricesData.referenceUrl">{{ cryptocurrencyPricesData.dataSource }}</a>
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
                            <span class="text-subtitle-2 mt-4">{{ tt('Base Currency') }}</span>
                            <p class="text-body-1 mt-1 mb-3">
                                <span v-if="!loading && cryptocurrencyPricesData">{{ cryptocurrencyPricesData.baseCurrency }}</span>
                                <span v-else-if="loading">
                                    <v-skeleton-loader class="skeleton-no-margin mt-3 mb-4" type="text" :loading="true"></v-skeleton-loader>
                                </span>
                            </p>
                        </div>
                        <v-divider />
                        <div class="mx-6 mt-4">
                            <span class="text-subtitle-2">{{ tt('Base Amount') }}</span>
                            <amount-input class="mt-2" density="compact"
                                          :currency="baseSymbol"
                                          :disabled="loading || !cryptocurrencyPricesData || !cryptocurrencyPricesData.prices || !cryptocurrencyPricesData.prices.length"
                                          v-model="baseAmount"/>
                        </div>
                        <div class="mx-6 mt-4">
                            <span class="text-subtitle-2">{{ tt('Base Cryptocurrency') }}</span>
                        </div>
                        <v-tabs show-arrows class="mb-4" direction="vertical"
                                :disabled="loading" v-model="baseSymbol"
                                v-if="cryptocurrencyPricesData && cryptocurrencyPricesData.prices && cryptocurrencyPricesData.prices.length">
                            <v-tab class="tab-text-truncate" :key="price.symbol" :value="price.symbol"
                                   v-for="price in availableCryptocurrencyPrices">
                                <div class="d-flex w-100">
                                    <span class="d-block text-truncate">{{ price.symbolDisplayName }}</span>
                                    <small class="smaller ms-1">{{ price.symbol }}</small>
                                </div>
                            </v-tab>
                        </v-tabs>
                        <div class="mx-6 mt-2 mb-4"
                             v-else-if="!cryptocurrencyPricesData || !cryptocurrencyPricesData.prices || !cryptocurrencyPricesData.prices.length">
                            <span v-if="!loading">{{ tt('None') }}</span>
                            <span v-else-if="loading">
                                <v-skeleton-loader class="skeleton-no-margin pt-2 pb-5" type="text"
                                                   :key="itemIdx" :loading="loading"
                                                   v-for="itemIdx in [ 1, 2, 3, 4, 5, 6, 7, 8, 9, 10 ]"></v-skeleton-loader>
                            </span>
                        </div>
                    </v-navigation-drawer>
                    <v-main>
                        <v-window class="d-flex flex-grow-1 disable-tab-transition w-100-window-container" v-model="activeTab">
                            <v-window-item value="cryptocurrencyPricesPage">
                                <v-card variant="flat" min-height="680">
                                    <template #title>
                                        <div class="title-and-toolbar d-flex align-center">
                                            <v-btn class="me-3 d-md-none" density="compact" color="default" variant="plain"
                                                   :ripple="false" :icon="true" @click="showNav = !showNav">
                                                <v-icon :icon="mdiMenu" size="24" />
                                            </v-btn>
                                            <span>{{ tt('Cryptocurrency Prices Data') }}</span>
                                            <v-btn density="compact" color="default" variant="text" size="24"
                                                   class="ms-2" :icon="true" :loading="loading" @click="reload(true)">
                                                <template #loader>
                                                    <v-progress-circular indeterminate size="20"/>
                                                </template>
                                                <v-icon :icon="mdiRefresh" size="24" />
                                                <v-tooltip activator="parent">{{ tt('Refresh') }}</v-tooltip>
                                            </v-btn>
                                        </div>
                                    </template>

                                    <v-table class="cryptocurrency-prices-table table-striped" :hover="!loading">
                                        <thead>
                                        <tr>
                                            <th>
                                                <div class="d-flex align-center">
                                                    <span>{{ tt('Cryptocurrency') }}</span>
                                                    <v-spacer/>
                                                    <span>{{ tt('Price in USDT') }}</span>
                                                    <v-spacer class="mx-2"/>
                                                    <span>{{ tt('Amount') }}</span>
                                                </div>
                                            </th>
                                        </tr>
                                        </thead>

                                        <tbody>
                                        <tr :key="itemIdx"
                                            v-for="itemIdx in (loading && (!cryptocurrencyPricesData || !cryptocurrencyPricesData.prices || cryptocurrencyPricesData.prices.length < 1) ? [ 1, 2, 3, 4, 5, 6, 7, 8, 9, 10, 11, 12 ] : [])">
                                            <td class="px-0">
                                                <v-skeleton-loader type="text" :loading="true"></v-skeleton-loader>
                                            </td>
                                        </tr>

                                        <tr v-if="!loading && (!cryptocurrencyPricesData || !cryptocurrencyPricesData.prices || !cryptocurrencyPricesData.prices.length)">
                                            <td>{{ tt('No cryptocurrency prices data') }}</td>
                                        </tr>

                                        <tr class="cryptocurrency-prices-table-row-data" :key="price.symbol"
                                            v-for="price in availableCryptocurrencyPrices">
                                            <td>
                                                <div class="d-flex align-center">
                                                    <span class="text-sm">{{ price.symbolDisplayName }}</span>
                                                    <span class="text-caption ms-1">{{ price.symbol }}</span>

                                                    <v-spacer/>

                                                    <span class="ms-3 text-sm">{{ formatPrice(price.price) }}</span>

                                                    <v-btn class="px-2 ms-2" color="default"
                                                           density="comfortable" variant="text"
                                                           :class="{ 'd-none': loading, 'hover-display': !loading }"
                                                           v-if="price.symbol !== baseSymbol"
                                                           @click="setAsBaseline(price.symbol, getFinalConvertedAmount(price, false))">
                                                        {{ tt('Set as Base') }}
                                                    </v-btn>
                                                    <span class="ms-3">{{ getFinalConvertedAmount(price, true) }}</span>
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

    <snack-bar ref="snackbar" />
</template>

<script setup lang="ts">
import SnackBar from '@/components/desktop/SnackBar.vue';

import { ref, computed, useTemplateRef, watch } from 'vue';
import { useDisplay } from 'vuetify';

import { useI18n } from '@/locales/helpers.ts';
import { useCryptocurrencyPricesPageBase } from '@/views/base/CryptocurrencyPricesPageBase.ts';

import { useCryptocurrencyPricesStore } from '@/stores/cryptocurrencyPrices.ts';

import { NumeralSystem } from '@/core/numeral.ts';

import type { LocalizedLatestCryptocurrencyPrice } from '@/views/base/CryptocurrencyPricesPageBase.ts';

import logger from '@/lib/logger.ts';

import {
    mdiRefresh,
    mdiMenu
} from '@mdi/js';

type SnackBarType = InstanceType<typeof SnackBar>;

const { mdAndUp } = useDisplay();

const { tt, getCurrentNumeralSystemType, formatExchangeRateAmountToWesternArabicNumerals } = useI18n();
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

const snackbar = useTemplateRef<SnackBarType>('snackbar');

const activeTab = ref<string>('cryptocurrencyPricesPage');
const loading = ref<boolean>(true);
const alwaysShowNav = ref<boolean>(mdAndUp.value);
const showNav = ref<boolean>(mdAndUp.value);

const numeralSystem = computed<NumeralSystem>(() => getCurrentNumeralSystemType());

function reload(force: boolean): void {
    loading.value = true;

    cryptocurrencyPricesStore.getLatestCryptocurrencyPrices({
        silent: false,
        force: force
    }).then(() => {
        loading.value = false;

        if (cryptocurrencyPricesData.value && cryptocurrencyPricesData.value.prices) {
            const prices = cryptocurrencyPricesData.value.prices;
            let foundBaseSymbol = false;

            for (const price of prices) {
                if (price.symbol === baseSymbol.value) {
                    foundBaseSymbol = true;
                    break;
                }
            }

            if (force) {
                snackbar.value?.showMessage('Cryptocurrency prices data has been updated');
            } else if (!foundBaseSymbol && prices.length > 0) {
                const firstPrice = prices[0];
                if (firstPrice) {
                    baseSymbol.value = firstPrice.symbol;
                }
            }
        }
    }).catch(error => {
        loading.value = false;

        if (!error.processed) {
            snackbar.value?.showError(error);
        }
    });
}

function getFinalConvertedAmount(toPrice: LocalizedLatestCryptocurrencyPrice, displayLocalizedDigits: boolean): string {
    if (!baseSymbol.value) {
        if (displayLocalizedDigits) {
            return numeralSystem.value.digitZero;
        } else {
            return NumeralSystem.WesternArabicNumerals.digitZero;
        }
    }

    const fromPrice = cryptocurrencyPricesStore.latestCryptocurrencyPriceMap[baseSymbol.value];
    let convertedAmount: number | '' | null = 0;

    try {
        convertedAmount = getConvertedAmount(baseAmount.value, fromPrice, toPrice);
    } catch (ex) {
        convertedAmount = 0;
        logger.warn('failed to convert amount by cryptocurrency prices, original base amount is ' + baseAmount.value, ex)
    }

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

watch(mdAndUp, (newValue) => {
    alwaysShowNav.value = newValue;

    if (!showNav.value) {
        showNav.value = newValue;
    }
});

reload(false);
</script>

<style>
.cryptocurrency-prices-table tr.cryptocurrency-prices-table-row-data .hover-display {
    display: none;
}

.cryptocurrency-prices-table tr.cryptocurrency-prices-table-row-data:hover .hover-display {
    display: grid;
}

</style>

