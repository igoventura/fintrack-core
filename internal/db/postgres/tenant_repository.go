package postgres

import (
	"context"
	"fmt"

	"github.com/igoventura/fintrack-core/domain"
)

type TenantRepository struct {
	db *DB
}

func NewTenantRepository(db *DB) *TenantRepository {
	return &TenantRepository{db: db}
}

func (r *TenantRepository) GetByID(ctx context.Context, id string) (*domain.Tenant, error) {
	query := `SELECT id, name, created_at, updated_at, deactivated_at FROM tenants WHERE id = $1 AND deactivated_at IS NULL`
	var t domain.Tenant
	err := r.db.Pool.QueryRow(ctx, query, id).Scan(
		&t.ID, &t.Name, &t.CreatedAt, &t.UpdatedAt, &t.DeactivatedAt,
	)
	if err != nil {
		return nil, fmt.Errorf("failed to get tenant by id: %w", err)
	}
	return &t, nil
}

func (r *TenantRepository) Create(ctx context.Context, t *domain.Tenant) error {
	query := `INSERT INTO tenants (name)
			  VALUES ($1)
			  RETURNING id, created_at, updated_at`
	row := r.db.Pool.QueryRow(ctx, query, t.Name)
	if err := row.Scan(&t.ID, &t.CreatedAt, &t.UpdatedAt); err != nil {
		return fmt.Errorf("failed to create tenant: %w", err)
	}
	return nil
}

func (r *TenantRepository) Update(ctx context.Context, t *domain.Tenant) error {
	query := `UPDATE tenants SET name = $2, updated_at = CURRENT_TIMESTAMP WHERE id = $1 RETURNING updated_at`
	row := r.db.Pool.QueryRow(ctx, query, t.ID, t.Name)
	if err := row.Scan(&t.UpdatedAt); err != nil {
		return fmt.Errorf("failed to update tenant: %w", err)
	}
	return nil
}

func (r *TenantRepository) Delete(ctx context.Context, id string) error {
	query := `UPDATE tenants SET deactivated_at = CURRENT_TIMESTAMP WHERE id = $1`
	_, err := r.db.Pool.Exec(ctx, query, id)
	if err != nil {
		return fmt.Errorf("failed to delete tenant: %w", err)
	}
	return nil
}

func (r *TenantRepository) ListByUserID(ctx context.Context, userID string) ([]domain.Tenant, error) {
	query := `SELECT t.id, t.name, t.created_at, t.updated_at, t.deactivated_at
			  FROM tenants t
			  JOIN users_tenants tu ON t.id = tu.tenant_id
			  WHERE tu.user_id = $1 AND t.deactivated_at IS NULL`
	rows, err := r.db.Pool.Query(ctx, query, userID)
	if err != nil {
		return nil, fmt.Errorf("failed to list tenants by user id: %w", err)
	}
	defer rows.Close()

	var tenants []domain.Tenant
	for rows.Next() {
		var t domain.Tenant
		if err := rows.Scan(&t.ID, &t.Name, &t.CreatedAt, &t.UpdatedAt, &t.DeactivatedAt); err != nil {
			return nil, fmt.Errorf("failed to scan tenant: %w", err)
		}
		tenants = append(tenants, t)
	}
	return tenants, nil
}
