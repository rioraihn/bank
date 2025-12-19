package persistence

import (
	"context"
	"errors"
	"sort"
	"sync"
	"time"

	"bank/internal/domain/entity"
	"bank/internal/domain/valueobject"
)

var (
	ErrTransactionNotFound = errors.New("transaction not found")
)

// MemoryTransactionRepository is an in-memory implementation of TransactionRepository
// For development and testing purposes only
type MemoryTransactionRepository struct {
	transactions map[string]*entity.Transaction
	mutex        sync.RWMutex
}

// NewMemoryTransactionRepository creates a new in-memory transaction repository
func NewMemoryTransactionRepository() *MemoryTransactionRepository {
	return &MemoryTransactionRepository{
		transactions: make(map[string]*entity.Transaction),
		mutex:        sync.RWMutex{},
	}
}

// Save saves a transaction to the repository
func (r *MemoryTransactionRepository) Save(ctx context.Context, transaction *entity.Transaction) error {
	r.mutex.Lock()
	defer r.mutex.Unlock()

	r.transactions[transaction.ID().String()] = r.copyTransaction(transaction)
	return nil
}

// FindByID finds a transaction by its ID
func (r *MemoryTransactionRepository) FindByID(ctx context.Context, transactionID valueobject.UserID) (*entity.Transaction, error) {
	r.mutex.RLock()
	defer r.mutex.RUnlock()

	transaction, exists := r.transactions[transactionID.String()]
	if !exists {
		return nil, ErrTransactionNotFound
	}

	return r.copyTransaction(transaction), nil
}

// FindByWalletID finds all transactions for a given wallet ID
func (r *MemoryTransactionRepository) FindByWalletID(ctx context.Context, walletID valueobject.UserID, limit, offset int) ([]*entity.Transaction, error) {
	r.mutex.RLock()
	defer r.mutex.RUnlock()

	var transactions []*entity.Transaction
	for _, transaction := range r.transactions {
		if transaction.WalletID().Equals(walletID) {
			transactions = append(transactions, r.copyTransaction(transaction))
		}
	}

	// Sort by creation time (newest first)
	sort.Slice(transactions, func(i, j int) bool {
		return transactions[i].CreatedAt().After(transactions[j].CreatedAt())
	})

	// Apply pagination
	if offset >= len(transactions) {
		return []*entity.Transaction{}, nil
	}

	end := offset + limit
	if end > len(transactions) {
		end = len(transactions)
	}

	return transactions[offset:end], nil
}

// FindByWalletIDAndType finds transactions by wallet ID and transaction type
func (r *MemoryTransactionRepository) FindByWalletIDAndType(
	ctx context.Context,
	walletID valueobject.UserID,
	transactionType entity.TransactionType,
	limit, offset int,
) ([]*entity.Transaction, error) {
	r.mutex.RLock()
	defer r.mutex.RUnlock()

	var transactions []*entity.Transaction
	for _, transaction := range r.transactions {
		if transaction.WalletID().Equals(walletID) && transaction.Type() == transactionType {
			transactions = append(transactions, r.copyTransaction(transaction))
		}
	}

	// Sort by creation time (newest first)
	sort.Slice(transactions, func(i, j int) bool {
		return transactions[i].CreatedAt().After(transactions[j].CreatedAt())
	})

	// Apply pagination
	if offset >= len(transactions) {
		return []*entity.Transaction{}, nil
	}

	end := offset + limit
	if end > len(transactions) {
		end = len(transactions)
	}

	return transactions[offset:end], nil
}

// CountByWalletID returns the total count of transactions for a wallet
func (r *MemoryTransactionRepository) CountByWalletID(ctx context.Context, walletID valueobject.UserID) (int, error) {
	r.mutex.RLock()
	defer r.mutex.RUnlock()

	count := 0
	for _, transaction := range r.transactions {
		if transaction.WalletID().Equals(walletID) {
			count++
		}
	}

	return count, nil
}

// copyTransaction creates a deep copy of a transaction to prevent mutations
func (r *MemoryTransactionRepository) copyTransaction(original *entity.Transaction) *entity.Transaction {
	// Reconstruct the transaction with the same data but new instance
	// This ensures that modifications to the returned transaction don't affect the stored one
	return entity.ReconstructTransaction(
		original.ID(),
		original.WalletID(),
		original.Type(),
		original.Amount(),
		original.Status(),
		original.CreatedAt().Format(time.RFC3339),
		original.FailureReason(),
	)
}

// GetAll returns all transactions (for testing/debugging purposes)
func (r *MemoryTransactionRepository) GetAll(ctx context.Context) []*entity.Transaction {
	r.mutex.RLock()
	defer r.mutex.RUnlock()

	transactions := make([]*entity.Transaction, 0, len(r.transactions))
	for _, transaction := range r.transactions {
		transactions = append(transactions, r.copyTransaction(transaction))
	}

	// Sort by creation time (newest first)
	sort.Slice(transactions, func(i, j int) bool {
		return transactions[i].CreatedAt().After(transactions[j].CreatedAt())
	})

	return transactions
}

// Clear removes all transactions (for testing purposes)
func (r *MemoryTransactionRepository) Clear() {
	r.mutex.Lock()
	defer r.mutex.Unlock()

	r.transactions = make(map[string]*entity.Transaction)
}