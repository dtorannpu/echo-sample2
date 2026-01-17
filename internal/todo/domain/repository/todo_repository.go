package repository

import (
	"context"
	"echo-sample2/internal/todo/domain"
)

type TodoRepository interface {
	Save(ctx context.Context, todo *domain.Todo) error
}
