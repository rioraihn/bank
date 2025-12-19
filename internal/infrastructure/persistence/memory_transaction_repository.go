package persistence

import (
	"context"
	"sync"
	"time"

	"bank/internal/domain/entity"
)

type MemoryTransactionRepository struct {
	transactions map[string]*entity.Transaction
	mutex        sync.RWMutex
}

func NewMemoryTransactionRepository() *MemoryTransactionRepository {
	return &MemoryTransactionRepository{
		transactions: make(map[string]*entity.Transaction),
		mutex:        sync.RWMutex{},
	}
}

func (r *MemoryTransactionRepository) Save(ctx context.Context, transaction *entity.Transaction) error {
	r.mutex.Lock()
	defer r.mutex.Unlock()

	r.transactions[transaction.ID().String()] = r.copyTransaction(transaction)
	return nil
}

func (r *MemoryTransactionRepository) copyTransaction(original *entity.Transaction) *entity.Transaction {
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
