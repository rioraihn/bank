package usecase

import (
	"context"

	"bank/internal/domain/entity"
	"bank/internal/domain/valueobject"
	"github.com/stretchr/testify/mock"
)

// MockWalletRepository is a mock for WalletRepository using testify/mock
type MockWalletRepository struct {
	mock.Mock
}

func (m *MockWalletRepository) FindByUserID(ctx context.Context, userID valueobject.UserID) (*entity.Wallet, error) {
	args := m.Called(ctx, userID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*entity.Wallet), args.Error(1)
}

func (m *MockWalletRepository) Save(ctx context.Context, wallet *entity.Wallet) error {
	args := m.Called(ctx, wallet)
	return args.Error(0)
}

func (m *MockWalletRepository) Update(ctx context.Context, wallet *entity.Wallet) error {
	args := m.Called(ctx, wallet)
	return args.Error(0)
}

func (m *MockWalletRepository) Create(ctx context.Context, wallet *entity.Wallet) error {
	args := m.Called(ctx, wallet)
	return args.Error(0)
}

func (m *MockWalletRepository) Exists(ctx context.Context, userID valueobject.UserID) (bool, error) {
	args := m.Called(ctx, userID)
	return args.Bool(0), args.Error(1)
}

// MockTransactionRepository is a mock for TransactionRepository using testify/mock
type MockTransactionRepository struct {
	mock.Mock
}

func (m *MockTransactionRepository) Save(ctx context.Context, transaction *entity.Transaction) error {
	args := m.Called(ctx, transaction)
	return args.Error(0)
}

func (m *MockTransactionRepository) FindByID(ctx context.Context, transactionID valueobject.UserID) (*entity.Transaction, error) {
	args := m.Called(ctx, transactionID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*entity.Transaction), args.Error(1)
}

func (m *MockTransactionRepository) FindByWalletID(ctx context.Context, walletID valueobject.UserID, limit, offset int) ([]*entity.Transaction, error) {
	args := m.Called(ctx, walletID, limit, offset)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]*entity.Transaction), args.Error(1)
}

func (m *MockTransactionRepository) FindByWalletIDAndType(ctx context.Context, walletID valueobject.UserID, transactionType entity.TransactionType, limit, offset int) ([]*entity.Transaction, error) {
	args := m.Called(ctx, walletID, transactionType, limit, offset)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]*entity.Transaction), args.Error(1)
}

func (m *MockTransactionRepository) CountByWalletID(ctx context.Context, walletID valueobject.UserID) (int, error) {
	args := m.Called(ctx, walletID)
	return args.Int(0), args.Error(1)
}