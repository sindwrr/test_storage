package metadata

import (
	"context"
	"database/sql"
	"fmt"
	"path/filepath"
	"time"

	"github.com/sindwrr/test_storage/internal/metadata/repository"
	"github.com/sindwrr/test_storage/internal/models"
)

type metadataService struct {
	db   *sql.DB
	repo repository.MetadataRepository
}

func NewMetadataService(db *sql.DB) MetadataService {
	return &metadataService{
		db:   db,
		repo: repository.NewPostgresRepo(db),
	}
}

func (s *metadataService) CreateArtifact(filePath string, fileSize int64, component string, build string, suite string) error {
	ext := filepath.Ext(filePath)
	fileTypeID := 1
	switch ext {
	case ".txt", ".log", ".json", ".xml", ".csv":
		fileTypeID = 1
	case ".png", ".jpg", ".jpeg", ".bmp", ".gif":
		fileTypeID = 2
	case ".mp4", ".avi", ".mov", ".mkv":
		fileTypeID = 3
	case ".pdf":
		fileTypeID = 4
	}

	now := time.Now().UTC()
	ctx := context.Background()

	tx, err := s.db.BeginTx(ctx, nil)
	if err != nil {
		return fmt.Errorf("begin tx: %w", err)
	}
	defer tx.Rollback()

	componentID, err := s.repo.GetOrCreateComponent(ctx, tx, component)
	if err != nil {
		return fmt.Errorf("component: %w", err)
	}

	buildID, err := s.repo.GetOrCreateBuild(ctx, tx, componentID, build)
	if err != nil {
		return fmt.Errorf("build: %w", err)
	}

	suiteID, err := s.repo.GetOrCreateSuite(ctx, tx, suite)
	if err != nil {
		return fmt.Errorf("suite: %w", err)
	}

	runStatusID, err := s.repo.GetOrCreateStatus(ctx, tx, "run_statuses", "Finished")
	if err != nil {
		return fmt.Errorf("run status: %w", err)
	}

	resultStatusID, err := s.repo.GetOrCreateStatus(ctx, tx, "result_statuses", "Passed")
	if err != nil {
		return fmt.Errorf("result status: %w", err)
	}

	runData := &models.TestRun{
		BuildID:    buildID,
		SuiteID:    suiteID,
		StatusID:   runStatusID,
		StartedAt:  now,
		FinishedAt: now,
	}
	runID, err := s.repo.CreateTestRun(ctx, tx, runData)
	if err != nil {
		return fmt.Errorf("test run: %w", err)
	}

	artifactData := &models.TestArtifact{
		RunID:      runID,
		StatusID:   resultStatusID,
		FileURL:    filePath,
		FileTypeID: fileTypeID,
		FileSize:   fileSize,
		CreatedAt:  now,
	}
	if err := s.repo.CreateTestArtifact(ctx, tx, artifactData); err != nil {
		return fmt.Errorf("artifact: %w", err)
	}

	return tx.Commit()
}

func (s *metadataService) GetArtifactInfo(component, build, suite string, fromTime, toTime time.Time) ([]models.ArtifactInfo, error) {
	return s.repo.GetArtifactInfo(component, build, suite, fromTime, toTime)
}
