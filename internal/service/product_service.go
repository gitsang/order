package service

import (
	"github.com/google/uuid"

	"github.com/gitsang/order/internal/model"
	"github.com/gitsang/order/internal/repository"
)

type ProductService struct {
	productRepo *repository.ProductRepository
}

func NewProductService(productRepo *repository.ProductRepository) *ProductService {
	return &ProductService{productRepo: productRepo}
}

func (s *ProductService) Get(id uuid.UUID) (*model.Product, error) {
	return s.productRepo.FindByID(id)
}

func (s *ProductService) List(categoryID *uuid.UUID) ([]model.Product, error) {
	return s.productRepo.List(categoryID, "active")
}

func (s *ProductService) Create(product *model.Product) error {
	return s.productRepo.Create(product)
}

func (s *ProductService) Update(product *model.Product) error {
	return s.productRepo.Update(product)
}

func (s *ProductService) Delete(id uuid.UUID) error {
	return s.productRepo.Delete(id)
}
