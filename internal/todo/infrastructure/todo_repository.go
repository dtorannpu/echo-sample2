package infrastructure

import (
	"context"
	"echo-sample2/internal/todo/domain"

	"gorm.io/gorm"
)

type TodoRepository struct {
	db *gorm.DB
}

func NewTodoRepository(db *gorm.DB) *TodoRepository {
	return &TodoRepository{db: db}
}

func (r *TodoRepository) getDB(ctx context.Context) *gorm.DB {
	tx, ok := ctx.Value("tx").(*gorm.DB)
	if ok && tx != nil {
		return tx
	}
	return r.db.WithContext(ctx)
}

func (r *TodoRepository) Save(ctx context.Context, todo *domain.Todo) error {
	entity := toEntity(todo)

	db := r.getDB(ctx)

	return db.Create(entity).Error
}

func toEntity(todo *domain.Todo) *TodoEntity {
	return &TodoEntity{
		ID:          todo.ID,
		Title:       todo.Title,
		Description: todo.Description,
	}
}
