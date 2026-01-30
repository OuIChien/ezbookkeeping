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
                                    :items="allCryptocurrencyDataSources"
                                    v-model="dataSource"
                                />
                            </v-col>
                            <v-col cols="12" md="6">
                                <v-text-field
                                    type="text"
                                    persistent-placeholder
                                    :label="tt('API Key')"
                                    :placeholder="tt('API Key')"
                                    v-model="apiKey"
                                    @change="saveConfig"
                                />
                            </v-col>
                            <v-col cols="12" md="6">
                                <v-text-field
                                    type="text"
                                    persistent-placeholder
                                    :label="tt('Proxy')"
                                    :placeholder="tt('Proxy')"
                                    v-model="proxy"
                                    @change="saveConfig"
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
                                    @change="saveConfig"
                                />
                            </v-col>
                            <v-col cols="12" md="6">
                                <v-text-field
                                    type="text"
                                    persistent-placeholder
                                    :label="tt('Update Frequency (minutes)')"
                                    :placeholder="tt('Update frequency placeholder')"
                                    v-model="updateFrequency"
                                    @change="saveConfig"
                                />
                            </v-col>
                        </v-row>
                    </v-card-text>
                </v-form>
            </v-card>
        </v-col>

        <v-col cols="12">
            <v-card :title="tt('Cryptocurrencies')">
                <template #append>
                    <v-btn class="ml-4" :icon="true" density="comfortable" variant="text" :disabled="loading" @click="add">
                        <v-icon :icon="icons.plus" />
                        <v-tooltip activator="parent">{{ tt('Add') }}</v-tooltip>
                    </v-btn>
                </template>

                <v-list rounded density="comfortable" class="pa-0">
                    <template v-for="(crypto, index) in allCryptocurrencies" :key="crypto.symbol">
                        <v-list-item :class="{ 'disabled': crypto.isHidden }">
                            <template #prepend>
                                <v-avatar class="mr-2" size="40" rounded color="transparent">
                                    <v-icon size="24" :icon="icons.bitcoin" />
                                </v-avatar>
                            </template>

                            <v-list-item-title>
                                <div class="d-flex align-center">
                                    <span>{{ crypto.name }}</span>
                                    <span class="ml-2 text-caption text-medium-emphasis">({{ crypto.symbol }})</span>
                                    <v-chip class="ml-2" size="x-small" v-if="crypto.isHidden">{{ tt('Hidden') }}</v-chip>
                                </div>
                            </v-list-item-title>

                            <template #append>
                                <v-btn class="ml-1" color="default" variant="text" density="comfortable" :icon="true" :disabled="loading" @click="edit(crypto)">
                                    <v-icon :icon="icons.pencil" />
                                    <v-tooltip activator="parent">{{ tt('Edit') }}</v-tooltip>
                                </v-btn>
                                <v-btn class="ml-1" color="default" variant="text" density="comfortable" :icon="true" :disabled="loading" v-if="!crypto.isHidden" @click="hide(crypto, true)">
                                    <v-icon :icon="icons.eyeOff" />
                                    <v-tooltip activator="parent">{{ tt('Hide') }}</v-tooltip>
                                </v-btn>
                                <v-btn class="ml-1" color="default" variant="text" density="comfortable" :icon="true" :disabled="loading" v-if="crypto.isHidden" @click="hide(crypto, false)">
                                    <v-icon :icon="icons.eye" />
                                    <v-tooltip activator="parent">{{ tt('Show') }}</v-tooltip>
                                </v-btn>
                                <v-btn class="ml-1" color="error" variant="text" density="comfortable" :icon="true" :disabled="loading" @click="remove(crypto)">
                                    <v-icon :icon="icons.delete" />
                                    <v-tooltip activator="parent">{{ tt('Delete') }}</v-tooltip>
                                </v-btn>
                            </template>
                        </v-list-item>
                        <v-divider v-if="index < allCryptocurrencies.length - 1" />
                    </template>

                    <div class="d-flex align-center justify-center pa-4" v-if="allCryptocurrencies.length < 1">
                        {{ tt('No cryptocurrency') }}
                    </div>
                </v-list>
            </v-card>
        </v-col>
    </v-row>

    <cryptocurrency-edit-dialog ref="editDialog" v-model="showEditDialog" @save="save" />

    <confirm-dialog ref="confirmDialog" />
    <snack-bar ref="snackbar" />
</template>

<script setup lang="ts">
import { ref, computed, useTemplateRef, watch } from 'vue';
import { useI18n } from '@/locales/helpers.ts';
import { useCryptocurrencyPricesStore } from '@/stores/cryptocurrencyPrices.ts';

import CryptocurrencyEditDialog from '@/views/desktop/app/settings/dialogs/CryptocurrencyEditDialog.vue';
import ConfirmDialog from '@/components/desktop/ConfirmDialog.vue';
import SnackBar from '@/components/desktop/SnackBar.vue';

import {
    mdiPlus,
    mdiPencil,
    mdiEye,
    mdiEyeOff,
    mdiDelete,
    mdiBitcoin
} from '@mdi/js';

import { ExternalDataSourceType } from '@/models/external_data_source.ts';
import type { CryptocurrencyInfoResponse } from '@/models/cryptocurrency.ts';

type CryptocurrencyEditDialogType = InstanceType<typeof CryptocurrencyEditDialog>;
type ConfirmDialogType = InstanceType<typeof ConfirmDialog>;
type SnackBarType = InstanceType<typeof SnackBar>;

const { tt } = useI18n();
const cryptocurrencyPricesStore = useCryptocurrencyPricesStore();

const editDialog = useTemplateRef<CryptocurrencyEditDialogType>('editDialog');
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
    bitcoin: mdiBitcoin
};

const allCryptocurrencyDataSources = [
    { name: 'CoinGecko', value: 'coingecko' }
];

const allCryptocurrencies = computed(() => cryptocurrencyPricesStore.allCryptocurrencies);

const apiKey = ref('');
const proxy = ref('');
const requestTimeout = ref(10000);
const updateFrequency = ref('');

watch(() => cryptocurrencyPricesStore.cryptocurrencyConfig, (newConfig) => {
    if (newConfig) {
        apiKey.value = newConfig.apiKey || '';
        proxy.value = newConfig.proxy || '';
        requestTimeout.value = newConfig.requestTimeout || 10000;
        updateFrequency.value = newConfig.updateFrequency || '';
    }
});

const dataSource = computed<string>({
    get: () => {
        return cryptocurrencyPricesStore.cryptocurrencyConfig?.dataSource || 'coingecko';
    },
    set: (value) => {
        if (value === cryptocurrencyPricesStore.cryptocurrencyConfig?.dataSource) {
            return;
        }
        saveConfig(value);
    }
});

function saveConfig(newDataSource?: string) {
    const config = {
        type: ExternalDataSourceType.Cryptocurrency,
        dataSource: newDataSource ?? cryptocurrencyPricesStore.cryptocurrencyConfig?.dataSource ?? 'coingecko',
        baseCurrency: cryptocurrencyPricesStore.cryptocurrencyConfig?.baseCurrency || 'USD',
        apiKey: apiKey.value,
        requestTimeout: requestTimeout.value || 10000,
        proxy: proxy.value || '',
        updateFrequency: updateFrequency.value || ''
    };

    cryptocurrencyPricesStore.saveCryptocurrencyConfig(config).then(() => {
        snackbar.value?.showMessage(tt('Data source saved'));
    }).catch(error => {
        snackbar.value?.showError(error);
    });
}

function init() {
    loading.value = true;
    Promise.all([
        cryptocurrencyPricesStore.loadAllCryptocurrencies({ force: false }),
        cryptocurrencyPricesStore.loadCryptocurrencyConfig()
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

function edit(crypto: CryptocurrencyInfoResponse) {
    showEditDialog.value = true;
    editDialog.value?.open(crypto);
}

function save(item: { symbol: string, name: string, isHidden: boolean, isEdit: boolean }) {
    editDialog.value?.setSubmitting(true);

    let promise;
    if (item.isEdit) {
        promise = cryptocurrencyPricesStore.modifyCryptocurrency({
            symbol: item.symbol,
            name: item.name,
            isHidden: item.isHidden,
            displayOrder: 0 // Keep order
        });
    } else {
        promise = cryptocurrencyPricesStore.addCryptocurrency({
            symbol: item.symbol,
            name: item.name,
            displayOrder: allCryptocurrencies.value.length + 1
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

function hide(crypto: CryptocurrencyInfoResponse, hidden: boolean) {
    loading.value = true;
    cryptocurrencyPricesStore.hideCryptocurrency({
        symbol: crypto.symbol,
        hidden: hidden
    }).then(() => {
        loading.value = false;
    }).catch(error => {
        loading.value = false;
        snackbar.value?.showError(error);
    });
}

function remove(crypto: CryptocurrencyInfoResponse) {
    confirmDialog.value?.open(tt('Are you sure you want to delete {name}?', { name: crypto.name })).then(() => {
        loading.value = true;
        cryptocurrencyPricesStore.deleteCryptocurrency({
            symbol: crypto.symbol
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
