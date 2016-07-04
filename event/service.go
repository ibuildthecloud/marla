package event

type Service struct {
}

func New() *Service {
	return &Service{}
}

func (s *Service) LogImageEvent(id, name, action string) {
}
