package repository

import (
	"context"
	"database/sql"
	"fmt"
	"time"

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
	return tx.QueryRowContext(ctx,
		`INSERT INTO test_artifacts (run_id, status_id, file_url, file_name, file_type_id, file_size, created_at, updated_at)
         VALUES ($1, $2, $3, $4, $5, $6, $7, $8)
         RETURNING id`,
		a.RunID, a.StatusID, a.FileURL, a.FileName, a.FileTypeID, a.FileSize, a.CreatedAt, a.CreatedAt,
	).Scan(&a.ID)
}

func (r *postgresRepo) GetArtifactInfo(component, build, suite string, fromTime, toTime time.Time) ([]models.ArtifactInfo, error) {
	query := `
        SELECT
			ta.id AS artifact_id,
            ta.file_url AS download_url,
            ta.file_name AS file_name,
            ta.file_size AS file_size,
            c.name AS component,
            b.name AS build,
            ts.name AS suite,
            ta.created_at AS upload_time,
			rs.name AS result
        FROM test_artifacts ta
        JOIN test_runs tr ON ta.run_id = tr.id
        JOIN builds b ON tr.build_id = b.id
        JOIN components c ON b.component_id = c.id
        JOIN test_suites ts ON tr.suite_id = ts.id
		JOIN result_statuses rs ON ta.status_id = rs.id
        WHERE 1=1
    `

	var args []interface{}
	argCounter := 1

	if component != "" {
		query += fmt.Sprintf(" AND c.name = $%d", argCounter)
		args = append(args, component)
		argCounter++
	}
	if build != "" {
		query += fmt.Sprintf(" AND b.name = $%d", argCounter)
		args = append(args, build)
		argCounter++
	}
	if suite != "" {
		query += fmt.Sprintf(" AND ts.name = $%d", argCounter)
		args = append(args, suite)
		argCounter++
	}
	if !fromTime.IsZero() {
		query += fmt.Sprintf(" AND ta.created_at >= $%d", argCounter)
		args = append(args, fromTime)
		argCounter++
	}
	if !toTime.IsZero() {
		query += fmt.Sprintf(" AND ta.created_at <= $%d", argCounter)
		args = append(args, toTime)
		argCounter++
	}

	query += " ORDER BY ta.created_at DESC;"

	rows, err := r.db.QueryContext(context.Background(), query, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var result []models.ArtifactInfo
	for rows.Next() {
		var a models.ArtifactInfo
		err := rows.Scan(
			&a.ID,
			&a.DownloadURL,
			&a.FileName,
			&a.FileSize,
			&a.Component,
			&a.Build,
			&a.Suite,
			&a.UploadTime,
			&a.Result,
		)
		if err != nil {
			return nil, err
		}
		result = append(result, a)
	}
	return result, rows.Err()
}

func (r *postgresRepo) GetFilePathByID(ctx context.Context, id int64) (string, error) {
	var filePath string
	err := r.db.QueryRowContext(ctx,
		`SELECT file_name FROM test_artifacts WHERE id = $1`, id,
	).Scan(&filePath)
	if err != nil {
		return "", fmt.Errorf("cannot get file path by id: %w", err)
	}
	return filePath, nil
}

func (r *postgresRepo) UpdateArtifactFileURL(ctx context.Context, tx DBTX, id int, fileURL string) error {
	_, err := tx.ExecContext(ctx,
		`UPDATE test_artifacts SET file_url = $1 WHERE id = $2`,
		fileURL, id)
	return err
}
