package auth

type Service struct{}

func NewService() *Service {
	return &Service{}
}

func (s *Service) Validate(username, password string) bool {
	// TODO: replace with DB later
	return username == "admin" && password == "123"
}
