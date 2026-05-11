package config

import (
	"testing"
)

func TestLoadDefaults(t *testing.T) {
	cfg := Load()
	if cfg.ArtifactVolume != "artifacts" {
		t.Errorf("expected 'artifacts', got %q", cfg.ArtifactVolume)
	}
	if cfg.MaxFileBytes != 50<<20 {
		t.Errorf("expected 52428800, got %d", cfg.MaxFileBytes)
	}
	if cfg.DatabaseURL != "" {
		t.Errorf("expected empty DatabaseURL, got %q", cfg.DatabaseURL)
	}
}

func TestLoadArtifactVolumeEnv(t *testing.T) {
	t.Setenv("ARTIFACT_VOLUME", "/custom/path")
	cfg := Load()
	if cfg.ArtifactVolume != "/custom/path" {
		t.Errorf("expected '/custom/path', got %q", cfg.ArtifactVolume)
	}
}

func TestLoadMaxFileBytesValid(t *testing.T) {
	t.Setenv("MAX_FILE_BYTES", "1048576") // 1 MB
	cfg := Load()
	if cfg.MaxFileBytes != 1048576 {
		t.Errorf("expected 1048576, got %d", cfg.MaxFileBytes)
	}
}

func TestLoadMaxFileBytesInvalid(t *testing.T) {
	t.Setenv("MAX_FILE_BYTES", "invalid")
	cfg := Load()
	if cfg.MaxFileBytes != 50<<20 {
		t.Errorf("expected default (52428800) for invalid input, got %d", cfg.MaxFileBytes)
	}
}

func TestLoadMaxFileBytesZero(t *testing.T) {
	t.Setenv("MAX_FILE_BYTES", "0")
	cfg := Load()
	if cfg.MaxFileBytes != 50<<20 {
		t.Errorf("expected default (52428800) for zero value, got %d", cfg.MaxFileBytes)
	}
}

func TestLoadDatabaseURLEnv(t *testing.T) {
	t.Setenv("DATABASE_URL", "postgres://user:pass@localhost/testdb")
	cfg := Load()
	if cfg.DatabaseURL != "postgres://user:pass@localhost/testdb" {
		t.Errorf("expected database URL, got %q", cfg.DatabaseURL)
	}
}
