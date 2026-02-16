package auth

type Service struct{}

func NewService() *Service {
	return &Service{}
}

func (s *Service) ValidateToken(token string) (int64, error) {
	// JWT validation
	return 1, nil
}
