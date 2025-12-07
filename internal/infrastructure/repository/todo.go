package repository

import (
	"echo-sample2/internal/domain/todo"

	"gorm.io/gorm"
)

type todoRepository struct {
	db *gorm.DB
}

func NewTodoRepository(db *gorm.DB) todo.TodoRepository {
	return &todoRepository{db: db}
}

// FindAll implements todo.TodoRepository.
func (t *todoRepository) FindAll() ([]todo.Todo, error) {
	panic("unimplemented")
}

// FindByID implements todo.TodoRepository.
func (t *todoRepository) FindByID(id uint) (*todo.Todo, error) {
	panic("unimplemented")
}
