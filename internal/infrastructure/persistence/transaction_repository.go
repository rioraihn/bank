package persistence

import (
	"bank/internal/domain/entity"
	"context"
	"database/sql"
)

type TransactionRepository struct {
	db *sql.DB
}

func NewTransactionRepository(db *sql.DB) *TransactionRepository {
	return &TransactionRepository{
		db: db,
	}
}

// InsertTransaction inserts transaction record inside a transaction
func (r *TransactionRepository) InsertTransaction(ctx context.Context, tx *sql.Tx, transaction *entity.Transaction) error {
	query := `
		INSERT INTO transactions (id, wallet_id, amount, transaction_type, created_at)
		VALUES ($1, $2, $3, $4, NOW());
	`

	_, err := tx.ExecContext(ctx, query,
		transaction.ID().String(),
		transaction.WalletID().String(),
		transaction.Amount().Amount(),
		string(transaction.Type()),
	)

	return err
}
