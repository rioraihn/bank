package http

import (
	"context"
	"net/http"
	"time"

	"bank/internal/domain/usecase"
	"bank/internal/domain/valueobject"

	"github.com/go-chi/render"
	"github.com/go-playground/validator/v10"
)

type WithdrawHandler struct {
	withdrawUseCase usecase.WithdrawUseCase
	validator       *validator.Validate
}

func NewWithdrawHandler(withdrawUseCase usecase.WithdrawUseCase) *WithdrawHandler {
	return &WithdrawHandler{
		withdrawUseCase: withdrawUseCase,
		validator:       validator.New(),
	}
}

// WithdrawRequest represents the HTTP request for withdraw
type WithdrawRequest struct {
	UserID string `json:"user_id" validate:"required,uuid"`
	Amount int64  `json:"amount" validate:"required,gt=0"`
}

// @Summary Withdraw funds from a wallet
// @Description Withdraw a specified amount from a user's wallet
// @Tags wallet
// @Accept json
// @Produce json
// @Param request body WithdrawRequest true "Withdraw request"
// @Success 200 {object} dto.WithdrawResponse
// @Failure 400 {object} ErrorResponse
// @Failure 404 {object} ErrorResponse
// @Failure 500 {object} ErrorResponse
// @Router /withdraw [post]
func (h *WithdrawHandler) HandleWithdraw(w http.ResponseWriter, r *http.Request) {
	var req WithdrawRequest

	// Parse and validate request
	if err := render.DecodeJSON(r.Body, &req); err != nil {
		render.Status(r, http.StatusBadRequest)
		render.JSON(w, r, ErrorResponse{
			Error:   "invalid_request",
			Message: "Invalid JSON format",
		})
		return
	}

	// Validate request struct
	if err := h.validator.Struct(&req); err != nil {
		render.Status(r, http.StatusBadRequest)
		render.JSON(w, r, ErrorResponse{
			Error:   "validation_error",
			Message: err.Error(),
		})
		return
	}

	// Create domain value objects
	userIDVO, err := valueobject.NewUserID(req.UserID)
	if err != nil {
		render.Status(r, http.StatusBadRequest)
		render.JSON(w, r, ErrorResponse{
			Error:   "validation_error",
			Message: "Invalid user ID format",
		})
		return
	}

	amountVO, err := valueobject.NewMoney(req.Amount)
	if err != nil {
		render.Status(r, http.StatusBadRequest)
		render.JSON(w, r, ErrorResponse{
			Error:   "validation_error",
			Message: "Invalid amount",
		})
		return
	}

	// Create context with timeout
	ctx, cancel := context.WithTimeout(r.Context(), 10*time.Second)
	defer cancel()

	// Execute use case
	response, err := h.withdrawUseCase.Withdraw(ctx, userIDVO, amountVO)
	if err != nil {
		// Handle different types of errors
		switch {
		case err.Error() == "invalid user ID format":
			fallthrough
		case err.Error() == "user ID is required":
			render.Status(r, http.StatusBadRequest)
			render.JSON(w, r, ErrorResponse{
				Error:   "validation_error",
				Message: err.Error(),
			})
			return

		case err.Error() == "wallet not found":
			render.Status(r, http.StatusNotFound)
			render.JSON(w, r, ErrorResponse{
				Error:   "wallet_not_found",
				Message: "No wallet found for this user",
			})
			return

		case err.Error() == "insufficient funds":
			render.Status(r, http.StatusBadRequest)
			render.JSON(w, r, ErrorResponse{
				Error:   "insufficient_funds",
				Message: "Insufficient funds for withdrawal",
			})
			return

		default:
			render.Status(r, http.StatusInternalServerError)
			render.JSON(w, r, ErrorResponse{
				Error:   "internal_error",
				Message: "An unexpected error occurred",
			})
			return
		}
	}

	// Return success response
	render.Status(r, http.StatusOK)
	render.JSON(w, r, response)
}
