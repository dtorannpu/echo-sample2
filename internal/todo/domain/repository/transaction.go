package repository

import (
	"context"
	"echo-sample2/internal/todo/domain"
)

type TodoRepository interface {
	Save(ctx context.Context, todo *domain.Todo) error
	Update(ctx context.Context, todo *domain.Todo) error
	Delete(ctx context.Context, id domain.TodoID) error
}

type Transaction interface {
	TodoRepo() TodoRepository
}

type TransactionManager interface {
	Do(ctx context.Context, fn func(tx Transaction) error) error
}
