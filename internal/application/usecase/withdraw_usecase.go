package usecase

import (
	"context"
	"database/sql"
	"log"

	"bank/internal/application/dto"
	"bank/internal/domain/entity"
	"bank/internal/domain/repository"
	domainusecase "bank/internal/domain/usecase"
	"bank/internal/domain/valueobject"
)

type withdrawUseCase struct {
	walletRepo      repository.WalletRepository
	transactionRepo repository.TransactionRepository
	db              *sql.DB
}

// NewWithdrawUseCase creates a new withdraw use case implementation
func NewWithdrawUseCase(walletRepo repository.WalletRepository, transactionRepo repository.TransactionRepository, db *sql.DB) domainusecase.WithdrawUseCase {
	return &withdrawUseCase{
		walletRepo:      walletRepo,
		transactionRepo: transactionRepo,
		db:              db,
	}
}

func (uc *withdrawUseCase) Withdraw(ctx context.Context, userID valueobject.UserID, amount valueobject.Money) (*dto.WithdrawResponse, error) {

	// Begin transaction
	tx, err := uc.db.BeginTx(ctx, nil)
	if err != nil {
		log.Printf("❌ Failed to begin transaction for user %s: %v", userID.String(), err)
		return &dto.WithdrawResponse{
			UserID:  userID.String(),
			Success: false,
			Message: "failed to begin transaction",
		}, err
	}

	defer func() {
		if err != nil {
			if rbErr := tx.Rollback(); rbErr != nil {
				log.Printf("❌ Failed to rollback transaction for user %s: %v", userID.String(), rbErr)
			}
		}
	}()

	wallet, err := uc.walletRepo.GetWalletForUpdate(ctx, tx, userID)
	if err != nil {
		log.Printf("❌ Wallet not found for user %s: %v", userID.String(), err)
		return &dto.WithdrawResponse{
			UserID:  userID.String(),
			Success: false,
			Message: "wallet not found",
		}, err
	}

	if wallet.Balance().Amount() < amount.Amount() {
		log.Printf("💸 Insufficient funds for user %s: attempted %d, available %d",
			userID.String(), amount.Amount(), wallet.Balance().Amount())
		return &dto.WithdrawResponse{
			UserID:  userID.String(),
			Success: false,
			Message: "insufficient funds",
		}, nil
	}

	newBalance := wallet.Balance().Amount() - amount.Amount()

	if err := uc.walletRepo.UpdateWalletBalance(ctx, tx, wallet.ID(), newBalance); err != nil {
		log.Printf("❌ Failed to update wallet balance for user %s: %v", userID.String(), err)
		return &dto.WithdrawResponse{
			UserID:  userID.String(),
			Success: false,
			Message: "failed to update wallet balance",
		}, err
	}

	transaction := entity.NewTransaction(
		wallet.ID(),
		entity.TransactionTypeWithdrawal,
		amount,
	)

	if err := uc.transactionRepo.InsertTransaction(ctx, tx, transaction); err != nil {
		log.Printf("❌ Failed to save transaction %s: %v", transaction.ID().String(), err)
		return &dto.WithdrawResponse{
			UserID:  userID.String(),
			Success: false,
			Message: "failed to record transaction",
		}, err
	}

	if err := tx.Commit(); err != nil {
		log.Printf("❌ Failed to commit transaction for user %s: %v", userID.String(), err)
		return &dto.WithdrawResponse{
			UserID:  userID.String(),
			Success: false,
			Message: "failed to commit transaction",
		}, err
	}

	return &dto.WithdrawResponse{
		UserID:          userID.String(),
		AmountWithdrawn: amount.Amount(),
		NewBalance:      newBalance,
		Success:         true,
		Message:         "withdrawal successful",
	}, nil
}
