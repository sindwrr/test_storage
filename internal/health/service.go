package health

import "context"

type HealthService interface {
	Alive(ctx context.Context) error
	Ready(ctx context.Context) error
}
