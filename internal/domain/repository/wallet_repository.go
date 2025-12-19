package repository

import (
	"context"

	"bank/internal/domain/entity"
	"bank/internal/domain/valueobject"
)

type WalletRepository interface {
	FindByUserID(ctx context.Context, userID valueobject.UserID) (*entity.Wallet, error)
	Update(ctx context.Context, wallet *entity.Wallet) error
}

type TransactionRepository interface {
	Save(ctx context.Context, transaction *entity.Transaction) error
}
