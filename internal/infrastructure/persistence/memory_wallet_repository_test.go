package persistence

import (
	"context"
	"testing"

	"bank/internal/domain/entity"
	"bank/internal/domain/valueobject"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestMemoryWalletRepository_CreateAndFind(t *testing.T) {
	// Arrange
	repo := NewMemoryWalletRepository()
	ctx := context.Background()

	userID, _ := valueobject.NewUserID("123e4567-e89b-12d3-a456-426614174000")
	initialBalance, _ := valueobject.NewMoney(5000) // $50.00
	wallet := entity.NewWalletWithBalance(userID, initialBalance)

	// Act & Assert - Create
	err := repo.Create(ctx, wallet)
	assert.NoError(t, err)

	// Act & Assert - Find
	found, err := repo.FindByUserID(ctx, userID)
	assert.NoError(t, err)
	assert.Equal(t, userID, found.UserID())
	assert.Equal(t, initialBalance, found.Balance())
	assert.Equal(t, wallet.ID(), found.ID())

	// Ensure they are different objects (copied)
	assert.NotSame(t, wallet, found)
}

func TestMemoryWalletRepository_CreateDuplicate(t *testing.T) {
	// Arrange
	repo := NewMemoryWalletRepository()
	ctx := context.Background()

	userID, _ := valueobject.NewUserID("123e4567-e89b-12d3-a456-426614174000")
	wallet1 := entity.NewWalletWithBalance(userID, valueobject.Money{})
	wallet2 := entity.NewWalletWithBalance(userID, valueobject.Money{})

	// Act & Assert
	err := repo.Create(ctx, wallet1)
	assert.NoError(t, err)

	err = repo.Create(ctx, wallet2)
	assert.Error(t, err)
	assert.Equal(t, ErrDuplicateWallet, err)
}

func TestMemoryWalletRepository_FindNotFound(t *testing.T) {
	// Arrange
	repo := NewMemoryWalletRepository()
	ctx := context.Background()

	userID, _ := valueobject.NewUserID("123e4567-e89b-12d3-a456-426614174000")

	// Act & Assert
	found, err := repo.FindByUserID(ctx, userID)
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

	err := repo.Create(ctx, wallet)
	require.NoError(t, err)

	// Modify the wallet (simulate withdrawal)
	withdrawAmount, _ := valueobject.NewMoney(2000)
	err = wallet.Withdraw(withdrawAmount)
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

func TestMemoryWalletRepository_SaveUpsert(t *testing.T) {
	// Arrange
	repo := NewMemoryWalletRepository()
	ctx := context.Background()

	userID, _ := valueobject.NewUserID("123e4567-e89b-12d3-a456-426614174000")
	wallet1 := entity.NewWalletWithBalance(userID, valueobject.Money{})

	// Act & Assert - Save (should create)
	err := repo.Save(ctx, wallet1)
	assert.NoError(t, err)

	found, err := repo.FindByUserID(ctx, userID)
	assert.NoError(t, err)
	assert.Equal(t, wallet1.ID(), found.ID())

	// Act & Assert - Save (should update)
	wallet2 := entity.NewWalletWithBalance(userID, valueobject.Money{})
	err = repo.Save(ctx, wallet2)
	assert.NoError(t, err)

	found, err = repo.FindByUserID(ctx, userID)
	assert.NoError(t, err)
	assert.Equal(t, wallet2.ID(), found.ID()) // Should be the new wallet
}

func TestMemoryWalletRepository_Exists(t *testing.T) {
	// Arrange
	repo := NewMemoryWalletRepository()
	ctx := context.Background()

	userID1, _ := valueobject.NewUserID("123e4567-e89b-12d3-a456-426614174000")
	userID2, _ := valueobject.NewUserID("123e4567-e89b-12d3-a456-426614174001")
	wallet := entity.NewWalletWithBalance(userID1, valueobject.Money{})

	err := repo.Create(ctx, wallet)
	require.NoError(t, err)

	// Act & Assert
	exists, err := repo.Exists(ctx, userID1)
	assert.NoError(t, err)
	assert.True(t, exists)

	exists, err = repo.Exists(ctx, userID2)
	assert.NoError(t, err)
	assert.False(t, exists)
}

func TestMemoryWalletRepository_GetAllAndClear(t *testing.T) {
	// Arrange
	repo := NewMemoryWalletRepository()
	ctx := context.Background()

	userID1, _ := valueobject.NewUserID("123e4567-e89b-12d3-a456-426614174000")
	userID2, _ := valueobject.NewUserID("123e4567-e89b-12d3-a456-426614174001")

	wallet1 := entity.NewWalletWithBalance(userID1, valueobject.Money{})
	wallet2 := entity.NewWalletWithBalance(userID2, valueobject.Money{})

	err := repo.Create(ctx, wallet1)
	require.NoError(t, err)
	err = repo.Create(ctx, wallet2)
	require.NoError(t, err)

	// Act & Assert - GetAll
	wallets := repo.GetAll(ctx)
	assert.Len(t, wallets, 2)

	// Act & Assert - Clear
	repo.Clear()
	wallets = repo.GetAll(ctx)
	assert.Len(t, wallets, 0)
}