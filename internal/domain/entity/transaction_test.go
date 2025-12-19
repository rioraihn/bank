package entity

import (
	"testing"

	"bank/internal/domain/valueobject"
)

func TestNewTransaction(t *testing.T) {
	t.Run("should create a new withdrawal transaction", func(t *testing.T) {
		// Arrange
		walletID := valueobject.NewUserIDRandom()
		amount, _ := valueobject.NewMoney(20000)

		// Act
		tx := NewTransaction(walletID, TransactionTypeWithdrawal, amount)

		// Assert
		if tx.ID().String() == "" {
			t.Error("expected transaction ID to be generated")
		}
		if !tx.WalletID().Equals(walletID) {
			t.Error("wallet ID mismatch")
		}
		if tx.Type() != TransactionTypeWithdrawal {
			t.Errorf("expected withdrawal type, got %s", tx.Type())
		}
		if tx.Amount().Amount() != amount.Amount() {
			t.Errorf("expected amount %d, got %d", amount.Amount(), tx.Amount().Amount())
		}
		if tx.Status() != TransactionStatusPending {
			t.Errorf("expected pending status, got %s", tx.Status())
		}
		if tx.CreatedAt().IsZero() {
			t.Error("expected created at to be set")
		}
	})

	t.Run("should generate unique transaction IDs", func(t *testing.T) {
		// Arrange
		walletID := valueobject.NewUserIDRandom()
		amount, _ := valueobject.NewMoney(10000)

		// Act
		tx1 := NewTransaction(walletID, TransactionTypeWithdrawal, amount)
		tx2 := NewTransaction(walletID, TransactionTypeWithdrawal, amount)

		// Assert
		if tx1.ID().String() == tx2.ID().String() {
			t.Error("expected different transaction IDs")
		}
	})
}


func TestReconstructTransaction(t *testing.T) {
	// Arrange
	id := valueobject.NewUserIDRandom()
	walletID := valueobject.NewUserIDRandom()
	amount, _ := valueobject.NewMoney(15000)
	createdAt := "2023-01-01T10:00:00Z"
	failureReason := "Insufficient funds"

	// Act
	tx := ReconstructTransaction(
		id,
		walletID,
		TransactionTypeWithdrawal,
		amount,
		TransactionStatusFailed,
		createdAt,
		failureReason,
	)

	// Assert
	if !tx.ID().Equals(id) {
		t.Error("transaction ID mismatch")
	}
	if !tx.WalletID().Equals(walletID) {
		t.Error("wallet ID mismatch")
	}
	if tx.Type() != TransactionTypeWithdrawal {
		t.Errorf("expected withdrawal type, got %s", tx.Type())
	}
	if tx.Amount().Amount() != amount.Amount() {
		t.Errorf("expected amount %d, got %d", amount.Amount(), tx.Amount().Amount())
	}
	if tx.Status() != TransactionStatusFailed {
		t.Errorf("expected failed status, got %s", tx.Status())
	}
	if tx.FailureReason() != failureReason {
		t.Errorf("expected failure reason '%s', got '%s'", failureReason, tx.FailureReason())
	}
}