package auth

type authService struct{}

func NewService() AuthService {
	return &authService{}
}

func (s *authService) Validate(username, password string) bool {
	// TODO: replace with full auth later
	return username == "admin" && password == "123"
}
