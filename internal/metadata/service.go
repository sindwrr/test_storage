package metadata

type metadataService struct{}

func NewMetadataService() MetadataService {
	return &metadataService{}
}
