package repository

import "context"

type UserRepository interface {
	EnsureUser(ctx context.Context, username string) error
	SetActive(ctx context.Context, username string, active bool) error
	GetGroupID(ctx context.Context, username string) (int, error)
}
