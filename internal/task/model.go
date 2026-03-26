package task

// Task описывает структуру таблицы в базе данных.
// Теги `gorm:"..."` подсказывают GORM, как именно создавать колонки.
type Task struct {
	ID     uint32 `gorm:"primaryKey"`    // Уникальный идентификатор задачи
	UserID uint32 `gorm:"index"`         // ID пользователя (индекс ускоряет поиск по этому полю)
	Title  string `gorm:"not null"`      // Текст задачи (например: "Найти безопасное укрытие")
	IsDone bool   `gorm:"default:false"` // Статус выполнения (по умолчанию задача не выполнена)
}
