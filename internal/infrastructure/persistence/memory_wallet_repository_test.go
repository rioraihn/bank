package persistence

import (
	"context"
	"testing"

	"bank/internal/domain/entity"
	"bank/internal/domain/valueobject"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestMemoryWalletRepository_FindByUserID(t *testing.T) {
	// Arrange
	repo := NewMemoryWalletRepository()
	ctx := context.Background()

	userID, _ := valueobject.NewUserID("123e4567-e89b-12d3-a456-426614174000")
	initialBalance, _ := valueobject.NewMoney(5000) // $50.00
	wallet := entity.NewWalletWithBalance(userID, initialBalance)

	// Manually add wallet to internal storage (simulating it was created elsewhere)
	repo.wallets[userID.String()] = repo.copyWallet(wallet)

	// Act
	found, err := repo.FindByUserID(ctx, userID)

	// Assert
	assert.NoError(t, err)
	assert.Equal(t, userID, found.UserID())
	assert.Equal(t, initialBalance, found.Balance())
	assert.Equal(t, wallet.ID(), found.ID())

	// Ensure they are different objects (copied)
	assert.NotSame(t, wallet, found)
}

func TestMemoryWalletRepository_FindByUserIDNotFound(t *testing.T) {
	// Arrange
	repo := NewMemoryWalletRepository()
	ctx := context.Background()

	userID, _ := valueobject.NewUserID("123e4567-e89b-12d3-a456-426614174000")

	// Act
	found, err := repo.FindByUserID(ctx, userID)

	// Assert
	assert.Error(t, err)
	assert.Equal(t, ErrWalletNotFound, err)
	assert.Nil(t, found)
}

func TestMemoryWalletRepository_Update(t *testing.T) {
	// Arrange
	repo := NewMemoryWalletRepository()
	ctx := context.Background()

	userID, _ := valueobject.NewUserID("123e4567-e89b-12d3-a456-426614174000")
	initialBalance, _ := valueobject.NewMoney(5000)
	wallet := entity.NewWalletWithBalance(userID, initialBalance)

	// Manually add wallet to storage (simulating it was created elsewhere)
	repo.wallets[userID.String()] = repo.copyWallet(wallet)

	// Modify the wallet (simulate withdrawal)
	withdrawAmount, _ := valueobject.NewMoney(2000)
	err := wallet.Withdraw(withdrawAmount)
	require.NoError(t, err)

	// Act
	err = repo.Update(ctx, wallet)
	assert.NoError(t, err)

	// Assert
	found, err := repo.FindByUserID(ctx, userID)
	assert.NoError(t, err)
	assert.Equal(t, int64(3000), found.Balance().Amount()) // 5000 - 2000 = 3000
}

func TestMemoryWalletRepository_UpdateNotFound(t *testing.T) {
	// Arrange
	repo := NewMemoryWalletRepository()
	ctx := context.Background()

	userID, _ := valueobject.NewUserID("123e4567-e89b-12d3-a456-426614174000")
	wallet := entity.NewWalletWithBalance(userID, valueobject.Money{})

	// Act & Assert
	err := repo.Update(ctx, wallet)
	assert.Error(t, err)
	assert.Equal(t, ErrWalletNotFound, err)
}

func TestMemoryWalletRepository_UpdatePreventsExternalMutation(t *testing.T) {
	// Arrange
	repo := NewMemoryWalletRepository()
	ctx := context.Background()

	userID, _ := valueobject.NewUserID("123e4567-e89b-12d3-a456-426614174000")
	initialBalance, _ := valueobject.NewMoney(1000)
	wallet := entity.NewWalletWithBalance(userID, initialBalance)

	// Manually add wallet to storage
	repo.wallets[userID.String()] = repo.copyWallet(wallet)

	// Get wallet from repository
	found, err := repo.FindByUserID(ctx, userID)
	require.NoError(t, err)

	// Modify the returned wallet
	withdrawAmount, _ := valueobject.NewMoney(500)
	err = found.Withdraw(withdrawAmount)
	require.NoError(t, err)

	// Update the modified wallet
	err = repo.Update(ctx, found)
	require.NoError(t, err)

	// Get wallet again to verify the update was applied
	foundAgain, err := repo.FindByUserID(ctx, userID)
	require.NoError(t, err)

	// Assert
	assert.Equal(t, int64(500), foundAgain.Balance().Amount()) // 1000 - 500 = 500
}
