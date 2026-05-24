package auth

import (
	"context"
	"errors"
	"testing"

	"github.com/go-ldap/ldap/v3"
	"github.com/stretchr/testify/assert"
)

type mockLdapConn struct {
	bindFunc   func(username, password string) error
	searchFunc func(*ldap.SearchRequest) (*ldap.SearchResult, error)
	closeCalls int
}

func (m *mockLdapConn) Bind(username, password string) error {
	return m.bindFunc(username, password)
}
func (m *mockLdapConn) Search(req *ldap.SearchRequest) (*ldap.SearchResult, error) {
	return m.searchFunc(req)
}
func (m *mockLdapConn) Close() { m.closeCalls++ }

func TestValidate_ServiceBind_Success(t *testing.T) {
	mockConn := &mockLdapConn{
		bindFunc: func(username, password string) error {
			if username == "cn=admin,dc=example,dc=org" && password == "adminpass" {
				return nil
			}
			if username == "uid=testuser,dc=example,dc=org" && password == "userpass" {
				return nil
			}
			return errors.New("invalid credentials")
		},
		searchFunc: func(req *ldap.SearchRequest) (*ldap.SearchResult, error) {
			return &ldap.SearchResult{
				Entries: []*ldap.Entry{
					{DN: "uid=testuser,dc=example,dc=org"},
				},
			}, nil
		},
	}
	svc := &authService{
		ldapAddr:     "ldap.example.com",
		ldapBaseDN:   "dc=example,dc=org",
		bindUser:     "cn=admin,dc=example,dc=org",
		bindPassword: "adminpass",
		dialer:       func(addr string) (ldapConn, error) { return mockConn, nil },
	}
	ok := svc.Validate("testuser", "userpass")
	assert.True(t, ok)
	assert.GreaterOrEqual(t, mockConn.closeCalls, 1)
}

func TestValidate_ServiceBind_SearchEmpty(t *testing.T) {
	mockConn := &mockLdapConn{
		bindFunc: func(username, password string) error {
			return nil
		},
		searchFunc: func(req *ldap.SearchRequest) (*ldap.SearchResult, error) {
			return &ldap.SearchResult{Entries: []*ldap.Entry{}}, nil // 0 записей
		},
	}
	svc := &authService{
		ldapAddr:     "ldap.example.com",
		ldapBaseDN:   "dc=example,dc=org",
		bindUser:     "cn=admin,dc=example,dc=org",
		bindPassword: "adminpass",
		dialer:       func(addr string) (ldapConn, error) { return mockConn, nil },
	}

	ok := svc.Validate("testuser", "userpass")
	assert.False(t, ok)
}

func TestValidate_ServiceBind_SearchError(t *testing.T) {
	mockConn := &mockLdapConn{
		bindFunc: func(username, password string) error { return nil },
		searchFunc: func(req *ldap.SearchRequest) (*ldap.SearchResult, error) {
			return nil, errors.New("search failed")
		},
	}
	svc := &authService{
		ldapAddr:     "ldap.example.com",
		ldapBaseDN:   "dc=example,dc=org",
		bindUser:     "cn=admin,dc=example,dc=org",
		bindPassword: "adminpass",
		dialer:       func(addr string) (ldapConn, error) { return mockConn, nil },
	}
	ok := svc.Validate("testuser", "userpass")
	assert.False(t, ok)
}

func TestValidate_DialerError(t *testing.T) {
	svc := &authService{
		ldapAddr: "invalid.addr",
		dialer: func(addr string) (ldapConn, error) {
			return nil, errors.New("dial failed")
		},
	}
	ok := svc.Validate("user", "pass")
	assert.False(t, ok)
}

type mockUserRepo struct {
	setActiveCalled  bool
	setActiveArgs    [2]interface{}
	setActiveErr     error
	getGroupIDFn     func(username string) (int, error) // новое
	getGroupIDCalled bool
}

func (m *mockUserRepo) EnsureUser(ctx context.Context, username string) error { return nil }
func (m *mockUserRepo) SetActive(ctx context.Context, username string, active bool) error {
	m.setActiveCalled = true
	m.setActiveArgs = [2]interface{}{username, active}
	return m.setActiveErr
}
func (m *mockUserRepo) GetGroupID(ctx context.Context, username string) (int, error) {
	m.getGroupIDCalled = true
	if m.getGroupIDFn != nil {
		return m.getGroupIDFn(username)
	}
	return 0, nil
}

func TestSetUserActive_CallsRepo(t *testing.T) {
	repo := &mockUserRepo{}
	svc := &authService{
		userRepo: repo,
	}
	svc.SetUserActive("testuser", true)
	assert.True(t, repo.setActiveCalled)
	assert.Equal(t, "testuser", repo.setActiveArgs[0])
	assert.Equal(t, true, repo.setActiveArgs[1])
}

func TestSetUserActive_NilRepo(t *testing.T) {
	svc := &authService{
		userRepo: nil,
	}

	assert.NotPanics(t, func() {
		svc.SetUserActive("testuser", false)
	})
}

func TestGetUserGroup_Success(t *testing.T) {
	repo := &mockUserRepo{getGroupIDFn: func(u string) (int, error) { return 2, nil }}
	svc := &authService{userRepo: repo}

	groupID, err := svc.GetUserGroup("admin")
	assert.NoError(t, err)
	assert.Equal(t, 2, groupID)
	assert.True(t, repo.getGroupIDCalled)
}

func TestGetUserGroup_NilRepo(t *testing.T) {
	svc := &authService{userRepo: nil}
	_, err := svc.GetUserGroup("admin")
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "not configured")
}
