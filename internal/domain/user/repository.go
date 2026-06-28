package user

import "context"

type Repository interface {
	Create(ctx context.Context, input CreateUserInput) (User, error)
	GetByID(ctx context.Context, id int64) (User, error)
	List(ctx context.Context, limit, offset int32) ([]User, error)
	Update(ctx context.Context, id int64, input UpdateUserInput) (User, error)
	Delete(ctx context.Context, id int64) error
}
