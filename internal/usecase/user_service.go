package usecase

import (
	"context"

	"service-app/internal/domain/user"
)

type UserService struct {
	repo user.Repository
}

func NewUserService(repo user.Repository) *UserService {
	return &UserService{repo: repo}
}

func (s *UserService) Create(ctx context.Context, input user.CreateUserInput) (user.User, error) {
	return s.repo.Create(ctx, input)
}

func (s *UserService) GetByID(ctx context.Context, id int64) (user.User, error) {
	return s.repo.GetByID(ctx, id)
}

func (s *UserService) List(ctx context.Context, limit, offset int32) ([]user.User, error) {
	return s.repo.List(ctx, limit, offset)
}

func (s *UserService) Update(ctx context.Context, id int64, input user.UpdateUserInput) (user.User, error) {
	return s.repo.Update(ctx, id, input)
}

func (s *UserService) Delete(ctx context.Context, id int64) error {
	return s.repo.Delete(ctx, id)
}
