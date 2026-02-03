package postgres

import (
	"context"
	"fmt"

	"github.com/igoventura/fintrack-core/domain"
)

type AccountRepository struct {
	db *DB
}

func NewAccountRepository(db *DB) *AccountRepository {
	return &AccountRepository{db: db}
}

func (r *AccountRepository) GetByID(ctx context.Context, id string) (*domain.Account, error) {
	query := `SELECT id, tenant_id, name, initial_balance, color, currency, icon, type, created_at, created_by, updated_at, updated_by, deactivated_at, deactivated_by FROM accounts WHERE id = $1 AND deactivated_at IS NULL`
	var a domain.Account
	err := r.db.Pool.QueryRow(ctx, query, id).Scan(
		&a.ID, &a.TenantID, &a.Name, &a.InitialBalance, &a.Color, &a.Currency, &a.Icon, &a.Type, &a.CreatedAt, &a.CreatedBy, &a.UpdatedAt, &a.UpdatedBy, &a.DeactivatedAt, &a.DeactivatedBy,
	)
	if err != nil {
		return nil, fmt.Errorf("failed to get account by id: %w", err)
	}
	return &a, nil
}

func (r *AccountRepository) List(ctx context.Context, tenantID string) ([]domain.Account, error) {
	query := `SELECT id, tenant_id, name, initial_balance, color, currency, icon, type, created_at, created_by, updated_at, updated_by, deactivated_at, deactivated_by FROM accounts WHERE tenant_id = $1 AND deactivated_at IS NULL`
	rows, err := r.db.Pool.Query(ctx, query, tenantID)
	if err != nil {
		return nil, fmt.Errorf("failed to list accounts: %w", err)
	}
	defer rows.Close()

	var accounts []domain.Account
	for rows.Next() {
		var a domain.Account
		if err := rows.Scan(&a.ID, &a.TenantID, &a.Name, &a.InitialBalance, &a.Color, &a.Currency, &a.Icon, &a.Type, &a.CreatedAt, &a.CreatedBy, &a.UpdatedAt, &a.UpdatedBy, &a.DeactivatedAt, &a.DeactivatedBy); err != nil {
			return nil, fmt.Errorf("failed to scan account: %w", err)
		}
		accounts = append(accounts, a)
	}
	return accounts, nil
}

func (r *AccountRepository) Create(ctx context.Context, a *domain.Account) error {
	query := `INSERT INTO accounts (tenant_id, name, initial_balance, color, currency, icon, type, created_by, updated_by)
			  VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9)
			  RETURNING id, created_at, updated_at`
	row := r.db.Pool.QueryRow(ctx, query, a.TenantID, a.Name, a.InitialBalance, a.Color, a.Currency, a.Icon, a.Type, a.CreatedBy, a.UpdatedBy)
	if err := row.Scan(&a.ID, &a.CreatedAt, &a.UpdatedAt); err != nil {
		return fmt.Errorf("failed to create account: %w", err)
	}
	return nil
}

func (r *AccountRepository) Update(ctx context.Context, a *domain.Account) error {
	query := `UPDATE accounts SET name = $2, initial_balance = $3, color = $4, currency = $5, icon = $6, type = $7, updated_at = CURRENT_TIMESTAMP, updated_by = $8 WHERE id = $1 AND tenant_id = $9 RETURNING updated_at`
	row := r.db.Pool.QueryRow(ctx, query, a.ID, a.Name, a.InitialBalance, a.Color, a.Currency, a.Icon, a.Type, a.UpdatedBy, a.TenantID)
	if err := row.Scan(&a.UpdatedAt); err != nil {
		return fmt.Errorf("failed to update account: %w", err)
	}
	return nil
}

func (r *AccountRepository) Delete(ctx context.Context, id string, userID string) error {
	query := `UPDATE accounts SET deactivated_at = CURRENT_TIMESTAMP, deactivated_by = $2 WHERE id = $1`
	_, err := r.db.Pool.Exec(ctx, query, id, userID)
	if err != nil {
		return fmt.Errorf("failed to delete account: %w", err)
	}
	return nil
}

func (r *AccountRepository) GetCreditCardInfo(ctx context.Context, accountID string) (*domain.CreditCardInfo, error) {
	query := `SELECT id, account_id, last_four, name, brand, closing_date, due_date, created_at, created_by, updated_at, updated_by, deactivated_at, deactivated_by FROM credit_card_info WHERE account_id = $1 AND deactivated_at IS NULL`
	var info domain.CreditCardInfo
	err := r.db.Pool.QueryRow(ctx, query, accountID).Scan(
		&info.ID, &info.AccountID, &info.LastFour, &info.Name, &info.Brand, &info.ClosingDate, &info.DueDate, &info.CreatedAt, &info.CreatedBy, &info.UpdatedAt, &info.UpdatedBy, &info.DeactivatedAt, &info.DeactivatedBy,
	)
	if err != nil {
		return nil, fmt.Errorf("failed to get credit card info: %w", err)
	}
	return &info, nil
}

func (r *AccountRepository) UpsertCreditCardInfo(ctx context.Context, info *domain.CreditCardInfo) error {
	query := `INSERT INTO credit_card_info (account_id, last_four, name, brand, closing_date, due_date, created_by, updated_by)
			  VALUES ($1, $2, $3, $4, $5, $6, $7, $8)
			  ON CONFLICT (account_id, deactivated_at) DO UPDATE SET
				last_four = EXCLUDED.last_four,
				name = EXCLUDED.name,
				brand = EXCLUDED.brand,
				closing_date = EXCLUDED.closing_date,
				due_date = EXCLUDED.due_date,
				updated_at = CURRENT_TIMESTAMP,
				updated_by = EXCLUDED.updated_by
			  RETURNING id, created_at, updated_at`
	row := r.db.Pool.QueryRow(ctx, query, info.AccountID, info.LastFour, info.Name, info.Brand, info.ClosingDate, info.DueDate, info.CreatedBy, info.UpdatedBy)
	if err := row.Scan(&info.ID, &info.CreatedAt, &info.UpdatedAt); err != nil {
		return fmt.Errorf("failed to upsert credit card info: %w", err)
	}
	return nil
}
