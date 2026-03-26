package grpc

import (
	"net"

	"google.golang.org/grpc"

	// Замени импорты
	taskpb "github.com/AlexDigitalART/project-protos/proto/tasks"
	userpb "github.com/AlexDigitalART/project-protos/proto/users"
	"github.com/AlexDigitalART/tasks-service/internal/task"
)

// RunGRPC инициализирует и запускает gRPC сервер на порту 50052
func RunGRPC(svc *task.Service, uc userpb.UserServiceClient) error {
	// 1. Указываем, что будем слушать сеть по протоколу TCP на порту 50052
	lis, err := net.Listen("tcp", ":50052")
	if err != nil {
		return err
	}

	// 2. Создаем пустой gRPC сервер
	grpcSrv := grpc.NewServer()

	// 3. Создаем наш хэндлер со всеми зависимостями
	handler := NewHandler(svc, uc)

	// 4. Регистрируем хэндлер в сервере (говорим gRPC, какой код должен отвечать на запросы)
	taskpb.RegisterTaskServiceServer(grpcSrv, handler)

	// 5. Блокируем горутину и начинаем слушать запросы.
	// Функция вернет ошибку только если сервер упадет или остановится.
	return grpcSrv.Serve(lis)
}
