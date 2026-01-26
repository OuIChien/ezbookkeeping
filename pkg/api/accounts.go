package api

import (
	"sort"

	"github.com/mayswind/ezbookkeeping/pkg/core"
	"github.com/mayswind/ezbookkeeping/pkg/cryptocurrency"
	"github.com/mayswind/ezbookkeeping/pkg/duplicatechecker"
	"github.com/mayswind/ezbookkeeping/pkg/errs"
	"github.com/mayswind/ezbookkeeping/pkg/exchangerates"
	"github.com/mayswind/ezbookkeeping/pkg/log"
	"github.com/mayswind/ezbookkeeping/pkg/models"
	"github.com/mayswind/ezbookkeeping/pkg/services"
	"github.com/mayswind/ezbookkeeping/pkg/settings"
	"github.com/mayswind/ezbookkeeping/pkg/stocks"
	"github.com/mayswind/ezbookkeeping/pkg/utils"
	"github.com/mayswind/ezbookkeeping/pkg/validators"
)

// AccountsApi represents account api
type AccountsApi struct {
	ApiUsingConfig
	ApiUsingDuplicateChecker
	accounts *services.AccountService
	users    *services.UserService
}

// Initialize an account api singleton instance
var (
	Accounts = &AccountsApi{
		ApiUsingConfig: ApiUsingConfig{
			container: settings.Container,
		},
		ApiUsingDuplicateChecker: ApiUsingDuplicateChecker{
			ApiUsingConfig: ApiUsingConfig{
				container: settings.Container,
			},
			container: duplicatechecker.Container,
		},
		accounts: services.Accounts,
		users:    services.Users,
	}
)

// AccountListHandler returns accounts list of current user
func (a *AccountsApi) AccountListHandler(c *core.WebContext) (any, *errs.Error) {
	var accountListReq models.AccountListRequest
	err := c.ShouldBindQuery(&accountListReq)

	if err != nil {
		log.Warnf(c, "[accounts.AccountListHandler] parse request failed, because %s", err.Error())
		return nil, errs.NewIncompleteOrIncorrectSubmissionError(err)
	}

	uid := c.GetCurrentUid()
	accounts, err := a.accounts.GetAllAccountsByUid(c, uid)

	if err != nil {
		log.Errorf(c, "[accounts.AccountListHandler] failed to get all accounts for user \"uid:%d\", because %s", uid, err.Error())
		return nil, errs.Or(err, errs.ErrOperationFailed)
	}

	userAllAccountResps := make([]*models.AccountInfoResponse, len(accounts))
	userAllAccountRespMap := make(map[int64]*models.AccountInfoResponse)

	for i := 0; i < len(accounts); i++ {
		userAllAccountResps[i] = accounts[i].ToAccountInfoResponse()
		userAllAccountRespMap[userAllAccountResps[i].Id] = userAllAccountResps[i]
	}

	for i := 0; i < len(userAllAccountResps); i++ {
		userAccountResp := userAllAccountResps[i]

		if accountListReq.VisibleOnly && userAccountResp.Hidden {
			continue
		}

		if userAccountResp.ParentId <= models.LevelOneAccountParentId {
			continue
		}

		parentAccount, parentExists := userAllAccountRespMap[userAccountResp.ParentId]

		if !parentExists || parentAccount == nil {
			continue
		}

		parentAccount.SubAccounts = append(parentAccount.SubAccounts, userAccountResp)
	}

	userFinalAccountResps := make(models.AccountInfoResponseSlice, 0, len(userAllAccountResps))

	for i := 0; i < len(userAllAccountResps); i++ {
		if userAllAccountResps[i].ParentId == models.LevelOneAccountParentId && (!accountListReq.VisibleOnly || !userAllAccountResps[i].Hidden) {
			sort.Sort(userAllAccountResps[i].SubAccounts)
			userFinalAccountResps = append(userFinalAccountResps, userAllAccountResps[i])
		}
	}

	sort.Sort(userFinalAccountResps)

	a.calculateAccountValuations(c, uid, userFinalAccountResps)

	return userFinalAccountResps, nil
}

// AccountGetHandler returns one specific account of current user
func (a *AccountsApi) AccountGetHandler(c *core.WebContext) (any, *errs.Error) {
	var accountGetReq models.AccountGetRequest
	err := c.ShouldBindQuery(&accountGetReq)

	if err != nil {
		log.Warnf(c, "[accounts.AccountGetHandler] parse request failed, because %s", err.Error())
		return nil, errs.NewIncompleteOrIncorrectSubmissionError(err)
	}

	uid := c.GetCurrentUid()
	accountAndSubAccounts, err := a.accounts.GetAccountAndSubAccountsByAccountId(c, uid, accountGetReq.Id)

	if err != nil {
		log.Errorf(c, "[accounts.AccountGetHandler] failed to get account \"id:%d\" for user \"uid:%d\", because %s", accountGetReq.Id, uid, err.Error())
		return nil, errs.Or(err, errs.ErrOperationFailed)
	}

	accountRespMap := make(map[int64]*models.AccountInfoResponse)

	for i := 0; i < len(accountAndSubAccounts); i++ {
		accountResp := accountAndSubAccounts[i].ToAccountInfoResponse()
		accountRespMap[accountResp.Id] = accountResp
	}

	accountResp, exists := accountRespMap[accountGetReq.Id]

	if !exists {
		return nil, errs.ErrAccountNotFound
	}

	for i := 0; i < len(accountAndSubAccounts); i++ {
		if accountAndSubAccounts[i].ParentAccountId == accountResp.Id {
			subAccountResp := accountAndSubAccounts[i].ToAccountInfoResponse()
			accountResp.SubAccounts = append(accountResp.SubAccounts, subAccountResp)
		}
	}

	sort.Sort(accountResp.SubAccounts)

	a.calculateAccountValuations(c, uid, []*models.AccountInfoResponse{accountResp})

	return accountResp, nil
}

// AccountCreateHandler saves a new account by request parameters for current user
func (a *AccountsApi) AccountCreateHandler(c *core.WebContext) (any, *errs.Error) {
	var accountCreateReq models.AccountCreateRequest
	err := c.ShouldBindJSON(&accountCreateReq)

	if err != nil {
		log.Warnf(c, "[accounts.AccountCreateHandler] parse request failed, because %s", err.Error())
		return nil, errs.NewIncompleteOrIncorrectSubmissionError(err)
	}

	clientTimezone, err := c.GetClientTimezone()

	if err != nil {
		log.Warnf(c, "[accounts.AccountCreateHandler] cannot get client timezone, because %s", err.Error())
		return nil, errs.ErrClientTimezoneOffsetInvalid
	}

	if accountCreateReq.Category < models.ACCOUNT_CATEGORY_CASH || accountCreateReq.Category > models.ACCOUNT_CATEGORY_CERTIFICATE_OF_DEPOSIT {
		log.Warnf(c, "[accounts.AccountCreateHandler] account category invalid, category is %d", accountCreateReq.Category)
		return nil, errs.ErrAccountCategoryInvalid
	}

	if accountCreateReq.Category != models.ACCOUNT_CATEGORY_CREDIT_CARD && accountCreateReq.CreditCardStatementDate != 0 {
		log.Warnf(c, "[accounts.AccountCreateHandler] cannot set statement date with category \"%d\"", accountCreateReq.Category)
		return nil, errs.ErrCannotSetStatementDateForNonCreditCard
	}

	if accountCreateReq.Type == models.ACCOUNT_TYPE_SINGLE_ACCOUNT {
		if len(accountCreateReq.SubAccounts) > 0 {
			log.Warnf(c, "[accounts.AccountCreateHandler] account cannot have any sub-accounts")
			return nil, errs.ErrAccountCannotHaveSubAccounts
		}

		if accountCreateReq.Currency == validators.ParentAccountCurrencyPlaceholder {
			log.Warnf(c, "[accounts.AccountCreateHandler] account cannot set currency placeholder")
			return nil, errs.ErrAccountCurrencyInvalid
		}

		if accountCreateReq.Balance != 0 && accountCreateReq.BalanceTime <= 0 {
			log.Warnf(c, "[accounts.AccountCreateHandler] account balance time is not set")
			return nil, errs.ErrAccountBalanceTimeNotSet
		}
	} else if accountCreateReq.Type == models.ACCOUNT_TYPE_MULTI_SUB_ACCOUNTS {
		if len(accountCreateReq.SubAccounts) < 1 {
			log.Warnf(c, "[accounts.AccountCreateHandler] account does not have any sub-accounts")
			return nil, errs.ErrAccountHaveNoSubAccount
		}

		if accountCreateReq.AssetType < models.ACCOUNT_ASSET_TYPE_FIAT || accountCreateReq.AssetType > models.ACCOUNT_ASSET_TYPE_STOCK {
			log.Warnf(c, "[accounts.AccountCreateHandler] parent account asset type is required for multi-sub-accounts")
			return nil, errs.ErrAccountAssetTypeRequiredForMultiSub
		}

		if accountCreateReq.Currency != validators.ParentAccountCurrencyPlaceholder {
			log.Warnf(c, "[accounts.AccountCreateHandler] parent account cannot set currency")
			return nil, errs.ErrParentAccountCannotSetCurrency
		}

		if accountCreateReq.Balance != 0 {
			log.Warnf(c, "[accounts.AccountCreateHandler] parent account cannot set balance")
			return nil, errs.ErrParentAccountCannotSetBalance
		}

		for i := 0; i < len(accountCreateReq.SubAccounts); i++ {
			subAccount := accountCreateReq.SubAccounts[i]

			if subAccount.Category != accountCreateReq.Category {
				log.Warnf(c, "[accounts.AccountCreateHandler] category of sub-account#%d not equals to parent", i)
				return nil, errs.ErrSubAccountCategoryNotEqualsToParent
			}

			if subAccount.Type != models.ACCOUNT_TYPE_SINGLE_ACCOUNT {
				log.Warnf(c, "[accounts.AccountCreateHandler] sub-account#%d type invalid", i)
				return nil, errs.ErrSubAccountTypeInvalid
			}

			if subAccount.AssetType != accountCreateReq.AssetType {
				log.Warnf(c, "[accounts.AccountCreateHandler] asset type of sub-account#%d not equals to parent", i)
				return nil, errs.ErrSubAccountAssetTypeNotEqualsToParent
			}

			if subAccount.Currency == validators.ParentAccountCurrencyPlaceholder {
				log.Warnf(c, "[accounts.AccountCreateHandler] sub-account#%d cannot set currency placeholder", i)
				return nil, errs.ErrAccountCurrencyInvalid
			}

			if subAccount.Balance != 0 && subAccount.BalanceTime <= 0 {
				log.Warnf(c, "[accounts.AccountCreateHandler] sub-account#%d balance time is not set", i)
				return nil, errs.ErrAccountBalanceTimeNotSet
			}

			if subAccount.CreditCardStatementDate != 0 {
				log.Warnf(c, "[accounts.AccountCreateHandler] sub-account#%d cannot set statement date", i)
				return nil, errs.ErrCannotSetStatementDateForSubAccount
			}
		}
	} else {
		log.Warnf(c, "[accounts.AccountCreateHandler] account type invalid, type is %d", accountCreateReq.Type)
		return nil, errs.ErrAccountTypeInvalid
	}

	uid := c.GetCurrentUid()
	maxOrderId, err := a.accounts.GetMaxDisplayOrder(c, uid, accountCreateReq.Category)

	if err != nil {
		log.Errorf(c, "[accounts.AccountCreateHandler] failed to get max display order for user \"uid:%d\", because %s", uid, err.Error())
		return nil, errs.Or(err, errs.ErrOperationFailed)
	}

	mainAccount := a.createNewAccountModel(uid, &accountCreateReq, false, maxOrderId+1)
	childrenAccounts, childrenAccountBalanceTimes := a.createSubAccountModels(uid, &accountCreateReq)

	if a.CurrentConfig().EnableDuplicateSubmissionsCheck && accountCreateReq.ClientSessionId != "" {
		found, remark := a.GetSubmissionRemark(duplicatechecker.DUPLICATE_CHECKER_TYPE_NEW_ACCOUNT, uid, accountCreateReq.ClientSessionId)

		if found {
			log.Infof(c, "[accounts.AccountCreateHandler] another account \"id:%s\" has been created for user \"uid:%d\"", remark, uid)
			accountId, err := utils.StringToInt64(remark)

			if err == nil {
				accountAndSubAccounts, err := a.accounts.GetAccountAndSubAccountsByAccountId(c, uid, accountId)

				if err != nil {
					log.Errorf(c, "[accounts.AccountCreateHandler] failed to get existed account \"id:%d\" for user \"uid:%d\", because %s", accountId, uid, err.Error())
					return nil, errs.Or(err, errs.ErrOperationFailed)
				}

				accountMap := a.accounts.GetAccountMapByList(accountAndSubAccounts)
				mainAccount, exists := accountMap[accountId]

				if !exists {
					return nil, errs.ErrOperationFailed
				}

				accountInfoResp := mainAccount.ToAccountInfoResponse()

				for i := 0; i < len(accountAndSubAccounts); i++ {
					if accountAndSubAccounts[i].ParentAccountId == mainAccount.AccountId {
						subAccountResp := accountAndSubAccounts[i].ToAccountInfoResponse()
						accountInfoResp.SubAccounts = append(accountInfoResp.SubAccounts, subAccountResp)
					}
				}

				a.calculateAccountValuations(c, uid, []*models.AccountInfoResponse{accountInfoResp})

				return accountInfoResp, nil
			}
		}
	}

	err = a.accounts.CreateAccounts(c, mainAccount, accountCreateReq.BalanceTime, childrenAccounts, childrenAccountBalanceTimes, clientTimezone)

	if err != nil {
		log.Errorf(c, "[accounts.AccountCreateHandler] failed to create account \"id:%d\" for user \"uid:%d\", because %s", mainAccount.AccountId, uid, err.Error())
		return nil, errs.Or(err, errs.ErrOperationFailed)
	}

	log.Infof(c, "[accounts.AccountCreateHandler] user \"uid:%d\" has created a new account \"id:%d\" successfully", uid, mainAccount.AccountId)

	a.SetSubmissionRemarkIfEnable(duplicatechecker.DUPLICATE_CHECKER_TYPE_NEW_ACCOUNT, uid, accountCreateReq.ClientSessionId, utils.Int64ToString(mainAccount.AccountId))
	accountInfoResp := mainAccount.ToAccountInfoResponse()

	if len(childrenAccounts) > 0 {
		accountInfoResp.SubAccounts = make([]*models.AccountInfoResponse, len(childrenAccounts))

		for i := 0; i < len(childrenAccounts); i++ {
			accountInfoResp.SubAccounts[i] = childrenAccounts[i].ToAccountInfoResponse()
		}
	}

	a.calculateAccountValuations(c, uid, []*models.AccountInfoResponse{accountInfoResp})

	return accountInfoResp, nil
}

// AccountModifyHandler saves an existed account by request parameters for current user
func (a *AccountsApi) AccountModifyHandler(c *core.WebContext) (any, *errs.Error) {
	var accountModifyReq models.AccountModifyRequest
	err := c.ShouldBindJSON(&accountModifyReq)

	if err != nil {
		log.Warnf(c, "[accounts.AccountModifyHandler] parse request failed, because %s", err.Error())
		return nil, errs.NewIncompleteOrIncorrectSubmissionError(err)
	}

	if accountModifyReq.Id <= 0 {
		return nil, errs.ErrAccountIdInvalid
	}

	clientTimezone, err := c.GetClientTimezone()

	if err != nil {
		log.Warnf(c, "[accounts.AccountModifyHandler] cannot get client timezone, because %s", err.Error())
		return nil, errs.ErrClientTimezoneOffsetInvalid
	}

	if accountModifyReq.Category < models.ACCOUNT_CATEGORY_CASH || accountModifyReq.Category > models.ACCOUNT_CATEGORY_CERTIFICATE_OF_DEPOSIT {
		log.Warnf(c, "[accounts.AccountModifyHandler] account category invalid, category is %d", accountModifyReq.Category)
		return nil, errs.ErrAccountCategoryInvalid
	}

	if accountModifyReq.Category != models.ACCOUNT_CATEGORY_CREDIT_CARD && accountModifyReq.CreditCardStatementDate != 0 {
		log.Warnf(c, "[accounts.AccountModifyHandler] cannot set statement date with category \"%d\"", accountModifyReq.Category)
		return nil, errs.ErrCannotSetStatementDateForNonCreditCard
	}

	uid := c.GetCurrentUid()
	accountAndSubAccounts, err := a.accounts.GetAccountAndSubAccountsByAccountId(c, uid, accountModifyReq.Id)

	if err != nil {
		log.Errorf(c, "[accounts.AccountModifyHandler] failed to get account \"id:%d\" for user \"uid:%d\", because %s", accountModifyReq.Id, uid, err.Error())
		return nil, errs.Or(err, errs.ErrOperationFailed)
	}

	accountMap := a.accounts.GetAccountMapByList(accountAndSubAccounts)
	mainAccount, exists := accountMap[accountModifyReq.Id]

	if !exists {
		return nil, errs.ErrAccountNotFound
	}

	if accountModifyReq.Currency != nil && mainAccount.Currency != *accountModifyReq.Currency {
		return nil, errs.ErrNotSupportedChangeCurrency
	}

	if accountModifyReq.Balance != nil {
		return nil, errs.ErrNotSupportedChangeBalance
	}

	if accountModifyReq.BalanceTime != nil {
		return nil, errs.ErrNotSupportedChangeBalanceTime
	}

	// Check if asset type is being changed (not allowed for existing accounts)
	if accountModifyReq.AssetType != nil {
		var oldAssetType models.AccountAssetType
		if mainAccount.Extend != nil {
			oldAssetType = mainAccount.Extend.AssetType
		}
		if *accountModifyReq.AssetType != oldAssetType {
			log.Warnf(c, "[accounts.AccountModifyHandler] cannot modify account asset type")
			return nil, errs.ErrNotSupportedChangeAssetType
		}
	}

	if mainAccount.Type == models.ACCOUNT_TYPE_SINGLE_ACCOUNT {
		if len(accountModifyReq.SubAccounts) > 0 {
			log.Warnf(c, "[accounts.AccountModifyHandler] account cannot have any sub-accounts")
			return nil, errs.ErrAccountCannotHaveSubAccounts
		}
	} else if mainAccount.Type == models.ACCOUNT_TYPE_MULTI_SUB_ACCOUNTS {
		if len(accountModifyReq.SubAccounts) < 1 {
			log.Warnf(c, "[accounts.AccountModifyHandler] account does not have any sub-accounts")
			return nil, errs.ErrAccountHaveNoSubAccount
		}

		for i := 0; i < len(accountModifyReq.SubAccounts); i++ {
			subAccountReq := accountModifyReq.SubAccounts[i]

			if subAccountReq.Category != accountModifyReq.Category {
				log.Warnf(c, "[accounts.AccountModifyHandler] category of sub-account#%d not equals to parent", i)
				return nil, errs.ErrSubAccountCategoryNotEqualsToParent
			}

			if subAccountReq.Id == 0 { // create new sub-account
				// Verify asset type matches parent account
				var parentAssetType models.AccountAssetType
				if mainAccount.Extend != nil {
					parentAssetType = mainAccount.Extend.AssetType
				}
				if subAccountReq.AssetType != nil && *subAccountReq.AssetType != parentAssetType {
					log.Warnf(c, "[accounts.AccountModifyHandler] asset type of sub-account#%d not equals to parent", i)
					return nil, errs.ErrSubAccountAssetTypeNotEqualsToParent
				}

				if subAccountReq.Currency == nil {
					log.Warnf(c, "[accounts.AccountModifyHandler] sub-account#%d not set currency", i)
					return nil, errs.ErrAccountCurrencyInvalid
				} else if subAccountReq.Currency != nil && *subAccountReq.Currency == validators.ParentAccountCurrencyPlaceholder {
					log.Warnf(c, "[accounts.AccountModifyHandler] sub-account#%d cannot set currency placeholder", i)
					return nil, errs.ErrAccountCurrencyInvalid
				}

				if subAccountReq.Balance == nil {
					defaultBalance := int64(0)
					subAccountReq.Balance = &defaultBalance
				}

				if *subAccountReq.Balance == 0 {
					defaultBalanceTime := int64(0)
					subAccountReq.BalanceTime = &defaultBalanceTime
				}

				if *subAccountReq.Balance != 0 && (subAccountReq.BalanceTime == nil || *subAccountReq.BalanceTime <= 0) {
					log.Warnf(c, "[accounts.AccountModifyHandler] sub-account#%d balance time is not set", i)
					return nil, errs.ErrAccountBalanceTimeNotSet
				}
			} else { // modify existed sub-account
				subAccount, exists := accountMap[subAccountReq.Id]

				if !exists {
					return nil, errs.ErrAccountNotFound
				}

				if subAccountReq.Currency != nil && subAccount.Currency != *subAccountReq.Currency {
					return nil, errs.ErrNotSupportedChangeCurrency
				}

				if subAccountReq.Balance != nil {
					return nil, errs.ErrNotSupportedChangeBalance
				}

				if subAccountReq.BalanceTime != nil {
					return nil, errs.ErrNotSupportedChangeBalanceTime
				}

				// Check if asset type is being changed (not allowed for existing sub-accounts)
				if subAccountReq.AssetType != nil {
					var oldAssetType models.AccountAssetType
					if subAccount.Extend != nil {
						oldAssetType = subAccount.Extend.AssetType
					}
					if *subAccountReq.AssetType != oldAssetType {
						log.Warnf(c, "[accounts.AccountModifyHandler] cannot modify sub-account#%d asset type", i)
						return nil, errs.ErrNotSupportedChangeAssetType
					}
				}
			}

			if subAccountReq.CreditCardStatementDate != 0 {
				log.Warnf(c, "[accounts.AccountModifyHandler] sub-account#%d cannot set statement date", i)
				return nil, errs.ErrCannotSetStatementDateForSubAccount
			}
		}
	}

	anythingUpdate := false
	var toUpdateAccounts []*models.Account
	var toAddAccounts []*models.Account
	var toAddAccountBalanceTimes []int64
	var toDeleteAccountIds []int64

	toUpdateAccount := a.getToUpdateAccount(uid, &accountModifyReq, mainAccount, false)

	if toUpdateAccount != nil {
		anythingUpdate = true
		toUpdateAccounts = append(toUpdateAccounts, toUpdateAccount)
	}

	toDeleteAccountIds = a.getToDeleteSubAccountIds(&accountModifyReq, mainAccount, accountAndSubAccounts)

	if len(toDeleteAccountIds) > 0 {
		anythingUpdate = true
	}

	maxOrderId := int32(0)

	for i := 0; i < len(accountAndSubAccounts); i++ {
		account := accountAndSubAccounts[i]

		if account.AccountId != mainAccount.AccountId && account.DisplayOrder > maxOrderId {
			maxOrderId = account.DisplayOrder
		}
	}

	for i := 0; i < len(accountModifyReq.SubAccounts); i++ {
		subAccountReq := accountModifyReq.SubAccounts[i]

		if _, exists := accountMap[subAccountReq.Id]; !exists {
			// For new sub-accounts, inherit asset type from parent account
			var parentAssetType models.AccountAssetType
			if mainAccount.Extend != nil {
				parentAssetType = mainAccount.Extend.AssetType
			}
			subAccountReq.AssetType = &parentAssetType

			anythingUpdate = true
			maxOrderId = maxOrderId + 1
			newSubAccount := a.createNewSubAccountModelForModify(uid, mainAccount.Type, subAccountReq, maxOrderId)
			toAddAccounts = append(toAddAccounts, newSubAccount)

			if subAccountReq.BalanceTime != nil {
				toAddAccountBalanceTimes = append(toAddAccountBalanceTimes, *subAccountReq.BalanceTime)
			} else {
				toAddAccountBalanceTimes = append(toAddAccountBalanceTimes, 0)
			}
		} else {
			toUpdateSubAccount := a.getToUpdateAccount(uid, subAccountReq, accountMap[subAccountReq.Id], true)

			if toUpdateSubAccount != nil {
				anythingUpdate = true
				toUpdateAccounts = append(toUpdateAccounts, toUpdateSubAccount)
			}
		}
	}

	if !anythingUpdate {
		return nil, errs.ErrNothingWillBeUpdated
	}

	if len(toAddAccounts) > 0 && a.CurrentConfig().EnableDuplicateSubmissionsCheck && accountModifyReq.ClientSessionId != "" {
		found, remark := a.GetSubmissionRemark(duplicatechecker.DUPLICATE_CHECKER_TYPE_NEW_SUBACCOUNT, uid, accountModifyReq.ClientSessionId)

		if found {
			log.Infof(c, "[accounts.AccountModifyHandler] another account \"id:%s\" modification has been created for user \"uid:%d\"", remark, uid)
			accountId, err := utils.StringToInt64(remark)

			if err == nil {
				accountAndSubAccounts, err := a.accounts.GetAccountAndSubAccountsByAccountId(c, uid, accountId)

				if err != nil {
					log.Errorf(c, "[accounts.AccountModifyHandler] failed to get existed account \"id:%d\" for user \"uid:%d\", because %s", accountId, uid, err.Error())
					return nil, errs.Or(err, errs.ErrOperationFailed)
				}

				accountMap := a.accounts.GetAccountMapByList(accountAndSubAccounts)
				mainAccount, exists := accountMap[accountId]

				if !exists {
					return nil, errs.ErrOperationFailed
				}

				accountInfoResp := mainAccount.ToAccountInfoResponse()

				for i := 0; i < len(accountAndSubAccounts); i++ {
					if accountAndSubAccounts[i].ParentAccountId == mainAccount.AccountId {
						subAccountResp := accountAndSubAccounts[i].ToAccountInfoResponse()
						accountInfoResp.SubAccounts = append(accountInfoResp.SubAccounts, subAccountResp)
					}
				}

				a.calculateAccountValuations(c, uid, []*models.AccountInfoResponse{accountInfoResp})

				return accountInfoResp, nil
			}
		}
	}

	err = a.accounts.ModifyAccounts(c, mainAccount, toUpdateAccounts, toAddAccounts, toAddAccountBalanceTimes, toDeleteAccountIds, clientTimezone)

	if err != nil {
		log.Errorf(c, "[accounts.AccountModifyHandler] failed to update account \"id:%d\" for user \"uid:%d\", because %s", accountModifyReq.Id, uid, err.Error())
		return nil, errs.Or(err, errs.ErrOperationFailed)
	}

	log.Infof(c, "[accounts.AccountModifyHandler] user \"uid:%d\" has updated account \"id:%d\" successfully", uid, accountModifyReq.Id)

	if len(toAddAccounts) > 0 {
		a.SetSubmissionRemarkIfEnable(duplicatechecker.DUPLICATE_CHECKER_TYPE_NEW_SUBACCOUNT, uid, accountModifyReq.ClientSessionId, utils.Int64ToString(mainAccount.AccountId))
	}

	accountRespMap := make(map[int64]*models.AccountInfoResponse)

	for i := 0; i < len(toUpdateAccounts); i++ {
		account := toUpdateAccounts[i]
		oldAccount := accountMap[account.AccountId]

		account.Type = oldAccount.Type
		account.ParentAccountId = oldAccount.ParentAccountId
		account.DisplayOrder = oldAccount.DisplayOrder
		account.Currency = oldAccount.Currency
		account.Balance = oldAccount.Balance

		accountResp := account.ToAccountInfoResponse()
		accountRespMap[accountResp.Id] = accountResp
	}

	for i := 0; i < len(toAddAccounts); i++ {
		account := toAddAccounts[i]
		accountResp := account.ToAccountInfoResponse()
		accountRespMap[accountResp.Id] = accountResp
	}

	deletedAccountIds := make(map[int64]bool)

	for i := 0; i < len(toDeleteAccountIds); i++ {
		deletedAccountIds[toDeleteAccountIds[i]] = true
	}

	for i := 0; i < len(accountAndSubAccounts); i++ {
		oldAccount := accountAndSubAccounts[i]
		_, exists := accountRespMap[oldAccount.AccountId]

		if !exists && !deletedAccountIds[oldAccount.AccountId] {
			oldAccountResp := oldAccount.ToAccountInfoResponse()
			accountRespMap[oldAccountResp.Id] = oldAccountResp
		}
	}

	accountResp := accountRespMap[accountModifyReq.Id]

	for i := 0; i < len(accountAndSubAccounts); i++ {
		account := accountAndSubAccounts[i]

		if account.ParentAccountId == accountResp.Id && !deletedAccountIds[account.AccountId] {
			subAccountResp := accountRespMap[account.AccountId]
			accountResp.SubAccounts = append(accountResp.SubAccounts, subAccountResp)
		}
	}

	for i := 0; i < len(toAddAccounts); i++ {
		account := toAddAccounts[i]

		if account.ParentAccountId == accountResp.Id {
			subAccountResp := accountRespMap[account.AccountId]
			accountResp.SubAccounts = append(accountResp.SubAccounts, subAccountResp)
		}
	}

	sort.Sort(accountResp.SubAccounts)

	a.calculateAccountValuations(c, uid, []*models.AccountInfoResponse{accountResp})

	return accountResp, nil
}

// AccountHideHandler hides an existed account by request parameters for current user
func (a *AccountsApi) AccountHideHandler(c *core.WebContext) (any, *errs.Error) {
	var accountHideReq models.AccountHideRequest
	err := c.ShouldBindJSON(&accountHideReq)

	if err != nil {
		log.Warnf(c, "[accounts.AccountHideHandler] parse request failed, because %s", err.Error())
		return nil, errs.NewIncompleteOrIncorrectSubmissionError(err)
	}

	uid := c.GetCurrentUid()
	err = a.accounts.HideAccount(c, uid, []int64{accountHideReq.Id}, accountHideReq.Hidden)

	if err != nil {
		log.Errorf(c, "[accounts.AccountHideHandler] failed to hide account \"id:%d\" for user \"uid:%d\", because %s", accountHideReq.Id, uid, err.Error())
		return nil, errs.Or(err, errs.ErrOperationFailed)
	}

	log.Infof(c, "[accounts.AccountHideHandler] user \"uid:%d\" has hidden account \"id:%d\"", uid, accountHideReq.Id)
	return true, nil
}

// AccountMoveHandler moves display order of existed accounts by request parameters for current user
func (a *AccountsApi) AccountMoveHandler(c *core.WebContext) (any, *errs.Error) {
	var accountMoveReq models.AccountMoveRequest
	err := c.ShouldBindJSON(&accountMoveReq)

	if err != nil {
		log.Warnf(c, "[accounts.AccountMoveHandler] parse request failed, because %s", err.Error())
		return nil, errs.NewIncompleteOrIncorrectSubmissionError(err)
	}

	uid := c.GetCurrentUid()
	accounts := make([]*models.Account, len(accountMoveReq.NewDisplayOrders))

	for i := 0; i < len(accountMoveReq.NewDisplayOrders); i++ {
		newDisplayOrder := accountMoveReq.NewDisplayOrders[i]
		account := &models.Account{
			Uid:          uid,
			AccountId:    newDisplayOrder.Id,
			DisplayOrder: newDisplayOrder.DisplayOrder,
		}

		accounts[i] = account
	}

	err = a.accounts.ModifyAccountDisplayOrders(c, uid, accounts)

	if err != nil {
		log.Errorf(c, "[accounts.AccountMoveHandler] failed to move accounts for user \"uid:%d\", because %s", uid, err.Error())
		return nil, errs.Or(err, errs.ErrOperationFailed)
	}

	log.Infof(c, "[accounts.AccountMoveHandler] user \"uid:%d\" has moved accounts", uid)
	return true, nil
}

// AccountDeleteHandler deletes an existed account by request parameters for current user
func (a *AccountsApi) AccountDeleteHandler(c *core.WebContext) (any, *errs.Error) {
	var accountDeleteReq models.AccountDeleteRequest
	err := c.ShouldBindJSON(&accountDeleteReq)

	if err != nil {
		log.Warnf(c, "[accounts.AccountDeleteHandler] parse request failed, because %s", err.Error())
		return nil, errs.NewIncompleteOrIncorrectSubmissionError(err)
	}

	uid := c.GetCurrentUid()
	err = a.accounts.DeleteAccount(c, uid, accountDeleteReq.Id)

	if err != nil {
		log.Errorf(c, "[accounts.AccountDeleteHandler] failed to delete account \"id:%d\" for user \"uid:%d\", because %s", accountDeleteReq.Id, uid, err.Error())
		return nil, errs.Or(err, errs.ErrOperationFailed)
	}

	log.Infof(c, "[accounts.AccountDeleteHandler] user \"uid:%d\" has deleted account \"id:%d\"", uid, accountDeleteReq.Id)
	return true, nil
}

// SubAccountDeleteHandler deletes an existed sub-account by request parameters for current user
func (a *AccountsApi) SubAccountDeleteHandler(c *core.WebContext) (any, *errs.Error) {
	var accountDeleteReq models.AccountDeleteRequest
	err := c.ShouldBindJSON(&accountDeleteReq)

	if err != nil {
		log.Warnf(c, "[accounts.SubAccountDeleteHandler] parse request failed, because %s", err.Error())
		return nil, errs.NewIncompleteOrIncorrectSubmissionError(err)
	}

	uid := c.GetCurrentUid()
	err = a.accounts.DeleteSubAccount(c, uid, accountDeleteReq.Id)

	if err != nil {
		log.Errorf(c, "[accounts.SubAccountDeleteHandler] failed to delete sub-account \"id:%d\" for user \"uid:%d\", because %s", accountDeleteReq.Id, uid, err.Error())
		return nil, errs.Or(err, errs.ErrOperationFailed)
	}

	log.Infof(c, "[accounts.SubAccountDeleteHandler] user \"uid:%d\" has deleted sub-account \"id:%d\"", uid, accountDeleteReq.Id)
	return true, nil
}

func (a *AccountsApi) createNewAccountModel(uid int64, accountCreateReq *models.AccountCreateRequest, isSubAccount bool, order int32) *models.Account {
	accountExtend := &models.AccountExtend{
		AssetType: accountCreateReq.AssetType,
	}

	if !isSubAccount && accountCreateReq.Category == models.ACCOUNT_CATEGORY_CREDIT_CARD {
		accountExtend.CreditCardStatementDate = &accountCreateReq.CreditCardStatementDate
	}

	return &models.Account{
		Uid:          uid,
		Name:         accountCreateReq.Name,
		DisplayOrder: order,
		Category:     accountCreateReq.Category,
		Type:         accountCreateReq.Type,
		Icon:         accountCreateReq.Icon,
		Color:        accountCreateReq.Color,
		Currency:     accountCreateReq.Currency,
		Balance:      accountCreateReq.Balance,
		Comment:      accountCreateReq.Comment,
		Extend:       accountExtend,
	}
}

func (a *AccountsApi) createNewSubAccountModelForModify(uid int64, accountType models.AccountType, accountModifyReq *models.AccountModifyRequest, order int32) *models.Account {
	var assetType models.AccountAssetType

	if accountModifyReq.AssetType != nil {
		assetType = *accountModifyReq.AssetType
	}

	accountExtend := &models.AccountExtend{
		AssetType: assetType,
	}

	return &models.Account{
		Uid:          uid,
		Name:         accountModifyReq.Name,
		DisplayOrder: order,
		Category:     accountModifyReq.Category,
		Type:         accountType,
		Icon:         accountModifyReq.Icon,
		Color:        accountModifyReq.Color,
		Currency:     *accountModifyReq.Currency,
		Balance:      *accountModifyReq.Balance,
		Comment:      accountModifyReq.Comment,
		Extend:       accountExtend,
	}
}

func (a *AccountsApi) createSubAccountModels(uid int64, accountCreateReq *models.AccountCreateRequest) ([]*models.Account, []int64) {
	if len(accountCreateReq.SubAccounts) <= 0 {
		return nil, nil
	}

	childrenAccounts := make([]*models.Account, len(accountCreateReq.SubAccounts))
	childrenAccountBalanceTimes := make([]int64, len(accountCreateReq.SubAccounts))

	for i := int32(0); i < int32(len(accountCreateReq.SubAccounts)); i++ {
		// Create sub-account model, but ensure asset type is inherited from parent
		subAccountReq := accountCreateReq.SubAccounts[i]
		subAccountReq.AssetType = accountCreateReq.AssetType
		childrenAccounts[i] = a.createNewAccountModel(uid, subAccountReq, true, i+1)
		childrenAccountBalanceTimes[i] = accountCreateReq.SubAccounts[i].BalanceTime
	}

	return childrenAccounts, childrenAccountBalanceTimes
}

func (a *AccountsApi) getToUpdateAccount(uid int64, accountModifyReq *models.AccountModifyRequest, oldAccount *models.Account, isSubAccount bool) *models.Account {
	var assetType models.AccountAssetType

	if accountModifyReq.AssetType != nil {
		assetType = *accountModifyReq.AssetType
	} else if oldAccount.Extend != nil {
		assetType = oldAccount.Extend.AssetType
	}

	newAccountExtend := &models.AccountExtend{
		AssetType: assetType,
	}

	if !isSubAccount && accountModifyReq.Category == models.ACCOUNT_CATEGORY_CREDIT_CARD {
		newAccountExtend.CreditCardStatementDate = &accountModifyReq.CreditCardStatementDate
	}

	newAccount := &models.Account{
		AccountId: oldAccount.AccountId,
		Uid:       uid,
		Name:      accountModifyReq.Name,
		Category:  accountModifyReq.Category,
		Icon:      accountModifyReq.Icon,
		Color:     accountModifyReq.Color,
		Comment:   accountModifyReq.Comment,
		Extend:    newAccountExtend,
		Hidden:    accountModifyReq.Hidden,
	}

	if newAccount.Name != oldAccount.Name ||
		newAccount.Category != oldAccount.Category ||
		newAccount.Icon != oldAccount.Icon ||
		newAccount.Color != oldAccount.Color ||
		newAccount.Comment != oldAccount.Comment ||
		newAccount.Hidden != oldAccount.Hidden {
		return newAccount
	}

	if (newAccount.Extend != nil && oldAccount.Extend == nil) ||
		(newAccount.Extend == nil && oldAccount.Extend != nil) {
		return newAccount
	}

	oldAccountExtend := oldAccount.Extend

	if newAccountExtend.AssetType != oldAccountExtend.AssetType ||
		newAccountExtend.CreditCardStatementDate != oldAccountExtend.CreditCardStatementDate {
		return newAccount
	}

	return nil
}

func (a *AccountsApi) getToDeleteSubAccountIds(accountModifyReq *models.AccountModifyRequest, mainAccount *models.Account, accountAndSubAccounts []*models.Account) []int64 {
	newSubAccountIds := make(map[int64]bool, len(accountModifyReq.SubAccounts))

	for i := 0; i < len(accountModifyReq.SubAccounts); i++ {
		newSubAccountIds[accountModifyReq.SubAccounts[i].Id] = true
	}

	toDeleteAccountIds := make([]int64, 0)

	for i := 0; i < len(accountAndSubAccounts); i++ {
		subAccount := accountAndSubAccounts[i]

		if subAccount.AccountId == mainAccount.AccountId {
			continue
		}

		if _, exists := newSubAccountIds[subAccount.AccountId]; !exists {
			toDeleteAccountIds = append(toDeleteAccountIds, subAccount.AccountId)
		}
	}

	return toDeleteAccountIds
}

func (a *AccountsApi) calculateAccountValuations(c *core.WebContext, uid int64, accountResps []*models.AccountInfoResponse) {
	// 1. Get user's default currency
	user, err := a.users.GetUserById(c, uid)
	if err != nil {
		return
	}
	defaultCurrency := user.DefaultCurrency

	// 2. Fetch all rates/prices
	exchangeRateResponse, _ := exchangerates.Container.GetLatestExchangeRates(c, uid, a.CurrentConfig())
	cryptoPriceResponse, _ := cryptocurrency.Container.GetLatestCryptocurrencyPrices(c, uid, a.CurrentConfig())
	stockPriceResponse, _ := stocks.Container.GetLatestStockPrices(c, uid, a.CurrentConfig())

	// 3. Create maps for quick lookup
	exchangeRates := make(map[string]float64)
	exchangeBaseCurrency := ""
	if exchangeRateResponse != nil {
		exchangeBaseCurrency = exchangeRateResponse.BaseCurrency
		for _, rate := range exchangeRateResponse.ExchangeRates {
			val, _ := utils.StringToFloat64(rate.Rate)
			exchangeRates[rate.Currency] = val
		}
	}
	if _, ok := exchangeRates[exchangeBaseCurrency]; !ok && exchangeBaseCurrency != "" {
		exchangeRates[exchangeBaseCurrency] = 1.0
	}

	cryptoPrices := make(map[string]float64)
	cryptoBaseCurrency := ""
	if cryptoPriceResponse != nil {
		cryptoBaseCurrency = cryptoPriceResponse.BaseCurrency
		for _, price := range cryptoPriceResponse.Prices {
			val, _ := utils.StringToFloat64(price.Price)
			cryptoPrices[price.Symbol] = val
		}
	}

	stockPrices := make(map[string]struct {
		price    float64
		currency string
	})
	if stockPriceResponse != nil {
		for _, price := range stockPriceResponse.Prices {
			val, _ := utils.StringToFloat64(price.Price)
			stockPrices[price.Symbol] = struct {
				price    float64
				currency string
			}{
				price:    val,
				currency: price.Currency,
			}
		}
	}

	// 4. Calculate valuation for each account
	for _, account := range accountResps {
		a.calculateSingleAccountValuation(account, defaultCurrency, exchangeBaseCurrency, exchangeRates, cryptoPrices, cryptoBaseCurrency, stockPrices)
	}
}

func (a *AccountsApi) calculateSingleAccountValuation(account *models.AccountInfoResponse, defaultCurrency string, exchangeBaseCurrency string, exchangeRates map[string]float64, cryptoPrices map[string]float64, cryptoBaseCurrency string, stockPrices map[string]struct {
	price    float64
	currency string
}) {
	if account.AssetType == models.ACCOUNT_ASSET_TYPE_FIAT {
		if account.Currency == defaultCurrency {
			account.TotalBalance = account.Balance
		} else {
			rateSrc, okSrc := exchangeRates[account.Currency]
			rateDst, okDst := exchangeRates[defaultCurrency]

			if okSrc && okDst && rateSrc != 0 {
				sourceFraction := a.getCurrencyFraction(account.Currency)
				targetFraction := a.getCurrencyFraction(defaultCurrency)
				account.TotalBalance = int64(float64(account.Balance) / rateSrc * rateDst / utils.Pow10(sourceFraction-targetFraction))
			} else {
				account.TotalBalance = 0
			}
		}
	} else if account.AssetType == models.ACCOUNT_ASSET_TYPE_CRYPTO {
		if price, ok := cryptoPrices[account.Currency]; ok {
			sourceFraction := a.getCurrencyFraction(account.Currency)
			totalInCryptoBase := float64(account.Balance) * price / utils.Pow10(sourceFraction)
			rateSrc, okSrc := exchangeRates[cryptoBaseCurrency]
			rateDst, okDst := exchangeRates[defaultCurrency]

			if cryptoBaseCurrency == defaultCurrency {
				targetFraction := a.getCurrencyFraction(defaultCurrency)
				account.TotalBalance = int64(totalInCryptoBase * utils.Pow10(targetFraction))
			} else if okSrc && okDst && rateSrc != 0 {
				targetFraction := a.getCurrencyFraction(defaultCurrency)
				account.TotalBalance = int64(totalInCryptoBase / rateSrc * rateDst * utils.Pow10(targetFraction))
			} else {
				account.TotalBalance = 0
			}
		} else {
			account.TotalBalance = 0
		}
	} else if account.AssetType == models.ACCOUNT_ASSET_TYPE_STOCK {
		if stockData, ok := stockPrices[account.Currency]; ok {
			sourceFraction := a.getCurrencyFraction(account.Currency)
			totalInStockCurrency := float64(account.Balance) * stockData.price / utils.Pow10(sourceFraction)
			rateSrc, okSrc := exchangeRates[stockData.currency]
			rateDst, okDst := exchangeRates[defaultCurrency]

			if stockData.currency == defaultCurrency {
				targetFraction := a.getCurrencyFraction(defaultCurrency)
				account.TotalBalance = int64(totalInStockCurrency * utils.Pow10(targetFraction))
			} else if okSrc && okDst && rateSrc != 0 {
				targetFraction := a.getCurrencyFraction(defaultCurrency)
				account.TotalBalance = int64(totalInStockCurrency / rateSrc * rateDst * utils.Pow10(targetFraction))
			} else {
				account.TotalBalance = 0
			}
		} else {
			account.TotalBalance = 0
		}
	} else {
		account.TotalBalance = account.Balance
	}

	// Recurse for sub-accounts
	for _, subAccount := range account.SubAccounts {
		a.calculateSingleAccountValuation(subAccount, defaultCurrency, exchangeBaseCurrency, exchangeRates, cryptoPrices, cryptoBaseCurrency, stockPrices)
	}

	// If it's a multi-sub-accounts parent, the total balance should be the sum of sub-accounts
	if account.Type == models.ACCOUNT_TYPE_MULTI_SUB_ACCOUNTS {
		account.TotalBalance = 0
		for _, subAccount := range account.SubAccounts {
			account.TotalBalance += subAccount.TotalBalance
		}
	}
}

func (a *AccountsApi) getCurrencyFraction(currency string) int {
	if fraction, ok := commonCurrencyFractions[currency]; ok {
		// For cryptocurrencies with fraction > 6, limit to 6 to avoid int64 overflow
		// int64 max value: 9223372036854775807
		// With fraction=8: can store up to ~92,233,720,368 BTC (more than enough)
		// With fraction=18: can only store ~9.22 ETH (not enough)
		if fraction > 8 {
			return 8
		}
		return fraction
	}
	return 2
}

var commonCurrencyFractions = map[string]int{
	"BIF": 0, "CLP": 0, "DJF": 0, "GNF": 0, "ISK": 0, "JPY": 0, "KMF": 0, "KRW": 0, "PYG": 0, "RWF": 0, "UGX": 0, "VND": 0, "VUV": 0, "XAF": 0, "XOF": 0, "XPF": 0,
	"BHD": 3, "IQD": 3, "JOD": 3, "KWD": 3, "LYD": 3, "OMR": 3, "TND": 3,
	"BTC": 8, "ETH": 5, "BNB": 5, "SOL": 5, "ADA": 4, "XRP": 4, "DOT": 3, "DOGE": 2, "MATIC": 4, "USDT": 2, "USDC": 2, "DAI": 2, "LTC": 4, "BCH": 4, "LINK": 4, "XLM": 4, "UNI": 4, "ATOM": 4, "XMR": 4, "ETC": 4,
}
