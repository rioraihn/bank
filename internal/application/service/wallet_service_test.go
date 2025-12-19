package service

import (
	"context"
	"testing"

	"bank/internal/application/dto"
	"bank/internal/application/mocks"
	"bank/internal/domain/entity"
	"bank/internal/domain/valueobject"

	"github.com/stretchr/testify/assert"
)

func TestWalletService_GetBalanceByUserID_Success(t *testing.T) {
	// Arrange
	ctx := context.Background()
	userID, _ := valueobject.NewUserID("123e4567-e89b-12d3-a456-426614174000")
	balance, _ := valueobject.NewMoney(5000) // $50.00
	wallet := entity.NewWalletWithBalance(userID, balance)

	// Explicitly reference the dto type to ensure the import is used
	var _ *dto.BalanceResponse = nil

	mockWalletRepo := new(mocks.MockWalletRepository)
	service := NewWalletService(mockWalletRepo)

	// Set up mock expectations
	mockWalletRepo.On("FindByUserID", ctx, userID).Return(wallet, nil)

	// Act
	resp, err := service.GetBalanceByUserID(ctx, userID)

	// Assert
	assert.NoError(t, err)
	assert.NotNil(t, resp)
	assert.Equal(t, userID.String(), resp.UserID)
	assert.Equal(t, int64(5000), resp.Balance)

	mockWalletRepo.AssertExpectations(t)
}

func TestWalletService_GetBalanceByUserID_WalletNotFound(t *testing.T) {
	// Arrange
	ctx := context.Background()
	userID, _ := valueobject.NewUserID("123e4567-e89b-12d3-a456-426614174000")

	mockWalletRepo := new(mocks.MockWalletRepository)
	service := NewWalletService(mockWalletRepo)

	// Set up mock expectations
	mockWalletRepo.On("FindByUserID", ctx, userID).Return(nil, assert.AnError)

	// Act
	resp, err := service.GetBalanceByUserID(ctx, userID)

	// Assert
	assert.Error(t, err)
	assert.Nil(t, resp)

	mockWalletRepo.AssertExpectations(t)
}
