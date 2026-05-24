package auth

type AuthService interface {
	Validate(username, password string) bool
	SetUserActive(username string, active bool)
}
