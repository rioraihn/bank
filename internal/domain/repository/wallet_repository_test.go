package repository

import (
	"context"
	"testing"

	"bank/internal/domain/entity"
	"bank/internal/domain/repository/mocks"
	"bank/internal/domain/valueobject"
)

// Test that our mock implementation satisfies the interface
func TestMockWalletRepositoryImplementsInterface(t *testing.T) {
	// This is just a compile-time check to ensure our mock implements the interface
	var _ WalletRepository = (*mocks.MockWalletRepository)(nil)
}

// Test that our mock implementation behaves correctly for basic scenarios
func TestMockWalletRepository_BasicBehavior(t *testing.T) {
	ctx := context.Background()
	repo := mocks.NewMockWalletRepository()

	userID := valueobject.NewUserIDRandom()
	initialBalance, _ := valueobject.NewMoney(50000)
	wallet := entity.NewWalletWithBalance(userID, initialBalance)

	// Test Create
	err := repo.Create(ctx, wallet)
	if err != nil {
		t.Fatalf("expected no error creating wallet, got %v", err)
	}

	// Test FindByUserID
	found, err := repo.FindByUserID(ctx, userID)
	if err != nil {
		t.Fatalf("expected no error finding wallet, got %v", err)
	}
	if !found.UserID().Equals(userID) {
		t.Error("userID mismatch")
	}

	// Test Exists
	exists, err := repo.Exists(ctx, userID)
	if err != nil {
		t.Fatalf("expected no error checking existence, got %v", err)
	}
	if !exists {
		t.Error("expected wallet to exist")
	}
}