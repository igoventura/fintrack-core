package postgres

import (
	"context"
	"fmt"

	"github.com/igoventura/fintrack-api/domain"
)

type TransactionRepository struct {
	db *DB
}

func NewTransactionRepository(db *DB) *TransactionRepository {
	return &TransactionRepository{db: db}
}

func (r *TransactionRepository) GetByID(ctx context.Context, tenantID, id string) (*domain.Transaction, error) {
	query := `SELECT id, parent_transaction_id, tenant_id, from_account_id, to_account_id, currency, amount, accrual_month, transaction_type, category_id, comments, due_date, payment_date, created_at, created_by, updated_at, updated_by, deactivated_at, deactivated_by FROM transactions WHERE id = $1 AND tenant_id = $2 AND deactivated_at IS NULL`
	var t domain.Transaction
	err := r.db.Pool.QueryRow(ctx, query, id, tenantID).Scan(
		&t.ID, &t.ParentTransactionID, &t.TenantID, &t.FromAccountID, &t.ToAccountID, &t.Currency, &t.Amount, &t.AccrualMonth, &t.TransactionType, &t.CategoryID, &t.Comments, &t.DueDate, &t.PaymentDate, &t.CreatedAt, &t.CreatedBy, &t.UpdatedAt, &t.UpdatedBy, &t.DeactivatedAt, &t.DeactivatedBy,
	)
	if err != nil {
		return nil, fmt.Errorf("failed to get transaction by id: %w", err)
	}
	return &t, nil
}

func (r *TransactionRepository) List(ctx context.Context, tenantID string, filter domain.TransactionFilter) ([]domain.Transaction, error) {
	query := `SELECT id, parent_transaction_id, tenant_id, from_account_id, to_account_id, currency, amount, accrual_month, transaction_type, category_id, comments, due_date, payment_date, created_at, created_by, updated_at, updated_by, deactivated_at, deactivated_by FROM transactions WHERE tenant_id = $1 AND deactivated_at IS NULL`
	args := []interface{}{tenantID}
	argIdx := 2

	if filter.AccrualMonth != "" {
		query += fmt.Sprintf(" AND accrual_month = $%d", argIdx)
		args = append(args, filter.AccrualMonth)
		argIdx++
	}
	if filter.AccountID != "" {
		query += fmt.Sprintf(" AND from_account_id = $%d", argIdx)
		args = append(args, filter.AccountID)
		argIdx++
	}
	if filter.TransactionType != "" {
		query += fmt.Sprintf(" AND transaction_type = $%d", argIdx)
		args = append(args, filter.TransactionType)
		argIdx++
	}

	rows, err := r.db.Pool.Query(ctx, query, args...)
	if err != nil {
		return nil, fmt.Errorf("failed to list transactions: %w", err)
	}
	defer rows.Close()

	var transactions []domain.Transaction
	for rows.Next() {
		var t domain.Transaction
		if err := rows.Scan(&t.ID, &t.ParentTransactionID, &t.TenantID, &t.FromAccountID, &t.ToAccountID, &t.Currency, &t.Amount, &t.AccrualMonth, &t.TransactionType, &t.CategoryID, &t.Comments, &t.DueDate, &t.PaymentDate, &t.CreatedAt, &t.CreatedBy, &t.UpdatedAt, &t.UpdatedBy, &t.DeactivatedAt, &t.DeactivatedBy); err != nil {
			return nil, fmt.Errorf("failed to scan transaction: %w", err)
		}
		transactions = append(transactions, t)
	}
	return transactions, nil
}

func (r *TransactionRepository) Create(ctx context.Context, t *domain.Transaction) error {
	query := `INSERT INTO transactions (parent_transaction_id, tenant_id, from_account_id, to_account_id, currency, amount, accrual_month, transaction_type, category_id, comments, due_date, payment_date, created_by, updated_by)
			  VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13, $14)
			  RETURNING id, created_at, updated_at`
	row := r.db.Pool.QueryRow(ctx, query, t.ParentTransactionID, t.TenantID, t.FromAccountID, t.ToAccountID, t.Currency, t.Amount, t.AccrualMonth, t.TransactionType, t.CategoryID, t.Comments, t.DueDate, t.PaymentDate, t.CreatedBy, t.UpdatedBy)
	if err := row.Scan(&t.ID, &t.CreatedAt, &t.UpdatedAt); err != nil {
		return fmt.Errorf("failed to create transaction: %w", err)
	}
	t.UpdatedBy = t.CreatedBy // Initial state
	return nil
}

func (r *TransactionRepository) Update(ctx context.Context, t *domain.Transaction) error {
	query := `UPDATE transactions SET parent_transaction_id = $2, from_account_id = $3, to_account_id = $4, currency = $5, amount = $6, accrual_month = $7, transaction_type = $8, category_id = $9, comments = $10, due_date = $11, payment_date = $12, updated_at = CURRENT_TIMESTAMP, updated_by = $13 WHERE id = $1 AND tenant_id = $14 RETURNING updated_at`
	row := r.db.Pool.QueryRow(ctx, query, t.ID, t.ParentTransactionID, t.FromAccountID, t.ToAccountID, t.Currency, t.Amount, t.AccrualMonth, t.TransactionType, t.CategoryID, t.Comments, t.DueDate, t.PaymentDate, t.UpdatedBy, t.TenantID)
	if err := row.Scan(&t.UpdatedAt); err != nil {
		return fmt.Errorf("failed to update transaction: %w", err)
	}
	return nil
}

func (r *TransactionRepository) Delete(ctx context.Context, tenantID, id string, userID string) error {
	query := `UPDATE transactions SET deactivated_at = CURRENT_TIMESTAMP, deactivated_by = $2 WHERE id = $1 AND tenant_id = $3`
	_, err := r.db.Pool.Exec(ctx, query, id, userID, tenantID)
	if err != nil {
		return fmt.Errorf("failed to delete transaction: %w", err)
	}
	return nil
}

func (r *TransactionRepository) AddTagsToTransaction(ctx context.Context, transactionID string, tagIDs []string) error {
	if len(tagIDs) == 0 {
		return nil
	}

	query := `INSERT INTO transactions_tags (transaction_id, tag_id) VALUES `
	values := []interface{}{}
	for i, tagID := range tagIDs {
		n := i * 2
		query += fmt.Sprintf("($%d, $%d),", n+1, n+2)
		values = append(values, transactionID, tagID)
	}
	query = query[:len(query)-1] // Remove trailing comma

	_, err := r.db.Pool.Exec(ctx, query, values...)
	if err != nil {
		return fmt.Errorf("failed to add tags to transaction: %w", err)
	}
	return nil
}

func (r *TransactionRepository) ReplaceTags(ctx context.Context, transactionID string, tagIDs []string) error {
	tx, err := r.db.Pool.Begin(ctx)
	if err != nil {
		return fmt.Errorf("failed to begin transaction: %w", err)
	}
	defer tx.Rollback(ctx)

	// 1. Delete existing tags
	deleteQuery := `DELETE FROM transactions_tags WHERE transaction_id = $1`
	if _, err := tx.Exec(ctx, deleteQuery, transactionID); err != nil {
		return fmt.Errorf("failed to delete existing tags: %w", err)
	}

	// 2. Insert new tags (if any)
	if len(tagIDs) > 0 {
		insertQuery := `INSERT INTO transactions_tags (transaction_id, tag_id) VALUES `
		values := []interface{}{}
		for i, tagID := range tagIDs {
			n := i * 2
			insertQuery += fmt.Sprintf("($%d, $%d),", n+1, n+2)
			values = append(values, transactionID, tagID)
		}
		insertQuery = insertQuery[:len(insertQuery)-1] // Remove trailing comma

		if _, err := tx.Exec(ctx, insertQuery, values...); err != nil {
			return fmt.Errorf("failed to insert new tags: %w", err)
		}
	}

	if err := tx.Commit(ctx); err != nil {
		return fmt.Errorf("failed to commit transaction: %w", err)
	}
	return nil
}

func (r *TransactionRepository) RemoveTagFromTransaction(ctx context.Context, transactionID, tagID string) error {
	query := `DELETE FROM transactions_tags WHERE transaction_id = $1 AND tag_id = $2`
	_, err := r.db.Pool.Exec(ctx, query, transactionID, tagID)
	if err != nil {
		return fmt.Errorf("failed to remove tag from transaction: %w", err)
	}
	return nil
}

func (r *TransactionRepository) ListTransactionTags(ctx context.Context, transactionID string) ([]domain.Tag, error) {
	query := `SELECT t.id, t.tenant_id, t.name, t.deactivated_at FROM tags t
			  JOIN transactions_tags tt ON t.id = tt.tag_id
			  WHERE tt.transaction_id = $1 AND t.deactivated_at IS NULL`
	rows, err := r.db.Pool.Query(ctx, query, transactionID)
	if err != nil {
		return nil, fmt.Errorf("failed to list transaction tags: %w", err)
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

func (r *TransactionRepository) AddAttachment(ctx context.Context, a *domain.TransactionAttachment) error {
	query := `INSERT INTO transaction_attachments (transaction_id, name, path, created_by, updated_by)
			  VALUES ($1, $2, $3, $4, $5)
			  RETURNING id, created_at, updated_at`
	row := r.db.Pool.QueryRow(ctx, query, a.TransactionID, a.Name, a.Path, a.CreatedBy, a.UpdatedBy)
	if err := row.Scan(&a.ID, &a.CreatedAt, &a.UpdatedAt); err != nil {
		return fmt.Errorf("failed to add attachment: %w", err)
	}
	return nil
}

func (r *TransactionRepository) RemoveAttachment(ctx context.Context, id string, userID string) error {
	query := `UPDATE transaction_attachments SET deactivated_at = CURRENT_TIMESTAMP, deactivated_by = $2 WHERE id = $1`
	_, err := r.db.Pool.Exec(ctx, query, id, userID)
	if err != nil {
		return fmt.Errorf("failed to remove attachment: %w", err)
	}
	return nil
}

func (r *TransactionRepository) ListAttachments(ctx context.Context, transactionID string) ([]domain.TransactionAttachment, error) {
	query := `SELECT id, transaction_id, name, path, created_at, created_by, updated_at, updated_by, deactivated_at, deactivated_by FROM transaction_attachments WHERE transaction_id = $1 AND deactivated_at IS NULL`
	rows, err := r.db.Pool.Query(ctx, query, transactionID)
	if err != nil {
		return nil, fmt.Errorf("failed to list attachments: %w", err)
	}
	defer rows.Close()

	var attachments []domain.TransactionAttachment
	for rows.Next() {
		var a domain.TransactionAttachment
		if err := rows.Scan(&a.ID, &a.TransactionID, &a.Name, &a.Path, &a.CreatedAt, &a.CreatedBy, &a.UpdatedAt, &a.UpdatedBy, &a.DeactivatedAt, &a.DeactivatedBy); err != nil {
			return nil, fmt.Errorf("failed to scan attachment: %w", err)
		}
		attachments = append(attachments, a)
	}
	return attachments, nil
}
