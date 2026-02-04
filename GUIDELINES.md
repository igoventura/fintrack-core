# Fintrack Core Development Guidelines

These guidelines are mandatory for all contributions to `fintrack-core`. They ensure consistency, maintainability, and adherence to the project's architectural standards.

## 1. Project Structure

We follow **Clean Architecture** principles. See [PROJECT_STRUCTURE.md](./PROJECT_STRUCTURE.md) for a detailed breakdown of the directory tree.

### Quick Reference
- **`/domain`**: Pure Go entities and repository interfaces. NO external dependencies.
- **`/internal/service`**: Business logic. Orchestrates domain entities.
- **`/internal/api`**: HTTP transport (Handlers, DTOs, Router).
- **`/internal/db`**: Database implementations (Postgres).

---

## 2. Implementation Workflow (How to Add a Feature)

Follow this order when implementing a new feature (e.g., "Add Expense"):

1.  **Domain Layer** (`/domain`)
    *   Define the Entity struct (e.g., `Expense`) with JSON tags excluded (or managed carefully).
    *   Add validation logic: `func (e *Expense) IsValid() (bool, map[string]string)`.
    *   Define the Repository Interface (e.g., `ExpenseRepository`).

2.  **Database Layer** (`/internal/db/postgres`)
    *   Implement the Repository Interface.
    *   Write SQL queries (use `pgx`).
    *   **Rule**: Ensure Soft Deletes and Timestamp handling (see Section 4).

3.  **Service Layer** (`/internal/service`)
    *   Define the Service struct and constructor.
    *   Implement business methods (e.g., `CreateExpense`).
    *   Inject dependencies (Repositories) via interfaces.

4.  **API Layer** (`/internal/api`)
    *   **DTOs**: Create request/response structs in `/api/dto`. Add `binding` tags for validation.
    *   **Handlers**: Create the handler in `/api/handler`. Map DTO -> Entity -> Service -> Entity -> DTO.
    *   **Router**: Register the new route in `/api/router/router.go`.

5.  **Documentation**
    *   Add Swagger annotations to the Handler.
    *   Run `make swagger` to update `docs/swagger.yaml`.

---

## 3. Architecture Rules

### Domain Layer
*   **Purity**: Must NOT import `internal/...`.
*   **Entities**: Should support `IsValid()` for self-validation.

### Service Layer
*   **Agnostic**: NO knowledge of HTTP, DTOs, or specific database implementations.
*   **Input/Output**: Receives Keys/Entities -> Returns Entities/Errors.

### API Layer
*   **Validation**: Handled first by Gin binding tags (in DTOs), then by Domain validation.
*   **Mapping**: Explicitly map DTOs to Entities and vice-versa. Do not pass DTOs to the service.
*   **Context**:
    *   **Authentication**: Use `domain.GetUserID(ctx)` to get the authenticated User ID.
    *   **Multi-Tenancy**: Use `domain.GetTenantID(ctx)` to get the Tenant ID.
    *   *Note*: The middleware filters requests without valid Tenant IDs, so handlers/services can assume it exists. Do NOT manually check for empty ID.

### Tenant Isolation Protocol
*   **Repository Layer**: ALL methods (Get, List, Update, Delete) MUST accept `tenantID` as an argument and filter by `tenant_id = $x` in the SQL query.
    *   *Exception*: Global lookups (if any) or user-scoped queries that are not tenant-bound.
*   **Service Layer**: Retrieve `tenantID` from `domain.GetTenantID(ctx)` and pass it to the Repository. Do not accept `tenantID` as a function argument from the Handler (unless it's a specific requirement).

---

## 4. Coding Standards

### Database & Persistence

#### Soft Deletes
*   **Filtering**: All `Select`, `Get`, `List` queries **MUST** filter out soft-deleted records.
    ```sql
    WHERE deactivated_at IS NULL
    ```
*   **Deleting**: Never use `DELETE FROM`. Update the record:
    ```sql
    UPDATE table SET deactivated_at = NOW() WHERE id = $1
    ```

#### Timestamps
*   **Return Values**: `Create` and `Update` operations **MUST** return the generated `created_at` and `updated_at` timestamps from the database (using `RETURNING` clause) to ensure the client has the authoritative time.

#### Migrations
*   **Immutable History**: NEVER modify an already applied migration file.
*   **New Changes**: Create a new migration using `make new-migration name=what_changed`.

### Error Handling
*   **Sentinel Errors**: Define standard errors in `domain/errors.go` (if applicable) or usage `errors.New`.
*   **Wrapping**: Wrap errors in Repositories and Services to provide context.
    ```go
    if err != nil {
        return fmt.Errorf("repository failed to get user: %w", err)
    }
    ```
*   **Checking**: Use `errors.Is(err, target)` to check for specific error types (e.g., `pgx.ErrNoRows`).

### Configuration
*   Use environment variables.
*   Structure configurations in `/internal/config`.

---

## 5. API Documentation Standards (Swagger)

We use `swag` to generate OpenAPI documentation.

### General Tags
*   `@Summary`: Short description.
*   `@Description`: Detailed description.
*   `@Tags`: Resource group (e.g., `users`, `accounts`).
*   `@Accept`, `@Produce`: Usually `json`.

### Security & Multi-Tenancy
*   **Authenticated Endpoints**: MUST include:
    ```go
    // @Security AuthPassword
    ```
*   **Tenant-Specific Endpoints**: MUST include:
    ```go
    // @Param X-Tenant-ID header string true "Tenant ID"
    ```

### Example
```go
// CreateAccount creates a new account
// @Summary Create account
// @Tags accounts
// @Accept json
// @Produce json
// @Security AuthPassword
// @Param X-Tenant-ID header string true "Tenant ID"
// @Param request body dto.CreateAccountRequest true "Account Payload"
// @Success 201 {object} dto.AccountResponse
// @Router /accounts [post]
func (h *AccountHandler) CreateAccount(c *gin.Context) { ... }
```

---

## 6. Testing Strategy

*(Currently in early stages, but aim for the following)*

*   **Unit Tests**: Internal logic (Service layer). Mock repositories.
*   **Integration Tests**: Handlers and Database.

To run tests:
```bash
make test
```

---

## 7. Task & Git Discipline

*   **Roadmap**: Before starting, check `ROADMAP.md` to see what needs to be done.
*   **Scope**: Only implement what is requested.
*   **Docs**:
    *   Update `PROJECT_STRUCTURE.md` if you add files.
    *   Update `README.md` if you add env vars.
