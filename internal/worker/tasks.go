package worker

import (
	"context"
	"database/sql"
	"fmt"
	"log"
	"os"
)

type IntegrityCheckTask struct {
	DB       *sql.DB
	BasePath string
}

func (t IntegrityCheckTask) Name() string { return "integrity-check" }

func (t IntegrityCheckTask) Execute(ctx context.Context) error {
	if err := t.ensureMissingStatus(ctx); err != nil {
		return fmt.Errorf("ensure missing status: %w", err)
	}

	rows, err := t.DB.QueryContext(ctx,
		`SELECT id, file_url FROM test_artifacts WHERE status_id != (
            SELECT id FROM result_statuses WHERE name = 'Missing')`)
	if err != nil {
		return fmt.Errorf("query artifacts: %w", err)
	}
	defer rows.Close()

	var (
		artifactID int
		fileURL    string
	)

	for rows.Next() {
		if err := rows.Scan(&artifactID, &fileURL); err != nil {
			return fmt.Errorf("scan row: %w", err)
		}

		if _, err := os.Stat(fileURL); os.IsNotExist(err) {
			_, err := t.DB.ExecContext(ctx,
				`UPDATE test_artifacts 
                 SET status_id = (SELECT id FROM result_statuses WHERE name = 'Missing'),
                     updated_at = NOW()
                 WHERE id = $1`, artifactID)
			if err != nil {
				log.Printf("integrity check: failed to update artifact %d: %v\n", artifactID, err)
			} else {
				log.Printf("integrity check: marked artifact %d as missing (file %s not found)\n", artifactID, fileURL)
			}
		}
	}

	return rows.Err()
}

func (t IntegrityCheckTask) ensureMissingStatus(ctx context.Context) error {
	var id int
	err := t.DB.QueryRowContext(ctx,
		`INSERT INTO result_statuses (name) VALUES ('Missing')
         ON CONFLICT (name) DO NOTHING
         RETURNING id`).Scan(&id)
	if err != nil && err != sql.ErrNoRows {
		return err
	}
	return nil
}
