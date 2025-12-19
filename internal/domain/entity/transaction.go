package entity

import (
	"time"

	"bank/internal/domain/valueobject"
)

type TransactionType string

const (
	TransactionTypeWithdrawal TransactionType = "WITHDRAWAL"
	TransactionTypeDeposit    TransactionType = "DEPOSIT"
)

type TransactionStatus string

const (
	TransactionStatusPending   TransactionStatus = "PENDING"
	TransactionStatusCompleted TransactionStatus = "COMPLETED"
	TransactionStatusFailed    TransactionStatus = "FAILED"
)

type Transaction struct {
	id              valueobject.UserID
	walletID        valueobject.UserID
	transactionType TransactionType
	amount          valueobject.Money
	status          TransactionStatus
	failureReason   string
	createdAt       time.Time
}

func NewTransaction(walletID valueobject.UserID, txType TransactionType, amount valueobject.Money) *Transaction {
	return &Transaction{
		id:              valueobject.NewUserIDRandom(),
		walletID:        walletID,
		transactionType: txType,
		amount:          amount,
		status:          TransactionStatusPending,
		failureReason:   "",
		createdAt:       time.Now().UTC(),
	}
}

func ReconstructTransaction(
	id valueobject.UserID,
	walletID valueobject.UserID,
	txType TransactionType,
	amount valueobject.Money,
	status TransactionStatus,
	createdAt string,
	failureReason string,
) *Transaction {
	parsedTime, _ := time.Parse(time.RFC3339, createdAt)

	return &Transaction{
		id:              id,
		walletID:        walletID,
		transactionType: txType,
		amount:          amount,
		status:          status,
		failureReason:   failureReason,
		createdAt:       parsedTime,
	}
}

func (t *Transaction) ID() valueobject.UserID {
	return t.id
}

func (t *Transaction) WalletID() valueobject.UserID {
	return t.walletID
}

func (t *Transaction) Type() TransactionType {
	return t.transactionType
}

func (t *Transaction) Amount() valueobject.Money {
	return t.amount
}

func (t *Transaction) Status() TransactionStatus {
	return t.status
}

func (t *Transaction) FailureReason() string {
	return t.failureReason
}

func (t *Transaction) CreatedAt() time.Time {
	return t.createdAt
}

