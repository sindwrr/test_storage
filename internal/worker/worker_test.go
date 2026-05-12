package worker

import (
	"context"
	"database/sql"
	"errors"
	"os"
	"path/filepath"
	"sync/atomic"
	"testing"
	"time"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/stretchr/testify/assert"
)

type mockTask struct {
	name    string
	execute func(ctx context.Context) error
}

func (m *mockTask) Name() string                      { return m.name }
func (m *mockTask) Execute(ctx context.Context) error { return m.execute(ctx) }

// ------------------------------------------------------------
// Pool tests
// ------------------------------------------------------------

func TestPool_SubmitAndExecute(t *testing.T) {
	var executed int32
	task := &mockTask{
		name: "test-task",
		execute: func(ctx context.Context) error {
			atomic.AddInt32(&executed, 1)
			return nil
		},
	}

	pool := NewPool(2)
	pool.Start()
	defer pool.Shutdown()

	pool.Submit(task)
	time.Sleep(50 * time.Millisecond)

	if atomic.LoadInt32(&executed) != 1 {
		t.Errorf("expected task to be executed once, got %d", executed)
	}
}

func TestPool_Shutdown(t *testing.T) {
	var executed int32
	task := &mockTask{
		name: "drop-task",
		execute: func(ctx context.Context) error {
			atomic.AddInt32(&executed, 1)
			return nil
		},
	}

	pool := NewPool(1)
	pool.Start()
	pool.Shutdown()
	pool.Submit(task)

	time.Sleep(50 * time.Millisecond)
	if atomic.LoadInt32(&executed) != 0 {
		t.Errorf("task should not be executed after shutdown")
	}
}

// ------------------------------------------------------------
// IntegrityCheckTask tests
// ------------------------------------------------------------

func newMockDBAndTempDir(t *testing.T) (*sql.DB, sqlmock.Sqlmock, string) {
	t.Helper()
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("sqlmock.New: %v", err)
	}
	tmpDir := t.TempDir()
	return db, mock, tmpDir
}

func TestIntegrityCheckTask_FileExists(t *testing.T) {
	db, mock, basePath := newMockDBAndTempDir(t)
	defer db.Close()

	existingFile := filepath.Join(basePath, "existing.log")
	os.WriteFile(existingFile, []byte("data"), 0644)

	mock.ExpectQuery(
		`(?s)INSERT INTO result_statuses \(name\) VALUES \('Missing'\)` +
			`.*ON CONFLICT \(name\) DO NOTHING.*RETURNING id`,
	).WillReturnRows(sqlmock.NewRows([]string{"id"}).AddRow(100))

	rows := sqlmock.NewRows([]string{"id", "file_url"}).
		AddRow(1, existingFile)
	mock.ExpectQuery(
		`(?s)SELECT id, file_url FROM test_artifacts WHERE status_id != \(` +
			`.*SELECT id FROM result_statuses WHERE name = 'Missing'\)`,
	).WillReturnRows(rows)

	task := IntegrityCheckTask{DB: db, BasePath: basePath}
	err := task.Execute(context.Background())

	assert.NoError(t, err)
	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestIntegrityCheckTask_FileMissing(t *testing.T) {
	db, mock, basePath := newMockDBAndTempDir(t)
	defer db.Close()

	missingFile := filepath.Join(basePath, "gone.log")

	mock.ExpectQuery(
		`(?s)INSERT INTO result_statuses \(name\) VALUES \('Missing'\)` +
			`.*ON CONFLICT \(name\) DO NOTHING.*RETURNING id`,
	).WillReturnRows(sqlmock.NewRows([]string{"id"}).AddRow(100))

	rows := sqlmock.NewRows([]string{"id", "file_url"}).
		AddRow(2, missingFile)
	mock.ExpectQuery(
		`(?s)SELECT id, file_url FROM test_artifacts WHERE status_id != \(` +
			`.*SELECT id FROM result_statuses WHERE name = 'Missing'\)`,
	).WillReturnRows(rows)

	mock.ExpectExec(
		`(?s)UPDATE test_artifacts.*SET status_id = \(SELECT id FROM result_statuses WHERE name = 'Missing'\)` +
			`.*updated_at = NOW\(\).*WHERE id = \$1`,
	).WithArgs(2).WillReturnResult(sqlmock.NewResult(0, 1))

	task := IntegrityCheckTask{DB: db, BasePath: basePath}
	err := task.Execute(context.Background())

	assert.NoError(t, err)
	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestIntegrityCheckTask_EnsureMissingStatus_AlreadyExists(t *testing.T) {
	db, mock, basePath := newMockDBAndTempDir(t)
	defer db.Close()

	mock.ExpectQuery(
		`(?s)INSERT INTO result_statuses \(name\) VALUES \('Missing'\)` +
			`.*ON CONFLICT \(name\) DO NOTHING.*RETURNING id`,
	).WillReturnError(sql.ErrNoRows)

	task := IntegrityCheckTask{DB: db, BasePath: basePath}
	err := task.ensureMissingStatus(context.Background())
	assert.NoError(t, err)
	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestIntegrityCheckTask_QueryError(t *testing.T) {
	db, mock, basePath := newMockDBAndTempDir(t)
	defer db.Close()

	mock.ExpectQuery(
		`(?s)INSERT INTO result_statuses \(name\) VALUES \('Missing'\)` +
			`.*ON CONFLICT \(name\) DO NOTHING.*RETURNING id`,
	).WillReturnRows(sqlmock.NewRows([]string{"id"}).AddRow(100))

	mock.ExpectQuery(
		`(?s)SELECT id, file_url FROM test_artifacts`,
	).WillReturnError(errors.New("db timeout"))

	task := IntegrityCheckTask{DB: db, BasePath: basePath}
	err := task.Execute(context.Background())
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "query artifacts")
	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestIntegrityCheckTask_Name(t *testing.T) {
	task := IntegrityCheckTask{DB: nil, BasePath: "/tmp"}
	if name := task.Name(); name != "integrity-check" {
		t.Errorf("expected 'integrity-check', got %q", name)
	}
}
