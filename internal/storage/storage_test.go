package storage

import (
	"bytes"
	"io"
	"mime/multipart"
	"net/http"
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func createMultipartFileHelper(t *testing.T, fileName, content string) (multipart.File, *multipart.FileHeader) {
	t.Helper()

	body := &bytes.Buffer{}
	writer := multipart.NewWriter(body)
	part, err := writer.CreateFormFile("file", fileName)
	if err != nil {
		t.Fatalf("failed to create form file: %v", err)
	}
	if _, err := part.Write([]byte(content)); err != nil {
		t.Fatalf("failed to write to part: %v", err)
	}
	writer.Close()

	req, err := http.NewRequest("POST", "/", body)
	if err != nil {
		t.Fatalf("failed to create request: %v", err)
	}
	req.Header.Set("Content-Type", writer.FormDataContentType())

	file, header, err := req.FormFile("file")
	if err != nil {
		t.Fatalf("failed to parse form file: %v", err)
	}
	return file, header
}

func TestNewStorageService_CreatesDir(t *testing.T) {
	tmpDir := t.TempDir()
	basePath := filepath.Join(tmpDir, "artifacts")

	svc, err := NewStorageService(basePath)
	if err != nil {
		t.Fatalf("newStorageService failed: %v", err)
	}
	if svc == nil {
		t.Fatal("service is nil")
	}

	info, err := os.Stat(basePath)
	if err != nil {
		t.Fatalf("directory not created: %v", err)
	}
	if !info.IsDir() {
		t.Fatal("path is not a directory")
	}
}

func TestSaveAndOpen(t *testing.T) {
	tmpDir := t.TempDir()
	svc, err := NewStorageService(tmpDir)
	if err != nil {
		t.Fatalf("failed to create service: %v", err)
	}

	fileName := "testfile.txt"
	fileContent := "Hello, storage!"

	file, header := createMultipartFileHelper(t, fileName, fileContent)
	defer file.Close()

	savedPath, err := svc.Save(file, header)
	if err != nil {
		t.Fatalf("save failed: %v", err)
	}

	if _, err := os.Stat(savedPath); os.IsNotExist(err) {
		t.Fatalf("file does not exist: %s", savedPath)
	}

	if !strings.HasPrefix(savedPath, tmpDir) {
		t.Errorf("saved file is outside basePath: %s", savedPath)
	}

	if filepath.Base(savedPath)[len(filepath.Base(savedPath))-len(fileName):] != fileName {
		t.Errorf("file name does not contain original name: %s", savedPath)
	}

	reader, err := svc.Open(savedPath)
	if err != nil {
		t.Fatalf("open failed: %v", err)
	}
	defer reader.Close()

	readContent, err := io.ReadAll(reader)
	if err != nil {
		t.Fatalf("failed to read opened file: %v", err)
	}
	if string(readContent) != fileContent {
		t.Errorf("content mismatch: got %q, wanted %q", readContent, fileContent)
	}
}

func TestOpen_NonExistent(t *testing.T) {
	tmpDir := t.TempDir()
	svc, err := NewStorageService(tmpDir)
	if err != nil {
		t.Fatalf("failed to create service: %v", err)
	}

	_, err = svc.Open(filepath.Join(tmpDir, "random.bin"))
	if err == nil {
		t.Fatal("expected error when opening non-existent file, got nil")
	}
}
