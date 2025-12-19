package service

import (
	"context"

	"bank/internal/application/dto"
	"bank/internal/domain/valueobject"
)

type WalletService interface {
	GetBalanceByUserID(ctx context.Context, userID valueobject.UserID) (*dto.BalanceResponse, error)
}
