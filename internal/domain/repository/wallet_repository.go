package repository

import (
	"context"

	"bank/internal/domain/entity"
	"bank/internal/domain/valueobject"
)

type WalletRepository interface {
	FindByUserID(ctx context.Context, userID valueobject.UserID) (*entity.Wallet, error)
	Save(ctx context.Context, wallet *entity.Wallet) error
	Create(ctx context.Context, wallet *entity.Wallet) error
	Update(ctx context.Context, wallet *entity.Wallet) error
	Exists(ctx context.Context, userID valueobject.UserID) (bool, error)
}

type TransactionRepository interface {
	Save(ctx context.Context, transaction *entity.Transaction) error
	FindByID(ctx context.Context, transactionID valueobject.UserID) (*entity.Transaction, error)
	FindByWalletID(ctx context.Context, walletID valueobject.UserID, limit, offset int) ([]*entity.Transaction, error)
	FindByWalletIDAndType(
		ctx context.Context,
		walletID valueobject.UserID,
		transactionType entity.TransactionType,
		limit, offset int,
	) ([]*entity.Transaction, error)
	CountByWalletID(ctx context.Context, walletID valueobject.UserID) (int, error)
}
