package postgres

import (
	"context"
	"fmt"

	"github.com/igoventura/fintrack-api/domain"
)

type TagRepository struct {
	db *DB
}

func NewTagRepository(db *DB) *TagRepository {
	return &TagRepository{db: db}
}

func (r *TagRepository) GetByID(ctx context.Context, id, tenantID string) (*domain.Tag, error) {
	query := `SELECT id, tenant_id, name, created_at, created_by, updated_at, updated_by, deactivated_at, deactivated_by FROM tags WHERE id = $1 AND tenant_id = $2 AND deactivated_at IS NULL`
	var t domain.Tag
	err := r.db.Pool.QueryRow(ctx, query, id, tenantID).Scan(
		&t.ID, &t.TenantID, &t.Name, &t.CreatedAt, &t.CreatedBy, &t.UpdatedAt, &t.UpdatedBy, &t.DeactivatedAt, &t.DeactivatedBy,
	)
	if err != nil {
		return nil, fmt.Errorf("failed to get tag by id: %w", err)
	}
	return &t, nil
}

func (r *TagRepository) List(ctx context.Context, tenantID string) ([]domain.Tag, error) {
	query := `SELECT id, tenant_id, name, created_at, created_by, updated_at, updated_by, deactivated_at, deactivated_by FROM tags WHERE tenant_id = $1 AND deactivated_at IS NULL`
	rows, err := r.db.Pool.Query(ctx, query, tenantID)
	if err != nil {
		return nil, fmt.Errorf("failed to list tags: %w", err)
	}
	defer rows.Close()

	var tags []domain.Tag
	for rows.Next() {
		var t domain.Tag
		if err := rows.Scan(&t.ID, &t.TenantID, &t.Name, &t.CreatedAt, &t.CreatedBy, &t.UpdatedAt, &t.UpdatedBy, &t.DeactivatedAt, &t.DeactivatedBy); err != nil {
			return nil, fmt.Errorf("failed to scan tag: %w", err)
		}
		tags = append(tags, t)
	}
	return tags, nil
}

func (r *TagRepository) Create(ctx context.Context, t *domain.Tag) error {
	query := `INSERT INTO tags (tenant_id, name, created_by, updated_by)
			  VALUES ($1, $2, $3, $3)
			  RETURNING id, created_at, updated_at`
	row := r.db.Pool.QueryRow(ctx, query, t.TenantID, t.Name, t.CreatedBy)
	if err := row.Scan(&t.ID, &t.CreatedAt, &t.UpdatedAt); err != nil {
		return fmt.Errorf("failed to create tag: %w", err)
	}
	// Initial state setup
	t.UpdatedBy = t.CreatedBy
	return nil
}

func (r *TagRepository) Update(ctx context.Context, t *domain.Tag) error {
	query := `UPDATE tags SET name = $2, updated_by = $3, updated_at = CURRENT_TIMESTAMP WHERE id = $1 AND tenant_id = $4 RETURNING updated_at`
	if err := r.db.Pool.QueryRow(ctx, query, t.ID, t.Name, t.UpdatedBy, t.TenantID).Scan(&t.UpdatedAt); err != nil {
		return fmt.Errorf("failed to update tag: %w", err)
	}
	return nil
}

func (r *TagRepository) Delete(ctx context.Context, id, tenantID, userID string) error {
	query := `UPDATE tags SET deactivated_at = CURRENT_TIMESTAMP, deactivated_by = $3 WHERE id = $1 AND tenant_id = $2`
	_, err := r.db.Pool.Exec(ctx, query, id, tenantID, userID)
	if err != nil {
		return fmt.Errorf("failed to delete tag: %w", err)
	}
	return nil
}
func (r *TagRepository) ValidateTags(ctx context.Context, tenantID string, tagIDs []string) (bool, error) {
	if len(tagIDs) == 0 {
		return true, nil
	}
	query := `SELECT COUNT(*) FROM tags WHERE tenant_id = $1 AND id = ANY($2) AND deactivated_at IS NULL`
	var count int
	err := r.db.Pool.QueryRow(ctx, query, tenantID, tagIDs).Scan(&count)
	if err != nil {
		return false, fmt.Errorf("failed to validate tags: %w", err)
	}

	return count == len(tagIDs), nil
}
