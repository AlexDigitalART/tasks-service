package database

import (
	"log"

	"github.com/AlexDigitalART/tasks-service/internal/task"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

// DB — глобальная переменная, которая хранит подключение к базе.
// В более сложных архитектурах избегают глобальных переменных, но для старта это окей.
var DB *gorm.DB

// InitDB открывает соединение и автоматически создает таблицы.
func InitDB() {
	var err error
	// Подключаемся к файловой базе данных SQLite
	DB, err = gorm.Open(sqlite.Open("tasks.db"), &gorm.Config{})
	if err != nil {
		// Если БД недоступна — приложение дальше работать не может, вызываем фатальную ошибку
		log.Fatalf("failed to connect to database: %v", err)
	}

	// AutoMigrate сверяет нашу структуру Task с реальной таблицей в БД.
	// Если таблицы нет — создаст, если добавились новые поля — допишет их.
	err = DB.AutoMigrate(&task.Task{})
	if err != nil {
		log.Fatalf("failed to migrate database: %v", err)
	}

	log.Println("Database connected and migrated successfully")
}
