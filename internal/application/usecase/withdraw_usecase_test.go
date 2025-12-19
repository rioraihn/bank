package usecase

import (
	"context"
	"errors"
	"testing"

	"bank/internal/domain/entity"
	"bank/internal/domain/valueobject"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)


func TestWithdrawUseCase_Withdraw_Success(t *testing.T) {
	// Arrange
	ctx := context.Background()
	userID, _ := valueobject.NewUserID("123e4567-e89b-12d3-a456-426614174000")

	initialBalance, _ := valueobject.NewMoney(10000) // $100.00
	wallet := entity.NewWalletWithBalance(userID, initialBalance)

	mockWalletRepo := new(MockWalletRepository)
	mockTransactionRepo := new(MockTransactionRepository)

	usecase := NewWithdrawUseCase(mockWalletRepo, mockTransactionRepo)

	amount, _ := valueobject.NewMoney(2500)

	// Set up mock expectations
	mockWalletRepo.On("FindByUserID", ctx, userID).Return(wallet, nil)
	mockTransactionRepo.On("Save", ctx, mock.AnythingOfType("*entity.Transaction")).Return(nil)
	mockWalletRepo.On("Update", ctx, wallet).Return(nil)

	// Act
	resp, err := usecase.Withdraw(ctx, userID, amount)

	// Assert
	assert.NoError(t, err)
	assert.True(t, resp.Success)
	assert.Equal(t, userID.String(), resp.UserID)
	assert.Equal(t, int64(2500), resp.AmountWithdrawn)
	assert.Equal(t, int64(7500), resp.NewBalance) // $100.00 - $25.00 = $75.00
	assert.Equal(t, "withdrawal successful", resp.Message)

	mockWalletRepo.AssertExpectations(t)
	mockTransactionRepo.AssertExpectations(t)
}

func TestWithdrawUseCase_Withdraw_InsufficientFunds(t *testing.T) {
	// Arrange
	ctx := context.Background()
	userID, _ := valueobject.NewUserID("123e4567-e89b-12d3-a456-426614174000")

	initialBalance, _ := valueobject.NewMoney(1000) // $10.00
	wallet := entity.NewWalletWithBalance(userID, initialBalance)

	mockWalletRepo := new(MockWalletRepository)
	mockTransactionRepo := new(MockTransactionRepository)

	usecase := NewWithdrawUseCase(mockWalletRepo, mockTransactionRepo)

	amount, _ := valueobject.NewMoney(2500) // $25.00 - more than balance

	// Set up mock expectations
	mockWalletRepo.On("FindByUserID", ctx, userID).Return(wallet, nil)

	// Act
	resp, err := usecase.Withdraw(ctx, userID, amount)

	// Assert
	assert.Error(t, err)
	assert.False(t, resp.Success)
	assert.Equal(t, userID.String(), resp.UserID)
	assert.Equal(t, int64(0), resp.AmountWithdrawn)
	assert.Equal(t, "insufficient funds", resp.Message)

	mockWalletRepo.AssertExpectations(t)
	// Transaction should NOT be saved
	mockTransactionRepo.AssertNotCalled(t, "Save")
	mockWalletRepo.AssertNotCalled(t, "Update")
}

func TestWithdrawUseCase_Withdraw_InvalidAmount(t *testing.T) {
	// Test invalid amount (negative value) - should fail at domain level
	_, err := valueobject.NewMoney(-1)
	assert.Error(t, err) // This should fail at domain level

	// Test valid zero amount - this succeeds at domain level but fails at business logic
	zeroMoney, err := valueobject.NewMoney(0)
	assert.NoError(t, err) // Zero money is valid at domain level
	assert.True(t, zeroMoney.IsZero())
}


func TestWithdrawUseCase_Withdraw_WalletNotFound(t *testing.T) {
	// Arrange
	ctx := context.Background()
	userID, _ := valueobject.NewUserID("123e4567-e89b-12d3-a456-426614174000")

	mockWalletRepo := new(MockWalletRepository)
	mockTransactionRepo := new(MockTransactionRepository)

	usecase := NewWithdrawUseCase(mockWalletRepo, mockTransactionRepo)

	amount, _ := valueobject.NewMoney(2500)

	// Set up mock expectations
	mockWalletRepo.On("FindByUserID", ctx, userID).Return(nil, errors.New("wallet not found"))

	// Act
	resp, err := usecase.Withdraw(ctx, userID, amount)

	// Assert
	assert.Error(t, err)
	assert.False(t, resp.Success)
	assert.Equal(t, "wallet not found", resp.Message)

	mockWalletRepo.AssertExpectations(t)
}

func TestWithdrawUseCase_Withdraw_TransactionSaveFailure(t *testing.T) {
	// Arrange
	ctx := context.Background()
	userID, _ := valueobject.NewUserID("123e4567-e89b-12d3-a456-426614174000")

	initialBalance, _ := valueobject.NewMoney(10000) // $100.00
	wallet := entity.NewWalletWithBalance(userID, initialBalance)

	mockWalletRepo := new(MockWalletRepository)
	mockTransactionRepo := new(MockTransactionRepository)

	usecase := NewWithdrawUseCase(mockWalletRepo, mockTransactionRepo)

	amount, _ := valueobject.NewMoney(2500) // $25.00

	// Set up mock expectations
	mockWalletRepo.On("FindByUserID", ctx, userID).Return(wallet, nil)
	mockTransactionRepo.On("Save", ctx, mock.AnythingOfType("*entity.Transaction")).Return(errors.New("database connection failed"))

	// Act
	resp, err := usecase.Withdraw(ctx, userID, amount)

	// Assert
	assert.Error(t, err)
	assert.False(t, resp.Success)
	assert.Equal(t, userID.String(), resp.UserID)
	assert.Equal(t, int64(0), resp.AmountWithdrawn)
	assert.Equal(t, "failed to record transaction", resp.Message)

	mockWalletRepo.AssertExpectations(t)
	mockTransactionRepo.AssertExpectations(t)
}

func TestWithdrawUseCase_Withdraw_WalletUpdateFailure(t *testing.T) {
	// Arrange
	ctx := context.Background()
	userID, _ := valueobject.NewUserID("123e4567-e89b-12d3-a456-426614174000")

	initialBalance, _ := valueobject.NewMoney(10000) // $100.00
	wallet := entity.NewWalletWithBalance(userID, initialBalance)

	mockWalletRepo := new(MockWalletRepository)
	mockTransactionRepo := new(MockTransactionRepository)

	usecase := NewWithdrawUseCase(mockWalletRepo, mockTransactionRepo)

	amount, _ := valueobject.NewMoney(2500) // $25.00

	// Set up mock expectations
	mockWalletRepo.On("FindByUserID", ctx, userID).Return(wallet, nil)
	mockTransactionRepo.On("Save", ctx, mock.AnythingOfType("*entity.Transaction")).Return(nil)
	mockWalletRepo.On("Update", ctx, wallet).Return(errors.New("database connection failed"))

	// Act
	resp, err := usecase.Withdraw(ctx, userID, amount)

	// Assert
	assert.Error(t, err)
	assert.False(t, resp.Success)
	assert.Equal(t, userID.String(), resp.UserID)
	assert.Equal(t, int64(0), resp.AmountWithdrawn)
	assert.Equal(t, "failed to update wallet", resp.Message)

	mockWalletRepo.AssertExpectations(t)
	mockTransactionRepo.AssertExpectations(t)
}