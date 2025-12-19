package entity

import (
	"testing"

	"bank/internal/domain/valueobject"
)

func TestNewTransaction(t *testing.T) {
	tests := []struct {
		name        string
		txType      TransactionType
		amount      int64
		expectError bool
		description string
	}{
		{
			name:        "withdrawal transaction",
			txType:      TransactionTypeWithdrawal,
			amount:      20000,
			expectError: false,
			description: "should create a new withdrawal transaction with default pending status",
		},
		{
			name:        "deposit transaction",
			txType:      TransactionTypeDeposit,
			amount:      50000,
			expectError: false,
			description: "should create a new deposit transaction with default pending status",
		},
		{
			name:        "zero amount transaction",
			txType:      TransactionTypeWithdrawal,
			amount:      0,
			expectError: false,
			description: "should create transaction even with zero amount (validation at usecase level)",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Arrange
			walletID := valueobject.NewUserIDRandom()
			amount, _ := valueobject.NewMoney(tt.amount)

			// Act
			tx := NewTransaction(walletID, tt.txType, amount)

			// Assert
			if tx.ID().String() == "" {
				t.Error("expected transaction ID to be generated")
			}
			if !tx.WalletID().Equals(walletID) {
				t.Error("wallet ID mismatch")
			}
			if tx.Type() != tt.txType {
				t.Errorf("expected transaction type %s, got %s", tt.txType, tx.Type())
			}
			if tx.Amount().Amount() != tt.amount {
				t.Errorf("expected amount %d, got %d", tt.amount, tx.Amount().Amount())
			}
			if tx.Status() != TransactionStatusPending {
				t.Errorf("expected pending status, got %s", tx.Status())
			}
			if tx.CreatedAt().IsZero() {
				t.Error("expected created at to be set")
			}
			if tx.FailureReason() != "" {
				t.Errorf("expected empty failure reason, got '%s'", tx.FailureReason())
			}
		})
	}
}

func TestNewTransaction_UniqueIDs(t *testing.T) {
	t.Run("should generate unique transaction IDs", func(t *testing.T) {
		// Arrange
		walletID := valueobject.NewUserIDRandom()
		amount, _ := valueobject.NewMoney(10000)

		// Act
		tx1 := NewTransaction(walletID, TransactionTypeWithdrawal, amount)
		tx2 := NewTransaction(walletID, TransactionTypeWithdrawal, amount)
		tx3 := NewTransaction(walletID, TransactionTypeDeposit, amount)

		// Assert
		if tx1.ID().String() == tx2.ID().String() {
			t.Error("expected different transaction IDs for same type transactions")
		}
		if tx1.ID().String() == tx3.ID().String() {
			t.Error("expected different transaction IDs for different type transactions")
		}
		if tx2.ID().String() == tx3.ID().String() {
			t.Error("expected different transaction IDs for different type transactions")
		}
	})
}

func TestReconstructTransaction(t *testing.T) {
	tests := []struct {
		name          string
		txType        TransactionType
		status        TransactionStatus
		amount        int64
		failureReason string
		createdAt     string
		description   string
	}{
		{
			name:          "successful withdrawal",
			txType:        TransactionTypeWithdrawal,
			status:        TransactionStatusCompleted,
			amount:        15000,
			failureReason: "",
			createdAt:     "2023-01-01T10:00:00Z",
			description:   "should reconstruct completed withdrawal transaction",
		},
		{
			name:          "failed withdrawal",
			txType:        TransactionTypeWithdrawal,
			status:        TransactionStatusFailed,
			amount:        20000,
			failureReason: "Insufficient funds",
			createdAt:     "2023-01-01T11:00:00Z",
			description:   "should reconstruct failed withdrawal transaction with failure reason",
		},
		{
			name:          "pending deposit",
			txType:        TransactionTypeDeposit,
			status:        TransactionStatusPending,
			amount:        50000,
			failureReason: "",
			createdAt:     "2023-01-01T12:00:00Z",
			description:   "should reconstruct pending deposit transaction",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Arrange
			id := valueobject.NewUserIDRandom()
			walletID := valueobject.NewUserIDRandom()
			amount, _ := valueobject.NewMoney(tt.amount)

			// Act
			tx := ReconstructTransaction(
				id,
				walletID,
				tt.txType,
				amount,
				tt.status,
				tt.createdAt,
				tt.failureReason,
			)

			// Assert
			if !tx.ID().Equals(id) {
				t.Error("transaction ID mismatch")
			}
			if !tx.WalletID().Equals(walletID) {
				t.Error("wallet ID mismatch")
			}
			if tx.Type() != tt.txType {
				t.Errorf("expected transaction type %s, got %s", tt.txType, tx.Type())
			}
			if tx.Amount().Amount() != tt.amount {
				t.Errorf("expected amount %d, got %d", tt.amount, tx.Amount().Amount())
			}
			if tx.Status() != tt.status {
				t.Errorf("expected status %s, got %s", tt.status, tx.Status())
			}
			if tx.FailureReason() != tt.failureReason {
				t.Errorf("expected failure reason '%s', got '%s'", tt.failureReason, tx.FailureReason())
			}
			if tx.CreatedAt().IsZero() {
				t.Error("expected created at to be set")
			}
		})
	}
}

func TestTransaction_Constants(t *testing.T) {
	t.Run("should have correct transaction type constants", func(t *testing.T) {
		tests := []struct {
			name     string
			value    TransactionType
			expected string
		}{
			{"withdrawal constant", TransactionTypeWithdrawal, "WITHDRAWAL"},
			{"deposit constant", TransactionTypeDeposit, "DEPOSIT"},
		}

		for _, tt := range tests {
			t.Run(tt.name, func(t *testing.T) {
				if string(tt.value) != tt.expected {
					t.Errorf("expected %s, got %s", tt.expected, string(tt.value))
				}
			})
		}
	})

	t.Run("should have correct transaction status constants", func(t *testing.T) {
		tests := []struct {
			name     string
			value    TransactionStatus
			expected string
		}{
			{"pending constant", TransactionStatusPending, "PENDING"},
			{"completed constant", TransactionStatusCompleted, "COMPLETED"},
			{"failed constant", TransactionStatusFailed, "FAILED"},
		}

		for _, tt := range tests {
			t.Run(tt.name, func(t *testing.T) {
				if string(tt.value) != tt.expected {
					t.Errorf("expected %s, got %s", tt.expected, string(tt.value))
				}
			})
		}
	})
}
