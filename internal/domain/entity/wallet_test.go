package entity

import (
	"testing"

	"bank/internal/domain/valueobject"
)

func TestNewWallet(t *testing.T) {
	t.Run("should create a new wallet with zero balance", func(t *testing.T) {
		// Arrange
		userID := valueobject.NewUserIDRandom()

		// Act
		wallet := NewWallet(userID)

		// Assert
		if !wallet.UserID().Equals(userID) {
			t.Error("userID mismatch")
		}
		if !wallet.Balance().IsZero() {
			t.Errorf("expected zero balance, got %d", wallet.Balance().Amount())
		}
	})

	t.Run("should create a wallet with initial balance", func(t *testing.T) {
		// Arrange
		userID := valueobject.NewUserIDRandom()
		initialBalance, _ := valueobject.NewMoney(50000)

		// Act
		wallet := NewWalletWithBalance(userID, initialBalance)

		// Assert
		if !wallet.UserID().Equals(userID) {
			t.Error("userID mismatch")
		}
		if wallet.Balance().Amount() != initialBalance.Amount() {
			t.Errorf("expected balance %d, got %d", initialBalance.Amount(), wallet.Balance().Amount())
		}
	})
}

func TestWalletWithdraw(t *testing.T) {
	t.Run("should withdraw money successfully", func(t *testing.T) {
		// Arrange
		userID := valueobject.NewUserIDRandom()
		initialBalance, _ := valueobject.NewMoney(50000)
		withdrawAmount, _ := valueobject.NewMoney(20000)
		wallet := NewWalletWithBalance(userID, initialBalance)

		// Act
		err := wallet.Withdraw(withdrawAmount)

		// Assert
		if err != nil {
			t.Fatalf("expected no error, got %v", err)
		}

		expectedBalance, _ := valueobject.NewMoney(30000)
		if wallet.Balance().Amount() != expectedBalance.Amount() {
			t.Errorf("expected balance %d, got %d", expectedBalance.Amount(), wallet.Balance().Amount())
		}
	})

	t.Run("should fail when withdrawing more than balance", func(t *testing.T) {
		// Arrange
		userID := valueobject.NewUserIDRandom()
		initialBalance, _ := valueobject.NewMoney(10000)
		withdrawAmount, _ := valueobject.NewMoney(20000)
		wallet := NewWalletWithBalance(userID, initialBalance)

		// Act
		err := wallet.Withdraw(withdrawAmount)

		// Assert
		if err == nil {
			t.Error("expected error for insufficient funds, got nil")
		}
		if err.Error() != "insufficient funds" {
			t.Errorf("expected 'insufficient funds', got %v", err)
		}

		// Balance should remain unchanged
		if wallet.Balance().Amount() != initialBalance.Amount() {
			t.Errorf("balance should remain unchanged, expected %d, got %d", initialBalance.Amount(), wallet.Balance().Amount())
		}
	})

	t.Run("should fail when withdrawing zero amount", func(t *testing.T) {
		// Arrange
		userID := valueobject.NewUserIDRandom()
		initialBalance, _ := valueobject.NewMoney(10000)
		withdrawAmount, _ := valueobject.NewMoney(0)
		wallet := NewWalletWithBalance(userID, initialBalance)

		// Act
		err := wallet.Withdraw(withdrawAmount)

		// Assert
		if err == nil {
			t.Error("expected error for zero amount, got nil")
		}
		if err.Error() != "withdraw amount must be greater than zero" {
			t.Errorf("expected 'withdraw amount must be greater than zero', got %v", err)
		}
	})
}

func TestWalletCanWithdraw(t *testing.T) {
	t.Run("should return true when sufficient balance", func(t *testing.T) {
		// Arrange
		userID := valueobject.NewUserIDRandom()
		initialBalance, _ := valueobject.NewMoney(50000)
		withdrawAmount, _ := valueobject.NewMoney(20000)
		wallet := NewWalletWithBalance(userID, initialBalance)

		// Act
		canWithdraw := wallet.CanWithdraw(withdrawAmount)

		// Assert
		if !canWithdraw {
			t.Error("expected to be able to withdraw")
		}
	})

	t.Run("should return false when insufficient balance", func(t *testing.T) {
		// Arrange
		userID := valueobject.NewUserIDRandom()
		initialBalance, _ := valueobject.NewMoney(10000)
		withdrawAmount, _ := valueobject.NewMoney(20000)
		wallet := NewWalletWithBalance(userID, initialBalance)

		// Act
		canWithdraw := wallet.CanWithdraw(withdrawAmount)

		// Assert
		if canWithdraw {
			t.Error("expected not to be able to withdraw")
		}
	})

	t.Run("should return false when zero amount", func(t *testing.T) {
		// Arrange
		userID := valueobject.NewUserIDRandom()
		initialBalance, _ := valueobject.NewMoney(10000)
		withdrawAmount, _ := valueobject.NewMoney(0)
		wallet := NewWalletWithBalance(userID, initialBalance)

		// Act
		canWithdraw := wallet.CanWithdraw(withdrawAmount)

		// Assert
		if canWithdraw {
			t.Error("expected not to be able to withdraw zero amount")
		}
	})
}

func TestWalletID(t *testing.T) {
	t.Run("should generate unique wallet IDs", func(t *testing.T) {
		// Act
		wallet1 := NewWallet(valueobject.NewUserIDRandom())
		wallet2 := NewWallet(valueobject.NewUserIDRandom())

		// Assert
		if wallet1.ID().String() == wallet2.ID().String() {
			t.Error("expected different wallet IDs")
		}
	})
}