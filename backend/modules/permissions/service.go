package permissions

import (
	"context"

	"github.com/google/uuid"
)

type Service interface {
	List(ctx context.Context) ([]*Permission, error)
	Get(ctx context.Context, id uuid.UUID) (*Permission, error)
	Create(ctx context.Context, dto CreatePermissionDTO) (*Permission, error)
	Update(ctx context.Context, id uuid.UUID, dto UpdatePermissionDTO) (*Permission, error)
	Delete(ctx context.Context, id uuid.UUID) error
}

type service struct{ repo Repository }

func NewService(r Repository) Service { return &service{repo: r} }

func (s *service) List(ctx context.Context) ([]*Permission, error) {
	return s.repo.List(ctx)
}

func (s *service) Get(ctx context.Context, id uuid.UUID) (*Permission, error) {
	return s.repo.FindByID(ctx, id)
}

func (s *service) Create(ctx context.Context, dto CreatePermissionDTO) (*Permission, error) {
	entity := &Permission{
		Email: dto.Email,
		Name:  dto.Name,
	}
	if err := s.repo.Create(ctx, entity); err != nil {
		return nil, err
	}
	return entity, nil
}

func (s *service) Update(ctx context.Context, id uuid.UUID, dto UpdatePermissionDTO) (*Permission, error) {
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
