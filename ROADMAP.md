# Feature Roadmap

This document tracks the implementation status of FinTrack Core features.

## Phase 1: Foundation (Current Status: ~90%)
Core infrastructure and user management.

- [x] **Project Structure**: Clean Architecture setup, folders, Makefiles.
- [x] **Database**: Migration system (tern), Initial Schema.
- [x] **Authentication**: Supabase integration, JWT Validator, Auth Middleware.
- [ ] **Multi-Tenancy** (~80%)
  - [x] Tenant Middleware (Strict Validation)
  - [x] Tenant Repository
  - [x] API: Create Tenant Endpoint (Onboarding)
- [ ] **User Management** (~80%)
  - [x] User Repository
  - [x] Auth Service (Register/Login)
  - [ ] API: Get/Update Profile
  - [ ] API: List User Tenants
- [x] **Accounts**: Account Domain, Repository, Service, API Handlers.

## Phase 2: Core Classification (Next Priorities)
Organizing financial data.

- [ ] **Categories**
  - [x] Repository
  - [ ] Service Layer (Business Logic)
  - [ ] API Handlers & Routes
  - [ ] Integration Tests

- [ ] **Tags**
  - [x] Repository
  - [ ] Service Layer (Business Logic)
  - [ ] API Handlers & Routes

## Phase 3: Financial Core
The heart of the application: tracking money movement.

- [ ] **Transactions**
  - [x] Repository (Basic CRUD)
  - [ ] Service Layer (Logic for Types: Credit, Debit, Transfer, Payment)
  - [ ] API Handlers & Routes
  - [ ] Link to Categories and Tags
  - [ ] Validation (Balance checks, same currency logic - *To Be Decided*)

- [ ] **Attachments**
  - [x] Schema Support
  - [ ] File Upload Logic (Service)
  - [ ] Storage Provider Integration (e.g. S3/Supabase Storage) - *Pending Design*
  - [ ] API Handlers

## Phase 4: Extensions
Advanced features.

- [ ] **Credit Card Management**
  - [x] Schema Support (`credit_card_info`)
  - [ ] Domain & Repository
  - [ ] Service & API
  - [ ] Statement Closing/Due Date Logic

- [ ] **Reporting**
  - [ ] Aggregation Queries (Monthly Spend, Income vs Expense)
  - [ ] Dashboard Endpoints

- [ ] **Invitations** (New)
  - [ ] Schema (`invitations` table)
  - [ ] Domain & Repository
  - [ ] Service (Send, Accept, Revoke logic)
  - [ ] API Handlers

## Infrastructure & Quality
Ongoing improvements.

- [x] Docker Composition (DB + App)
- [x] API Documentation (Swagger/Scalar)
- [ ] CI/CD Pipelines
- [ ] Unit Test Coverage (>80%)
