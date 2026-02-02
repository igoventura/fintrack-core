package postgres

import (
	"context"
	"fmt"

	"github.com/igoventura/fintrack-core/domain"
)

type TransactionRepository struct {
	db *DB
}

func NewTransactionRepository(db *DB) *TransactionRepository {
	return &TransactionRepository{db: db}
}

func (r *TransactionRepository) GetByID(ctx context.Context, id string) (*domain.Transaction, error) {
	query := `SELECT id, previous_sibling_transaction_id, next_sibling_transaction_id, tenant_id, from_account_id, to_account_id, amount, accrual_month, transaction_type, category_id, comments, due_date, payment_date, created_at, created_by, updated_at, updated_by, deactivated_at, deactivated_by FROM transactions WHERE id = $1`
	var t domain.Transaction
	err := r.db.Pool.QueryRow(ctx, query, id).Scan(
		&t.ID, &t.PreviousSiblingTransactionID, &t.NextSiblingTransactionID, &t.TenantID, &t.FromAccountID, &t.ToAccountID, &t.Amount, &t.AccrualMonth, &t.TransactionType, &t.CategoryID, &t.Comments, &t.DueDate, &t.PaymentDate, &t.CreatedAt, &t.CreatedBy, &t.UpdatedAt, &t.UpdatedBy, &t.DeactivatedAt, &t.DeactivatedBy,
	)
	if err != nil {
		return nil, fmt.Errorf("failed to get transaction by id: %w", err)
	}
	return &t, nil
}

func (r *TransactionRepository) List(ctx context.Context, tenantID string) ([]domain.Transaction, error) {
	query := `SELECT id, previous_sibling_transaction_id, next_sibling_transaction_id, tenant_id, from_account_id, to_account_id, amount, accrual_month, transaction_type, category_id, comments, due_date, payment_date, created_at, created_by, updated_at, updated_by, deactivated_at, deactivated_by FROM transactions WHERE tenant_id = $1`
	rows, err := r.db.Pool.Query(ctx, query, tenantID)
	if err != nil {
		return nil, fmt.Errorf("failed to list transactions: %w", err)
	}
	defer rows.Close()

	var transactions []domain.Transaction
	for rows.Next() {
		var t domain.Transaction
		if err := rows.Scan(&t.ID, &t.PreviousSiblingTransactionID, &t.NextSiblingTransactionID, &t.TenantID, &t.FromAccountID, &t.ToAccountID, &t.Amount, &t.AccrualMonth, &t.TransactionType, &t.CategoryID, &t.Comments, &t.DueDate, &t.PaymentDate, &t.CreatedAt, &t.CreatedBy, &t.UpdatedAt, &t.UpdatedBy, &t.DeactivatedAt, &t.DeactivatedBy); err != nil {
			return nil, fmt.Errorf("failed to scan transaction: %w", err)
		}
		transactions = append(transactions, t)
	}
	return transactions, nil
}

func (r *TransactionRepository) Create(ctx context.Context, t *domain.Transaction) error {
	query := `INSERT INTO transactions (id, previous_sibling_transaction_id, next_sibling_transaction_id, tenant_id, from_account_id, to_account_id, amount, accrual_month, transaction_type, category_id, comments, due_date, payment_date, created_at, created_by, updated_at, updated_by, deactivated_at, deactivated_by)
			  VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13, $14, $15, $16, $17, $18, $19)`
	_, err := r.db.Pool.Exec(ctx, query, t.ID, t.PreviousSiblingTransactionID, t.NextSiblingTransactionID, t.TenantID, t.FromAccountID, t.ToAccountID, t.Amount, t.AccrualMonth, t.TransactionType, t.CategoryID, t.Comments, t.DueDate, t.PaymentDate, t.CreatedAt, t.CreatedBy, t.UpdatedAt, t.UpdatedBy, t.DeactivatedAt, t.DeactivatedBy)
	if err != nil {
		return fmt.Errorf("failed to create transaction: %w", err)
	}
	return nil
}

func (r *TransactionRepository) Update(ctx context.Context, t *domain.Transaction) error {
	query := `UPDATE transactions SET previous_sibling_transaction_id = $2, next_sibling_transaction_id = $3, from_account_id = $4, to_account_id = $5, amount = $6, accrual_month = $7, transaction_type = $8, category_id = $9, comments = $10, due_date = $11, payment_date = $12, updated_at = $13, updated_by = $14, deactivated_at = $15, deactivated_by = $16 WHERE id = $1`
	_, err := r.db.Pool.Exec(ctx, query, t.ID, t.PreviousSiblingTransactionID, t.NextSiblingTransactionID, t.FromAccountID, t.ToAccountID, t.Amount, t.AccrualMonth, t.TransactionType, t.CategoryID, t.Comments, t.DueDate, t.PaymentDate, t.UpdatedAt, t.UpdatedBy, t.DeactivatedAt, t.DeactivatedBy)
	if err != nil {
		return fmt.Errorf("failed to update transaction: %w", err)
	}
	return nil
}

func (r *TransactionRepository) Delete(ctx context.Context, id string, userID string) error {
	query := `UPDATE transactions SET deactivated_at = CURRENT_TIMESTAMP, deactivated_by = $2 WHERE id = $1`
	_, err := r.db.Pool.Exec(ctx, query, id, userID)
	if err != nil {
		return fmt.Errorf("failed to delete transaction: %w", err)
	}
	return nil
}

func (r *TransactionRepository) AddTagToTransaction(ctx context.Context, transactionID, tagID string) error {
	query := `INSERT INTO transactions_tags (transaction_id, tag_id) VALUES ($1, $2)`
	_, err := r.db.Pool.Exec(ctx, query, transactionID, tagID)
	if err != nil {
		return fmt.Errorf("failed to add tag to transaction: %w", err)
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
			  WHERE tt.transaction_id = $1`
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
	query := `INSERT INTO transaction_attachments (id, transaction_id, name, path, created_at, created_by, updated_at, updated_by, deactivated_at, deactivated_by)
			  VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10)`
	_, err := r.db.Pool.Exec(ctx, query, a.ID, a.TransactionID, a.Name, a.Path, a.CreatedAt, a.CreatedBy, a.UpdatedAt, a.UpdatedBy, a.DeactivatedAt, a.DeactivatedBy)
	if err != nil {
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
	query := `SELECT id, transaction_id, name, path, created_at, created_by, updated_at, updated_by, deactivated_at, deactivated_by FROM transaction_attachments WHERE transaction_id = $1`
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
