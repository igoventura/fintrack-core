# Feature Roadmap

This document tracks the implementation status of FinTrack Core features.

## Phase 1: Foundation (Current Status: 100%)
Core infrastructure and user management.

- [x] **Project Structure**: Clean Architecture setup, folders, Makefiles.
- [x] **Database**: Migration system (tern), Initial Schema.
- [x] **Authentication**: Supabase integration, JWT Validator, Auth Middleware.
- [x] **Multi-Tenancy**
  - [x] Tenant Middleware (Strict Validation)
  - [x] Tenant Repository
  - [x] API: Create Tenant Endpoint (Onboarding)
- [x] **User Management**
  - [x] User Repository
  - [x] Auth Service (Register/Login)
  - [x] API: Get/Update Profile
  - [x] API: List User Tenants
- [x] **Accounts**: Account Domain, Repository, Service, API Handlers.

## Phase 2: Core Classification (Next Priorities)
Organizing financial data.

- [x] **Categories** (authenticated endpoint and tenant-scoped)
  - [x] Repository
  - [x] Service Layer (Business Logic)
  - [x] API Handlers & Routes

- [x] **Tags** (authenticated endpoint and tenant-scoped)
  - [x] Repository
  - [x] Service Layer (Business Logic)
  - [x] API Handlers & Routes

## Phase 3: Financial Core
The heart of the application: tracking money movement.

- [ ] **Transactions** (authenticated endpoint and tenant-scoped)
  - [x] Repository (Basic CRUD)
  - [ ] Service Layer (Logic for Types: Credit, Debit, Transfer, Payment)
  - [ ] API Handlers & Routes
  - [ ] Link to Categories and Tags (Note: Must update `transactions_tags` schema to include `created_at`)
  - [ ] **Transaction Logic Rules**:
    - [ ] **Schema Migration**:
      - Add `currency` column (varchar(3)): Default to `FromAccount` currency.
      - Remove `PreviousSiblingTransactionID` and `NextSiblingTransactionID`.
      - Add `ParentTransactionID` (UUID, Nullable):
        - For Transfers: links credit transaction to original debit.
        - For Installments: links subsequent installments to first one.
      - Update `transactions_tags` schema to include `created_at`.
    - [ ] **Fields**:
      - `FromAccountID`: Source account.
      - `ToAccountID`: Recipient account (triggers credit transaction for transfers/payments).
      - `TenantID`: From context.
      - `AccrualMonth` (YYYYMM): Defaults to due date's month, or explicit from frontend.
      - `PaymentDate`: Tracks payment status. For CC, defaults to Due Date (except Payment transaction).
    - [ ] **Credit Card Installments**:
      - Input: `installments` (int) + `accrual_month`.
      - Logic: Splits value, generates N transactions.
      - Rounding: First installment absorbs remainder (e.g. 10/3 -> 3.34, 3.33, 3.33). Sum must match Amount.
      - Due Dates: Calculated based on accrual month + 1 month. Handles month-end logic (31st -> 28th/29th).

- [ ] **Attachments** (authenticated endpoint and tenant-scoped)
  - [x] Schema Support
  - [ ] File Upload Logic (Service)
  - [ ] Storage Provider Integration (e.g. S3/Supabase Storage) - *Pending Design*
  - [ ] API Handlers

## Phase 4: Extensions
Advanced features.

- [ ] **Credit Card Management** (authenticated endpoint and tenant-scoped)
  - [x] Schema Support (`credit_card_info`)
  - [ ] Domain & Repository
  - [ ] Service & API
  - [ ] Statement Closing/Due Date Logic

- [ ] **Reporting** (authenticated endpoint and tenant-scoped)
  - [ ] Aggregation Queries (Monthly Spend, Income vs Expense)
  - [ ] Dashboard Endpoints

- [ ] **Invitations** (authenticated and tenant-scoped create endpoint, unauthenticated accept endpoint, tenant is not required for accept endpoint)
  - [ ] Schema (`invitations` table)
    - id, inviter_user_id, email, tenant_id, status, expires_at, created_at, updated_at
  - [ ] Domain & Repository
  - [ ] Service (Send, Accept, Revoke logic)
  - [ ] API Handlers

## Infrastructure & Quality
Ongoing improvements.

- [x] Docker Composition (DB + App)
- [x] API Documentation (Swagger/Scalar)
- [ ] CI/CD Pipelines
- [ ] Unit Test Coverage (>80%)
