package repository

import (
	"github.com/google/uuid"
	"gorm.io/gorm"

	"github.com/gitsang/order/internal/model"
)

type CategoryRepository struct {
	db *gorm.DB
}

func NewCategoryRepository(db *gorm.DB) *CategoryRepository {
	return &CategoryRepository{db: db}
}

func (r *CategoryRepository) Create(category *model.Category) error {
	return r.db.Create(category).Error
}

func (r *CategoryRepository) FindByID(id uuid.UUID) (*model.Category, error) {
	var category model.Category
	err := r.db.Where("id = ?", id).First(&category).Error
	return &category, err
}

func (r *CategoryRepository) List() ([]model.Category, error) {
	var categories []model.Category
	err := r.db.Order("sort_order ASC, created_at DESC").Find(&categories).Error
	return categories, err
}

func (r *CategoryRepository) Update(category *model.Category) error {
	return r.db.Save(category).Error
}

func (r *CategoryRepository) Delete(id uuid.UUID) error {
	return r.db.Delete(&model.Category{}, id).Error
}
