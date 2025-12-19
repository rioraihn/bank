package dto

type WithdrawRequest struct {
	UserID string `json:"user_id" validate:"required,uuid"`
	Amount int64  `json:"amount" validate:"required,gt=0"`
}

type WithdrawResponse struct {
	UserID          string `json:"user_id"`
	AmountWithdrawn int64  `json:"amount_withdrawn"`
	NewBalance      int64  `json:"new_balance"`
	Success         bool   `json:"success"`
	Message         string `json:"message,omitempty"`
}

type BalanceResponse struct {
	UserID  string `json:"user_id"`
	Balance int64  `json:"balance"`
}

type ErrorResponse struct {
	Error   string `json:"error"`
	Message string `json:"message,omitempty"`
}
