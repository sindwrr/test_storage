package storage

import "mime/multipart"

type StorageService interface {
	Save(file multipart.File, header *multipart.FileHeader) (filePath string, err error)
}
