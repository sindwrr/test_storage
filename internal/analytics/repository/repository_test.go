package repository

import (
	"context"
	"regexp"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/stretchr/testify/assert"
)

func TestArtifactsPerDay_Success(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("failed to create sqlmock: %v", err)
	}
	defer db.Close()

	repo := &postgresRepo{db: db}

	rows := sqlmock.NewRows([]string{"day", "count"}).
		AddRow("2026-05-10", 5).
		AddRow("2026-05-11", 3)

	mock.ExpectQuery(regexp.QuoteMeta(
		`SELECT date(created_at) AS day, count(*)
         FROM test_artifacts
         GROUP BY day
         ORDER BY day`)).
		WillReturnRows(rows)

	got, err := repo.ArtifactsPerDay(context.Background())

	assert.NoError(t, err)
	assert.Len(t, got, 2)
	assert.Equal(t, "2026-05-10", got[0].Date)
	assert.Equal(t, 5, got[0].Count)
	assert.Equal(t, "2026-05-11", got[1].Date)
	assert.Equal(t, 3, got[1].Count)

	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestStatusDistribution_Success(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("failed to create sqlmock: %v", err)
	}
	defer db.Close()

	repo := &postgresRepo{db: db}

	rows := sqlmock.NewRows([]string{"name", "count"}).
		AddRow("passed", 10).
		AddRow("failed", 2)

	mock.ExpectQuery(regexp.QuoteMeta(
		`SELECT rs.name, count(*)
		FROM test_artifacts ta
		JOIN result_statuses rs ON ta.status_id = rs.id
		GROUP BY rs.name
		ORDER BY rs.name`)).
		WillReturnRows(rows)

	got, err := repo.StatusDistribution(context.Background())

	assert.NoError(t, err)
	assert.Len(t, got, 2)
	assert.Equal(t, "passed", got[0].Status)
	assert.Equal(t, 10, got[0].Count)
	assert.Equal(t, "failed", got[1].Status)
	assert.Equal(t, 2, got[1].Count)

	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestArtifactsPerDay_QueryError(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("failed to create sqlmock: %v", err)
	}
	defer db.Close()

	repo := &postgresRepo{db: db}

	mock.ExpectQuery(regexp.QuoteMeta(
		`SELECT date(created_at) AS day, count(*)
         FROM test_artifacts
         GROUP BY day
         ORDER BY day`)).
		WillReturnError(sqlmock.ErrCancelled)

	_, err = repo.ArtifactsPerDay(context.Background())

	assert.Error(t, err)
	assert.NoError(t, mock.ExpectationsWereMet())
}
