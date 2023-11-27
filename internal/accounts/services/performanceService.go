package services

import (
	"context"
	repository "thyra/internal/accounts/repositories" // Replace with your actual package path
	"time"

	"github.com/google/uuid"
)

type AccountPerformanceService struct {
	repo *repository.AccountPerformanceRepository
}

func NewAccountPerformanceService(repo *repository.AccountPerformanceRepository) *AccountPerformanceService {
	return &AccountPerformanceService{
		repo: repo,
	}
}

// GetAccountValueChange calculates the value change for an account between two dates
func (s *AccountPerformanceService) GetAccountValueChange(ctx context.Context, accountID uuid.UUID, startDate, endDate time.Time) (repository.ValueChange, error) {
	return s.repo.GetAccountPerformanceChange(ctx, accountID, startDate, endDate)
}

func (s *AccountPerformanceService) GetUserValueChange(ctx context.Context, userID uuid.UUID, startDate, endDate time.Time) (repository.ValueChange, error) {
	return s.repo.GetUserPerformanceChange(ctx, userID, startDate, endDate)
}
