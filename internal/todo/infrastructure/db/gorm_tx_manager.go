package db

import (
	"context"
	"echo-sample2/internal/todo/domain/repository"

	"gorm.io/gorm"
)

type GormTxManager struct {
	db *gorm.DB
}

func NewGormTxManager(db *gorm.DB) *GormTxManager {
	return &GormTxManager{db: db}
}

func (g GormTxManager) Do(ctx context.Context, fn func(tx repository.Transaction) error) error {
	return g.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		return fn(&GormTx{tx: tx})
	})
}
