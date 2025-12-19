package service

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