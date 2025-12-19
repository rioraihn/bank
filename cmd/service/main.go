package main

import (
	"context"
	"flag"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"strconv"
	"syscall"
	"time"

	"github.com/joho/godotenv"

	appservice "bank/internal/application/service"
	appusecase "bank/internal/application/usecase"
	"bank/internal/domain/service"
	"bank/internal/domain/usecase"
	infrahttp "bank/internal/infrastructure/http"
	"bank/internal/infrastructure/persistence"
)

const (
	DefaultServerPort = "8080"

	DefaultServerHost = "0.0.0.0"

	ShutdownTimeout = 30 * time.Second
)

// AppConfig holds the application configuration
type AppConfig struct {
	ServerHost string
	ServerPort string
	Debug      bool
}

// Container holds all application dependencies
type Container struct {
	WalletRepo      *persistence.MemoryWalletRepository
	TransactionRepo *persistence.MemoryTransactionRepository
	WithdrawUseCase usecase.WithdrawUseCase
	WalletService   service.WalletService
	Server          *infrahttp.Server
}

func main() {
	// Load .env file if it exists
	if err := godotenv.Load(); err != nil {
		log.Println("No .env file found or error loading .env file, using environment variables or defaults")
	}

	config := parseFlags()

	setupLogging(config.Debug)

	container := setupContainer()

	if err := runApplication(container, config); err != nil {
		log.Fatalf("Failed to run application: %v", err)
	}
}

func parseFlags() *AppConfig {
	config := &AppConfig{}

	hostFlag := flag.String("host", "", "Server host")
	portFlag := flag.String("port", "", "Server port")
	debugFlag := flag.Bool("debug", false, "Enable debug logging")

	flag.Parse()

	config.ServerHost = getStringValue(*hostFlag, "SERVER_HOST", DefaultServerHost)
	config.ServerPort = getStringValue(*portFlag, "SERVER_PORT", DefaultServerPort)
	config.Debug = *debugFlag || getEnvBool("DEBUG", false)

	return config
}

func getStringValue(flagValue, envKey, defaultValue string) string {
	if flagValue != "" {
		return flagValue
	}
	if value := os.Getenv(envKey); value != "" {
		return value
	}
	return defaultValue
}

func getEnvBool(key string, defaultValue bool) bool {
	if value := os.Getenv(key); value != "" {
		if parsed, err := strconv.ParseBool(value); err == nil {
			return parsed
		}
	}
	return defaultValue
}

func setupLogging(debug bool) {
	if debug {
		log.SetFlags(log.LstdFlags | log.Lshortfile)
	} else {
		log.SetFlags(log.LstdFlags)
	}
}

func setupContainer() *Container {
	walletRepo := persistence.NewMemoryWalletRepository()
	transactionRepo := persistence.NewMemoryTransactionRepository()

	withdrawUseCase := appusecase.NewWithdrawUseCase(walletRepo, transactionRepo)
	walletService := appservice.NewWalletService(walletRepo)

	server := infrahttp.NewServer(withdrawUseCase, walletService)

	return &Container{
		WalletRepo:      walletRepo,
		TransactionRepo: transactionRepo,
		WithdrawUseCase: withdrawUseCase,
		WalletService:   walletService,
		Server:          server,
	}
}

func runApplication(container *Container, config *AppConfig) error {
	serverAddr := fmt.Sprintf("%s:%s", config.ServerHost, config.ServerPort)
	httpServer := &http.Server{
		Addr:         serverAddr,
		Handler:      container.Server.GetRouter(),
		ReadTimeout:  15 * time.Second,
		WriteTimeout: 15 * time.Second,
		IdleTimeout:  60 * time.Second,
	}

	serverErrors := make(chan error, 1)

	go func() {
		log.Printf("Starting wallet service on %s", serverAddr)
		log.Printf("Health check: http://%s/health", serverAddr)
		log.Printf("API Documentation:")
		log.Printf("  Withdraw: POST http://%s/withdraw", serverAddr)
		log.Printf("  Balance:  GET  http://%s/balance?user_id=<uuid>", serverAddr)
		log.Printf("Versioned API:")
		log.Printf("  Withdraw: POST http://%s/api/v1/withdraw", serverAddr)
		log.Printf("  Balance:  GET  http://%s/api/v1/balance?user_id=<uuid>", serverAddr)

		if err := httpServer.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			serverErrors <- fmt.Errorf("server failed to start: %w", err)
		}
	}()

	shutdown := make(chan os.Signal, 1)
	signal.Notify(shutdown, os.Interrupt, syscall.SIGTERM)

	select {
	case err := <-serverErrors:
		return err

	case sig := <-shutdown:
		log.Printf("Received shutdown signal: %v", sig)
		return gracefulShutdown(httpServer)
	}
}

func gracefulShutdown(server *http.Server) error {
	ctx, cancel := context.WithTimeout(context.Background(), ShutdownTimeout)
	defer cancel()

	log.Printf("Shutting down server gracefully (timeout: %v)...", ShutdownTimeout)

	if err := server.Shutdown(ctx); err != nil {
		log.Printf("Server shutdown error: %v", err)

		if err := server.Close(); err != nil {
			return fmt.Errorf("server forced shutdown failed: %w", err)
		}

		return fmt.Errorf("server graceful shutdown failed: %w", err)
	}

	log.Println("Server shutdown complete")
	return nil
}
