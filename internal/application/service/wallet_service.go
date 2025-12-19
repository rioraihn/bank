package service

import (
	"context"

	"bank/internal/application/dto"
	"bank/internal/domain/repository"
	domainservice "bank/internal/domain/service"
	"bank/internal/domain/valueobject"
)

type walletService struct {
	walletRepo repository.WalletRepository
}

func NewWalletService(walletRepo repository.WalletRepository) domainservice.WalletService {
	return &walletService{
		walletRepo: walletRepo,
	}
}

func (s *walletService) GetBalanceByUserID(ctx context.Context, userID valueobject.UserID) (*dto.BalanceResponse, error) {
	wallet, err := s.walletRepo.FindByUserID(ctx, userID)
	if err != nil {
		return nil, err
	}

	return &dto.BalanceResponse{
		UserID:  userID.String(),
		Balance: wallet.Balance().Amount(),
	}, nil
}
