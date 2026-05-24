package repository

import (
	"context"
	"database/sql"
)

type DBTX interface {
	ExecContext(ctx context.Context, query string, args ...any) (sql.Result, error)
}

type UserRepo struct {
	db DBTX
}

func NewUserRepo(db DBTX) *UserRepo {
	return &UserRepo{db: db}
}

func (r *UserRepo) EnsureUser(ctx context.Context, username string) error {
	const query = `
        INSERT INTO users (username, group_id, is_ldap, is_active, created_at)
        VALUES ($1, 1, true, true, NOW())
        ON CONFLICT (username) DO NOTHING;
    `
	_, err := r.db.ExecContext(ctx, query, username)
	return err
}

func (r *UserRepo) SetActive(ctx context.Context, username string, active bool) error {
	const query = `UPDATE users SET is_active = $1 WHERE username = $2`
	_, err := r.db.ExecContext(ctx, query, active, username)
	return err
}
