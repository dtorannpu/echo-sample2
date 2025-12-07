package repository

import (
	"context"
	"echo-sample2/internal/domain/todo"

	"gorm.io/gorm"
)

type todoRepository struct {
	db *gorm.DB
}

func NewTodoRepository(db *gorm.DB) todo.Repository {
	return &todoRepository{db: db}
}

func (t *todoRepository) FindAll(ctx context.Context) ([]*todo.Todo, error) {
	var models []Todo
	if err := t.db.WithContext(ctx).Find(&models).Error; err != nil {
		return nil, err
	}

	return modelsToTodos(models), nil
}

func (t *todoRepository) FindByID(ctx context.Context, id uint) (*todo.Todo, error) {
	var model Todo
	result := t.db.WithContext(ctx).First(&model, id)
	if result.Error != nil {
		return nil, result.Error
	}

	return modelToTodo(&model), nil
}

func modelsToTodos(models []Todo) []*todo.Todo {
	todos := make([]*todo.Todo, len(models))

	for i, model := range models {
		todos[i] = modelToTodo(&model)
	}

	return todos
}

func modelToTodo(model *Todo) *todo.Todo {
	return &todo.Todo{ID: model.ID, Title: model.Title, Description: model.Description}
}
