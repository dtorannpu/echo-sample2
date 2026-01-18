package usecase

import (
	"context"
	"echo-sample2/internal/todo/domain"
	"echo-sample2/internal/todo/domain/repository"
)

type CreateTodoUseCase struct {
	txManager repository.TransactionManager
}
type CreateTodoCommand struct {
	Title       string
	Description string
}

func NewCreateTodoUseCase(tm repository.TransactionManager) *CreateTodoUseCase {
	return &CreateTodoUseCase{txManager: tm}
}

func (u *CreateTodoUseCase) Execute(ctx context.Context, request CreateTodoCommand) error {
	return u.txManager.Do(ctx, func(tx repository.Transaction) error {
		return tx.TodoRepo().Save(ctx, createTodo(request))
	})
}

func createTodo(request CreateTodoCommand) *domain.Todo {
	return &domain.Todo{
		ID:          domain.NewTodoID(),
		Title:       request.Title,
		Description: request.Description,
	}
}
