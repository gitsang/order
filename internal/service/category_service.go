package service

import (
	"github.com/google/uuid"

	"github.com/gitsang/order/internal/model"
	"github.com/gitsang/order/internal/repository"
)

type CategoryService struct {
	categoryRepo *repository.CategoryRepository
}

func NewCategoryService(categoryRepo *repository.CategoryRepository) *CategoryService {
	return &CategoryService{categoryRepo: categoryRepo}
}

func (s *CategoryService) Get(id uuid.UUID) (*model.Category, error) {
	return s.categoryRepo.FindByID(id)
}

func (s *CategoryService) List() ([]model.Category, error) {
	return s.categoryRepo.List()
}

func (s *CategoryService) Create(category *model.Category) error {
	return s.categoryRepo.Create(category)
}

func (s *CategoryService) Update(category *model.Category) error {
	return s.categoryRepo.Update(category)
}

func (s *CategoryService) Delete(id uuid.UUID) error {
	return s.categoryRepo.Delete(id)
}
