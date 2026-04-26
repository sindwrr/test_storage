package config

import "os"

type Config struct {
	DatabaseURL    string
	ArtifactVolume string
}

func Load() Config {
	return Config{DatabaseURL: os.Getenv("DATABASE_URL"), ArtifactVolume: "artifacts"}
}
