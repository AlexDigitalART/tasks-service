package task

import "gorm.io/gorm"

// Repository структура хранит в себе подключение к БД.
type Repository struct {
	db *gorm.DB
}

// NewRepository — конструктор (создает новый репозиторий).
func NewRepository(db *gorm.DB) *Repository {
	return &Repository{db: db}
}

// Create сохраняет новую задачу в базу.
func (r *Repository) Create(t *Task) error {
	return r.db.Create(t).Error
}

// GetByID ищет задачу по её ID.
func (r *Repository) GetByID(id uint32) (*Task, error) {
	var t Task
	// First ищет первую попавшуюся запись с таким ID
	err := r.db.First(&t, id).Error
	return &t, err
}

// GetAll достает вообще все задачи из таблицы.
func (r *Repository) GetAll() ([]*Task, error) {
	var tasks []*Task
	err := r.db.Find(&tasks).Error
	return tasks, err
}

// GetByUserID ищет задачи конкретного пользователя.
func (r *Repository) GetByUserID(userID uint32) ([]*Task, error) {
	var tasks []*Task
	// Where добавляет условие SQL (WHERE user_id = ?)
	err := r.db.Where("user_id = ?", userID).Find(&tasks).Error
	return tasks, err
}

// Update обновляет существующую задачу.
func (r *Repository) Update(t *Task) error {
	// Save обновит все поля записи, у которой совпадает ID
	return r.db.Save(t).Error
}

// Delete удаляет задачу по ID.
func (r *Repository) Delete(id uint32) error {
	return r.db.Delete(&Task{}, id).Error
}
