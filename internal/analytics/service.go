package analytics

import (
	"context"
	"database/sql"

	"github.com/sindwrr/test_storage/internal/analytics/repository"
	"github.com/sindwrr/test_storage/internal/models/analytics"
)

type analyticsService struct {
	repo repository.AnalyticsRepository
}

func NewService(db *sql.DB) AnalyticsService {
	return &analyticsService{repo: repository.NewAnalyticsRepo(db)}
}

func (s *analyticsService) ArtifactsPerDay(ctx context.Context) ([]analytics.DayCount, error) {
	return s.repo.ArtifactsPerDay(ctx)
}

func (s *analyticsService) StatusDistribution(ctx context.Context) ([]analytics.StatusCount, error) {
	return s.repo.StatusDistribution(ctx)
}
