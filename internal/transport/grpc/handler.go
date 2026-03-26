package grpc

import (
	"context"
	"fmt"

	// Замени импорты на свои
	taskpb "github.com/AlexDigitalART/project-protos/proto/tasks"
	userpb "github.com/AlexDigitalART/project-protos/proto/users"
	"github.com/AlexDigitalART/tasks-service/internal/task"
)

// Handler — это структура, которая "слушает" запросы.
// В ней хранятся зависимости: бизнес-логика (svc) и клиент для похода в Users (userClient).
type Handler struct {
	svc        *task.Service
	userClient userpb.UserServiceClient

	// Встраиваем структуру для обратной совместимости (требование gRPC)
	taskpb.UnimplementedTaskServiceServer
}

// NewHandler — конструктор для создания нашего хэндлера.
func NewHandler(svc *task.Service, uc userpb.UserServiceClient) *Handler {
	return &Handler{svc: svc, userClient: uc}
}

// CreateTask создает новую задачу.
func (h *Handler) CreateTask(ctx context.Context, req *taskpb.CreateTaskRequest) (*taskpb.CreateTaskResponse, error) {
	// 1. Идем по сети в другой микросервис (Users), чтобы проверить, существует ли пользователь.
	// ctx (контекст) передаем дальше, чтобы при необходимости можно было отменить запрос.
	_, err := h.userClient.GetUser(ctx, &userpb.GetUserRequest{Id: req.Id})
	if err != nil {
		// Если вернулась ошибка — пользователя нет (или сервис недоступен). Прерываем создание.
		return nil, fmt.Errorf("user %d not found: %w", req.UserId, err)
	}

	// 2. Внутренняя логика: вызываем метод монолита для создания задачи в БД.
	t, err := h.svc.CreateTask(task.Task{UserID: req.UserId, Title: req.Title})
	if err != nil {
		return nil, err
	}

	// 3. Формируем успешный ответ, перекладывая данные из модели БД в модель Protobuf.
	return &taskpb.CreateTaskResponse{
		Task: &taskpb.Task{
			Id:     t.ID,
			UserId: t.UserID,
			Title:  t.Title,
			IsDone: t.IsDone,
		},
	}, nil
}

// GetTask получает задачу по ID (без проверки пользователя).
func (h *Handler) GetTask(ctx context.Context, req *taskpb.GetTaskRequest) (*taskpb.GetTaskResponse, error) {
	t, err := h.svc.GetTaskByID(req.Id) // Предполагаем, что такой метод есть в svc
	if err != nil {
		return nil, err
	}

	return &taskpb.GetTaskResponse{
		Task: &taskpb.Task{
			Id:     t.ID,
			UserId: t.UserID,
			Title:  t.Title,
			IsDone: t.IsDone,
		},
	}, nil
}

// ListTasks возвращает список всех задач (без проверки пользователя).
func (h *Handler) ListTasks(ctx context.Context, req *taskpb.ListTasksRequest) (*taskpb.ListTasksResponse, error) {
	tasks, err := h.svc.GetAllTasks() // Предполагаемый метод
	if err != nil {
		return nil, err
	}

	// Преобразуем слайс моделей БД в слайс моделей Protobuf
	var pbTasks []*taskpb.Task
	for _, t := range tasks {
		pbTasks = append(pbTasks, &taskpb.Task{
			Id:     t.ID,
			UserId: t.UserID,
			Title:  t.Title,
			IsDone: t.IsDone,
		})
	}

	return &taskpb.ListTasksResponse{Tasks: pbTasks}, nil
}

// ListTasksByUser возвращает задачи конкретного юзера.
func (h *Handler) ListTasksByUser(ctx context.Context, req *taskpb.ListTasksByUserRequest) (*taskpb.ListTasksResponse, error) {
	tasks, err := h.svc.GetTasksByUserID(req.UserId) // Предполагаемый метод
	if err != nil {
		return nil, err
	}

	var pbTasks []*taskpb.Task
	for _, t := range tasks {
		pbTasks = append(pbTasks, &taskpb.Task{
			Id:     t.ID,
			UserId: t.UserID,
			Title:  t.Title,
			IsDone: t.IsDone,
		})
	}

	return &taskpb.ListTasksResponse{Tasks: pbTasks}, nil
}

// UpdateTask обновляет задачу. Как сказано в задании, здесь мы тоже проверяем пользователя (если меняем владельца).
func (h *Handler) UpdateTask(ctx context.Context, req *taskpb.UpdateTaskRequest) (*taskpb.UpdateTaskResponse, error) {
	// Проверяем пользователя (если UserId передается в запросе на обновление)
	// В зависимости от твоей схемы Protobuf, ID пользователя может быть в UpdateTaskRequest.
	// Здесь я делаю проверку по аналогии с CreateTask.
	_, err := h.userClient.GetUser(ctx, &userpb.GetUserRequest{Id: req.UserId})
	if err != nil {
		return nil, fmt.Errorf("user %d not found: %w", req.UserId, err)
	}

	t, err := h.svc.UpdateTask(task.Task{
		ID:     req.Id,
		UserID: req.UserId,
		Title:  req.Title,
		IsDone: req.IsDone,
	})
	if err != nil {
		return nil, err
	}

	return &taskpb.UpdateTaskResponse{
		Task: &taskpb.Task{
			Id:     t.ID,
			UserId: t.UserID,
			Title:  t.Title,
			IsDone: t.IsDone,
		},
	}, nil
}

// DeleteTask удаляет задачу.
func (h *Handler) DeleteTask(ctx context.Context, req *taskpb.DeleteTaskRequest) (*taskpb.DeleteTaskResponse, error) {
	err := h.svc.DeleteTask(req.Id) // Предполагаемый метод
	if err != nil {
		return nil, err
	}

	return &taskpb.DeleteTaskResponse{Success: true}, nil
}
