package postgres

import (
	"context"
	"database/sql"
	"fmt"

	"github.com/gabrielcau/ledger-system/internal/domain"
)

type AccountRepository struct {
	db *sql.DB
}

func NewAccountRepository(db *sql.DB) *AccountRepository {
	return &AccountRepository{db: db}
}

// TransferTx performs a money transfer within a database transaction using Pessimistic Locking.
func (r *AccountRepository) TransferTx(ctx context.Context, transfer domain.Transfer) error {
	tx, err := r.db.BeginTx(ctx, nil)
	if err != nil {
		return fmt.Errorf("failed to begin transaction: %w", err)
	}
	defer tx.Rollback()

	// 1. Sort Account IDs to prevent Deadlocks
	// Always lock the lower ID first.
	// This ensures that two concurrent transfers (A -> B and B -> A) don't deadlock waiting for each other.
	firstID, secondID := transfer.FromAccountID, transfer.ToAccountID
	if firstID > secondID {
		firstID, secondID = secondID, firstID
	}

	// 2. Lock the accounts
	// We use SELECT ... FOR UPDATE to lock the rows.
	// This prevents other transactions from modifying these rows until this transaction commits/rolls back.

	// Lock first account
	if err := r.lockAccount(ctx, tx, firstID); err != nil {
		return err
	}

	// Lock second account
	if err := r.lockAccount(ctx, tx, secondID); err != nil {
		return err
	}

	// 3. Get current balances (re-querying inside the lock to be safe, though lockAccount creates the lock)
	// Actually, lockAccount simply locks. We need to verify sufficient funds for the sender.
	var fromBalance int64
	err = tx.QueryRowContext(ctx, "SELECT balance FROM accounts WHERE id = $1", transfer.FromAccountID).Scan(&fromBalance)
	if err != nil {
		return fmt.Errorf("failed to get sender balance: %w", err)
	}

	if fromBalance < transfer.Amount {
		return domain.ErrInsufficientFunds
	}

	// 4. Update Balances
	_, err = tx.ExecContext(ctx, "UPDATE accounts SET balance = balance - $1 WHERE id = $2", transfer.Amount, transfer.FromAccountID)
	if err != nil {
		return fmt.Errorf("failed to debit sender: %w", err)
	}

	_, err = tx.ExecContext(ctx, "UPDATE accounts SET balance = balance + $1 WHERE id = $2", transfer.Amount, transfer.ToAccountID)
	if err != nil {
		return fmt.Errorf("failed to credit receiver: %w", err)
	}

	// 5. Create Transfer Record
	_, err = tx.ExecContext(ctx,
		"INSERT INTO transfers (from_account_id, to_account_id, amount) VALUES ($1, $2, $3)",
		transfer.FromAccountID, transfer.ToAccountID, transfer.Amount)
	if err != nil {
		return fmt.Errorf("failed to create transfer record: %w", err)
	}

	return tx.Commit()
}

func (r *AccountRepository) lockAccount(ctx context.Context, tx *sql.Tx, id int) error {
	// The query verifies existence and acquires the lock.
	// The 'FOR UPDATE' clause is critical here.
	var dummyID int
	err := tx.QueryRowContext(ctx, "SELECT id FROM accounts WHERE id = $1 FOR UPDATE", id).Scan(&dummyID)
	if err != nil {
		if err == sql.ErrNoRows {
			return domain.ErrAccountNotFound
		}
		return fmt.Errorf("failed to lock account %d: %w", id, err)
	}
	return nil
}

func (r *AccountRepository) GetBalance(ctx context.Context, id int) (int64, error) {
	var balance int64
	err := r.db.QueryRowContext(ctx, "SELECT balance FROM accounts WHERE id = $1", id).Scan(&balance)
	if err != nil {
		return 0, err
	}
	return balance, nil
}
