<template>
    <f7-page @page:afterin="onPageAfterIn">
        <f7-navbar>
            <f7-nav-left :back-link="tt('Back')"></f7-nav-left>
            <f7-nav-title :title="tt('Stock Prices Settings')"></f7-nav-title>
            <f7-nav-right>
                <f7-link icon-f7="add" :class="{ 'disabled': loading }" @click="openAddPopup"></f7-link>
            </f7-nav-right>
        </f7-navbar>

        <f7-block-title>{{ tt('Data Source') }}</f7-block-title>
        <f7-list form strong inset dividers>
            <f7-list-item link="#" :title="tt('Data Source')" :after="getDataSourceName(dataSource)" @click="showDataSourcePopup = true">
                <list-item-selection-popup value-type="item"
                    key-field="value" value-field="value" title-field="name"
                    :title="tt('Data Source')"
                    :enable-filter="false"
                    :items="allStockDataSources"
                    v-model:show="showDataSourcePopup"
                    v-model="dataSource"
                    @update:model-value="saveConfig">
                </list-item-selection-popup>
            </f7-list-item>
            <f7-list-input v-if="dataSource === 'alphavantage' || dataSource === 'financial_modeling_prep'" type="text" clear-button :label="tt('API Key')" :placeholder="tt('API Key')" v-model:value="apiKey" @change="saveConfig"></f7-list-input>
            <f7-list-input type="text" clear-button :label="tt('Proxy')" :placeholder="tt('Proxy')" v-model:value="proxy" @change="saveConfig"></f7-list-input>
            <f7-list-input type="number" clear-button :label="tt('Request Timeout')" :placeholder="tt('Request Timeout')" v-model:value="requestTimeout" @change="saveConfig"></f7-list-input>
            <f7-list-input type="text" clear-button :label="tt('Update Frequency (minutes)')" :placeholder="tt('Update frequency placeholder')" v-model:value="updateFrequency" @change="saveConfig"></f7-list-input>
        </f7-list>

        <f7-block-title>{{ tt('Stock') }}</f7-block-title>
        <f7-list strong inset dividers v-if="allStocks.length">
            <f7-list-item v-for="stock in allStocks" :key="stock.symbol" link="#" @click="openEditPopup(stock)">
                <template #media>
                    <f7-icon f7="chart_line_fill" size="24"></f7-icon>
                </template>
                <template #title>
                    <span>{{ stock.name }}</span>
                    <span class="margin-horizontal-half text-color-gray">({{ stock.symbol }})</span>
                    <f7-badge v-if="stock.market" color="blue" class="margin-start">{{ stock.market }}</f7-badge>
                    <f7-badge v-if="stock.isHidden" color="gray" class="margin-start">{{ tt('Hidden') }}</f7-badge>
                </template>
                <template #after>
                    <f7-icon f7="chevron_right"></f7-icon>
                </template>
            </f7-list-item>
        </f7-list>
        <f7-list strong inset dividers v-else>
            <f7-list-item :title="tt('No stock')"></f7-list-item>
        </f7-list>

        <f7-popup class="stock-edit-popup" :opened="showEditPopup" @popup:closed="showEditPopup = false">
            <f7-page>
                <f7-navbar>
                    <f7-nav-left>
                        <f7-link @click="closeEditPopup">{{ tt('Cancel') }}</f7-link>
                    </f7-nav-left>
                    <f7-nav-title :title="editMode === 'add' ? tt('Add Stock') : tt('Edit Stock')"></f7-nav-title>
                    <f7-nav-right>
                        <f7-link :class="{ 'disabled': submitting || !editForm.symbol || !editForm.name }" @click="saveStock">{{ tt('Save') }}</f7-link>
                    </f7-nav-right>
                </f7-navbar>
                <f7-list form strong inset dividers class="margin-top">
                    <f7-list-input type="text" clear-button :label="tt('Symbol')" :placeholder="tt('Symbol')" v-model:value="editForm.symbol" :disabled="editMode === 'edit'"></f7-list-input>
                    <f7-list-input type="text" clear-button :label="tt('Name')" :placeholder="tt('Name')" v-model:value="editForm.name"></f7-list-input>
                    <f7-list-input type="text" clear-button :label="tt('Market')" :placeholder="tt('Market (Optional)')" v-model:value="editForm.market"></f7-list-input>
                    <f7-list-item v-if="editMode === 'edit'">
                        <template #title>{{ tt('Hidden') }}</template>
                        <template #after>
                            <f7-toggle :checked="editForm.isHidden" @toggle:change="editForm.isHidden = $event"></f7-toggle>
                        </template>
                    </f7-list-item>
                </f7-list>
            </f7-page>
        </f7-popup>
    </f7-page>
</template>

<script setup lang="ts">
import { ref, computed } from 'vue';
import { useI18n } from '@/locales/helpers.ts';
import { useI18nUIComponents } from '@/lib/ui/mobile.ts';
import { useStockPricesStore } from '@/stores/stockPrices.ts';
import { ExternalDataSourceType } from '@/models/external_data_source.ts';
import type { StockInfoResponse } from '@/models/stock.ts';

const { tt } = useI18n();
const { showToast } = useI18nUIComponents();
const stockPricesStore = useStockPricesStore();

const loading = ref(false);
const submitting = ref(false);
const showDataSourcePopup = ref(false);
const showEditPopup = ref(false);
const editMode = ref<'add' | 'edit'>('add');
const editForm = ref({ symbol: '', name: '', market: '', isHidden: false });

const allStockDataSources = [
    { name: 'Yahoo Finance', value: 'yahoo_finance' },
    { name: 'Alpha Vantage', value: 'alphavantage' },
    { name: 'Financial Modeling Prep', value: 'financial_modeling_prep' }
];

const dataSource = ref('yahoo_finance');
const apiKey = ref('');
const proxy = ref('');
const requestTimeout = ref(10000);
const updateFrequency = ref('');

const allStocks = computed(() => stockPricesStore.allStocks);

function getDataSourceName(value: string): string {
    const item = allStockDataSources.find(d => d.value === value);
    return item ? item.name : value;
}

function loadConfig(): void {
    const config = stockPricesStore.stockConfig;
    if (config) {
        dataSource.value = config.dataSource || 'yahoo_finance';
        apiKey.value = config.apiKey || '';
        proxy.value = config.proxy || '';
        requestTimeout.value = config.requestTimeout || 10000;
        updateFrequency.value = config.updateFrequency || '';
    }
}

function saveConfig(): void {
    const config = {
        type: ExternalDataSourceType.Stock,
        dataSource: dataSource.value,
        baseCurrency: stockPricesStore.stockConfig?.baseCurrency || 'USD',
        apiKey: apiKey.value,
        requestTimeout: typeof requestTimeout.value === 'number' ? requestTimeout.value : parseInt(String(requestTimeout.value), 10) || 10000,
        proxy: proxy.value,
        updateFrequency: updateFrequency.value
    };
    stockPricesStore.saveStockConfig(config).then(() => {
        showToast(tt('Settings saved'));
    }).catch(err => {
        showToast(err?.message || err);
    });
}

function openAddPopup(): void {
    editMode.value = 'add';
    editForm.value = { symbol: '', name: '', market: '', isHidden: false };
    showEditPopup.value = true;
}

function openEditPopup(stock: StockInfoResponse): void {
    editMode.value = 'edit';
    editForm.value = { symbol: stock.symbol, name: stock.name, market: stock.market || '', isHidden: stock.isHidden };
    showEditPopup.value = true;
}

function closeEditPopup(): void {
    showEditPopup.value = false;
}

function saveStock(): void {
    if (!editForm.value.symbol || !editForm.value.name || submitting.value) return;
    submitting.value = true;
    const isEdit = editMode.value === 'edit';
    const promise = isEdit
        ? stockPricesStore.modifyStock({
            symbol: editForm.value.symbol,
            name: editForm.value.name,
            market: editForm.value.market,
            isHidden: editForm.value.isHidden,
            displayOrder: 0
        })
        : stockPricesStore.addStock({
            symbol: editForm.value.symbol,
            name: editForm.value.name,
            market: editForm.value.market,
            displayOrder: allStocks.value.length + 1
        });
    promise.then(() => {
        submitting.value = false;
        showEditPopup.value = false;
        showToast(tt('Saved successfully'));
    }).catch(err => {
        submitting.value = false;
        showToast(err?.message || err);
    });
}

function onPageAfterIn(): void {
    loading.value = true;
    Promise.all([
        stockPricesStore.loadAllStocks({ force: false }),
        stockPricesStore.loadStockConfig()
    ]).then(() => {
        loadConfig();
        loading.value = false;
    }).catch(err => {
        loading.value = false;
        showToast(err?.message || err);
    });
}
</script>
