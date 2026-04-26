package metadata

type MetadataService interface {
	CreateArtifact(filePath string, component string, build string, suite string) error
}
