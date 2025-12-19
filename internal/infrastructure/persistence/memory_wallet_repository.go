package persistence

import (
	"context"
	"errors"
	"sync"

	"bank/internal/domain/entity"
	"bank/internal/domain/valueobject"
)

var (
	ErrWalletNotFound = errors.New("wallet not found")
)

type MemoryWalletRepository struct {
	wallets map[string]*entity.Wallet
	mutex   sync.RWMutex
}

func NewMemoryWalletRepository() *MemoryWalletRepository {
	return &MemoryWalletRepository{
		wallets: make(map[string]*entity.Wallet),
		mutex:   sync.RWMutex{},
	}
}

func (r *MemoryWalletRepository) FindByUserID(ctx context.Context, userID valueobject.UserID) (*entity.Wallet, error) {
	r.mutex.RLock()
	defer r.mutex.RUnlock()

	wallet, exists := r.wallets[userID.String()]
	if !exists {
		return nil, ErrWalletNotFound
	}

	// Return a copy to prevent external mutations
	return r.copyWallet(wallet), nil
}

func (r *MemoryWalletRepository) Update(ctx context.Context, wallet *entity.Wallet) error {
	r.mutex.Lock()
	defer r.mutex.Unlock()

	if _, exists := r.wallets[wallet.UserID().String()]; !exists {
		return ErrWalletNotFound
	}

	r.wallets[wallet.UserID().String()] = r.copyWallet(wallet)
	return nil
}

func (r *MemoryWalletRepository) copyWallet(original *entity.Wallet) *entity.Wallet {
	return entity.ReconstructWallet(
		original.ID(),
		original.UserID(),
		original.Balance(),
	)
}

func (r *MemoryWalletRepository) AddWalletForTesting(wallet *entity.Wallet) {
	r.mutex.Lock()
	defer r.mutex.Unlock()
	r.wallets[wallet.UserID().String()] = r.copyWallet(wallet)
}
