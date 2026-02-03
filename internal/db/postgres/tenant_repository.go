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
	query := `SELECT id, name, created_at, updated_at, deactivated_at FROM tenants WHERE id = $1`
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
	query := `INSERT INTO tenants (id, name, created_at, updated_at, deactivated_at)
			  VALUES ($1, $2, $3, $4, $5)`
	_, err := r.db.Pool.Exec(ctx, query, t.ID, t.Name, t.CreatedAt, t.UpdatedAt, t.DeactivatedAt)
	if err != nil {
		return fmt.Errorf("failed to create tenant: %w", err)
	}
	return nil
}

func (r *TenantRepository) Update(ctx context.Context, t *domain.Tenant) error {
	query := `UPDATE tenants SET name = $2, updated_at = $3, deactivated_at = $4 WHERE id = $1`
	_, err := r.db.Pool.Exec(ctx, query, t.ID, t.Name, t.UpdatedAt, t.DeactivatedAt)
	if err != nil {
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
			  JOIN tenant_users tu ON t.id = tu.tenant_id
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
