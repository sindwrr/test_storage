package repository

import (
	"context"
	"database/sql"
	"time"

	"github.com/sindwrr/test_storage/internal/models"
)

type DBTX interface {
	ExecContext(ctx context.Context, query string, args ...interface{}) (sql.Result, error)
	QueryRowContext(ctx context.Context, query string, args ...interface{}) *sql.Row
	QueryContext(ctx context.Context, query string, args ...interface{}) (*sql.Rows, error)
}

type MetadataRepository interface {
	GetOrCreateComponent(ctx context.Context, tx DBTX, name string) (int, error)
	GetOrCreateBuild(ctx context.Context, tx DBTX, componentID int, name string) (int, error)
	GetOrCreateSuite(ctx context.Context, tx DBTX, name string) (int, error)
	GetOrCreateStatus(ctx context.Context, tx DBTX, tableName, statusName string) (int, error)
	CreateTestRun(ctx context.Context, tx DBTX, run *models.TestRun) (int, error)
	CreateTestArtifact(ctx context.Context, tx DBTX, artifact *models.TestArtifact) error
	GetArtifactInfo(component string, build string, suite string, fromTime time.Time, toTime time.Time) ([]models.ArtifactInfo, error)
	GetFilePathByID(ctx context.Context, id int64) (string, error)
	UpdateArtifactFileURL(ctx context.Context, tx DBTX, id int, fileURL string) error
}
