package valueobject

import (
	"testing"
)

func TestNewMoney(t *testing.T) {
	t.Run("should create valid Money with positive amount", func(t *testing.T) {
		// Arrange
		amount := int64(10000)

		// Act
		money, err := NewMoney(amount)

		// Assert
		if err != nil {
			t.Fatalf("expected no error, got %v", err)
		}
		if money.Amount() != amount {
			t.Errorf("expected amount %d, got %d", amount, money.Amount())
		}
	})

	t.Run("should create valid Money with zero amount", func(t *testing.T) {
		// Arrange
		amount := int64(0)

		// Act
		money, err := NewMoney(amount)

		// Assert
		if err != nil {
			t.Fatalf("expected no error, got %v", err)
		}
		if money.Amount() != amount {
			t.Errorf("expected amount %d, got %d", amount, money.Amount())
		}
	})

	t.Run("should return error for negative amount", func(t *testing.T) {
		// Arrange
		amount := int64(-1000)

		// Act
		_, err := NewMoney(amount)

		// Assert
		if err == nil {
			t.Error("expected error for negative amount, got nil")
		}
		if err.Error() != "money amount cannot be negative" {
			t.Errorf("expected 'money amount cannot be negative', got %v", err)
		}
	})
}

func TestMoneySubtract(t *testing.T) {
	t.Run("should subtract money correctly", func(t *testing.T) {
		// Arrange
		money1, _ := NewMoney(15000)
		money2, _ := NewMoney(5000)

		// Act
		result, err := money1.Subtract(money2)

		// Assert
		expected := int64(10000)
		if err != nil {
			t.Fatalf("expected no error, got %v", err)
		}
		if result.Amount() != expected {
			t.Errorf("expected %d, got %d", expected, result.Amount())
		}
	})

	t.Run("should return error when subtracting more than available", func(t *testing.T) {
		// Arrange
		money1, _ := NewMoney(5000)
		money2, _ := NewMoney(10000)

		// Act
		_, err := money1.Subtract(money2)

		// Assert
		if err == nil {
			t.Error("expected error for insufficient funds, got nil")
		}
		if err.Error() != "insufficient funds" {
			t.Errorf("expected 'insufficient funds', got %v", err)
		}
	})
}

func TestMoneyComparisons(t *testing.T) {
	t.Run("should check if money is less than or equal", func(t *testing.T) {
		// Arrange
		money1, _ := NewMoney(10000)
		money2, _ := NewMoney(10000)
		money3, _ := NewMoney(15000)

		// Act & Assert
		if !money1.LessThanOrEqual(money2) {
			t.Error("expected money1 to be less than or equal to money2")
		}
		if !money1.LessThanOrEqual(money3) {
			t.Error("expected money1 to be less than or equal to money3")
		}
		if money3.LessThanOrEqual(money1) {
			t.Error("expected money3 to not be less than or equal to money1")
		}
	})

	t.Run("should check if money is equal", func(t *testing.T) {
		// Arrange
		money1, _ := NewMoney(10000)
		money2, _ := NewMoney(10000)
		money3, _ := NewMoney(15000)

		// Act & Assert
		if money1.Amount() != money2.Amount() {
			t.Error("expected money1 to equal money2")
		}
		if money1.Amount() == money3.Amount() {
			t.Error("expected money1 to not equal money3")
		}
	})

	t.Run("should check if money is zero", func(t *testing.T) {
		// Arrange
		money1, _ := NewMoney(0)
		money2, _ := NewMoney(10000)

		// Act & Assert
		if !money1.IsZero() {
			t.Error("expected money1 to be zero")
		}
		if money2.IsZero() {
			t.Error("expected money2 to not be zero")
		}
	})
}

func TestMoneyString(t *testing.T) {
	t.Run("should format money as string correctly", func(t *testing.T) {
		// Arrange
		money, _ := NewMoney(15000)

		// Act
		result := money.String()

		// Assert
		expected := "15000"
		if result != expected {
			t.Errorf("expected %s, got %s", expected, result)
		}
	})
}
