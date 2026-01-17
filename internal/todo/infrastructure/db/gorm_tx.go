package db

import (
	domainRepo "echo-sample2/internal/todo/domain/repository"
	"echo-sample2/internal/todo/infrastructure/repository"

	"gorm.io/gorm"
)

type GormTx struct {
	tx *gorm.DB
}

func (g *GormTx) TodoRepo() domainRepo.TodoRepository {
	return repository.NewTodoRepository(g.tx)
}
