package repository

import (
	"github.com/google/uuid"
	"gorm.io/gorm"

	"github.com/gitsang/order/internal/model"
)

type OrderRepository struct {
	db *gorm.DB
}

func NewOrderRepository(db *gorm.DB) *OrderRepository {
	return &OrderRepository{db: db}
}

func (r *OrderRepository) Create(order *model.Order) error {
	return r.db.Create(order).Error
}

func (r *OrderRepository) FindByID(id uuid.UUID) (*model.Order, error) {
	var order model.Order
	err := r.db.Preload("OrderItems").Preload("OrderItems.Product").Where("id = ?", id).First(&order).Error
	return &order, err
}

func (r *OrderRepository) FindByUserID(userID uuid.UUID, limit, offset int) ([]model.Order, error) {
	var orders []model.Order
	err := r.db.Preload("OrderItems").Preload("OrderItems.Product").
		Where("user_id = ?", userID).
		Order("created_at DESC").
		Limit(limit).Offset(offset).
		Find(&orders).Error
	return orders, err
}

func (r *OrderRepository) Update(order *model.Order) error {
	return r.db.Save(order).Error
}

func (r *OrderRepository) UpdateStatus(id uuid.UUID, status string) error {
	return r.db.Model(&model.Order{}).Where("id = ?", id).Update("status", status).Error
}
