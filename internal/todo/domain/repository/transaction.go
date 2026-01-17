package repository

import (
	"context"
)

type Transaction interface {
	TodoRepo() TodoRepository
}

type TransactionManager interface {
	Do(ctx context.Context, fn func(tx Transaction) error) error
}
