package repository

import "context"

type UserRepository interface {
	EnsureUser(ctx context.Context, username string) error
}
