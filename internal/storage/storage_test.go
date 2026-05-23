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

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
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

func TestNewStorageService_MkdirError(t *testing.T) {
	// Путь с нулевым байтом гарантированно вызывает ошибку
	_, err := NewStorageService("/invalid\x00path")
	assert.Error(t, err)
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

func TestSave_CreateFileError(t *testing.T) {
	svc, err := NewStorageService(t.TempDir())
	require.NoError(t, err)

	file, header := createMultipartFileHelper(t, "test.txt", "content")
	header.Filename = "bad\x00name.txt"
	_, err = svc.Save(file, header)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "create file")
}

type errReader struct{}

func (e errReader) Read(p []byte) (int, error) { return 0, io.ErrUnexpectedEOF }

type mockMultipartFile struct {
	io.Reader
}

func (m *mockMultipartFile) Close() error                                 { return nil }
func (m *mockMultipartFile) Seek(offset int64, whence int) (int64, error) { return 0, nil }
func (m *mockMultipartFile) ReadAt(p []byte, off int64) (n int, err error) {
	return m.Reader.Read(p)
}

func TestSave_CopyError(t *testing.T) {
	svc, err := NewStorageService(t.TempDir())
	require.NoError(t, err)

	_, header := createMultipartFileHelper(t, "test.txt", "content")
	badFile := &mockMultipartFile{Reader: &errReader{}}
	_, err = svc.Save(badFile, header)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "copy file")
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
