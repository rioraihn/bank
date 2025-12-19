package error

import (
	"errors"
	"fmt"
)

var (
	ErrWalletNotFound      = errors.New("wallet not found")
	ErrTransactionNotFound = errors.New("transaction not found")
	ErrInvalidAmount       = errors.New("invalid amount")
	ErrInsufficientFunds   = errors.New("insufficient funds")
	ErrDuplicateWallet     = errors.New("wallet already exists for user")
)

type DomainError struct {
	Code    string
	Message string
	Cause   error
}

func (e *DomainError) Error() string {
	if e.Cause != nil {
		return fmt.Sprintf("%s: %s (caused by: %v)", e.Code, e.Message, e.Cause)
	}
	return fmt.Sprintf("%s: %s", e.Code, e.Message)
}

func (e *DomainError) Unwrap() error {
	return e.Cause
}

func NewDomainError(code, message string, cause error) *DomainError {
	return &DomainError{
		Code:    code,
		Message: message,
		Cause:   cause,
	}
}

func WalletError(message string, cause error) *DomainError {
	return NewDomainError("WALLET_ERROR", message, cause)
}

func TransactionError(message string, cause error) *DomainError {
	return NewDomainError("TRANSACTION_ERROR", message, cause)
}

func ValidationError(message string, cause error) *DomainError {
	return NewDomainError("VALIDATION_ERROR", message, cause)
}

func ConflictError(message string, cause error) *DomainError {
	return NewDomainError("CONFLICT_ERROR", message, cause)
}

func NotFoundError(resource string) *DomainError {
	return NewDomainError("NOT_FOUND", fmt.Sprintf("%s not found", resource), nil)
}
