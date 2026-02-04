package postgres

import (
	"context"
	"fmt"

	"github.com/igoventura/fintrack-core/domain"
)

type CategoryRepository struct {
	db *DB
}

func NewCategoryRepository(db *DB) *CategoryRepository {
	return &CategoryRepository{db: db}
}

func (r *CategoryRepository) GetByID(ctx context.Context, id, tenantID string) (*domain.Category, error) {
	query := `SELECT id, parent_category, tenant_id, name, deactivated_at, color, icon, created_at, created_by, updated_at, updated_by, deactivated_by FROM categories WHERE id = $1 AND tenant_id = $2 AND deactivated_at IS NULL`
	var c domain.Category
	err := r.db.Pool.QueryRow(ctx, query, id, tenantID).Scan(
		&c.ID, &c.ParentCategoryID, &c.TenantID, &c.Name, &c.DeactivatedAt, &c.Color, &c.Icon,
		&c.CreatedAt, &c.CreatedBy, &c.UpdatedAt, &c.UpdatedBy, &c.DeactivatedBy,
	)
	if err != nil {
		return nil, fmt.Errorf("failed to get category by id: %w", err)
	}
	return &c, nil
}

func (r *CategoryRepository) List(ctx context.Context, tenantID string) ([]domain.Category, error) {
	query := `SELECT id, parent_category, tenant_id, name, deactivated_at, color, icon, created_at, created_by, updated_at, updated_by, deactivated_by FROM categories WHERE tenant_id = $1 AND deactivated_at IS NULL`
	rows, err := r.db.Pool.Query(ctx, query, tenantID)
	if err != nil {
		return nil, fmt.Errorf("failed to list categories: %w", err)
	}
	defer rows.Close()

	var categories []domain.Category
	for rows.Next() {
		var c domain.Category
		if err := rows.Scan(&c.ID, &c.ParentCategoryID, &c.TenantID, &c.Name, &c.DeactivatedAt, &c.Color, &c.Icon, &c.CreatedAt, &c.CreatedBy, &c.UpdatedAt, &c.UpdatedBy, &c.DeactivatedBy); err != nil {
			return nil, fmt.Errorf("failed to scan category: %w", err)
		}
		categories = append(categories, c)
	}
	return categories, nil
}

func (r *CategoryRepository) Create(ctx context.Context, c *domain.Category) error {
	query := `INSERT INTO categories (parent_category, tenant_id, name, color, icon, created_by, updated_by)
			  VALUES ($1, $2, $3, $4, $5, $6, $7)
			  RETURNING id, created_at, updated_at`
	row := r.db.Pool.QueryRow(ctx, query, c.ParentCategoryID, c.TenantID, c.Name, c.Color, c.Icon, c.CreatedBy, c.UpdatedBy)
	if err := row.Scan(&c.ID, &c.CreatedAt, &c.UpdatedAt); err != nil {
		return fmt.Errorf("failed to create category: %w", err)
	}
	return nil
}

func (r *CategoryRepository) Update(ctx context.Context, c *domain.Category) error {
	query := `UPDATE categories SET parent_category = $2, name = $3, color = $4, icon = $5, updated_by = $6, updated_at = CURRENT_TIMESTAMP WHERE id = $1 AND tenant_id = $7 RETURNING updated_at`
	err := r.db.Pool.QueryRow(ctx, query, c.ID, c.ParentCategoryID, c.Name, c.Color, c.Icon, c.UpdatedBy, c.TenantID).Scan(&c.UpdatedAt)
	if err != nil {
		return fmt.Errorf("failed to update category: %w", err)
	}
	return nil
}

func (r *CategoryRepository) Delete(ctx context.Context, id, tenantID, userID string) error {
	query := `UPDATE categories SET deactivated_at = CURRENT_TIMESTAMP, deactivated_by = $2 WHERE id = $1 AND tenant_id = $3`
	_, err := r.db.Pool.Exec(ctx, query, id, userID, tenantID)
	if err != nil {
		return fmt.Errorf("failed to delete category: %w", err)
	}
	return nil
}
