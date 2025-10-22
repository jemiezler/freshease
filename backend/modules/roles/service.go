package roles

import (
	"context"

	"github.com/google/uuid"
)

type Service interface {
	List(ctx context.Context) ([]*Role, error)
	Get(ctx context.Context, id uuid.UUID) (*Role, error)
	Create(ctx context.Context, dto CreateRoleDTO) (*Role, error)
	Update(ctx context.Context, id uuid.UUID, dto UpdateRoleDTO) (*Role, error)
	Delete(ctx context.Context, id uuid.UUID) error
}

type service struct{ repo Repository }

func NewService(r Repository) Service { return &service{repo: r} }

func (s *service) List(ctx context.Context) ([]*Role, error) {
	return s.repo.List(ctx)
}

func (s *service) Get(ctx context.Context, id uuid.UUID) (*Role, error) {
	return s.repo.FindByID(ctx, id)
}

func (s *service) Create(ctx context.Context, dto CreateRoleDTO) (*Role, error) {
	entity := &Role{
		Email: dto.Email,
		Name:  dto.Name,
	}
	if err := s.repo.Create(ctx, entity); err != nil {
		return nil, err
	}
	return entity, nil
}

func (s *service) Update(ctx context.Context, id uuid.UUID, dto UpdateRoleDTO) (*Role, error) {
	entity, err := s.repo.FindByID(ctx, id)
	if err != nil { return nil, err }
	if dto.Email != nil { entity.Email = *dto.Email }
	if dto.Name  != nil { entity.Name  = *dto.Name  }
	if err := s.repo.Update(ctx, entity); err != nil {
		return nil, err
	}
	return entity, nil
}

func (s *service) Delete(ctx context.Context, id uuid.UUID) error {
	return s.repo.Delete(ctx, id)
}
