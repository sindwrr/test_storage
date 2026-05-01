package metadata

type MetadataService interface {
	CreateArtifact(filePath string, fileSize int64, component string, build string, suite string) error
}
