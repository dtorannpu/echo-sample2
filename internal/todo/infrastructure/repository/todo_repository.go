package repository

import (
	"context"
	"echo-sample2/internal/todo/domain"
	"echo-sample2/internal/todo/infrastructure"
	"errors"

	"gorm.io/gorm"
)

type TodoRepository struct {
	db *gorm.DB
}

func NewTodoRepository(db *gorm.DB) *TodoRepository {
	return &TodoRepository{db: db}
}

func (r *TodoRepository) Save(ctx context.Context, todo *domain.Todo) error {
	entity := toEntity(todo)

	return gorm.G[infrastructure.TodoEntity](r.db).Create(ctx, entity)
}

func (r *TodoRepository) FindAll(ctx context.Context) ([]*domain.Todo, error) {
	todos, err := gorm.G[infrastructure.TodoEntity](r.db).Find(ctx)
	if err != nil {
		return nil, err
	}

	res := make([]*domain.Todo, len(todos))

	for i, todo := range todos {
		res[i] = &domain.Todo{
			ID:          todo.ID,
			Title:       todo.Title,
			Description: todo.Description,
		}
	}

	return res, nil
}

func (r *TodoRepository) FindById(ctx context.Context, id domain.TodoID) (*domain.Todo, error) {
	todo, err := gorm.G[infrastructure.TodoEntity](r.db).Where(&infrastructure.TodoEntity{ID: id}, id).First(ctx)

	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}

	return &domain.Todo{
		ID:          todo.ID,
		Title:       todo.Title,
		Description: todo.Description,
	}, nil
}

func toEntity(todo *domain.Todo) *infrastructure.TodoEntity {
	return &infrastructure.TodoEntity{
		ID:          todo.ID,
		Title:       todo.Title,
		Description: todo.Description,
	}
}
