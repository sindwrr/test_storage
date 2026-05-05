package health

import (
	"context"
	"database/sql"
	"fmt"
)

type healthService struct {
	db *sql.DB
}

func NewService(db *sql.DB) HealthService {
	return &healthService{db: db}
}

func (s *healthService) Alive(ctx context.Context) error {
	return nil
}

func (s *healthService) Ready(ctx context.Context) error {
	if s.db == nil {
		return fmt.Errorf("DB not configured!")
	}
	return s.db.PingContext(ctx)
}
