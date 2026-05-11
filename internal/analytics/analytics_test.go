package analytics

import (
	"context"
	"errors"
	"testing"

	"github.com/sindwrr/test_storage/internal/models/analytics"
)

type mockRepo struct {
	artifactsPerDayFn func(ctx context.Context) ([]analytics.DayCount, error)
	statusDistFn      func(ctx context.Context) ([]analytics.StatusCount, error)
}

func (m *mockRepo) ArtifactsPerDay(ctx context.Context) ([]analytics.DayCount, error) {
	if m.artifactsPerDayFn != nil {
		return m.artifactsPerDayFn(ctx)
	}
	return nil, nil
}

func (m *mockRepo) StatusDistribution(ctx context.Context) ([]analytics.StatusCount, error) {
	if m.statusDistFn != nil {
		return m.statusDistFn(ctx)
	}
	return nil, nil
}

func TestArtifactsPerDay_Success(t *testing.T) {
	expected := []analytics.DayCount{
		{Date: "2026-05-10", Count: 5},
		{Date: "2026-05-11", Count: 3},
	}
	repo := &mockRepo{
		artifactsPerDayFn: func(ctx context.Context) ([]analytics.DayCount, error) {
			return expected, nil
		},
	}
	svc := &analyticsService{repo: repo}

	got, err := svc.ArtifactsPerDay(context.Background())
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(got) != 2 || got[0].Count != 5 || got[1].Date != "2026-05-11" {
		t.Errorf("unexpected result: %+v", got)
	}
}

func TestArtifactsPerDay_Error(t *testing.T) {
	repo := &mockRepo{
		artifactsPerDayFn: func(ctx context.Context) ([]analytics.DayCount, error) {
			return nil, errors.New("db down")
		},
	}
	svc := &analyticsService{repo: repo}

	_, err := svc.ArtifactsPerDay(context.Background())
	if err == nil || err.Error() != "db down" {
		t.Fatalf("expected error 'db down', got %v", err)
	}
}

func TestStatusDistribution_Success(t *testing.T) {
	expected := []analytics.StatusCount{
		{Status: "passed", Count: 10},
		{Status: "failed", Count: 2},
	}
	repo := &mockRepo{
		statusDistFn: func(ctx context.Context) ([]analytics.StatusCount, error) {
			return expected, nil
		},
	}
	svc := &analyticsService{repo: repo}

	got, err := svc.StatusDistribution(context.Background())
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(got) != 2 || got[0].Status != "passed" || got[1].Count != 2 {
		t.Errorf("unexpected result: %+v", got)
	}
}

func TestStatusDistribution_Error(t *testing.T) {
	repo := &mockRepo{
		statusDistFn: func(ctx context.Context) ([]analytics.StatusCount, error) {
			return nil, errors.New("timeout")
		},
	}
	svc := &analyticsService{repo: repo}

	_, err := svc.StatusDistribution(context.Background())
	if err == nil || err.Error() != "timeout" {
		t.Fatalf("expected error 'timeout', got %v", err)
	}
}
