package models

import "time"

type TestArtifact struct {
	ID         int       `json:"id"`
	RunID      int       `json:"run_id"`
	StatusID   int       `json:"status_id"`
	FileURL    string    `json:"file_url"`
	FileTypeID int       `json:"file_type_id"`
	CreatedAt  time.Time `json:"created_at"`
	UpdatedAt  time.Time `json:"updated_at"`
}
