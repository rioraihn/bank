package mocks

import (
	"context"

	"bank/internal/domain/entity"
	domainerr "bank/internal/domain/error"
	"bank/internal/domain/valueobject"
)

type MockTransactionRepository struct {
	transactions map[valueobject.UserID]*entity.Transaction
	walletTxs    map[valueobject.UserID][]valueobject.UserID // walletID -> transactionIDs
}

func NewMockTransactionRepository() *MockTransactionRepository {
	return &MockTransactionRepository{
		transactions: make(map[valueobject.UserID]*entity.Transaction),
		walletTxs:    make(map[valueobject.UserID][]valueobject.UserID),
	}
}

func (m *MockTransactionRepository) Save(ctx context.Context, transaction *entity.Transaction) error {
	m.transactions[transaction.ID()] = transaction
	m.walletTxs[transaction.WalletID()] = append(m.walletTxs[transaction.WalletID()], transaction.ID())
	return nil
}

func (m *MockTransactionRepository) FindByID(ctx context.Context, transactionID valueobject.UserID) (*entity.Transaction, error) {
	if tx, exists := m.transactions[transactionID]; exists {
		return tx, nil
	}
	return nil, domainerr.ErrTransactionNotFound
}

func (m *MockTransactionRepository) FindByWalletID(ctx context.Context, walletID valueobject.UserID, limit, offset int) ([]*entity.Transaction, error) {
	txs, exists := m.walletTxs[walletID]
	if !exists {
		return []*entity.Transaction{}, nil
	}

	var result []*entity.Transaction
	start := offset
	end := offset + limit

	if start > len(txs) {
		start = len(txs)
	}
	if end > len(txs) {
		end = len(txs)
	}

	for i := start; i < end; i++ {
		tx := m.transactions[txs[i]]
		result = append(result, tx)
	}

	return result, nil
}

func (m *MockTransactionRepository) FindByWalletIDAndType(
	ctx context.Context,
	walletID valueobject.UserID,
	transactionType entity.TransactionType,
	limit, offset int,
) ([]*entity.Transaction, error) {
	allTxs, err := m.FindByWalletID(ctx, walletID, 1000, 0) // Get all transactions
	if err != nil {
		return nil, err
	}

	var filtered []*entity.Transaction
	for _, tx := range allTxs {
		if tx.Type() == transactionType {
			filtered = append(filtered, tx)
		}
	}

	// Apply pagination
	start := offset
	end := offset + limit

	if start > len(filtered) {
		start = len(filtered)
	}
	if end > len(filtered) {
		end = len(filtered)
	}

	if start >= end {
		return []*entity.Transaction{}, nil
	}

	return filtered[start:end], nil
}

func (m *MockTransactionRepository) CountByWalletID(ctx context.Context, walletID valueobject.UserID) (int, error) {
	txs, exists := m.walletTxs[walletID]
	if !exists {
		return 0, nil
	}
	return len(txs), nil
}