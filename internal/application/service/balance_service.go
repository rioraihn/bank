package service

import (
	domainService "bank/internal/domain/service"
	"context"
	"log"

	"bank/internal/application/dto"
	"bank/internal/domain/repository"
	"bank/internal/domain/valueobject"
)

type balanceService struct {
	walletRepo repository.WalletRepository
}

// NewBalanceUseCase creates a new balance use case implementation
func NewBalanceUseCase(walletRepo repository.WalletRepository) domainService.BalanceService {
	return &balanceService{
		walletRepo: walletRepo,
	}
}

func (uc *balanceService) GetBalance(ctx context.Context, userID valueobject.UserID) (*dto.BalanceResponse, error) {

	wallet, err := uc.walletRepo.GetWallet(ctx, userID)
	if err != nil {
		log.Printf("❌ Wallet not found for user %s: %v", userID.String(), err)
		return &dto.BalanceResponse{
			UserID:  userID.String(),
			Balance: 0,
		}, err
	}

	return &dto.BalanceResponse{
		UserID:  userID.String(),
		Balance: wallet.Balance().Amount(),
	}, nil
}
