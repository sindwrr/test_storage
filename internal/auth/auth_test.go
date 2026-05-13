package auth

import (
	"testing"
)

func TestValidateFreeIPA(t *testing.T) {

	svc := NewService(
		"ipa.demo1.freeipa.org:389",
		"dc=demo1,dc=freeipa,dc=org",
		"",
		"",
	)

	if !svc.Validate("uid=employee,cn=users,cn=accounts,dc=demo1,dc=freeipa,dc=org", "Secret123") {
		t.Error("expected successful authentication")
	}

	if svc.Validate("uid=employee,cn=users,cn=accounts,dc=demo1,dc=freeipa,dc=org", "wrong") {
		t.Error("expected authentication failure")
	}
}
