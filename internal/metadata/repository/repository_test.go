package repository

import (
	"context"
	"database/sql"
	"errors"
	"regexp"
	"testing"
	"time"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/sindwrr/test_storage/internal/models"
	"github.com/stretchr/testify/assert"
)

func newMockDB(t *testing.T) (*sql.DB, sqlmock.Sqlmock) {
	t.Helper()
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("sqlmock.New: %v", err)
	}
	return db, mock
}

// -------------------------------------------
// GetOrCreateComponent
// -------------------------------------------
func TestGetOrCreateComponent_Success(t *testing.T) {
	db, mock := newMockDB(t)
	defer db.Close()
	repo := NewPostgresRepo(db).(*postgresRepo)

	mock.ExpectQuery(regexp.QuoteMeta(
		`INSERT INTO components (name) VALUES ($1)
         ON CONFLICT (name) DO UPDATE SET name = EXCLUDED.name
         RETURNING id`)).
		WithArgs("core").
		WillReturnRows(sqlmock.NewRows([]string{"id"}).AddRow(1))

	id, err := repo.GetOrCreateComponent(context.Background(), db, "core")
	assert.NoError(t, err)
	assert.Equal(t, 1, id)
	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestGetOrCreateComponent_Error(t *testing.T) {
	db, mock := newMockDB(t)
	defer db.Close()
	repo := NewPostgresRepo(db).(*postgresRepo)

	mock.ExpectQuery(regexp.QuoteMeta(
		`INSERT INTO components (name) VALUES ($1)
         ON CONFLICT (name) DO UPDATE SET name = EXCLUDED.name
         RETURNING id`)).
		WithArgs("core").
		WillReturnError(errors.New("connection lost"))

	_, err := repo.GetOrCreateComponent(context.Background(), db, "core")
	assert.Error(t, err)
	assert.NoError(t, mock.ExpectationsWereMet())
}

// -------------------------------------------
// GetOrCreateBuild
// -------------------------------------------
func TestGetOrCreateBuild_Success(t *testing.T) {
	db, mock := newMockDB(t)
	defer db.Close()
	repo := NewPostgresRepo(db).(*postgresRepo)

	mock.ExpectQuery(regexp.QuoteMeta(
		`INSERT INTO builds (component_id, name) VALUES ($1, $2)
         RETURNING id`)).
		WithArgs(1, "v1.0").
		WillReturnRows(sqlmock.NewRows([]string{"id"}).AddRow(10))

	id, err := repo.GetOrCreateBuild(context.Background(), db, 1, "v1.0")
	assert.NoError(t, err)
	assert.Equal(t, 10, id)
	assert.NoError(t, mock.ExpectationsWereMet())
}

// -------------------------------------------
// GetOrCreateSuite
// -------------------------------------------
func TestGetOrCreateSuite_Success(t *testing.T) {
	db, mock := newMockDB(t)
	defer db.Close()
	repo := NewPostgresRepo(db).(*postgresRepo)

	mock.ExpectQuery(regexp.QuoteMeta(
		`INSERT INTO test_suites (name) VALUES ($1)
         ON CONFLICT (name) DO UPDATE SET name = EXCLUDED.name
         RETURNING id`)).
		WithArgs("smoke").
		WillReturnRows(sqlmock.NewRows([]string{"id"}).AddRow(5))

	id, err := repo.GetOrCreateSuite(context.Background(), db, "smoke")
	assert.NoError(t, err)
	assert.Equal(t, 5, id)
	assert.NoError(t, mock.ExpectationsWereMet())
}

// -------------------------------------------
// GetOrCreateStatus
// -------------------------------------------
func TestGetOrCreateStatus_Existing(t *testing.T) {
	db, mock := newMockDB(t)
	defer db.Close()
	repo := NewPostgresRepo(db).(*postgresRepo)

	mock.ExpectQuery(regexp.QuoteMeta("SELECT id FROM run_statuses WHERE name = $1")).
		WithArgs("passed").
		WillReturnRows(sqlmock.NewRows([]string{"id"}).AddRow(2))

	id, err := repo.GetOrCreateStatus(context.Background(), db, "run_statuses", "passed")
	assert.NoError(t, err)
	assert.Equal(t, 2, id)
	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestGetOrCreateStatus_NotFound_Inserts(t *testing.T) {
	db, mock := newMockDB(t)
	defer db.Close()
	repo := NewPostgresRepo(db).(*postgresRepo)

	mock.ExpectQuery(regexp.QuoteMeta("SELECT id FROM run_statuses WHERE name = $1")).
		WithArgs("skipped").
		WillReturnError(sql.ErrNoRows)

	mock.ExpectQuery(regexp.QuoteMeta("INSERT INTO run_statuses (name) VALUES ($1) RETURNING id")).
		WithArgs("skipped").
		WillReturnRows(sqlmock.NewRows([]string{"id"}).AddRow(7))

	id, err := repo.GetOrCreateStatus(context.Background(), db, "run_statuses", "skipped")
	assert.NoError(t, err)
	assert.Equal(t, 7, id)
	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestGetOrCreateStatus_QueryError(t *testing.T) {
	db, mock := newMockDB(t)
	defer db.Close()
	repo := NewPostgresRepo(db).(*postgresRepo)

	mock.ExpectQuery(regexp.QuoteMeta("SELECT id FROM run_statuses WHERE name = $1")).
		WithArgs("passed").
		WillReturnError(errors.New("timeout"))

	_, err := repo.GetOrCreateStatus(context.Background(), db, "run_statuses", "passed")
	assert.Error(t, err)
	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestCreateTestRun_Success(t *testing.T) {
	db, mock := newMockDB(t)
	defer db.Close()
	repo := NewPostgresRepo(db).(*postgresRepo)

	run := &models.TestRun{
		BuildID:    10,
		SuiteID:    20,
		StatusID:   1,
		StartedAt:  time.Date(2026, 5, 11, 12, 0, 0, 0, time.UTC),
		FinishedAt: time.Date(2026, 5, 11, 12, 5, 0, 0, time.UTC),
	}

	mock.ExpectQuery(regexp.QuoteMeta(
		`INSERT INTO test_runs (build_id, suite_id, status_id, started_at, finished_at, updated_at)
         VALUES ($1, $2, $3, $4, $5, $6)
         RETURNING id`)).
		WithArgs(run.BuildID, run.SuiteID, run.StatusID, run.StartedAt, run.FinishedAt, run.StartedAt).
		WillReturnRows(sqlmock.NewRows([]string{"id"}).AddRow(100))

	id, err := repo.CreateTestRun(context.Background(), db, run)
	assert.NoError(t, err)
	assert.Equal(t, 100, id)
	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestCreateTestRun_Error(t *testing.T) {
	db, mock := newMockDB(t)
	defer db.Close()
	repo := NewPostgresRepo(db).(*postgresRepo)

	run := &models.TestRun{BuildID: 1}
	mock.ExpectQuery(`INSERT INTO test_runs`).
		WillReturnError(errors.New("duplicate key"))

	_, err := repo.CreateTestRun(context.Background(), db, run)
	assert.Error(t, err)
	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestCreateTestArtifact_Success(t *testing.T) {
	db, mock := newMockDB(t)
	defer db.Close()
	repo := NewPostgresRepo(db).(*postgresRepo)

	artifact := &models.TestArtifact{
		RunID:      1,
		StatusID:   2,
		FileURL:    "/download/1",
		FileName:   "log.txt",
		FileTypeID: 1,
		FileSize:   1024,
		CreatedAt:  time.Date(2026, 5, 11, 12, 0, 0, 0, time.UTC),
	}

	mock.ExpectQuery(regexp.QuoteMeta(
		`INSERT INTO test_artifacts (run_id, status_id, file_url, file_name, file_type_id, file_size, created_at, updated_at)
         VALUES ($1, $2, $3, $4, $5, $6, $7, $8)
         RETURNING id`)).
		WithArgs(
			artifact.RunID, artifact.StatusID, artifact.FileURL, artifact.FileName,
			artifact.FileTypeID, artifact.FileSize, artifact.CreatedAt, artifact.CreatedAt,
		).
		WillReturnRows(sqlmock.NewRows([]string{"id"}).AddRow(55))

	err := repo.CreateTestArtifact(context.Background(), db, artifact)
	assert.NoError(t, err)
	assert.Equal(t, 55, artifact.ID)
	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestGetArtifactInfo_NoFilters(t *testing.T) {
	db, mock := newMockDB(t)
	defer db.Close()
	repo := NewPostgresRepo(db).(*postgresRepo)

	baseQuery := `SELECT ta.id AS artifact_id, ta.file_url AS download_url, ` +
		`ta.file_name AS file_name, ta.file_size AS file_size, c.name AS component, ` +
		`b.name AS build, ts.name AS suite, ta.created_at AS upload_time, rs.name AS result ` +
		`FROM test_artifacts ta ` +
		`JOIN test_runs tr ON ta.run_id = tr.id ` +
		`JOIN builds b ON tr.build_id = b.id ` +
		`JOIN components c ON b.component_id = c.id ` +
		`JOIN test_suites ts ON tr.suite_id = ts.id ` +
		`JOIN result_statuses rs ON ta.status_id = rs.id ` +
		`WHERE 1=1 ORDER BY ta.created_at DESC;`

	mock.ExpectQuery(baseQuery).
		WillReturnRows(sqlmock.NewRows([]string{
			"artifact_id", "download_url", "file_name", "file_size",
			"component", "build", "suite", "upload_time", "result",
		}).AddRow(1, "/dl/1", "a.txt", 500, "core", "v1", "smoke", time.Now(), "passed"))

	results, err := repo.GetArtifactInfo("", "", "", time.Time{}, time.Time{})
	assert.NoError(t, err)
	assert.Len(t, results, 1)
	assert.Equal(t, "a.txt", results[0].FileName)
	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestGetArtifactInfo_WithFilters(t *testing.T) {
	db, mock := newMockDB(t)
	defer db.Close()
	repo := NewPostgresRepo(db).(*postgresRepo)

	from := time.Date(2026, 1, 1, 0, 0, 0, 0, time.UTC)
	to := time.Date(2026, 5, 1, 0, 0, 0, 0, time.UTC)

	expectedQuery := `(?s)` +
		`SELECT.*ta\.id AS artifact_id.*` +
		`FROM test_artifacts ta.*` +
		`WHERE 1=1 AND c\.name = \$1 AND ta\.created_at >= \$2 AND ta\.created_at <= \$3.*` +
		`ORDER BY ta\.created_at DESC;`

	mock.ExpectQuery(expectedQuery).
		WithArgs("core", from, to).
		WillReturnRows(sqlmock.NewRows([]string{
			"artifact_id", "download_url", "file_name", "file_size",
			"component", "build", "suite", "upload_time", "result", // добавлено "result"
		}).AddRow(2, "/dl/2", "b.log", 1200, "core", "v2", "regression", time.Now(), "passed"))

	results, err := repo.GetArtifactInfo("core", "", "", from, to)
	assert.NoError(t, err)
	assert.Len(t, results, 1)
	assert.Equal(t, "b.log", results[0].FileName)
	assert.NoError(t, mock.ExpectationsWereMet())
}

// -------------------------------------------
// GetFilePathByID
// -------------------------------------------
func TestGetFilePathByID_Success(t *testing.T) {
	db, mock := newMockDB(t)
	defer db.Close()
	repo := NewPostgresRepo(db).(*postgresRepo)

	mock.ExpectQuery(regexp.QuoteMeta(
		`SELECT file_name FROM test_artifacts WHERE id = $1`)).
		WithArgs(int64(5)).
		WillReturnRows(sqlmock.NewRows([]string{"file_name"}).AddRow("/data/log.txt"))

	path, err := repo.GetFilePathByID(context.Background(), 5)
	assert.NoError(t, err)
	assert.Equal(t, "/data/log.txt", path)
	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestGetFilePathByID_NotFound(t *testing.T) {
	db, mock := newMockDB(t)
	defer db.Close()
	repo := NewPostgresRepo(db).(*postgresRepo)

	mock.ExpectQuery(regexp.QuoteMeta(
		`SELECT file_name FROM test_artifacts WHERE id = $1`)).
		WithArgs(int64(999)).
		WillReturnError(sql.ErrNoRows)

	_, err := repo.GetFilePathByID(context.Background(), 999)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "cannot get file path by id")
	assert.NoError(t, mock.ExpectationsWereMet())
}

// -------------------------------------------
// UpdateArtifactFileURL
// -------------------------------------------
func TestUpdateArtifactFileURL_Success(t *testing.T) {
	db, mock := newMockDB(t)
	defer db.Close()
	repo := NewPostgresRepo(db).(*postgresRepo)

	mock.ExpectExec(regexp.QuoteMeta(
		`UPDATE test_artifacts SET file_url = $1 WHERE id = $2`)).
		WithArgs("/new/url", 10).
		WillReturnResult(sqlmock.NewResult(0, 1))

	err := repo.UpdateArtifactFileURL(context.Background(), db, 10, "/new/url")
	assert.NoError(t, err)
	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestUpdateArtifactFileURL_Error(t *testing.T) {
	db, mock := newMockDB(t)
	defer db.Close()
	repo := NewPostgresRepo(db).(*postgresRepo)

	mock.ExpectExec(regexp.QuoteMeta(
		`UPDATE test_artifacts SET file_url = $1 WHERE id = $2`)).
		WillReturnError(errors.New("disk full"))

	err := repo.UpdateArtifactFileURL(context.Background(), db, 10, "/url")
	assert.Error(t, err)
	assert.NoError(t, mock.ExpectationsWereMet())
}
