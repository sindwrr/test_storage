package repository

import (
	"context"
	"database/sql"

	"github.com/sindwrr/test_storage/internal/models/analytics"
)

type postgresRepo struct {
	db *sql.DB
}

func NewAnalyticsRepo(db *sql.DB) AnalyticsRepository {
	return &postgresRepo{db: db}
}

func (r *postgresRepo) ArtifactsPerDay(ctx context.Context) ([]analytics.DayCount, error) {
	rows, err := r.db.QueryContext(ctx,
		`SELECT date(created_at) AS day, count(*)
		 FROM test_artifacts
		 GROUP BY day
		 ORDER BY day`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var result []analytics.DayCount
	for rows.Next() {
		var dc analytics.DayCount
		if err := rows.Scan(&dc.Date, &dc.Count); err != nil {
			return nil, err
		}
		result = append(result, dc)
	}
	return result, rows.Err()
}

func (r *postgresRepo) StatusDistribution(ctx context.Context) ([]analytics.StatusCount, error) {
	rows, err := r.db.QueryContext(ctx,
		`SELECT rs.name, count(*)
         FROM test_runs tr
         JOIN run_statuses rs ON tr.status_id = rs.id
         GROUP BY rs.name
         ORDER BY rs.name`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var result []analytics.StatusCount
	for rows.Next() {
		var sc analytics.StatusCount
		if err := rows.Scan(&sc.Status, &sc.Count); err != nil {
			return nil, err
		}
		result = append(result, sc)
	}
	return result, rows.Err()
}
