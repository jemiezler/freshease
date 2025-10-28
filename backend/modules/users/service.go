package users

import (
	"context"

	"github.com/google/uuid"
)

type Service interface {
	List(ctx context.Context) ([]*GetUserDTO, error)
	Get(ctx context.Context, id uuid.UUID) (*GetUserDTO, error)
	Create(ctx context.Context, dto CreateUserDTO) (*GetUserDTO, error)
	Update(ctx context.Context, id uuid.UUID, dto UpdateUserDTO) (*GetUserDTO, error)
	Delete(ctx context.Context, id uuid.UUID) error
}

type service struct {
	repo Repository
}

func NewService(r Repository) Service { return &service{repo: r} }

func (s *service) List(ctx context.Context) ([]*GetUserDTO, error) {
	return s.repo.List(ctx)
}

func (s *service) Get(ctx context.Context, id uuid.UUID) (*GetUserDTO, error) {
	return s.repo.FindByID(ctx, id)
}

func (s *service) Create(ctx context.Context, dto CreateUserDTO) (*GetUserDTO, error) {
	// pass the inbound DTO straight to the repo
	return s.repo.Create(ctx, &dto)
}

func (s *service) Update(ctx context.Context, id uuid.UUID, dto UpdateUserDTO) (*GetUserDTO, error) {
	// ensure the DTO has the ID (path param is source of truth)
	dto.ID = id
	return s.repo.Update(ctx, &dto)
}

func (s *service) Delete(ctx context.Context, id uuid.UUID) error {
	return s.repo.Delete(ctx, id)
}
