package config

import (
	"os"
	"strconv"
)

type Config struct {
	DatabaseURL    string
	ArtifactVolume string
	MaxFileBytes   int64
	LDAPAddr       string
	LDAPBaseDN     string
	LDAPUser       string
	LDAPPassword   string
}

func Load() Config {
	vol := os.Getenv("ARTIFACT_VOLUME")
	if vol == "" {
		vol = "artifacts"
	}

	maxBytes := int64(50 << 20) // default 50 MB
	if raw := os.Getenv("MAX_FILE_BYTES"); raw != "" {
		if val, err := strconv.ParseInt(raw, 10, 64); err == nil && val > 0 {
			maxBytes = val
		}
	}

	return Config{
		DatabaseURL:    os.Getenv("DATABASE_URL"),
		ArtifactVolume: vol,
		MaxFileBytes:   maxBytes,
		LDAPAddr:       os.Getenv("LDAP_ADDR"),
		LDAPBaseDN:     os.Getenv("LDAP_BASE_DN"),
		LDAPUser:       os.Getenv("LDAP_USER"),
		LDAPPassword:   os.Getenv("LDAP_PASSWORD"),
	}
}
