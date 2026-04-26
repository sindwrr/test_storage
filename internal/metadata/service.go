package metadata

import "database/sql"

type metadataService struct {
	db *sql.DB
}

func NewMetadataService(db *sql.DB) MetadataService {
	return &metadataService{db: db}
}

func (s *metadataService) CreateArtifact(filePath string, component string, build string, suite string) error {
	return nil
}
