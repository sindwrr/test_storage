package metadata

import (
	"context"
	"time"

	"github.com/sindwrr/test_storage/internal/models"
)

type MetadataService interface {
	CreateArtifact(filePath string, fileSize int64, component string, build string, suite string) error
	GetArtifactInfo(component string, build string, suite string, fromTime time.Time, toTime time.Time) ([]models.ArtifactInfo, error)
	GetFilePathByID(ctx context.Context, id int64) (string, error)
}
