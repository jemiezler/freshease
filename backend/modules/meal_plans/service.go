package meal_plans

import (
	"context"

	"github.com/google/uuid"
)

type Service interface {
	List(ctx context.Context) ([]*GetMeal_planDTO, error)
	Get(ctx context.Context, id uuid.UUID) (*GetMeal_planDTO, error)
	Create(ctx context.Context, dto CreateMeal_planDTO) (*GetMeal_planDTO, error)
	Update(ctx context.Context, id uuid.UUID, dto UpdateMeal_planDTO) (*GetMeal_planDTO, error)
	Delete(ctx context.Context, id uuid.UUID) error
}

type service struct {
	repo Repository
}

func NewService(r Repository) Service { return &service{repo: r} }

func (s *service) List(ctx context.Context) ([]*GetMeal_planDTO, error) {
	return s.repo.List(ctx)
}

func (s *service) Get(ctx context.Context, id uuid.UUID) (*GetMeal_planDTO, error) {
	return s.repo.FindByID(ctx, id)
}

func (s *service) Create(ctx context.Context, dto CreateMeal_planDTO) (*GetMeal_planDTO, error) {
	return s.repo.Create(ctx, &dto)
}

func (s *service) Update(ctx context.Context, id uuid.UUID, dto UpdateMeal_planDTO) (*GetMeal_planDTO, error) {
	dto.ID = id
	return s.repo.Update(ctx, &dto)
}

func (s *service) Delete(ctx context.Context, id uuid.UUID) error {
	return s.repo.Delete(ctx, id)
}
