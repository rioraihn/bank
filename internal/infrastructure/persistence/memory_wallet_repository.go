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
	ErrDuplicateWallet = errors.New("wallet already exists for user")
)

// MemoryWalletRepository is an in-memory implementation of WalletRepository
// For development and testing purposes only
type MemoryWalletRepository struct {
	wallets map[string]*entity.Wallet
	mutex   sync.RWMutex
}

// NewMemoryWalletRepository creates a new in-memory wallet repository
func NewMemoryWalletRepository() *MemoryWalletRepository {
	return &MemoryWalletRepository{
		wallets: make(map[string]*entity.Wallet),
		mutex:   sync.RWMutex{},
	}
}

// FindByUserID finds a wallet by user ID
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

// Save saves a wallet to the repository (upsert operation)
func (r *MemoryWalletRepository) Save(ctx context.Context, wallet *entity.Wallet) error {
	r.mutex.Lock()
	defer r.mutex.Unlock()

	r.wallets[wallet.UserID().String()] = r.copyWallet(wallet)
	return nil
}

// Create creates a new wallet
func (r *MemoryWalletRepository) Create(ctx context.Context, wallet *entity.Wallet) error {
	r.mutex.Lock()
	defer r.mutex.Unlock()

	if _, exists := r.wallets[wallet.UserID().String()]; exists {
		return ErrDuplicateWallet
	}

	r.wallets[wallet.UserID().String()] = r.copyWallet(wallet)
	return nil
}

// Update updates an existing wallet
func (r *MemoryWalletRepository) Update(ctx context.Context, wallet *entity.Wallet) error {
	r.mutex.Lock()
	defer r.mutex.Unlock()

	if _, exists := r.wallets[wallet.UserID().String()]; !exists {
		return ErrWalletNotFound
	}

	r.wallets[wallet.UserID().String()] = r.copyWallet(wallet)
	return nil
}

// Exists checks if a wallet exists for the given user ID
func (r *MemoryWalletRepository) Exists(ctx context.Context, userID valueobject.UserID) (bool, error) {
	r.mutex.RLock()
	defer r.mutex.RUnlock()

	_, exists := r.wallets[userID.String()]
	return exists, nil
}

// copyWallet creates a deep copy of a wallet to prevent mutations
func (r *MemoryWalletRepository) copyWallet(original *entity.Wallet) *entity.Wallet {
	// Reconstruct the wallet with the same data but new instance
	// This ensures that modifications to the returned wallet don't affect the stored one
	return entity.ReconstructWallet(
		original.ID(),
		original.UserID(),
		original.Balance(),
	)
}

// GetAll returns all wallets (for testing/debugging purposes)
func (r *MemoryWalletRepository) GetAll(ctx context.Context) []*entity.Wallet {
	r.mutex.RLock()
	defer r.mutex.RUnlock()

	wallets := make([]*entity.Wallet, 0, len(r.wallets))
	for _, wallet := range r.wallets {
		wallets = append(wallets, r.copyWallet(wallet))
	}

	return wallets
}

// Clear removes all wallets (for testing purposes)
func (r *MemoryWalletRepository) Clear() {
	r.mutex.Lock()
	defer r.mutex.Unlock()

	r.wallets = make(map[string]*entity.Wallet)
}