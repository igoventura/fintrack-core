# Development Guidelines

These guidelines are mandatory for all contributions to `fintrack-core`.

## 1. Architecture Rules

We follow **Clean Architecture**.

- **Domain Layer (`/domain`)**:
    - **NO** imports from `internal/...`.
    - Must contain pure Go structs and interfaces.
    - Entities must implement `IsValid()`.

- **Service Layer (`/internal/service`)**:
    - Orchestrates `domain` logic.
    - **NO** knowledge of HTTP, DTOs, or JSON tags.
    - Returns `domain` entities or standard errors.

- **API Layer (`/internal/api`)**:
    - Handles HTTP, Parsing, and Validation (Validation tags go here in DTOs).
    - Maps **DTO <-> Entity** explicitly (Mappers belong here).
    - Calls Service Layer.

## 2. Coding Standards

### Database & Persistence
- **Soft Deletes**:
    - **READ**: All `List` and `Get` queries **MUST** filter out deleted records:
      ```sql
      WHERE deactivated_at IS NULL
      ```
    - **DELETE**: Never `DELETE FROM`. Update `deactivated_at = NOW()`.

- **Timestamps**:
    - Mutations (`Create`/`Update`) **MUST** return the generated `created_at` and `updated_at` timestamps to the caller.

### Error Handling
- Use `errors.Is`/`errors.As` for checking sentinel errors.
- Wrap errors in the Service/Repository layers with context (e.g., `fmt.Errorf("failed to get user: %w", err)`).

## 3. Workflow & Documentation

- **Task Discipline**:
    - Only implement what is explicitly requested. If unsure, ask first.
    - Before starting, check `ROADMAP.md` and `PROJECT_STRUCTURE.md`.

- **Documentation Updates**:
    - **MUST** update `PROJECT_STRUCTURE.md` if you add/move files.
    - **MUST** update `README.md` if you add new features or env vars.
    - **Swagger**: Run `make swagger` after modifying API Handlers or DTOs.

- **Migrations**:
    - **NEVER** modify an existing migration file that has been applied (unless in local dev and you reset DB).
    - Create new migrations for schema changes: `make new-migration name=description`.
