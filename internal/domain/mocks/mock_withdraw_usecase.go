package mocks

import (
	"context"

	"bank/internal/application/dto"
	"bank/internal/domain/valueobject"

	"github.com/stretchr/testify/mock"
)

// MockWithdrawUseCase is a mock implementation of usecase.WithdrawUseCase
type MockWithdrawUseCase struct {
	mock.Mock
}

func (m *MockWithdrawUseCase) Withdraw(ctx context.Context, userID valueobject.UserID, amount valueobject.Money) (*dto.WithdrawResponse, error) {
	args := m.Called(ctx, userID, amount)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*dto.WithdrawResponse), args.Error(1)
}

func NewMockWithdrawUseCase() *MockWithdrawUseCase {
	return &MockWithdrawUseCase{}
}