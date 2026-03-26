package task

// Service хранит в себе ссылку на репозиторий, чтобы просить его сохранять данные.
type Service struct {
	repo *Repository
}

// NewService — конструктор для сервиса.
func NewService(repo *Repository) *Service {
	return &Service{repo: repo}
}

// CreateTask принимает данные, создает структуру и передает в репозиторий.
func (s *Service) CreateTask(t Task) (*Task, error) {
	// Здесь могла бы быть бизнес-логика: например, проверка, что Title не пустой
	err := s.repo.Create(&t)
	if err != nil {
		return nil, err
	}
	return &t, nil
}

// GetTaskByID запрашивает задачу у репозитория.
func (s *Service) GetTaskByID(id uint32) (*Task, error) {
	return s.repo.GetByID(id)
}

// GetAllTasks получает список всех задач.
func (s *Service) GetAllTasks() ([]*Task, error) {
	return s.repo.GetAll()
}

// GetTasksByUserID получает задачи по ID юзера.
func (s *Service) GetTasksByUserID(userID uint32) ([]*Task, error) {
	return s.repo.GetByUserID(userID)
}

// UpdateTask передает обновленные данные в базу.
func (s *Service) UpdateTask(t Task) (*Task, error) {
	err := s.repo.Update(&t)
	if err != nil {
		return nil, err
	}
	return &t, nil
}

// DeleteTask удаляет задачу.
func (s *Service) DeleteTask(id uint32) error {
	return s.repo.Delete(id)
}
