import { ref, computed } from 'vue';

import { useI18n } from '@/locales/helpers.ts';

import { useSettingsStore } from '@/stores/setting.ts';
import { useUserStore } from '@/stores/user.ts';
import { useAccountsStore } from '@/stores/account.ts';

import type { HiddenAmount, NumberWithSuffix } from '@/core/numeral.ts';
import type { WeekDayValue } from '@/core/datetime.ts';
import { AccountCategory, AccountType, AccountAssetType as AccountAssetTypeClass } from '@/core/account.ts';
import type { Account, CategorizedAccount } from '@/models/account.ts';

import { isObject, isNumber, isString } from '@/lib/common.ts';

export function useAccountListPageBase() {
    const { formatAmountToLocalizedNumeralsWithCurrency } = useI18n();

    const settingsStore = useSettingsStore();
    const userStore = useUserStore();
    const accountsStore = useAccountsStore();

    const loading = ref<boolean>(true);
    const showHidden = ref<boolean>(false);
    const displayOrderModified = ref<boolean>(false);

    const showAccountBalance = computed<boolean>({
        get: () => settingsStore.appSettings.showAccountBalance,
        set: (value) => settingsStore.setShowAccountBalance(value)
    });

    const customAccountCategoryOrder = computed<string>(() => settingsStore.appSettings.accountCategoryOrders);
    const defaultAccountCategory = computed<AccountCategory>(() => AccountCategory.values(customAccountCategoryOrder.value)[0] ?? AccountCategory.Default);

    const firstDayOfWeek = computed<WeekDayValue>(() => userStore.currentUserFirstDayOfWeek);
    const fiscalYearStart = computed<number>(() => userStore.currentUserFiscalYearStart);
    const defaultCurrency = computed<string>(() => userStore.currentUserDefaultCurrency);

    const allAccounts = computed<Account[]>(() => accountsStore.allAccounts);
    const allCategorizedAccountsMap = computed<Record<number, CategorizedAccount>>(() => accountsStore.allCategorizedAccountsMap);
    const allAccountCount = computed<number>(() => accountsStore.allAvailableAccountsCount);

    const netAssets = computed<string>(() => {
        const netAssets: number | HiddenAmount | NumberWithSuffix = accountsStore.getNetAssets(showAccountBalance.value);
        return formatAmountToLocalizedNumeralsWithCurrency(netAssets, defaultCurrency.value);
    });

    const totalAssets = computed<string>(() => {
        const totalAssets: number | HiddenAmount | NumberWithSuffix = accountsStore.getTotalAssets(showAccountBalance.value);
        return formatAmountToLocalizedNumeralsWithCurrency(totalAssets, defaultCurrency.value);
    });

    const totalLiabilities = computed<string>(() => {
        const totalLiabilities: number | HiddenAmount | NumberWithSuffix = accountsStore.getTotalLiabilities(showAccountBalance.value);
        return formatAmountToLocalizedNumeralsWithCurrency(totalLiabilities, defaultCurrency.value);
    });

    function accountCategoryTotalBalance(accountCategory?: AccountCategory): string {
        if (!accountCategory) {
            return '';
        }

        const totalBalance: number | HiddenAmount | NumberWithSuffix = accountsStore.getAccountCategoryTotalBalance(showAccountBalance.value, accountCategory);
        return formatAmountToLocalizedNumeralsWithCurrency(totalBalance, defaultCurrency.value);
    }

    function accountBalance(account: Account, currentSubAccountId?: string): string | null {
        const defaultCurrency = userStore.currentUserDefaultCurrency;

        if (account.type === AccountType.SingleAccount.type) {
            const balance: number| HiddenAmount | null = accountsStore.getAccountBalance(showAccountBalance.value, account);

            if (!isNumber(balance) && !isString(balance)) {
                return '';
            }

            const displayBalance = formatAmountToLocalizedNumeralsWithCurrency(balance, account.currency);

            if (showAccountBalance.value) {
                const totalBalance = formatAmountToLocalizedNumeralsWithCurrency(account.totalBalance, defaultCurrency);
                if (account.currency != defaultCurrency) {
                    return `${displayBalance} (≈ ${totalBalance})`;
                } else {
                    return `≈ ${totalBalance}`;
                }
            }

            return displayBalance;
        } else if (account.type === AccountType.MultiSubAccounts.type) {
            const balanceResult = accountsStore.getAccountSubAccountBalance(showAccountBalance.value, showHidden.value, account, currentSubAccountId);

            if (!isObject(balanceResult)) {
                return '';
            }

            const displayBalance = formatAmountToLocalizedNumeralsWithCurrency(balanceResult.balance, balanceResult.currency);

            if (showAccountBalance.value && currentSubAccountId) {
                const subAccount = account.getSubAccount(currentSubAccountId);

                if (subAccount) {
                    const totalBalance = formatAmountToLocalizedNumeralsWithCurrency(subAccount.totalBalance, defaultCurrency);
                    if (balanceResult.currency != defaultCurrency) {
                        return `${displayBalance} (≈ ${totalBalance})`;
                    } else {
                        return `≈ ${totalBalance}`;
                    }
                }
            } else if (showAccountBalance.value && !currentSubAccountId) {
                let hasNonFiatSubAccount = false;

                if (account.subAccounts) {
                    for (const subAccount of account.subAccounts) {
                        if (subAccount.assetType !== AccountAssetTypeClass.Fiat.type) {
                            hasNonFiatSubAccount = true;
                            break;
                        }
                    }
                }

                if (hasNonFiatSubAccount) {
                    const totalBalance = formatAmountToLocalizedNumeralsWithCurrency(account.totalBalance, defaultCurrency);
                    if (balanceResult.currency != defaultCurrency) {
                        return `${displayBalance} (≈ ${totalBalance})`;
                    } else {
                        return `≈ ${totalBalance}`;
                    }
                }
            }

            return displayBalance;
        } else {
            return null;
        }
    }

    return {
        // states
        loading,
        showHidden,
        displayOrderModified,
        // computed states
        showAccountBalance,
        customAccountCategoryOrder,
        defaultAccountCategory,
        firstDayOfWeek,
        fiscalYearStart,
        defaultCurrency,
        allAccounts,
        allCategorizedAccountsMap,
        allAccountCount,
        netAssets,
        totalAssets,
        totalLiabilities,
        // functions
        accountCategoryTotalBalance,
        accountBalance
    };
}
