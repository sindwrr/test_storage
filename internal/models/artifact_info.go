package models

import "time"

type ArtifactInfo struct {
	ID          int64
	DownloadURL string
	FileName    string
	FileSize    int64
	Component   string
	Build       string
	Suite       string
	UploadTime  time.Time
}
