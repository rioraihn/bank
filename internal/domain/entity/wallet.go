package entity

import (
	"errors"

	"bank/internal/domain/valueobject"
)

type Wallet struct {
	id      valueobject.UserID // Using UserID as wallet ID for simplicity
	userID  valueobject.UserID
	balance valueobject.Money
}

func NewWallet(userID valueobject.UserID) *Wallet {
	balance, _ := valueobject.NewMoney(0)
	return &Wallet{
		id:      valueobject.NewUserIDRandom(),
		userID:  userID,
		balance: balance,
	}
}

func NewWalletWithBalance(userID valueobject.UserID, initialBalance valueobject.Money) *Wallet {
	return &Wallet{
		id:      valueobject.NewUserIDRandom(),
		userID:  userID,
		balance: initialBalance,
	}
}

func ReconstructWallet(id, userID valueobject.UserID, balance valueobject.Money) *Wallet {
	return &Wallet{
		id:      id,
		userID:  userID,
		balance: balance,
	}
}

func (w *Wallet) ID() valueobject.UserID {
	return w.id
}

func (w *Wallet) UserID() valueobject.UserID {
	return w.userID
}

func (w *Wallet) Balance() valueobject.Money {
	return w.balance
}

func (w *Wallet) Withdraw(amount valueobject.Money) error {
	if amount.IsZero() {
		return errors.New("withdraw amount must be greater than zero")
	}

	if !w.CanWithdraw(amount) {
		return errors.New("insufficient funds")
	}

	newBalance, err := w.balance.Subtract(amount)
	if err != nil {
		return err
	}

	w.balance = newBalance
	return nil
}

func (w *Wallet) CanWithdraw(amount valueobject.Money) bool {
	if amount.IsZero() {
		return false
	}
	return amount.LessThanOrEqual(w.balance)
}
