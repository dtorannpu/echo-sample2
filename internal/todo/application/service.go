package application

import (
	"context"
	"echo-sample2/internal/todo/domain"
	"echo-sample2/internal/todo/domain/repository"
	"echo-sample2/internal/todo/dto"
)

type TodoService struct {
	txManager repository.TransactionManager
}

func NewTodoService(tm repository.TransactionManager) *TodoService {
	return &TodoService{txManager: tm}
}

func (s *TodoService) Create(ctx context.Context, request dto.CreateTodoInput) error {
	return s.txManager.Do(ctx, func(tx repository.Transaction) error {
		return tx.TodoRepo().Save(ctx, createTodo(request))
	})
}

func createTodo(request dto.CreateTodoInput) *domain.Todo {
	return &domain.Todo{
		ID:          domain.NewTodoID(),
		Title:       request.Title,
		Description: request.Description,
	}
}
