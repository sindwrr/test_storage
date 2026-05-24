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

func TestSetActive_Success(t *testing.T) {
	db, mock, err := sqlmock.New()
	assert.NoError(t, err)
	defer db.Close()

	repo := NewUserRepo(db)

	mock.ExpectExec(regexp.QuoteMeta("UPDATE users SET is_active = $1 WHERE username = $2")).
		WithArgs(true, "testuser").
		WillReturnResult(sqlmock.NewResult(0, 1))

	err = repo.SetActive(context.Background(), "testuser", true)
	assert.NoError(t, err)
	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestSetActive_Error(t *testing.T) {
	db, mock, err := sqlmock.New()
	assert.NoError(t, err)
	defer db.Close()

	repo := NewUserRepo(db)

	mock.ExpectExec(regexp.QuoteMeta("UPDATE users SET is_active = $1 WHERE username = $2")).
		WithArgs(false, "testuser").
		WillReturnError(sql.ErrConnDone)

	err = repo.SetActive(context.Background(), "testuser", false)
	assert.Error(t, err)
	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestGetGroupID_Success(t *testing.T) {
	db, mock, err := sqlmock.New()
	assert.NoError(t, err)
	defer db.Close()

	repo := NewUserRepo(db)
	rows := sqlmock.NewRows([]string{"group_id"}).AddRow(2)
	mock.ExpectQuery(regexp.QuoteMeta(`SELECT group_id FROM users WHERE username = $1`)).
		WithArgs("admin").
		WillReturnRows(rows)

	groupID, err := repo.GetGroupID(context.Background(), "admin")
	assert.NoError(t, err)
	assert.Equal(t, 2, groupID)
	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestGetGroupID_NotFound(t *testing.T) {
	db, mock, err := sqlmock.New()
	assert.NoError(t, err)
	defer db.Close()

	repo := NewUserRepo(db)
	mock.ExpectQuery(regexp.QuoteMeta(`SELECT group_id FROM users WHERE username = $1`)).
		WithArgs("nobody").
		WillReturnError(sql.ErrNoRows)

	_, err = repo.GetGroupID(context.Background(), "nobody")
	assert.ErrorIs(t, err, sql.ErrNoRows)
	assert.NoError(t, mock.ExpectationsWereMet())
}
