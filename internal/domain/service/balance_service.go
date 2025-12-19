package service

import (
	"context"

	"bank/internal/application/dto"
	"bank/internal/domain/valueobject"
)

type BalanceService interface {
	GetBalance(ctx context.Context, userID valueobject.UserID) (*dto.BalanceResponse, error)
}
