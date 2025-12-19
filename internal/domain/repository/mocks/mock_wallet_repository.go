package mocks

import (
	"context"

	"bank/internal/domain/entity"
	domainerr "bank/internal/domain/error"
	"bank/internal/domain/valueobject"
)

type MockWalletRepository struct {
	wallets map[valueobject.UserID]*entity.Wallet
}

func NewMockWalletRepository() *MockWalletRepository {
	return &MockWalletRepository{
		wallets: make(map[valueobject.UserID]*entity.Wallet),
	}
}

func (m *MockWalletRepository) FindByUserID(ctx context.Context, userID valueobject.UserID) (*entity.Wallet, error) {
	if wallet, exists := m.wallets[userID]; exists {
		return wallet, nil
	}
	return nil, domainerr.ErrWalletNotFound
}

func (m *MockWalletRepository) Save(ctx context.Context, wallet *entity.Wallet) error {
	m.wallets[wallet.UserID()] = wallet
	return nil
}

func (m *MockWalletRepository) Create(ctx context.Context, wallet *entity.Wallet) error {
	if _, exists := m.wallets[wallet.UserID()]; exists {
		return domainerr.ErrDuplicateWallet
	}
	m.wallets[wallet.UserID()] = wallet
	return nil
}

func (m *MockWalletRepository) Update(ctx context.Context, wallet *entity.Wallet) error {
	if _, exists := m.wallets[wallet.UserID()]; !exists {
		return domainerr.ErrWalletNotFound
	}
	m.wallets[wallet.UserID()] = wallet
	return nil
}

func (m *MockWalletRepository) Exists(ctx context.Context, userID valueobject.UserID) (bool, error) {
	_, exists := m.wallets[userID]
	return exists, nil
}