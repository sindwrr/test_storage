package auth

type AuthService interface {
	Validate(username, password string) bool
}
