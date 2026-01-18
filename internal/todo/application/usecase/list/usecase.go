package list

import (
	"context"
	"echo-sample2/internal/todo/domain"
)

type TodoRepository interface {
	FindAll(ctx context.Context) ([]*domain.Todo, error)
}

type UseCase struct {
	repo TodoRepository
}

type Todo struct {
	ID   int64
	Name string
}

type TodoItem struct {
	ID          domain.TodoID
	Title       string
	Description string
}

type Result struct {
	Todos []*TodoItem
}

func New(repo TodoRepository) *UseCase {
	return &UseCase{repo: repo}
}

func (u *UseCase) Execute(ctx context.Context) (*Result, error) {
	todos, err := u.repo.FindAll(ctx)
	if err != nil {
		return nil, err
	}

	res := make([]*TodoItem, len(todos))

	for i, todo := range todos {
		res[i] = &TodoItem{
			ID:          todo.ID,
			Title:       todo.Title,
			Description: todo.Description,
		}
	}

	return &Result{Todos: res}, nil
}
