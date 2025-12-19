package http

import (
	"net/http"
	"time"

	"bank/internal/domain/service"
	"bank/internal/domain/usecase"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/render"
)

// Server represents the HTTP server
type Server struct {
	router          *chi.Mux
	withdrawHandler *WithdrawHandler
	balanceHandler  *BalanceHandler
}

// NewServer creates a new HTTP server
func NewServer(
	withdrawUseCase usecase.WithdrawUseCase,
	walletService service.WalletService,
) *Server {
	server := &Server{
		router:          chi.NewRouter(),
		withdrawHandler: NewWithdrawHandler(withdrawUseCase),
		balanceHandler:  NewBalanceHandler(walletService),
	}

	server.setupRoutes()
	return server
}

// setupRoutes configures the HTTP routes
func (s *Server) setupRoutes() {
	// Middleware
	s.router.Use(middleware.Logger)
	s.router.Use(middleware.Recoverer)
	s.router.Use(middleware.RequestID)
	s.router.Use(middleware.Timeout(60 * time.Second))
	s.router.Use(middleware.AllowContentType("application/json"))
	s.router.Use(middleware.SetHeader("Content-Type", "application/json"))

	// Health check endpoint
	s.router.Get("/health", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		render.JSON(w, r, HealthResponse{
			Status:  "ok",
			Message: "Wallet service is running",
		})
	})

	// API routes
	s.router.Route("/api/v1", func(r chi.Router) {
		// Wallet endpoints
		r.Post("/withdraw", s.withdrawHandler.HandleWithdraw)
		r.Get("/balance", s.balanceHandler.HandleGetBalance)
	})

	// Legacy routes (for backwards compatibility with PRD)
	s.router.Post("/withdraw", s.withdrawHandler.HandleWithdraw)
	s.router.Get("/balance", s.balanceHandler.HandleGetBalance)
}

// GetRouter returns the chi router
func (s *Server) GetRouter() *chi.Mux {
	return s.router
}
