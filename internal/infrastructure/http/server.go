package http

import (
	"log"
	"net/http"
	"time"

	"bank/internal/domain/service"
	"bank/internal/domain/usecase"

	"github.com/go-chi/render"
	"github.com/gorilla/mux"
)

type Server struct {
	router          *mux.Router
	withdrawHandler *WithdrawHandler
	balanceHandler  *BalanceHandler
}

func NewServer(
	withdrawUseCase usecase.WithdrawUseCase,
	walletService service.WalletService,
) *Server {
	server := &Server{
		router:          mux.NewRouter(),
		withdrawHandler: NewWithdrawHandler(withdrawUseCase),
		balanceHandler:  NewBalanceHandler(walletService),
	}

	server.setupRoutes()
	return server
}

func (s *Server) setupRoutes() {
	// Apply middleware
	s.router.Use(s.loggingMiddleware)
	s.router.Use(s.recoveryMiddleware)
	s.router.Use(s.requestIDMiddleware)
	s.router.Use(s.timeoutMiddleware)
	s.router.Use(s.contentTypeMiddleware)

	// Health check endpoint
	s.router.HandleFunc("/health", s.healthHandler).Methods("GET")

	// API routes with subrouter
	api := s.router.PathPrefix("/api/v1").Subrouter()
	api.HandleFunc("/withdraw", s.withdrawHandler.HandleWithdraw).Methods("POST")
	api.HandleFunc("/balance", s.balanceHandler.HandleGetBalance).Methods("GET")

	// Legacy routes (for backwards compatibility with PRD)
	s.router.HandleFunc("/withdraw", s.withdrawHandler.HandleWithdraw).Methods("POST")
	s.router.HandleFunc("/balance", s.balanceHandler.HandleGetBalance).Methods("GET")
}

// GetRouter returns the gorilla mux router
func (s *Server) GetRouter() *mux.Router {
	return s.router
}

// Middleware functions
func (s *Server) loggingMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()
		next.ServeHTTP(w, r)
		log.Printf("%s %s %s", r.Method, r.RequestURI, time.Since(start))
	})
}

func (s *Server) recoveryMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		defer func() {
			if err := recover(); err != nil {
				log.Printf("Panic: %v", err)
				http.Error(w, "Internal Server Error", http.StatusInternalServerError)
			}
		}()
		next.ServeHTTP(w, r)
	})
}

func (s *Server) requestIDMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		requestID := r.Header.Get("X-Request-ID")
		if requestID == "" {
			requestID = generateRequestID()
		}
		w.Header().Set("X-Request-ID", requestID)
		next.ServeHTTP(w, r)
	})
}

func (s *Server) timeoutMiddleware(next http.Handler) http.Handler {
	return http.TimeoutHandler(next, 60*time.Second, "Request timeout")
}

func (s *Server) contentTypeMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method == "POST" || r.Method == "PUT" || r.Method == "PATCH" {
			contentType := r.Header.Get("Content-Type")
			if contentType != "application/json" {
				http.Error(w, "Content-Type must be application/json", http.StatusUnsupportedMediaType)
				return
			}
		}
		w.Header().Set("Content-Type", "application/json")
		next.ServeHTTP(w, r)
	})
}

func (s *Server) healthHandler(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)
	render.JSON(w, r, HealthResponse{
		Status:  "ok",
		Message: "Wallet service is running",
	})
}

func generateRequestID() string {
	return time.Now().Format("20060102150405") + "-" + randomString(8)
}

func randomString(length int) string {
	const charset = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"
	b := make([]byte, length)
	for i := range b {
		b[i] = charset[time.Now().UnixNano()%int64(len(charset))]
	}
	return string(b)
}
