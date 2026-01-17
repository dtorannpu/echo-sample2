package domain

import "context"

type Repository interface {
	Save(ctx context.Context, todo *Todo) error
}
