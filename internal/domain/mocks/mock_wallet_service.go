package mocks

import (
	"context"

	"bank/internal/application/dto"
	"bank/internal/domain/valueobject"

	"github.com/stretchr/testify/mock"
)

// MockWalletService is a mock implementation of service.WalletService
type MockWalletService struct {
	mock.Mock
}

func (m *MockWalletService) GetBalanceByUserID(ctx context.Context, userID valueobject.UserID) (*dto.BalanceResponse, error) {
	args := m.Called(ctx, userID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*dto.BalanceResponse), args.Error(1)
}

func NewMockWalletService() *MockWalletService {
	return &MockWalletService{}
}