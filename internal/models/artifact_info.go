package models

import "time"

type ArtifactInfo struct {
	DownloadURL string
	FileName    string
	FileSize    int64
	Component   string
	Build       string
	Suite       string
	UploadTime  time.Time
}
