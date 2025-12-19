package usecase

import (
	"context"

	"bank/internal/application/dto"
	"bank/internal/domain/entity"
	"bank/internal/domain/repository"
	domainusecase "bank/internal/domain/usecase"
	"bank/internal/domain/valueobject"
)

type withdrawUseCase struct {
	walletRepo      repository.WalletRepository
	transactionRepo repository.TransactionRepository
}

// NewWithdrawUseCase creates a new withdraw use case implementation
func NewWithdrawUseCase(walletRepo repository.WalletRepository, transactionRepo repository.TransactionRepository) domainusecase.WithdrawUseCase {
	return &withdrawUseCase{
		walletRepo:      walletRepo,
		transactionRepo: transactionRepo,
	}
}

func (uc *withdrawUseCase) Withdraw(ctx context.Context, userID valueobject.UserID, amount valueobject.Money) (*dto.WithdrawResponse, error) {

	wallet, err := uc.walletRepo.FindByUserID(ctx, userID)
	if err != nil {
		return &dto.WithdrawResponse{
			Success: false,
			Message: "wallet not found",
		}, err
	}

	err = wallet.Withdraw(amount)
	if err != nil {
		return &dto.WithdrawResponse{
			UserID:  userID.String(),
			Success: false,
			Message: "insufficient funds",
		}, err
	}

	transaction := entity.NewTransaction(
		wallet.ID(),
		entity.TransactionTypeWithdrawal,
		amount,
	)

	if err := uc.transactionRepo.Save(ctx, transaction); err != nil {
		return &dto.WithdrawResponse{
			UserID:  userID.String(),
			Success: false,
			Message: "failed to record transaction",
		}, err
	}

	if err := uc.walletRepo.Update(ctx, wallet); err != nil {
		return &dto.WithdrawResponse{
			UserID:  userID.String(),
			Success: false,
			Message: "failed to update wallet",
		}, err
	}

	return &dto.WithdrawResponse{
		UserID:          userID.String(),
		AmountWithdrawn: amount.Amount(),
		NewBalance:      wallet.Balance().Amount(),
		Success:         true,
		Message:         "withdrawal successful",
	}, nil
}
