package metadata

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/sindwrr/test_storage/internal/metadata/repository"
	"github.com/sindwrr/test_storage/internal/models"
	"github.com/stretchr/testify/assert"
)

type mockMetadataRepo struct {
	GetOrCreateComponentFn  func(ctx context.Context, tx repository.DBTX, name string) (int, error)
	GetOrCreateBuildFn      func(ctx context.Context, tx repository.DBTX, componentID int, name string) (int, error)
	GetOrCreateSuiteFn      func(ctx context.Context, tx repository.DBTX, name string) (int, error)
	GetOrCreateStatusFn     func(ctx context.Context, tx repository.DBTX, tableName, statusName string) (int, error)
	CreateTestRunFn         func(ctx context.Context, tx repository.DBTX, run *models.TestRun) (int, error)
	CreateTestArtifactFn    func(ctx context.Context, tx repository.DBTX, artifact *models.TestArtifact) error
	GetArtifactInfoFn       func(component, build, suite string, fromTime, toTime time.Time) ([]models.ArtifactInfo, error)
	GetFilePathByIDFn       func(ctx context.Context, id int64) (string, error)
	UpdateArtifactFileURLFn func(ctx context.Context, tx repository.DBTX, id int, fileURL string) error
}

func (m *mockMetadataRepo) GetOrCreateComponent(ctx context.Context, tx repository.DBTX, name string) (int, error) {
	if m.GetOrCreateComponentFn != nil {
		return m.GetOrCreateComponentFn(ctx, tx, name)
	}
	return 1, nil
}

func (m *mockMetadataRepo) GetOrCreateBuild(ctx context.Context, tx repository.DBTX, componentID int, name string) (int, error) {
	if m.GetOrCreateBuildFn != nil {
		return m.GetOrCreateBuildFn(ctx, tx, componentID, name)
	}
	return 1, nil
}

func (m *mockMetadataRepo) GetOrCreateSuite(ctx context.Context, tx repository.DBTX, name string) (int, error) {
	if m.GetOrCreateSuiteFn != nil {
		return m.GetOrCreateSuiteFn(ctx, tx, name)
	}
	return 1, nil
}

func (m *mockMetadataRepo) GetOrCreateStatus(ctx context.Context, tx repository.DBTX, tableName, statusName string) (int, error) {
	if m.GetOrCreateStatusFn != nil {
		return m.GetOrCreateStatusFn(ctx, tx, tableName, statusName)
	}
	return 1, nil
}

func (m *mockMetadataRepo) CreateTestRun(ctx context.Context, tx repository.DBTX, run *models.TestRun) (int, error) {
	if m.CreateTestRunFn != nil {
		return m.CreateTestRunFn(ctx, tx, run)
	}
	return 1, nil
}

func (m *mockMetadataRepo) CreateTestArtifact(ctx context.Context, tx repository.DBTX, artifact *models.TestArtifact) error {
	if m.CreateTestArtifactFn != nil {
		return m.CreateTestArtifactFn(ctx, tx, artifact)
	}
	artifact.ID = 42
	return nil
}

func (m *mockMetadataRepo) GetArtifactInfo(component, build, suite string, fromTime, toTime time.Time) ([]models.ArtifactInfo, error) {
	if m.GetArtifactInfoFn != nil {
		return m.GetArtifactInfoFn(component, build, suite, fromTime, toTime)
	}
	return nil, nil
}

func (m *mockMetadataRepo) GetFilePathByID(ctx context.Context, id int64) (string, error) {
	if m.GetFilePathByIDFn != nil {
		return m.GetFilePathByIDFn(ctx, id)
	}
	return "", nil
}

func (m *mockMetadataRepo) UpdateArtifactFileURL(ctx context.Context, tx repository.DBTX, id int, fileURL string) error {
	if m.UpdateArtifactFileURLFn != nil {
		return m.UpdateArtifactFileURLFn(ctx, tx, id, fileURL)
	}
	return nil
}

func TestGetArtifactInfo_Success(t *testing.T) {
	expected := []models.ArtifactInfo{
		{ID: 1, FileName: "log1.txt", Component: "core"},
		{ID: 2, FileName: "log2.txt", Component: "ui"},
	}
	repo := &mockMetadataRepo{
		GetArtifactInfoFn: func(component, build, suite string, fromTime, toTime time.Time) ([]models.ArtifactInfo, error) {
			return expected, nil
		},
	}
	svc := &metadataService{repo: repo}
	got, err := svc.GetArtifactInfo("core", "", "", time.Time{}, time.Time{})
	assert.NoError(t, err)
	assert.Equal(t, expected, got)
}

func TestGetArtifactInfo_Error(t *testing.T) {
	repo := &mockMetadataRepo{
		GetArtifactInfoFn: func(component, build, suite string, fromTime, toTime time.Time) ([]models.ArtifactInfo, error) {
			return nil, errors.New("db error")
		},
	}
	svc := &metadataService{repo: repo}
	_, err := svc.GetArtifactInfo("", "", "", time.Time{}, time.Time{})
	assert.Error(t, err)
}

func TestGetFilePathByID_Success(t *testing.T) {
	repo := &mockMetadataRepo{
		GetFilePathByIDFn: func(ctx context.Context, id int64) (string, error) {
			return "/artifacts/log.txt", nil
		},
	}
	svc := &metadataService{repo: repo}
	path, err := svc.GetFilePathByID(context.Background(), 1)
	assert.NoError(t, err)
	assert.Equal(t, "/artifacts/log.txt", path)
}

func TestGetFilePathByID_Error(t *testing.T) {
	repo := &mockMetadataRepo{
		GetFilePathByIDFn: func(ctx context.Context, id int64) (string, error) {
			return "", errors.New("not found")
		},
	}
	svc := &metadataService{repo: repo}
	_, err := svc.GetFilePathByID(context.Background(), 1)
	assert.Error(t, err)
}
