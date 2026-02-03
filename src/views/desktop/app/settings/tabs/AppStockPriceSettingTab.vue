<template>
    <v-row>
        <v-col cols="12">
            <v-card :title="tt('Data Source')">
                <v-form>
                    <v-card-text>
                        <v-row>
                            <v-col cols="12" md="6">
                                <v-select
                                    item-title="name"
                                    item-value="value"
                                    persistent-placeholder
                                    :label="tt('Data Source')"
                                    :placeholder="tt('Data Source')"
                                    :items="allStockDataSources"
                                    v-model="dataSource"
                                />
                            </v-col>
                            <v-col cols="12" md="6" v-if="dataSource === 'alphavantage' || dataSource === 'financial_modeling_prep'">
                                <v-text-field
                                    type="text"
                                    persistent-placeholder
                                    :label="tt('API Key')"
                                    :placeholder="tt('API Key')"
                                    v-model="apiKey"
                                    @change="() => saveConfig()"
                                />
                            </v-col>
                            <v-col cols="12" md="6">
                                <v-text-field
                                    type="text"
                                    persistent-placeholder
                                    :label="tt('Proxy')"
                                    :placeholder="tt('Proxy')"
                                    v-model="proxy"
                                    @change="() => saveConfig()"
                                />
                            </v-col>
                            <v-col cols="12" md="6">
                                <v-text-field
                                    type="number"
                                    persistent-placeholder
                                    :label="tt('Request Timeout')"
                                    :placeholder="tt('Request Timeout')"
                                    v-model.number="requestTimeout"
                                    min="1000"
                                    step="1000"
                                    @change="() => saveConfig()"
                                />
                            </v-col>
                            <v-col cols="12" md="6">
                                <v-text-field
                                    type="text"
                                    persistent-placeholder
                                    :label="tt('Update Frequency (minutes)')"
                                    :placeholder="tt('Update frequency placeholder')"
                                    v-model="updateFrequency"
                                    @change="() => saveConfig()"
                                />
                            </v-col>
                        </v-row>
                    </v-card-text>
                </v-form>
            </v-card>
        </v-col>

        <v-col cols="12">
            <v-card :title="tt('Stocks')">
                <template #append>
                    <v-btn class="ml-4" :icon="true" density="comfortable" variant="text" :disabled="loading" @click="add">
                        <v-icon :icon="icons.plus" />
                        <v-tooltip activator="parent">{{ tt('Add') }}</v-tooltip>
                    </v-btn>
                </template>

                <v-list rounded density="comfortable" class="pa-0">
                    <template v-for="(stock, index) in allStocks" :key="stock.symbol">
                        <v-list-item :class="{ 'disabled': stock.isHidden }">
                            <template #prepend>
                                <v-avatar class="mr-2" size="40" rounded color="transparent">
                                    <v-icon size="24" :icon="icons.chartLine" />
                                </v-avatar>
                            </template>

                            <v-list-item-title>
                                <div class="d-flex align-center">
                                    <span>{{ stock.name }}</span>
                                    <span class="ml-2 text-caption text-medium-emphasis">({{ stock.symbol }})</span>
                                    <v-chip class="ml-2" size="x-small" v-if="stock.market">{{ stock.market }}</v-chip>
                                    <v-chip class="ml-2" size="x-small" v-if="stock.isHidden">{{ tt('Hidden') }}</v-chip>
                                </div>
                            </v-list-item-title>

                            <template #append>
                                <v-btn class="ml-1" color="default" variant="text" density="comfortable" :icon="true" :disabled="loading" @click="edit(stock)">
                                    <v-icon :icon="icons.pencil" />
                                    <v-tooltip activator="parent">{{ tt('Edit') }}</v-tooltip>
                                </v-btn>
                                <v-btn class="ml-1" color="default" variant="text" density="comfortable" :icon="true" :disabled="loading" v-if="!stock.isHidden" @click="hide(stock, true)">
                                    <v-icon :icon="icons.eyeOff" />
                                    <v-tooltip activator="parent">{{ tt('Hide') }}</v-tooltip>
                                </v-btn>
                                <v-btn class="ml-1" color="default" variant="text" density="comfortable" :icon="true" :disabled="loading" v-if="stock.isHidden" @click="hide(stock, false)">
                                    <v-icon :icon="icons.eye" />
                                    <v-tooltip activator="parent">{{ tt('Show') }}</v-tooltip>
                                </v-btn>
                                <v-btn class="ml-1" color="error" variant="text" density="comfortable" :icon="true" :disabled="loading" @click="remove(stock)">
                                    <v-icon :icon="icons.delete" />
                                    <v-tooltip activator="parent">{{ tt('Delete') }}</v-tooltip>
                                </v-btn>
                            </template>
                        </v-list-item>
                        <v-divider v-if="index < allStocks.length - 1" />
                    </template>

                    <div class="d-flex align-center justify-center pa-4" v-if="allStocks.length < 1">
                        {{ tt('No stock') }}
                    </div>
                </v-list>
            </v-card>
        </v-col>
    </v-row>

    <stock-edit-dialog ref="editDialog" v-model="showEditDialog" @save="save" />

    <confirm-dialog ref="confirmDialog" />
    <snack-bar ref="snackbar" />
</template>

<script setup lang="ts">
import { ref, computed, useTemplateRef, watch } from 'vue';
import { useI18n } from '@/locales/helpers.ts';
import { useStockPricesStore } from '@/stores/stockPrices.ts';

import StockEditDialog from '@/views/desktop/app/settings/dialogs/StockEditDialog.vue';
import ConfirmDialog from '@/components/desktop/ConfirmDialog.vue';
import SnackBar from '@/components/desktop/SnackBar.vue';

import {
    mdiPlus,
    mdiPencil,
    mdiEye,
    mdiEyeOff,
    mdiDelete,
    mdiChartLine
} from '@mdi/js';

import { ExternalDataSourceType } from '@/models/external_data_source.ts';
import type { StockInfoResponse } from '@/models/stock.ts';

type StockEditDialogType = InstanceType<typeof StockEditDialog>;
type ConfirmDialogType = InstanceType<typeof ConfirmDialog>;
type SnackBarType = InstanceType<typeof SnackBar>;

const { tt } = useI18n();
const stockPricesStore = useStockPricesStore();

const editDialog = useTemplateRef<StockEditDialogType>('editDialog');
const confirmDialog = useTemplateRef<ConfirmDialogType>('confirmDialog');
const snackbar = useTemplateRef<SnackBarType>('snackbar');

const loading = ref(false);
const showEditDialog = ref(false);

const icons = {
    plus: mdiPlus,
    pencil: mdiPencil,
    eye: mdiEye,
    eyeOff: mdiEyeOff,
    delete: mdiDelete,
    chartLine: mdiChartLine
};

const allStockDataSources = [
    { name: 'Yahoo Finance', value: 'yahoo_finance' },
    { name: 'Alpha Vantage', value: 'alphavantage' },
    { name: 'Financial Modeling Prep', value: 'financial_modeling_prep' }
];

const allStocks = computed(() => stockPricesStore.allStocks);

const dataSource = computed<string>({
    get: () => {
        return stockPricesStore.stockConfig?.dataSource || 'yahoo_finance';
    },
    set: (value) => {
        if (value === stockPricesStore.stockConfig?.dataSource) {
            return;
        }
        saveConfig(value);
    }
});

const apiKey = ref('');
const proxy = ref('');
const requestTimeout = ref(10000);
const updateFrequency = ref('');

watch(() => stockPricesStore.stockConfig, (newConfig) => {
    if (newConfig) {
        apiKey.value = newConfig.apiKey || '';
        proxy.value = newConfig.proxy || '';
        requestTimeout.value = newConfig.requestTimeout || 10000;
        updateFrequency.value = newConfig.updateFrequency || '';
    }
});

function saveConfig(newDataSource?: string) {
    const dataSourceValue = typeof newDataSource === 'string'
        ? newDataSource
        : dataSource.value;
    const config = {
        type: ExternalDataSourceType.Stock,
        dataSource: dataSourceValue,
        baseCurrency: stockPricesStore.stockConfig?.baseCurrency || 'USD',
        apiKey: apiKey.value,
        requestTimeout: requestTimeout.value || 10000,
        proxy: proxy.value || '',
        updateFrequency: updateFrequency.value || ''
    };

    stockPricesStore.saveStockConfig(config).then(() => {
        snackbar.value?.showMessage(tt('Settings saved'));
    }).catch(error => {
        snackbar.value?.showError(error);
    });
}

function init() {
    loading.value = true;
    Promise.all([
        stockPricesStore.loadAllStocks({ force: false }),
        stockPricesStore.loadStockConfig()
    ]).then(() => {
        loading.value = false;
    }).catch(error => {
        loading.value = false;
        snackbar.value?.showError(error);
    });
}

function add() {
    showEditDialog.value = true;
    editDialog.value?.open();
}

function edit(stock: StockInfoResponse) {
    showEditDialog.value = true;
    editDialog.value?.open(stock);
}

function save(item: { symbol: string, name: string, market: string, isHidden: boolean, isEdit: boolean }) {
    editDialog.value?.setSubmitting(true);

    let promise;
    if (item.isEdit) {
        promise = stockPricesStore.modifyStock({
            symbol: item.symbol,
            name: item.name,
            market: item.market,
            isHidden: item.isHidden,
            displayOrder: 0 // Keep order
        });
    } else {
        promise = stockPricesStore.addStock({
            symbol: item.symbol,
            name: item.name,
            market: item.market,
            displayOrder: allStocks.value.length + 1
        });
    }

    promise.then(() => {
        editDialog.value?.setSubmitting(false);
        editDialog.value?.close();
        showEditDialog.value = false;
        snackbar.value?.showMessage(tt('Saved successfully'));
    }).catch(error => {
        editDialog.value?.setSubmitting(false);
        snackbar.value?.showError(error);
    });
}

function hide(stock: StockInfoResponse, hidden: boolean) {
    loading.value = true;
    stockPricesStore.hideStock({
        symbol: stock.symbol,
        hidden: hidden
    }).then(() => {
        loading.value = false;
    }).catch(error => {
        loading.value = false;
        snackbar.value?.showError(error);
    });
}

function remove(stock: StockInfoResponse) {
    confirmDialog.value?.open(tt('Are you sure you want to delete {name}?', { name: stock.name })).then(() => {
        loading.value = true;
        stockPricesStore.deleteStock({
            symbol: stock.symbol
        }).then(() => {
            loading.value = false;
            snackbar.value?.showMessage(tt('Deleted successfully'));
        }).catch(error => {
            loading.value = false;
            snackbar.value?.showError(error);
        });
    });
}

init();
</script>
