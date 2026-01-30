<template>
    <v-dialog :width="600" v-model="showState">
        <v-card class="pa-2 pa-sm-4 pa-md-8">
            <template #title>
                <div class="d-flex align-center justify-center">
                    <div class="d-flex align-center">
                        <v-icon class="mr-2" :icon="icons.pencil" v-if="mode === 'edit'" />
                        <v-icon class="mr-2" :icon="icons.plus" v-else />
                        <span>{{ title }}</span>
                    </div>
                </div>
            </template>
            <v-card-text>
                <v-form ref="form">
                    <v-row>
                        <v-col cols="12" md="6">
                            <v-text-field
                                type="text"
                                persistent-placeholder
                                :disabled="mode === 'edit'"
                                :label="tt('Symbol')"
                                :placeholder="tt('Symbol')"
                                v-model="crypto.symbol"
                                :rules="[
                                    (v) => !!v || tt('Symbol is required'),
                                    (v) => (v && v.length <= 20) || tt('Symbol must be less than 20 characters')
                                ]"
                            />
                        </v-col>
                        <v-col cols="12" md="6">
                            <v-text-field
                                type="text"
                                persistent-placeholder
                                :label="tt('Name')"
                                :placeholder="tt('Name')"
                                v-model="crypto.name"
                                :rules="[
                                    (v) => !!v || tt('Name is required'),
                                    (v) => (v && v.length <= 100) || tt('Name must be less than 100 characters')
                                ]"
                            />
                        </v-col>
                        <v-col cols="12" v-if="mode === 'edit'">
                            <v-checkbox
                                :label="tt('Hidden')"
                                v-model="crypto.isHidden"
                            />
                        </v-col>
                    </v-row>
                </v-form>
            </v-card-text>
            <v-card-actions>
                <v-spacer />
                <v-btn :disabled="submitting" @click="cancel">{{ tt('Cancel') }}</v-btn>
                <v-btn color="primary" :disabled="submitting" :loading="submitting" @click="confirm">{{ tt('OK') }}</v-btn>
            </v-card-actions>
        </v-card>
    </v-dialog>
</template>

<script setup lang="ts">
import { ref, computed } from 'vue';
import { useI18n } from '@/locales/helpers.ts';
import { mdiPencil, mdiPlus } from '@mdi/js';

import type { CryptocurrencyInfoResponse } from '@/models/cryptocurrency.ts';

const props = defineProps<{
    modelValue: boolean;
}>();

const emit = defineEmits<{
    (e: 'update:modelValue', value: boolean): void;
    (e: 'save', item: { symbol: string, name: string, isHidden: boolean, isEdit: boolean }): void;
}>();

const { tt } = useI18n();

const showState = computed({
    get: () => props.modelValue,
    set: (value) => emit('update:modelValue', value)
});

const mode = ref<'add' | 'edit'>('add');
const submitting = ref(false);
const crypto = ref({
    symbol: '',
    name: '',
    isHidden: false
});

const title = computed(() => {
    return mode.value === 'edit' ? tt('Edit Cryptocurrency') : tt('Add Cryptocurrency');
});

const icons = {
    pencil: mdiPencil,
    plus: mdiPlus
};

function open(item?: CryptocurrencyInfoResponse) {
    if (item) {
        mode.value = 'edit';
        crypto.value = {
            symbol: item.symbol,
            name: item.name,
            isHidden: item.isHidden
        };
    } else {
        mode.value = 'add';
        crypto.value = {
            symbol: '',
            name: '',
            isHidden: false
        };
    }
    showState.value = true;
}

function cancel() {
    showState.value = false;
}

function confirm() {
    if (!crypto.value.symbol || !crypto.value.name) {
        return;
    }

    emit('save', {
        symbol: crypto.value.symbol,
        name: crypto.value.name,
        isHidden: crypto.value.isHidden,
        isEdit: mode.value === 'edit'
    });
}

function setSubmitting(val: boolean) {
    submitting.value = val;
}

function close() {
    showState.value = false;
}

defineExpose({
    open,
    setSubmitting,
    close
});
</script>
