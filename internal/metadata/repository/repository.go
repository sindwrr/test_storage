package repository

import (
	"context"
	"database/sql"

	"github.com/sindwrr/test_storage/internal/models"
)

type postgresRepo struct {
	db *sql.DB
}

func NewPostgresRepo(db *sql.DB) MetadataRepository {
	return &postgresRepo{db: db}
}

func (r *postgresRepo) GetOrCreateComponent(ctx context.Context, tx DBTX, name string) (int, error) {
	var id int
	err := tx.QueryRowContext(ctx,
		`INSERT INTO components (name) VALUES ($1)
         ON CONFLICT (name) DO UPDATE SET name = EXCLUDED.name
         RETURNING id`, name).Scan(&id)
	return id, err
}

func (r *postgresRepo) GetOrCreateBuild(ctx context.Context, tx DBTX, componentID int, name string) (int, error) {
	var id int
	err := tx.QueryRowContext(ctx,
		`INSERT INTO builds (component_id, name) VALUES ($1, $2)
         RETURNING id`, componentID, name).Scan(&id)
	return id, err
}

func (r *postgresRepo) GetOrCreateSuite(ctx context.Context, tx DBTX, name string) (int, error) {
	var id int
	err := tx.QueryRowContext(ctx,
		`INSERT INTO test_suites (name) VALUES ($1)
         ON CONFLICT (name) DO UPDATE SET name = EXCLUDED.name
         RETURNING id`, name).Scan(&id)
	return id, err
}

func (r *postgresRepo) GetOrCreateStatus(ctx context.Context, tx DBTX, tableName string, statusName string) (int, error) {
	var id int
	query := "SELECT id FROM " + tableName + " WHERE name = $1"
	err := tx.QueryRowContext(ctx, query, statusName).Scan(&id)
	if err == nil {
		return id, nil
	}
	if err != sql.ErrNoRows {
		return 0, err
	}
	insertQuery := "INSERT INTO " + tableName + " (name) VALUES ($1) RETURNING id"
	err = tx.QueryRowContext(ctx, insertQuery, statusName).Scan(&id)
	return id, err
}

func (r *postgresRepo) CreateTestRun(ctx context.Context, tx DBTX, run *models.TestRun) (int, error) {
	var id int
	err := tx.QueryRowContext(ctx,
		`INSERT INTO test_runs (build_id, suite_id, status_id, started_at, finished_at, updated_at)
         VALUES ($1, $2, $3, $4, $5, $6)
         RETURNING id`,
		run.BuildID, run.SuiteID, run.StatusID, run.StartedAt, run.FinishedAt, run.StartedAt).Scan(&id)
	return id, err
}

func (r *postgresRepo) CreateTestArtifact(ctx context.Context, tx DBTX, a *models.TestArtifact) error {
	_, err := tx.ExecContext(ctx,
		`INSERT INTO test_artifacts (run_id, status_id, file_url, file_type_id, file_size, created_at, updated_at)
         VALUES ($1, $2, $3, $4, $5, $6, $7)`,
		a.RunID, a.StatusID, a.FileURL, a.FileTypeID, a.FileSize, a.CreatedAt, a.CreatedAt)
	return err
}
