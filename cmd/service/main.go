package main

import (
	"context"
	"database/sql"
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
	"bank/internal/domain/repository"
	"bank/internal/domain/service"
	"bank/internal/domain/usecase"
	"bank/internal/infrastructure/database"
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
	ServerHost             string
	ServerPort             string
	Debug                  bool
	FailFastOnDBConnection bool // If true, app fails to start if DB is not connected
}

// Container holds all application dependencies
type Container struct {
	DB              *sql.DB
	WalletRepo      repository.WalletRepository
	TransactionRepo repository.TransactionRepository
	WithdrawUseCase usecase.WithdrawUseCase
	BalanceService  service.BalanceService
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

	log.Printf("✅ Database connection established and migrations completed")

	if err := runApplication(container, config); err != nil {
		log.Fatalf("Failed to run application: %v", err)
	}
}

func parseFlags() *AppConfig {
	config := &AppConfig{}

	hostFlag := flag.String("host", "", "Server host")
	portFlag := flag.String("port", "", "Server port")
	debugFlag := flag.Bool("debug", false, "Enable debug logging")
	failFastFlag := flag.Bool("fail-fast-db", true, "Fail to start if database connection fails")

	flag.Parse()

	config.ServerHost = getStringValue(*hostFlag, "SERVER_HOST", DefaultServerHost)
	config.ServerPort = getStringValue(*portFlag, "SERVER_PORT", DefaultServerPort)
	config.Debug = *debugFlag || getEnvBool("DEBUG", false)
	config.FailFastOnDBConnection = *failFastFlag || getEnvBool("FAIL_FAST_DB", true)

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
	// Connect to real database
	dbConfig := database.NewDatabaseConfig()

	log.Printf("🔌 Connecting to database: %s:%s/%s", dbConfig.Host, dbConfig.Port, dbConfig.DBName)
	db, err := database.ConnectToDatabase(dbConfig)
	if err != nil {
		log.Fatalf("❌ Failed to connect to database: %v", err)
	}

	// Use real database repositories with SQL query execution
	walletRepo := persistence.NewWalletRepository(db)
	transactionRepo := persistence.NewTransactionRepository(db)

	withdrawUseCase := appusecase.NewWithdrawUseCase(walletRepo, transactionRepo, db)
	BalanceService := appservice.NewBalanceUseCase(walletRepo)

	server := infrahttp.NewServer(withdrawUseCase, BalanceService)

	return &Container{
		DB:              db,
		WalletRepo:      walletRepo,
		TransactionRepo: transactionRepo,
		WithdrawUseCase: withdrawUseCase,
		BalanceService:  BalanceService,
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
		log.Printf("  Withdraw: POST http://%s/withdraw", serverAddr)
		log.Printf("  Balance:  GET  http://%s/balance?user_id=<uuid>", serverAddr)

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
