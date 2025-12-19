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

func (h *BalanceHandler) HandleGetBalance(w http.ResponseWriter, r *http.Request) {
	userID := r.URL.Query().Get("user_id")
	if userID == "" {
		render.Status(r, http.StatusBadRequest)
		render.JSON(w, r, ErrorResponse{
			Error:   "missing_parameter",
			Message: "user_id query parameter is required",
		})
		return
	}

	if err := h.validator.Var(userID, "required,uuid"); err != nil {
		render.Status(r, http.StatusBadRequest)
		render.JSON(w, r, ErrorResponse{
			Error:   "validation_error",
			Message: "Invalid user ID format",
		})
		return
	}

	userIDVO, err := valueobject.NewUserID(userID)
	if err != nil {
		render.Status(r, http.StatusBadRequest)
		render.JSON(w, r, ErrorResponse{
			Error:   "validation_error",
			Message: "Invalid user ID format",
		})
		return
	}

	ctx, cancel := context.WithTimeout(r.Context(), 10*time.Second)
	defer cancel()

	response, err := h.walletService.GetBalanceByUserID(ctx, userIDVO)
	if err != nil {
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

	render.Status(r, http.StatusOK)
	render.JSON(w, r, response)
}
