package roles

import (
	"context"

	"github.com/google/uuid"
)

type Service interface {
	List(ctx context.Context) ([]*GetRoleDTO, error)
	Get(ctx context.Context, id uuid.UUID) (*GetRoleDTO, error)
	Create(ctx context.Context, dto CreateRoleDTO) (*GetRoleDTO, error)
	Update(ctx context.Context, id uuid.UUID, dto UpdateRoleDTO) (*GetRoleDTO, error)
	Delete(ctx context.Context, id uuid.UUID) error
}

type service struct {
	repo Repository
}

func NewService(r Repository) Service { return &service{repo: r} }

func (s *service) List(ctx context.Context) ([]*GetRoleDTO, error) {
	return s.repo.List(ctx)
}

func (s *service) Get(ctx context.Context, id uuid.UUID) (*GetRoleDTO, error) {
	return s.repo.FindByID(ctx, id)
}

func (s *service) Create(ctx context.Context, dto CreateRoleDTO) (*GetRoleDTO, error) {
	// pass the inbound DTO straight to the repo
	return s.repo.Create(ctx, &dto)
}

func (s *service) Update(ctx context.Context, id uuid.UUID, dto UpdateRoleDTO) (*GetRoleDTO, error) {
	// ensure the DTO has the ID (path param is source of truth)
	dto.ID = id
	return s.repo.Update(ctx, &dto)
}

func (s *service) Delete(ctx context.Context, id uuid.UUID) error {
	return s.repo.Delete(ctx, id)
}
