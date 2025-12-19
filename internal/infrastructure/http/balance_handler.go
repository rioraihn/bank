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
	balanceService service.BalanceService
	validator      *validator.Validate
}

func NewBalanceHandler(walletService service.BalanceService) *BalanceHandler {
	return &BalanceHandler{
		balanceService: walletService,
		validator:      validator.New(),
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

	response, err := h.balanceService.GetBalance(ctx, userIDVO)
	if err != nil {
		render.Status(r, http.StatusInternalServerError)
		render.JSON(w, r, ErrorResponse{
			Error:   err.Error(),
			Message: "An unexpected error occurred",
		})
		return
	}

	render.Status(r, http.StatusOK)
	render.JSON(w, r, response)
}
