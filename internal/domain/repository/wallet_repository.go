package repository

import (
	"context"
	"database/sql"

	"bank/internal/domain/entity"
	"bank/internal/domain/valueobject"
)

type WalletRepository interface {
	GetWallet(ctx context.Context, userID valueobject.UserID) (*entity.Wallet, error)
	GetWalletForUpdate(ctx context.Context, tx *sql.Tx, userID valueobject.UserID) (*entity.Wallet, error)
	UpdateWalletBalance(ctx context.Context, tx *sql.Tx, walletID valueobject.UserID, newBalance int64) error
}

type TransactionRepository interface {
	InsertTransaction(ctx context.Context, tx *sql.Tx, transaction *entity.Transaction) error
}
