package create

import (
	"context"
	"echo-sample2/internal/todo/domain"
	"echo-sample2/internal/todo/domain/repository"
)

type UseCase struct {
	txManager repository.TransactionManager
}
type Command struct {
	Title       string
	Description string
}

func New(tm repository.TransactionManager) *UseCase {
	return &UseCase{txManager: tm}
}

func (u *UseCase) Execute(ctx context.Context, request Command) error {
	return u.txManager.Do(ctx, func(tx repository.Transaction) error {
		return tx.TodoRepo().Save(ctx, createTodo(request))
	})
}

func createTodo(request Command) *domain.Todo {
	return &domain.Todo{
		ID:          domain.NewTodoID(),
		Title:       request.Title,
		Description: request.Description,
	}
}
