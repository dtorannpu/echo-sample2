package get

import (
	"context"
	"echo-sample2/internal/todo/domain"
)

type TodoRepository interface {
	FindById(ctx context.Context, id domain.TodoID) (*domain.Todo, error)
}

type Todo struct {
	ID          domain.TodoID
	Title       string
	Description string
}

type Command struct {
	ID domain.TodoID
}

type UseCase struct {
	repo TodoRepository
}

func New(repo TodoRepository) *UseCase {
	return &UseCase{repo: repo}
}

func (u *UseCase) Execute(ctx context.Context, command Command) (*Todo, error) {
	todo, err := u.repo.FindById(ctx, command.ID)
	if err != nil {
		return nil, err
	}

	if todo == nil {
		return nil, nil
	}

	return &Todo{
		ID:          todo.ID,
		Title:       todo.Title,
		Description: todo.Description,
	}, nil
}
