package product_categories

import (
	"context"

	"github.com/google/uuid"
)

type Service interface {
	List(ctx context.Context) ([]*GetProductCategoryDTO, error)
	Get(ctx context.Context, id uuid.UUID) (*GetProductCategoryDTO, error)
	Create(ctx context.Context, dto CreateProductCategoryDTO) (*GetProductCategoryDTO, error)
	Update(ctx context.Context, id uuid.UUID, dto UpdateProductCategoryDTO) (*GetProductCategoryDTO, error)
	Delete(ctx context.Context, id uuid.UUID) error
}

type service struct {
	repo Repository
}

func NewService(r Repository) Service { return &service{repo: r} }

func (s *service) List(ctx context.Context) ([]*GetProductCategoryDTO, error) {
	return s.repo.List(ctx)
}

func (s *service) Get(ctx context.Context, id uuid.UUID) (*GetProductCategoryDTO, error) {
	return s.repo.FindByID(ctx, id)
}

func (s *service) Create(ctx context.Context, dto CreateProductCategoryDTO) (*GetProductCategoryDTO, error) {
	return s.repo.Create(ctx, &dto)
}

func (s *service) Update(ctx context.Context, id uuid.UUID, dto UpdateProductCategoryDTO) (*GetProductCategoryDTO, error) {
	dto.ID = id
	return s.repo.Update(ctx, &dto)
}

func (s *service) Delete(ctx context.Context, id uuid.UUID) error {
	return s.repo.Delete(ctx, id)
}
