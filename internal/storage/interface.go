package storage

import (
	"io"
	"mime/multipart"
)

type StorageService interface {
	Save(file multipart.File, header *multipart.FileHeader) (filePath string, err error)
	Open(filePath string) (io.ReadCloser, error)
}
