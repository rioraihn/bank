package http

import (
	"bytes"
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	appservice "bank/internal/application/service"
	appusecase "bank/internal/application/usecase"
	"bank/internal/domain/entity"
	"bank/internal/domain/service"
	"bank/internal/domain/usecase"
	"bank/internal/domain/valueobject"
	"bank/internal/infrastructure/persistence"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

type TestSetup struct {
	Server          *Server
	WalletRepo      *persistence.MemoryWalletRepository
	TransactionRepo *persistence.MemoryTransactionRepository
	WithdrawUseCase usecase.WithdrawUseCase
	WalletService   service.WalletService
}

func setupTestServer() *TestSetup {
	// Create in-memory repositories
	walletRepo := persistence.NewMemoryWalletRepository()
	transactionRepo := persistence.NewMemoryTransactionRepository()

	// Create use cases and services
	withdrawUseCase := appusecase.NewWithdrawUseCase(walletRepo, transactionRepo)
	walletService := appservice.NewWalletService(walletRepo)

	// Create server
	server := NewServer(withdrawUseCase, walletService)

	return &TestSetup{
		Server:          server,
		WalletRepo:      walletRepo,
		TransactionRepo: transactionRepo,
		WithdrawUseCase: withdrawUseCase,
		WalletService:   walletService,
	}
}

func (ts *TestSetup) createTestWallet(t *testing.T, userID string, balance int64) {
	userIDVO, err := valueobject.NewUserID(userID)
	require.NoError(t, err)

	balanceVO, err := valueobject.NewMoney(balance)
	require.NoError(t, err)

	wallet := entity.NewWalletWithBalance(userIDVO, balanceVO)
	err = ts.WalletRepo.Create(context.Background(), wallet)
	require.NoError(t, err)
}

func TestHealthEndpoint(t *testing.T) {
	// Arrange
	setup := setupTestServer()
	req := httptest.NewRequest("GET", "/health", nil)
	w := httptest.NewRecorder()

	// Act
	setup.Server.GetRouter().ServeHTTP(w, req)

	// Assert
	assert.Equal(t, http.StatusOK, w.Code)

	var response HealthResponse
	err := json.Unmarshal(w.Body.Bytes(), &response)
	assert.NoError(t, err)
	assert.Equal(t, "ok", response.Status)
}

func TestWithdrawEndpoint_Success(t *testing.T) {
	// Arrange
	setup := setupTestServer()

	// Create a test wallet with initial balance
	setup.createTestWallet(t, "550e8400-e29b-41d4-a716-446655440000", 10000)

	requestBody := map[string]interface{}{
		"user_id": "550e8400-e29b-41d4-a716-446655440000",
		"amount":  2500,
	}

	bodyBytes, _ := json.Marshal(requestBody)
	req := httptest.NewRequest("POST", "/withdraw", bytes.NewBuffer(bodyBytes))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()

	// Act
	setup.Server.GetRouter().ServeHTTP(w, req)

	// Assert
	assert.Equal(t, http.StatusOK, w.Code)

	var response map[string]interface{}
	err := json.Unmarshal(w.Body.Bytes(), &response)
	assert.NoError(t, err)
	assert.True(t, response["success"].(bool))
	assert.Equal(t, float64(2500), response["amount_withdrawn"])
	assert.Equal(t, float64(7500), response["new_balance"])
}

func TestWithdrawEndpoint_InvalidRequest(t *testing.T) {
	// Arrange
	setup := setupTestServer()

	testCases := []struct {
		name          string
		requestBody   map[string]interface{}
		expectedError string
	}{
		{
			name: "missing user_id",
			requestBody: map[string]interface{}{
				"amount": 1000,
			},
			expectedError: "validation_error",
		},
		{
			name: "invalid user_id format",
			requestBody: map[string]interface{}{
				"user_id": "invalid-uuid",
				"amount":  1000,
			},
			expectedError: "validation_error",
		},
		{
			name: "missing amount",
			requestBody: map[string]interface{}{
				"user_id": "550e8400-e29b-41d4-a716-446655440000",
			},
			expectedError: "validation_error",
		},
		{
			name: "zero amount",
			requestBody: map[string]interface{}{
				"user_id": "550e8400-e29b-41d4-a716-446655440000",
				"amount":  0,
			},
			expectedError: "validation_error",
		},
		{
			name: "negative amount",
			requestBody: map[string]interface{}{
				"user_id": "550e8400-e29b-41d4-a716-446655440000",
				"amount":  -100,
			},
			expectedError: "validation_error",
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			bodyBytes, _ := json.Marshal(tc.requestBody)
			req := httptest.NewRequest("POST", "/withdraw", bytes.NewBuffer(bodyBytes))
			req.Header.Set("Content-Type", "application/json")
			w := httptest.NewRecorder()

			// Act
			setup.Server.GetRouter().ServeHTTP(w, req)

			// Assert
			assert.Equal(t, http.StatusBadRequest, w.Code)

			var response ErrorResponse
			err := json.Unmarshal(w.Body.Bytes(), &response)
			assert.NoError(t, err)
			assert.Equal(t, tc.expectedError, response.Error)
		})
	}
}

func TestWithdrawEndpoint_InsufficientFunds(t *testing.T) {
	// Arrange
	setup := setupTestServer()
	// Create a test wallet with low balance
	setup.createTestWallet(t, "550e8400-e29b-41d4-a716-446655440000", 1000)

	requestBody := map[string]interface{}{
		"user_id": "550e8400-e29b-41d4-a716-446655440000",
		"amount":  2500, // More than balance
	}

	bodyBytes, _ := json.Marshal(requestBody)
	req := httptest.NewRequest("POST", "/withdraw", bytes.NewBuffer(bodyBytes))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()

	// Act
	setup.Server.GetRouter().ServeHTTP(w, req)

	// Assert
	assert.Equal(t, http.StatusBadRequest, w.Code)

	var response ErrorResponse
	err := json.Unmarshal(w.Body.Bytes(), &response)
	assert.NoError(t, err)
	assert.Equal(t, "insufficient_funds", response.Error)
}

func TestWithdrawEndpoint_WalletNotFound(t *testing.T) {
	// Arrange
	setup := setupTestServer()

	requestBody := map[string]interface{}{
		"user_id": "550e8400-e29b-41d4-a716-446655440000",
		"amount":  1000,
	}

	bodyBytes, _ := json.Marshal(requestBody)
	req := httptest.NewRequest("POST", "/withdraw", bytes.NewBuffer(bodyBytes))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()

	// Act
	setup.Server.GetRouter().ServeHTTP(w, req)

	// Assert
	assert.Equal(t, http.StatusNotFound, w.Code)

	var response ErrorResponse
	err := json.Unmarshal(w.Body.Bytes(), &response)
	assert.NoError(t, err)
	assert.Equal(t, "wallet_not_found", response.Error)
}

func TestBalanceEndpoint_Success(t *testing.T) {
	// Arrange
	setup := setupTestServer()
	// Create a test wallet
	setup.createTestWallet(t, "550e8400-e29b-41d4-a716-446655440000", 5000)

	req := httptest.NewRequest("GET", "/balance?user_id=550e8400-e29b-41d4-a716-446655440000", nil)
	w := httptest.NewRecorder()

	// Act
	setup.Server.GetRouter().ServeHTTP(w, req)

	// Assert
	assert.Equal(t, http.StatusOK, w.Code)

	var response map[string]interface{}
	err := json.Unmarshal(w.Body.Bytes(), &response)
	assert.NoError(t, err)
	assert.Equal(t, "550e8400-e29b-41d4-a716-446655440000", response["user_id"])
	assert.Equal(t, float64(5000), response["balance"])
}

func TestBalanceEndpoint_InvalidRequest(t *testing.T) {
	// Arrange
	setup := setupTestServer()

	testCases := []struct {
		name           string
		url            string
		expectedError  string
		expectedStatus int
	}{
		{
			name:           "missing user_id",
			url:            "/balance",
			expectedError:  "missing_parameter",
			expectedStatus: http.StatusBadRequest,
		},
		{
			name:           "invalid user_id format",
			url:            "/balance?user_id=invalid-uuid",
			expectedError:  "validation_error",
			expectedStatus: http.StatusBadRequest,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			req := httptest.NewRequest("GET", tc.url, nil)
			w := httptest.NewRecorder()

			// Act
			setup.Server.GetRouter().ServeHTTP(w, req)

			// Assert
			assert.Equal(t, tc.expectedStatus, w.Code)

			var response ErrorResponse
			err := json.Unmarshal(w.Body.Bytes(), &response)
			assert.NoError(t, err)
			assert.Equal(t, tc.expectedError, response.Error)
		})
	}
}

func TestBalanceEndpoint_WalletNotFound(t *testing.T) {
	// Arrange
	setup := setupTestServer()

	req := httptest.NewRequest("GET", "/balance?user_id=550e8400-e29b-41d4-a716-446655440000", nil)
	w := httptest.NewRecorder()

	// Act
	setup.Server.GetRouter().ServeHTTP(w, req)

	// Assert
	assert.Equal(t, http.StatusNotFound, w.Code)

	var response ErrorResponse
	err := json.Unmarshal(w.Body.Bytes(), &response)
	assert.NoError(t, err)
	assert.Equal(t, "wallet_not_found", response.Error)
}
