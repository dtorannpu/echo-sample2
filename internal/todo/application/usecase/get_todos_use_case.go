package usecase

import (
	"context"
	"echo-sample2/internal/todo/domain"
)

type TodoRepository interface {
	FindAll(ctx context.Context) ([]*domain.Todo, error)
}

type GetTodosUseCase struct {
	repo TodoRepository
}

type Todo struct {
	ID   int64
	Name string
}

type TodoResult struct {
	ID          domain.TodoID
	Title       string
	Description string
}

type GetTodosResult struct {
	Todos []*TodoResult
}

func NewGetTodosUseCase(repo TodoRepository) *GetTodosUseCase {
	return &GetTodosUseCase{repo: repo}
}

func (u *GetTodosUseCase) Execute(ctx context.Context) (*GetTodosResult, error) {
	todos, err := u.repo.FindAll(ctx)
	if err != nil {
		return nil, err
	}

	res := make([]*TodoResult, len(todos))

	for i, todo := range todos {
		res[i] = &TodoResult{
			ID:          todo.ID,
			Title:       todo.Title,
			Description: todo.Description,
		}
	}

	return &GetTodosResult{Todos: res}, nil
}
