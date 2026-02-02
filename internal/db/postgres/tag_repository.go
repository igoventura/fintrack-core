package postgres

import (
	"context"
	"fmt"

	"github.com/igoventura/fintrack-core/domain"
)

type TagRepository struct {
	db *DB
}

func NewTagRepository(db *DB) *TagRepository {
	return &TagRepository{db: db}
}

func (r *TagRepository) GetByID(ctx context.Context, id string) (*domain.Tag, error) {
	query := `SELECT id, tenant_id, name, deactivated_at FROM tags WHERE id = $1`
	var t domain.Tag
	err := r.db.Pool.QueryRow(ctx, query, id).Scan(
		&t.ID, &t.TenantID, &t.Name, &t.DeactivatedAt,
	)
	if err != nil {
		return nil, fmt.Errorf("failed to get tag by id: %w", err)
	}
	return &t, nil
}

func (r *TagRepository) List(ctx context.Context, tenantID string) ([]domain.Tag, error) {
	query := `SELECT id, tenant_id, name, deactivated_at FROM tags WHERE tenant_id = $1`
	rows, err := r.db.Pool.Query(ctx, query, tenantID)
	if err != nil {
		return nil, fmt.Errorf("failed to list tags: %w", err)
	}
	defer rows.Close()

	var tags []domain.Tag
	for rows.Next() {
		var t domain.Tag
		if err := rows.Scan(&t.ID, &t.TenantID, &t.Name, &t.DeactivatedAt); err != nil {
			return nil, fmt.Errorf("failed to scan tag: %w", err)
		}
		tags = append(tags, t)
	}
	return tags, nil
}

func (r *TagRepository) Create(ctx context.Context, t *domain.Tag) error {
	query := `INSERT INTO tags (id, tenant_id, name, deactivated_at)
			  VALUES ($1, $2, $3, $4)`
	_, err := r.db.Pool.Exec(ctx, query, t.ID, t.TenantID, t.Name, t.DeactivatedAt)
	if err != nil {
		return fmt.Errorf("failed to create tag: %w", err)
	}
	return nil
}

func (r *TagRepository) Update(ctx context.Context, t *domain.Tag) error {
	query := `UPDATE tags SET name = $2, deactivated_at = $3 WHERE id = $1`
	_, err := r.db.Pool.Exec(ctx, query, t.ID, t.Name, t.DeactivatedAt)
	if err != nil {
		return fmt.Errorf("failed to update tag: %w", err)
	}
	return nil
}

func (r *TagRepository) Delete(ctx context.Context, id string) error {
	query := `UPDATE tags SET deactivated_at = CURRENT_TIMESTAMP WHERE id = $1`
	_, err := r.db.Pool.Exec(ctx, query, id)
	if err != nil {
		return fmt.Errorf("failed to delete tag: %w", err)
	}
	return nil
}
