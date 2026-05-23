package auth

import (
	"context"
	"database/sql"
	"fmt"
	"log"

	"github.com/go-ldap/ldap/v3"
	authrepo "github.com/sindwrr/test_storage/internal/auth/repository"
)

type ldapConn interface {
	Bind(username, password string) error
	Search(*ldap.SearchRequest) (*ldap.SearchResult, error)
	Close()
}

type realLdapConn struct {
	conn *ldap.Conn
}

func (r *realLdapConn) Bind(username, password string) error { return r.conn.Bind(username, password) }
func (r *realLdapConn) Search(req *ldap.SearchRequest) (*ldap.SearchResult, error) {
	return r.conn.Search(req)
}
func (r *realLdapConn) Close() { r.conn.Close() }

type authService struct {
	ldapAddr     string
	ldapBaseDN   string
	bindUser     string
	bindPassword string
	dialer       func(addr string) (ldapConn, error)
	userRepo     authrepo.UserRepository
}

func NewService(ldapAddr, ldapBaseDN, bindUser, bindPassword string, db *sql.DB) AuthService {
	return &authService{
		ldapAddr:     ldapAddr,
		ldapBaseDN:   ldapBaseDN,
		bindUser:     bindUser,
		bindPassword: bindPassword,
		userRepo:     authrepo.NewUserRepo(db),
		dialer: func(addr string) (ldapConn, error) {
			conn, err := ldap.DialURL(fmt.Sprintf("ldap://%s", addr))
			if err != nil {
				return nil, err
			}
			return &realLdapConn{conn: conn}, nil
		},
	}
}

func (s *authService) Validate(username, password string) bool {
	conn, err := s.dialer(s.ldapAddr)
	if err != nil {
		log.Printf("LDAP dial error: %v", err)
		return false
	}
	defer conn.Close()

	authenticated := false
	if s.bindUser != "" {
		err = conn.Bind(s.bindUser, s.bindPassword)
		if err != nil {
			log.Printf("LDAP initial bind error: %v", err)
			return false
		}
		searchRequest := ldap.NewSearchRequest(
			s.ldapBaseDN,
			ldap.ScopeWholeSubtree, ldap.NeverDerefAliases, 0, 0, false,
			fmt.Sprintf("(uid=%s)", username),
			[]string{"dn"},
			nil,
		)
		sr, err := conn.Search(searchRequest)
		if err != nil || len(sr.Entries) == 0 {
			log.Printf("LDAP user search error: %v", err)
			return false
		}
		err = conn.Bind(sr.Entries[0].DN, password)
		authenticated = (err == nil)
	} else {
		err = conn.Bind(username, password)
		authenticated = (err == nil)
	}

	if authenticated && s.userRepo != nil {
		if err := s.userRepo.EnsureUser(context.Background(), username); err != nil {
			log.Printf("failed to ensure user in DB: %v", err)
		}
	}

	return authenticated
}
