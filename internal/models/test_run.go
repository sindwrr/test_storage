package models

import "time"

type TestRun struct {
	ID         int       `json:"id"`
	BuildID    int       `json:"build_id"`
	SuiteID    int       `json:"suite_id"`
	StatusID   int       `json:"status_id"`
	StartedAt  time.Time `json:"started_at"`
	FinishedAt time.Time `json:"finished_at,omitempty"`
	UpdatedAt  time.Time `json:"updated_at"`
}
