package http

import (
	"context"
	"net/http"
	"time"

	"bank/internal/domain/service"
	"bank/internal/domain/valueobject"

	"github.com/go-chi/render"
	"github.com/go-playground/validator/v10"
)

type BalanceHandler struct {
	walletService service.WalletService
	validator     *validator.Validate
}

func NewBalanceHandler(walletService service.WalletService) *BalanceHandler {
	return &BalanceHandler{
		walletService: walletService,
		validator:     validator.New(),
	}
}

// @Summary Get wallet balance
// @Description Get the current balance for a user's wallet
// @Tags wallet
// @Accept json
// @Produce json
// @Param user_id query string true "User ID" Format(uuid)
// @Success 200 {object} dto.BalanceResponse
// @Failure 400 {object} ErrorResponse
// @Failure 404 {object} ErrorResponse
// @Failure 500 {object} ErrorResponse
// @Router /balance [get]
func (h *BalanceHandler) HandleGetBalance(w http.ResponseWriter, r *http.Request) {
	// Get user_id from query parameter
	userID := r.URL.Query().Get("user_id")
	if userID == "" {
		render.Status(r, http.StatusBadRequest)
		render.JSON(w, r, ErrorResponse{
			Error:   "missing_parameter",
			Message: "user_id query parameter is required",
		})
		return
	}

	// Validate user ID
	if err := h.validator.Var(userID, "required,uuid"); err != nil {
		render.Status(r, http.StatusBadRequest)
		render.JSON(w, r, ErrorResponse{
			Error:   "validation_error",
			Message: "Invalid user ID format",
		})
		return
	}

	// Create domain value object
	userIDVO, err := valueobject.NewUserID(userID)
	if err != nil {
		render.Status(r, http.StatusBadRequest)
		render.JSON(w, r, ErrorResponse{
			Error:   "validation_error",
			Message: "Invalid user ID format",
		})
		return
	}

	// Create context with timeout
	ctx, cancel := context.WithTimeout(r.Context(), 10*time.Second)
	defer cancel()

	// Execute service
	response, err := h.walletService.GetBalanceByUserID(ctx, userIDVO)
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
