package analytics

import (
	"context"

	"github.com/sindwrr/test_storage/internal/models/analytics"
)

type AnalyticsService interface {
	ArtifactsPerDay(ctx context.Context) ([]analytics.DayCount, error)
	StatusDistribution(ctx context.Context) ([]analytics.StatusCount, error)
}
