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

func (r *CategoryRepository) GetByID(ctx context.Context, id string) (*domain.Category, error) {
	query := `SELECT id, parent_category, tenant_id, name, deactivated_at, color, icon FROM categories WHERE id = $1 AND deactivated_at IS NULL`
	var c domain.Category
	err := r.db.Pool.QueryRow(ctx, query, id).Scan(
		&c.ID, &c.ParentCategoryID, &c.TenantID, &c.Name, &c.DeactivatedAt, &c.Color, &c.Icon,
	)
	if err != nil {
		return nil, fmt.Errorf("failed to get category by id: %w", err)
	}
	return &c, nil
}

func (r *CategoryRepository) List(ctx context.Context, tenantID string) ([]domain.Category, error) {
	query := `SELECT id, parent_category, tenant_id, name, deactivated_at, color, icon FROM categories WHERE tenant_id = $1 AND deactivated_at IS NULL`
	rows, err := r.db.Pool.Query(ctx, query, tenantID)
	if err != nil {
		return nil, fmt.Errorf("failed to list categories: %w", err)
	}
	defer rows.Close()

	var categories []domain.Category
	for rows.Next() {
		var c domain.Category
		if err := rows.Scan(&c.ID, &c.ParentCategoryID, &c.TenantID, &c.Name, &c.DeactivatedAt, &c.Color, &c.Icon); err != nil {
			return nil, fmt.Errorf("failed to scan category: %w", err)
		}
		categories = append(categories, c)
	}
	return categories, nil
}

func (r *CategoryRepository) Create(ctx context.Context, c *domain.Category) error {
	query := `INSERT INTO categories (parent_category, tenant_id, name, color, icon)
			  VALUES ($1, $2, $3, $4, $5)
			  RETURNING id`
	row := r.db.Pool.QueryRow(ctx, query, c.ParentCategoryID, c.TenantID, c.Name, c.Color, c.Icon)
	if err := row.Scan(&c.ID); err != nil {
		return fmt.Errorf("failed to create category: %w", err)
	}
	return nil
}

func (r *CategoryRepository) Update(ctx context.Context, c *domain.Category) error {
	query := `UPDATE categories SET parent_category = $2, name = $3, color = $4, icon = $5 WHERE id = $1`
	_, err := r.db.Pool.Exec(ctx, query, c.ID, c.ParentCategoryID, c.Name, c.Color, c.Icon)
	if err != nil {
		return fmt.Errorf("failed to update category: %w", err)
	}
	return nil
}

func (r *CategoryRepository) Delete(ctx context.Context, id string) error {
	query := `UPDATE categories SET deactivated_at = CURRENT_TIMESTAMP WHERE id = $1`
	_, err := r.db.Pool.Exec(ctx, query, id)
	if err != nil {
		return fmt.Errorf("failed to delete category: %w", err)
	}
	return nil
}
