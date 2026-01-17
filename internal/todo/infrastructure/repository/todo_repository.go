package repository

import (
	"context"
	"echo-sample2/internal/todo/domain"
	"echo-sample2/internal/todo/infrastructure"

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

func toEntity(todo *domain.Todo) *infrastructure.TodoEntity {
	return &infrastructure.TodoEntity{
		ID:          todo.ID,
		Title:       todo.Title,
		Description: todo.Description,
	}
}
