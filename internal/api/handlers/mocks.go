package handlers

import (
	"bytes"
	"context"
	"errors"
	"io"
	"mime/multipart"
	"time"

	"github.com/sindwrr/test_storage/internal/models"
	"github.com/sindwrr/test_storage/internal/models/analytics"
)

var errSome = errors.New("test error")

// -------------------------------------------
// Metadata
// -------------------------------------------
type mockMetadataService struct {
	createArtifactFn  func(filePath string, fileSize int64, component, build, suite, result string) error
	getArtifactInfoFn func(component, build, suite string, fromTime, toTime time.Time) ([]models.ArtifactInfo, error)
	getFilePathByIDFn func(ctx context.Context, id int64) (string, error)
}

func (m *mockMetadataService) CreateArtifact(filePath string, fileSize int64, component, build, suite, result string) error {
	if m.createArtifactFn != nil {
		return m.createArtifactFn(filePath, fileSize, component, build, suite, result)
	}
	return nil
}
func (m *mockMetadataService) GetArtifactInfo(component, build, suite string, fromTime, toTime time.Time) ([]models.ArtifactInfo, error) {
	if m.getArtifactInfoFn != nil {
		return m.getArtifactInfoFn(component, build, suite, fromTime, toTime)
	}
	return nil, nil
}
func (m *mockMetadataService) GetFilePathByID(ctx context.Context, id int64) (string, error) {
	if m.getFilePathByIDFn != nil {
		return m.getFilePathByIDFn(ctx, id)
	}
	return "", nil
}

// -------------------------------------------
// Storage
// -------------------------------------------
type mockStorageService struct {
	saveFn func(file multipart.File, header *multipart.FileHeader) (string, error)
	openFn func(path string) (io.ReadCloser, error)
}

func (m *mockStorageService) Save(file multipart.File, header *multipart.FileHeader) (string, error) {
	if m.saveFn != nil {
		return m.saveFn(file, header)
	}
	return "", nil
}
func (m *mockStorageService) Open(path string) (io.ReadCloser, error) {
	if m.openFn != nil {
		return m.openFn(path)
	}
	return io.NopCloser(bytes.NewBufferString("test")), nil
}

// -------------------------------------------
// Analytics
// -------------------------------------------
type mockAnalyticsService struct {
	artifactsPerDayFn func(ctx context.Context) ([]analytics.DayCount, error)
	statusDistFn      func(ctx context.Context) ([]analytics.StatusCount, error)
}

func (m *mockAnalyticsService) ArtifactsPerDay(ctx context.Context) ([]analytics.DayCount, error) {
	if m.artifactsPerDayFn != nil {
		return m.artifactsPerDayFn(ctx)
	}
	return nil, nil
}
func (m *mockAnalyticsService) StatusDistribution(ctx context.Context) ([]analytics.StatusCount, error) {
	if m.statusDistFn != nil {
		return m.statusDistFn(ctx)
	}
	return nil, nil
}

// -------------------------------------------
// Health
// -------------------------------------------
type mockHealthService struct {
	aliveErr error
	readyErr error
}

func (m *mockHealthService) Alive(ctx context.Context) error {
	return m.aliveErr
}

func (m *mockHealthService) Ready(ctx context.Context) error {
	return m.readyErr
}

// -------------------------------------------
// Auth
// -------------------------------------------
type mockAuthService struct {
	validateFn func(username, password string) bool
}

func (m *mockAuthService) Validate(username, password string) bool {
	if m.validateFn != nil {
		return m.validateFn(username, password)
	}
	return false
}
