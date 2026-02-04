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

- [ ] **Transactions** (authenticated endpoint and tenant-scoped) - **IN PROGRESS**
  - [x] Schema Migration
      - [x] Add `currency` column to `transactions` table.
      - [x] Remove `previous_sibling_transaction_id` and `next_sibling_transaction_id`.
      - [x] Add `parent_transaction_id` column for splitting transactions.
      - [x] Add `created_at`, `updated_at`, `deactivated_at` to `transactions_tags` table.
  - [x] Domain & Repository
      - [x] Update `Transaction` entity.
      - [x] Update `TransactionRepository` (CRUD + Filters).
      - [x] Update `TransactionRepository` for bulk tag insertion (`AddTagsToTransaction`).
  - [x] Service Layer
      - [x] Implement `CreateTransaction` with:
          - [x] Default currency logic (from Account).
          - [x] Tag association (list of IDs).
          - [x] Basic validation.
  - [x] API Layer
      - [x] Create DTOs.
      - [x] Implement `TransactionHandler`.
      - [x] Register Routes.
  - [x] Link to Categories and Tags
  - [ ] **Transaction Logic Rules**:
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
