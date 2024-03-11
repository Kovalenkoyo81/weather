package services

// Close закрывает ресурсы, связанные с сервисом для работы с пользователями.
func (s *Service) Close() error {

	return s.repo.Close()
}
