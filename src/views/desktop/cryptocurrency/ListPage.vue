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
                                                    <span>{{ tt('Price') }}</span>
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

                                                    <span class="ms-3 text-sm">{{ formatCryptocurrencyPrice(price.symbol) }}</span>
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

import { ref, useTemplateRef, watch } from 'vue';
import { useDisplay } from 'vuetify';

import { useI18n } from '@/locales/helpers.ts';
import { useCryptocurrencyPricesPageBase } from '@/views/base/CryptocurrencyPricesPageBase.ts';

import { useCryptocurrencyPricesStore } from '@/stores/cryptocurrencyPrices.ts';


import {
    mdiRefresh,
    mdiMenu
} from '@mdi/js';

type SnackBarType = InstanceType<typeof SnackBar>;

const { mdAndUp } = useDisplay();

const { tt } = useI18n();
const {
    cryptocurrencyPricesData,
    cryptocurrencyPricesDataUpdateTime,
    availableCryptocurrencyPrices,
    formatCryptocurrencyPrice
} = useCryptocurrencyPricesPageBase();

const cryptocurrencyPricesStore = useCryptocurrencyPricesStore();

const snackbar = useTemplateRef<SnackBarType>('snackbar');

const activeTab = ref<string>('cryptocurrencyPricesPage');
const loading = ref<boolean>(true);
const alwaysShowNav = ref<boolean>(mdAndUp.value);
const showNav = ref<boolean>(mdAndUp.value);

function reload(force: boolean): void {
    loading.value = true;

    cryptocurrencyPricesStore.getLatestCryptocurrencyPrices({
        silent: false,
        force: force
    }).then(() => {
        loading.value = false;

        if (force) {
            snackbar.value?.showMessage('Cryptocurrency prices data has been updated');
        }
    }).catch(error => {
        loading.value = false;

        if (!error.processed) {
            snackbar.value?.showError(error);
        }
    });
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
</style>

