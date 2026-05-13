package auth

import (
	"fmt"
	"log"

	"github.com/go-ldap/ldap/v3"
)

type authService struct {
	ldapAddr     string
	ldapBaseDN   string
	bindUser     string
	bindPassword string
}

func NewService(ldapAddr, ldapBaseDN, bindUser, bindPassword string) AuthService {
	return &authService{
		ldapAddr:     ldapAddr,
		ldapBaseDN:   ldapBaseDN,
		bindUser:     bindUser,
		bindPassword: bindPassword,
	}
}

func (s *authService) Validate(username, password string) bool {
	conn, err := ldap.DialURL(fmt.Sprintf("ldap://%s", s.ldapAddr))
	if err != nil {
		log.Printf("LDAP dial error: %v", err)
		return false
	}
	defer conn.Close()

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
		return err == nil
	}

	err = conn.Bind(username, password)
	return err == nil
}
