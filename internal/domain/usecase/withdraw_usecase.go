package usecase

import (
	"context"

	"bank/internal/application/dto"
	"bank/internal/domain/valueobject"
)

type WithdrawUseCase interface {
	Withdraw(ctx context.Context, userID valueobject.UserID, amount valueobject.Money) (*dto.WithdrawResponse, error)
}
