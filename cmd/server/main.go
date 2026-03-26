package main

import (
	"log"

	// Замени импорты на свои
	"github.com/AlexDigitalART/tasks-service/internal/database"
	"github.com/AlexDigitalART/tasks-service/internal/task"
	transportgrpc "github.com/AlexDigitalART/tasks-service/internal/transport/grpc"
)

func main() {
	// 1. Инициализация БД
	// Предполагаем, что этот код уже был перенесен из монолита
	database.InitDB()

	// 2. Репозиторий и сервис задач
	repo := task.NewRepository(database.DB)
	svc := task.NewService(repo)

	// 3. Подключаем клиент к Users-сервису (который крутится на 50051 порту)
	userClient, conn, err := transportgrpc.NewUserClient("localhost:50051")
	if err != nil {
		log.Fatalf("failed to connect to users: %v", err)
	}
	// defer гарантирует, что соединение с Users-сервисом закроется,
	// когда завершится функция main (программа остановится)
	defer conn.Close()

	// 4. Запуск gRPC Tasks-сервиса
	log.Println("Starting Tasks gRPC server on port 50052...")
	if err := transportgrpc.RunGRPC(svc, userClient); err != nil {
		log.Fatalf("Tasks gRPC server error: %v", err)
	}
}
