# ezbookkeeping Account Feature Analysis Report

## 1. Overview
The account feature in ezbookkeeping is the core of the entire financial management system. It not only records the state of assets and liabilities but also supports flexible bookkeeping needs through a hierarchical structure and various account categories. Accounts are tightly integrated with transactions to dynamically reflect the user's financial status.

## 2. Account Structure

### 2.1 Account Types
The system supports two primary forms of account organization:
- **Single Account**: The most basic unit that directly holds a balance and records transactions.
- **Multi-sub accounts**: A parent account used as a container for multiple sub-accounts. The parent account itself doesn't usually record transactions directly; instead, transactions are recorded in its sub-accounts.

### 2.2 Hierarchy
- **Level-one Account**: An account with `ParentAccountId` set to 0.
- **Sub-account**: An account linked to a level-one account. Currently, the system supports a two-level structure (Parent-Child).

### 2.3 Account Categories
Accounts are divided into different categories and classified as either **Asset** or **Liability**:

| Category | Classification | Description |
| :--- | :--- | :--- |
| Cash | Asset | Physical currency |
| Checking Account | Asset | Bank demand deposits |
| Credit Card | Liability | Credit limit that requires repayment |
| Virtual | Asset | Digital balances like PayPal, etc. |
| Debt | Liability | Money owed to others |
| Receivables | Asset | Money others owe to the user |
| Investment | Asset | Stocks, funds, etc. |
| Savings Account | Asset | Fixed-term or dedicated savings |
| Certificate of Deposit | Asset | Term deposits |

### 2.4 Key Attributes
- **Basic Info**: Name, Icon, Color, Currency, Comment.
- **State Control**: Visibility (Hidden), Display Order.
- **Balance**: Stored in real-time in the database, usually in the smallest currency unit (e.g., cents).

---

## 3. Business Logic

### 3.1 Balance Management
- **Initial Balance**: Set during account creation. The system automatically generates a system transaction of type `MODIFY_BALANCE` to record this initial state.
- **Dynamic Updates**:
    - **Income**: Increases the account balance.
    - **Expense**: Decreases the account balance.
    - **Transfer**: Decreases the source account balance and increases the destination account balance.
- **Atomicity**: Balance updates use atomic addition/subtraction at the database level to ensure data consistency under concurrent access.

### 3.2 Account Operations
- **Creation**: Supports bulk creation, such as initializing multiple sub-accounts when creating a "Multi-sub accounts" parent.
- **Modification**: Allows changing attributes like name and icon. For "Multi-sub accounts", it supports dynamically adding, modifying, or removing sub-accounts.
- **Hiding**: Hidden accounts do not appear in the selection list during bookkeeping but are retained for historical data and can still be included in asset statistics.
- **Deletion**:
    - Uses **Soft Delete** (marking the `Deleted` field).
    - **Deletion Constraints**: An account cannot be deleted if it has been used in actual transactions (beyond the initial balance adjustment) or is linked to scheduled tasks/transaction templates, ensuring financial data integrity.

### 3.3 Statistics and Reconciliation
- The system calculates total assets, total liabilities, and net assets based on account classifications.
- The account balance should logically match the sum of all associated transaction flows.

---

## 4. Technical Implementation Highlights
- **Backend (Golang)**: `pkg/services/accounts.go` handles business logic, interacting with the database via XORM and ensuring consistency with transactions (`DoTransaction`).
- **Frontend (Vue.js)**: Account state is managed in `src/stores/account.ts`, supporting multi-currency conversion for display.
- **Unique Identification**: Uses 64-bit integer UUIDs (generated via Snowflake or similar) as the `AccountId`.
