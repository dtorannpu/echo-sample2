package todo

import "context"

type Repository interface {
	FindAll(ctx context.Context) ([]*Todo, error)
	FindByID(ctx context.Context, id uint) (*Todo, error)
}
