package valueobject

import (
	"errors"
	"strconv"
)

type Money struct {
	amount int64
}

func NewMoney(amount int64) (Money, error) {
	if amount < 0 {
		return Money{}, errors.New("money amount cannot be negative")
	}

	return Money{amount: amount}, nil
}

func (m Money) Amount() int64 {
	return m.amount
}

func (m Money) Subtract(other Money) (Money, error) {
	if m.amount < other.amount {
		return Money{}, errors.New("insufficient funds")
	}

	return Money{amount: m.amount - other.amount}, nil
}

func (m Money) LessThanOrEqual(other Money) bool {
	return m.amount <= other.amount
}


func (m Money) IsZero() bool {
	return m.amount == 0
}

func (m Money) String() string {
	return strconv.FormatInt(m.amount, 10)
}
