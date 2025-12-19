package error

import (
	"errors"
	"testing"
)

func TestDomainError(t *testing.T) {
	t.Run("should create domain error with code and message", func(t *testing.T) {
		// Arrange
		code := "TEST_ERROR"
		message := "This is a test error"

		// Act
		err := NewDomainError(code, message, nil)

		// Assert
		expected := "TEST_ERROR: This is a test error"
		if err.Error() != expected {
			t.Errorf("expected '%s', got '%s'", expected, err.Error())
		}
		if err.Unwrap() != nil {
			t.Error("expected nil unwrap")
		}
	})

	t.Run("should create domain error with cause", func(t *testing.T) {
		// Arrange
		code := "TEST_ERROR"
		message := "This is a test error"
		cause := errors.New("underlying error")

		// Act
		err := NewDomainError(code, message, cause)

		// Assert
		expected := "TEST_ERROR: This is a test error (caused by: underlying error)"
		if err.Error() != expected {
			t.Errorf("expected '%s', got '%s'", expected, err.Error())
		}
		if err.Unwrap() != cause {
			t.Error("expected cause to be unwrapped")
		}
	})
}

func TestSpecializedErrors(t *testing.T) {
	t.Run("should create wallet error", func(t *testing.T) {
		// Arrange
		message := "Wallet operation failed"
		cause := ErrInsufficientFunds

		// Act
		err := WalletError(message, cause)

		// Assert
		if err.Code != "WALLET_ERROR" {
			t.Errorf("expected code 'WALLET_ERROR', got '%s'", err.Code)
		}
		if err.Message != message {
			t.Errorf("expected message '%s', got '%s'", message, err.Message)
		}
		if err.Unwrap() != cause {
			t.Error("expected cause to be unwrapped")
		}
	})

	t.Run("should create transaction error", func(t *testing.T) {
		// Arrange
		message := "Transaction failed"
		cause := ErrInvalidAmount

		// Act
		err := TransactionError(message, cause)

		// Assert
		if err.Code != "TRANSACTION_ERROR" {
			t.Errorf("expected code 'TRANSACTION_ERROR', got '%s'", err.Code)
		}
		if err.Message != message {
			t.Errorf("expected message '%s', got '%s'", message, err.Message)
		}
	})

	t.Run("should create validation error", func(t *testing.T) {
		// Arrange
		message := "Invalid input"

		// Act
		err := ValidationError(message, nil)

		// Assert
		if err.Code != "VALIDATION_ERROR" {
			t.Errorf("expected code 'VALIDATION_ERROR', got '%s'", err.Code)
		}
		if err.Message != message {
			t.Errorf("expected message '%s', got '%s'", message, err.Message)
		}
	})

	t.Run("should create conflict error", func(t *testing.T) {
		// Arrange
		message := "Resource conflict"

		// Act
		err := ConflictError(message, nil)

		// Assert
		if err.Code != "CONFLICT_ERROR" {
			t.Errorf("expected code 'CONFLICT_ERROR', got '%s'", err.Code)
		}
		if err.Message != message {
			t.Errorf("expected message '%s', got '%s'", message, err.Message)
		}
	})

	t.Run("should create not found error", func(t *testing.T) {
		// Arrange
		resource := "Wallet"

		// Act
		err := NotFoundError(resource)

		// Assert
		if err.Code != "NOT_FOUND" {
			t.Errorf("expected code 'NOT_FOUND', got '%s'", err.Code)
		}
		expectedMessage := "Wallet not found"
		if err.Message != expectedMessage {
			t.Errorf("expected message '%s', got '%s'", expectedMessage, err.Message)
		}
		if err.Unwrap() != nil {
			t.Error("expected nil unwrap")
		}
	})
}

func TestCommonErrors(t *testing.T) {
	t.Run("should have predefined errors", func(t *testing.T) {
		// Assert
		if ErrWalletNotFound == nil {
			t.Error("ErrWalletNotFound should not be nil")
		}
		if ErrTransactionNotFound == nil {
			t.Error("ErrTransactionNotFound should not be nil")
		}
		if ErrInvalidAmount == nil {
			t.Error("ErrInvalidAmount should not be nil")
		}
		if ErrInsufficientFunds == nil {
			t.Error("ErrInsufficientFunds should not be nil")
		}
		if ErrDuplicateWallet == nil {
			t.Error("ErrDuplicateWallet should not be nil")
		}

		// Test error messages
		if ErrWalletNotFound.Error() != "wallet not found" {
			t.Errorf("unexpected error message: %s", ErrWalletNotFound.Error())
		}
		if ErrInsufficientFunds.Error() != "insufficient funds" {
			t.Errorf("unexpected error message: %s", ErrInsufficientFunds.Error())
		}
	})
}