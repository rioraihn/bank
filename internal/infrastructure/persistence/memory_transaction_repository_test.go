package persistence

import (
	"context"
	"testing"

	"bank/internal/domain/entity"
	"bank/internal/domain/valueobject"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestMemoryTransactionRepository_SaveAndFind(t *testing.T) {
	// Arrange
	repo := NewMemoryTransactionRepository()
	ctx := context.Background()

	walletID, _ := valueobject.NewUserID("123e4567-e89b-12d3-a456-426614174000")
	amount, _ := valueobject.NewMoney(5000)
	transaction := entity.NewTransaction(walletID, entity.TransactionTypeWithdrawal, amount)

	// Act
	err := repo.Save(ctx, transaction)
	assert.NoError(t, err)

	// Assert
	found, err := repo.FindByID(ctx, transaction.ID())
	assert.NoError(t, err)
	assert.Equal(t, transaction.ID(), found.ID())
	assert.Equal(t, walletID, found.WalletID())
	assert.Equal(t, entity.TransactionTypeWithdrawal, found.Type())
	assert.Equal(t, amount, found.Amount())

	// Ensure they are different objects (copied)
	assert.NotSame(t, transaction, found)
}

func TestMemoryTransactionRepository_FindByIDNotFound(t *testing.T) {
	// Arrange
	repo := NewMemoryTransactionRepository()
	ctx := context.Background()

	transactionID, _ := valueobject.NewUserID("123e4567-e89b-12d3-a456-426614174000")

	// Act & Assert
	found, err := repo.FindByID(ctx, transactionID)
	assert.Error(t, err)
	assert.Equal(t, ErrTransactionNotFound, err)
	assert.Nil(t, found)
}

func TestMemoryTransactionRepository_FindByWalletID(t *testing.T) {
	// Arrange
	repo := NewMemoryTransactionRepository()
	ctx := context.Background()

	walletID1, _ := valueobject.NewUserID("123e4567-e89b-12d3-a456-426614174000")
	walletID2, _ := valueobject.NewUserID("123e4567-e89b-12d3-a456-426614174001")

	amount, _ := valueobject.NewMoney(1000)

	// Create transactions for wallet1
	tx1 := entity.NewTransaction(walletID1, entity.TransactionTypeWithdrawal, amount)
	tx2 := entity.NewTransaction(walletID1, entity.TransactionTypeDeposit, amount)

	// Create transaction for wallet2
	tx3 := entity.NewTransaction(walletID2, entity.TransactionTypeWithdrawal, amount)

	err := repo.Save(ctx, tx1)
	require.NoError(t, err)
	err = repo.Save(ctx, tx2)
	require.NoError(t, err)
	err = repo.Save(ctx, tx3)
	require.NoError(t, err)

	// Act
	transactions, err := repo.FindByWalletID(ctx, walletID1, 10, 0)

	// Assert
	assert.NoError(t, err)
	assert.Len(t, transactions, 2)

	// Should be sorted by creation time (newest first) or equal if created at same time
	assert.True(t, transactions[0].CreatedAt().After(transactions[1].CreatedAt()) ||
		transactions[0].CreatedAt().Equal(transactions[1].CreatedAt()))
}

func TestMemoryTransactionRepository_FindByWalletIDWithPagination(t *testing.T) {
	// Arrange
	repo := NewMemoryTransactionRepository()
	ctx := context.Background()

	walletID, _ := valueobject.NewUserID("123e4567-e89b-12d3-a456-426614174000")
	amount, _ := valueobject.NewMoney(1000)

	// Create multiple transactions
	var transactions []*entity.Transaction
	for i := 0; i < 5; i++ {
		tx := entity.NewTransaction(walletID, entity.TransactionTypeWithdrawal, amount)
		err := repo.Save(ctx, tx)
		require.NoError(t, err)
		transactions = append(transactions, tx)
	}

	// Test cases
	testCases := []struct {
		name     string
		limit    int
		offset   int
		expected int
	}{
		{"First page", 2, 0, 2},
		{"Second page", 2, 2, 2},
		{"Last page", 2, 4, 1},
		{"Beyond data", 2, 5, 0},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			// Act
			result, err := repo.FindByWalletID(ctx, walletID, tc.limit, tc.offset)

			// Assert
			assert.NoError(t, err)
			assert.Len(t, result, tc.expected)
		})
	}
}

func TestMemoryTransactionRepository_FindByWalletIDAndType(t *testing.T) {
	// Arrange
	repo := NewMemoryTransactionRepository()
	ctx := context.Background()

	walletID, _ := valueobject.NewUserID("123e4567-e89b-12d3-a456-426614174000")
	amount, _ := valueobject.NewMoney(1000)

	// Create transactions of different types
	withdrawal1 := entity.NewTransaction(walletID, entity.TransactionTypeWithdrawal, amount)
	deposit1 := entity.NewTransaction(walletID, entity.TransactionTypeDeposit, amount)
	withdrawal2 := entity.NewTransaction(walletID, entity.TransactionTypeWithdrawal, amount)

	err := repo.Save(ctx, withdrawal1)
	require.NoError(t, err)
	err = repo.Save(ctx, deposit1)
	require.NoError(t, err)
	err = repo.Save(ctx, withdrawal2)
	require.NoError(t, err)

	// Act
	withdrawals, err := repo.FindByWalletIDAndType(ctx, walletID, entity.TransactionTypeWithdrawal, 10, 0)
	deposits, err := repo.FindByWalletIDAndType(ctx, walletID, entity.TransactionTypeDeposit, 10, 0)

	// Assert
	assert.NoError(t, err)
	assert.Len(t, withdrawals, 2)
	assert.Len(t, deposits, 1)
}

func TestMemoryTransactionRepository_CountByWalletID(t *testing.T) {
	// Arrange
	repo := NewMemoryTransactionRepository()
	ctx := context.Background()

	walletID1, _ := valueobject.NewUserID("123e4567-e89b-12d3-a456-426614174000")
	walletID2, _ := valueobject.NewUserID("123e4567-e89b-12d3-a456-426614174001")
	amount, _ := valueobject.NewMoney(1000)

	// Create transactions for wallet1
	tx1 := entity.NewTransaction(walletID1, entity.TransactionTypeWithdrawal, amount)
	tx2 := entity.NewTransaction(walletID1, entity.TransactionTypeDeposit, amount)

	// Create transaction for wallet2
	tx3 := entity.NewTransaction(walletID2, entity.TransactionTypeWithdrawal, amount)

	err := repo.Save(ctx, tx1)
	require.NoError(t, err)
	err = repo.Save(ctx, tx2)
	require.NoError(t, err)
	err = repo.Save(ctx, tx3)
	require.NoError(t, err)

	// Act & Assert
	count, err := repo.CountByWalletID(ctx, walletID1)
	assert.NoError(t, err)
	assert.Equal(t, 2, count)

	count, err = repo.CountByWalletID(ctx, walletID2)
	assert.NoError(t, err)
	assert.Equal(t, 1, count)

	count, err = repo.CountByWalletID(ctx, valueobject.NewUserIDRandom())
	assert.NoError(t, err)
	assert.Equal(t, 0, count)
}

func TestMemoryTransactionRepository_GetAllAndClear(t *testing.T) {
	// Arrange
	repo := NewMemoryTransactionRepository()
	ctx := context.Background()

	walletID, _ := valueobject.NewUserID("123e4567-e89b-12d3-a456-426614174000")
	amount, _ := valueobject.NewMoney(1000)

	tx1 := entity.NewTransaction(walletID, entity.TransactionTypeWithdrawal, amount)
	tx2 := entity.NewTransaction(walletID, entity.TransactionTypeDeposit, amount)

	err := repo.Save(ctx, tx1)
	require.NoError(t, err)
	err = repo.Save(ctx, tx2)
	require.NoError(t, err)

	// Act & Assert - GetAll
	transactions := repo.GetAll(ctx)
	assert.Len(t, transactions, 2)

	// Act & Assert - Clear
	repo.Clear()
	transactions = repo.GetAll(ctx)
	assert.Len(t, transactions, 0)
}