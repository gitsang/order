package repository

import (
	"github.com/google/uuid"
	"gorm.io/gorm"

	"github.com/gitsang/order/internal/model"
)

type ProductRepository struct {
	db *gorm.DB
}

func NewProductRepository(db *gorm.DB) *ProductRepository {
	return &ProductRepository{db: db}
}

func (r *ProductRepository) Create(product *model.Product) error {
	return r.db.Create(product).Error
}

func (r *ProductRepository) FindByID(id uuid.UUID) (*model.Product, error) {
	var product model.Product
	err := r.db.Preload("Category").Where("id = ?", id).First(&product).Error
	return &product, err
}

func (r *ProductRepository) List(categoryID *uuid.UUID, status string) ([]model.Product, error) {
	var products []model.Product
	query := r.db.Preload("Category")

	if categoryID != nil {
		query = query.Where("category_id = ?", *categoryID)
	}
	if status != "" {
		query = query.Where("status = ?", status)
	}

	err := query.Order("sort_order ASC, created_at DESC").Find(&products).Error
	return products, err
}

func (r *ProductRepository) Update(product *model.Product) error {
	return r.db.Save(product).Error
}

func (r *ProductRepository) Delete(id uuid.UUID) error {
	return r.db.Delete(&model.Product{}, id).Error
}
