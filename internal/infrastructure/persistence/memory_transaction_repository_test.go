package persistence

import (
	"context"
	"testing"

	"bank/internal/domain/entity"
	"bank/internal/domain/valueobject"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestMemoryTransactionRepository_Save(t *testing.T) {
	// Arrange
	repo := NewMemoryTransactionRepository()
	ctx := context.Background()

	walletID, _ := valueobject.NewUserID("123e4567-e89b-12d3-a456-426614174000")
	amount, _ := valueobject.NewMoney(5000)
	transaction := entity.NewTransaction(walletID, entity.TransactionTypeWithdrawal, amount)

	// Act
	err := repo.Save(ctx, transaction)

	// Assert
	assert.NoError(t, err)

	// Verify transaction was stored by checking internal storage
	stored, exists := repo.transactions[transaction.ID().String()]
	assert.True(t, exists)
	assert.Equal(t, transaction.ID(), stored.ID())
	assert.Equal(t, walletID, stored.WalletID())
	assert.Equal(t, entity.TransactionTypeWithdrawal, stored.Type())
	assert.Equal(t, amount, stored.Amount())

	// Ensure they are different objects (copied)
	assert.NotSame(t, transaction, stored)
}

func TestMemoryTransactionRepository_SaveMultipleTransactions(t *testing.T) {
	// Arrange
	repo := NewMemoryTransactionRepository()
	ctx := context.Background()

	walletID, _ := valueobject.NewUserID("123e4567-e89b-12d3-a456-426614174000")
	amount, _ := valueobject.NewMoney(1000)

	// Create and save multiple transactions
	tx1 := entity.NewTransaction(walletID, entity.TransactionTypeWithdrawal, amount)
	tx2 := entity.NewTransaction(walletID, entity.TransactionTypeDeposit, amount)
	tx3 := entity.NewTransaction(walletID, entity.TransactionTypeWithdrawal, amount)

	// Act
	err1 := repo.Save(ctx, tx1)
	err2 := repo.Save(ctx, tx2)
	err3 := repo.Save(ctx, tx3)

	// Assert
	assert.NoError(t, err1)
	assert.NoError(t, err2)
	assert.NoError(t, err3)

	// Verify all transactions were stored
	assert.Equal(t, 3, len(repo.transactions))
	assert.Contains(t, repo.transactions, tx1.ID().String())
	assert.Contains(t, repo.transactions, tx2.ID().String())
	assert.Contains(t, repo.transactions, tx3.ID().String())
}

func TestMemoryTransactionRepository_SaveOverwritesExisting(t *testing.T) {
	// Arrange
	repo := NewMemoryTransactionRepository()
	ctx := context.Background()

	walletID, _ := valueobject.NewUserID("123e4567-e89b-12d3-a456-426614174000")
	amount1, _ := valueobject.NewMoney(1000)
	amount2, _ := valueobject.NewMoney(2000)

	transaction := entity.NewTransaction(walletID, entity.TransactionTypeWithdrawal, amount1)

	// Save first transaction
	err := repo.Save(ctx, transaction)
	require.NoError(t, err)

	// Create a new transaction with the same ID but different amount
	updatedTransaction := entity.ReconstructTransaction(
		transaction.ID(),
		transaction.WalletID(),
		transaction.Type(),
		amount2,
		transaction.Status(),
		transaction.CreatedAt().Format("RFC3339"),
		"",
	)

	// Act
	err = repo.Save(ctx, updatedTransaction)

	// Assert
	assert.NoError(t, err)

	// Verify the transaction was overwritten
	stored := repo.transactions[transaction.ID().String()]
	assert.Equal(t, amount2, stored.Amount())
	assert.Equal(t, amount2.Amount(), int64(2000))
}

func TestMemoryTransactionRepository_ConcurrentSave(t *testing.T) {
	// Arrange
	repo := NewMemoryTransactionRepository()
	ctx := context.Background()

	walletID, _ := valueobject.NewUserID("123e4567-e89b-12d3-a456-426614174000")
	amount, _ := valueobject.NewMoney(1000)

	// Test concurrent saves
	numTransactions := 10
	transactions := make([]*entity.Transaction, numTransactions)

	for i := 0; i < numTransactions; i++ {
		transactions[i] = entity.NewTransaction(walletID, entity.TransactionTypeWithdrawal, amount)
	}

	// Act - Save concurrently
	done := make(chan bool, numTransactions)

	for i := 0; i < numTransactions; i++ {
		go func(index int) {
			err := repo.Save(ctx, transactions[index])
			assert.NoError(t, err)
			done <- true
		}(i)
	}

	// Wait for all goroutines to complete
	for i := 0; i < numTransactions; i++ {
		<-done
	}

	// Assert
	assert.Equal(t, numTransactions, len(repo.transactions))

	// Verify all transactions are stored
	for _, tx := range transactions {
		stored, exists := repo.transactions[tx.ID().String()]
		assert.True(t, exists)
		assert.Equal(t, tx.ID(), stored.ID())
	}
}
