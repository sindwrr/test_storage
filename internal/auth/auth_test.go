package auth

import "testing"

func TestValidateSuccess(t *testing.T) {
	svc := NewService()
	if !svc.Validate("admin", "123") {
		t.Fatal("expected true for admin/123")
	}
}

func TestValidateBadPassword(t *testing.T) {
	svc := NewService()
	if svc.Validate("admin", "wrong") {
		t.Fatal("expected false for wrong password")
	}
}

func TestValidateBadUsername(t *testing.T) {
	svc := NewService()
	if svc.Validate("wrong", "123") {
		t.Fatal("expected false for wrong username")
	}
}

func TestValidateEmptyCredentials(t *testing.T) {
	svc := NewService()
	if svc.Validate("", "") {
		t.Fatal("expected false for empty credentials")
	}
}
