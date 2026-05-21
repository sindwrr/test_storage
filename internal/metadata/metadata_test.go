package metadata

import (
	"context"
	"errors"
	"fmt"
	"testing"
	"time"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/sindwrr/test_storage/internal/metadata/repository"
	"github.com/sindwrr/test_storage/internal/models"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
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

func setupServiceWithDB(t *testing.T) (*metadataService, *mockMetadataRepo, sqlmock.Sqlmock, func()) {
	t.Helper()

	db, mock, err := sqlmock.New()
	require.NoError(t, err)

	repo := &mockMetadataRepo{}
	svc := &metadataService{
		db:   db,
		repo: repo,
	}

	cleanup := func() { db.Close() }
	return svc, repo, mock, cleanup
}

func TestCreateArtifact_Success(t *testing.T) {
	svc, repo, mock, cleanup := setupServiceWithDB(t)
	defer cleanup()

	repo.GetOrCreateComponentFn = func(ctx context.Context, tx repository.DBTX, name string) (int, error) {
		assert.Equal(t, "core", name)
		return 1, nil
	}
	repo.GetOrCreateBuildFn = func(ctx context.Context, tx repository.DBTX, componentID int, name string) (int, error) {
		assert.Equal(t, 1, componentID)
		assert.Equal(t, "v1.0", name)
		return 2, nil
	}
	repo.GetOrCreateSuiteFn = func(ctx context.Context, tx repository.DBTX, name string) (int, error) {
		assert.Equal(t, "smoke", name)
		return 3, nil
	}
	repo.GetOrCreateStatusFn = func(ctx context.Context, tx repository.DBTX, tableName, statusName string) (int, error) {
		if tableName == "run_statuses" {
			assert.Equal(t, "Finished", statusName)
			return 4, nil
		}
		if tableName == "result_statuses" {
			assert.Equal(t, "Passed", statusName)
			return 5, nil
		}
		return 0, fmt.Errorf("unexpected table")
	}
	repo.CreateTestRunFn = func(ctx context.Context, tx repository.DBTX, run *models.TestRun) (int, error) {
		assert.Equal(t, 2, run.BuildID)
		assert.Equal(t, 3, run.SuiteID)
		assert.Equal(t, 4, run.StatusID)
		return 10, nil
	}
	repo.CreateTestArtifactFn = func(ctx context.Context, tx repository.DBTX, artifact *models.TestArtifact) error {
		assert.Equal(t, 10, artifact.RunID)
		assert.Equal(t, 5, artifact.StatusID)
		assert.Equal(t, "/tmp/test.log", artifact.FileName)
		artifact.ID = 100
		return nil
	}
	repo.UpdateArtifactFileURLFn = func(ctx context.Context, tx repository.DBTX, id int, fileURL string) error {
		assert.Equal(t, 100, id)
		assert.Equal(t, "/artifact/download/100", fileURL)
		return nil
	}

	mock.ExpectBegin()
	mock.ExpectCommit()

	err := svc.CreateArtifact("/tmp/test.log", 1024, "core", "v1.0", "smoke")
	assert.NoError(t, err)
	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestCreateArtifact_BeginTxError(t *testing.T) {
	svc, _, mock, cleanup := setupServiceWithDB(t)
	defer cleanup()

	mock.ExpectBegin().WillReturnError(errors.New("connection refused"))

	err := svc.CreateArtifact("/tmp/test.log", 1024, "core", "v1", "smoke")
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "begin tx")
	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestCreateArtifact_ComponentError(t *testing.T) {
	svc, repo, mock, cleanup := setupServiceWithDB(t)
	defer cleanup()

	repo.GetOrCreateComponentFn = func(ctx context.Context, tx repository.DBTX, name string) (int, error) {
		return 0, errors.New("component error")
	}

	mock.ExpectBegin()
	mock.ExpectRollback()

	err := svc.CreateArtifact("/tmp/test.log", 1024, "core", "v1", "smoke")
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "component")
	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestCreateArtifact_CommitError(t *testing.T) {
	svc, repo, mock, cleanup := setupServiceWithDB(t)
	defer cleanup()

	repo.GetOrCreateComponentFn = func(ctx context.Context, tx repository.DBTX, name string) (int, error) { return 1, nil }
	repo.GetOrCreateBuildFn = func(ctx context.Context, tx repository.DBTX, componentID int, name string) (int, error) {
		return 2, nil
	}
	repo.GetOrCreateSuiteFn = func(ctx context.Context, tx repository.DBTX, name string) (int, error) { return 3, nil }
	repo.GetOrCreateStatusFn = func(ctx context.Context, tx repository.DBTX, tableName, statusName string) (int, error) {
		return 1, nil
	}
	repo.CreateTestRunFn = func(ctx context.Context, tx repository.DBTX, run *models.TestRun) (int, error) { return 10, nil }
	repo.CreateTestArtifactFn = func(ctx context.Context, tx repository.DBTX, artifact *models.TestArtifact) error {
		artifact.ID = 100
		return nil
	}
	repo.UpdateArtifactFileURLFn = func(ctx context.Context, tx repository.DBTX, id int, fileURL string) error { return nil }

	mock.ExpectBegin()
	mock.ExpectCommit().WillReturnError(errors.New("commit error"))

	err := svc.CreateArtifact("/tmp/test.log", 1024, "core", "v1", "smoke")
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "commit")
	assert.NoError(t, mock.ExpectationsWereMet())
}
