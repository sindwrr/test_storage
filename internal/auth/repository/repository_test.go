package repository

import (
	"context"
	"database/sql"
	"regexp"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/stretchr/testify/assert"
)

func TestEnsureUser_Success(t *testing.T) {
	db, mock, err := sqlmock.New()
	assert.NoError(t, err)
	defer db.Close()

	repo := NewUserRepo(db)

	expectedQuery := regexp.QuoteMeta(
		`INSERT INTO users (username, group_id, is_ldap, is_active, created_at) ` +
			`VALUES ($1, 1, true, true, NOW()) ON CONFLICT (username) DO NOTHING;`,
	)

	mock.ExpectExec(expectedQuery).
		WithArgs("testuser").
		WillReturnResult(sqlmock.NewResult(0, 1))

	err = repo.EnsureUser(context.Background(), "testuser")
	assert.NoError(t, err)
	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestEnsureUser_ExecError(t *testing.T) {
	db, mock, err := sqlmock.New()
	assert.NoError(t, err)
	defer db.Close()

	repo := NewUserRepo(db)

	expectedQuery := regexp.QuoteMeta(
		`INSERT INTO users (username, group_id, is_ldap, is_active, created_at) ` +
			`VALUES ($1, 1, true, true, NOW()) ON CONFLICT (username) DO NOTHING;`,
	)

	mock.ExpectExec(expectedQuery).
		WithArgs("testuser").
		WillReturnError(sql.ErrConnDone)

	err = repo.EnsureUser(context.Background(), "testuser")
	assert.Error(t, err)
	assert.NoError(t, mock.ExpectationsWereMet())
}
