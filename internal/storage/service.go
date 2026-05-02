package storage

import (
	"fmt"
	"io"
	"mime/multipart"
	"os"
	"path/filepath"
	"strings"
	"time"
)

type storageService struct {
	basePath string
}

func NewStorageService(basePath string) (StorageService, error) {
	if err := os.MkdirAll(basePath, 0o755); err != nil {
		return nil, fmt.Errorf("create artifacts dir: %s", err)
	}
	return &storageService{basePath: basePath}, nil
}

func (s *storageService) Save(file multipart.File, header *multipart.FileHeader) (string, error) {
	fileName := filepath.Base(header.Filename)
	newName := fmt.Sprintf("%s_%s", time.Now().Format("2006-01-02_15-04-05.000000"), fileName)
	filePath := filepath.Join(s.basePath, newName)

	dst, err := os.Create(filePath)
	if err != nil {
		return "", fmt.Errorf("create file: %s", err)
	}
	defer dst.Close()

	if _, err := io.Copy(dst, file); err != nil {
		os.Remove(filePath)
		return "", fmt.Errorf("copy file: %s", err)
	}

	return filePath, nil
}

func (s *storageService) Open(filePath string) (io.ReadCloser, error) {
	cleanPath := filepath.Clean(filePath)

	absPath := filepath.Join(s.basePath, cleanPath)

	if !strings.HasPrefix(absPath, s.basePath) {
		return nil, fmt.Errorf("storage: file not found")
	}

	file, err := os.Open(absPath)
	if err != nil {
		return nil, fmt.Errorf("open file %s: %w", cleanPath, err)
	}

	return file, nil
}
