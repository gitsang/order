package service

import (
	"errors"
	"fmt"
	"time"

	"github.com/google/uuid"

	"github.com/gitsang/order/internal/model"
	"github.com/gitsang/order/internal/repository"
)

type OrderService struct {
	orderRepo   *repository.OrderRepository
	productRepo *repository.ProductRepository
}

func NewOrderService(orderRepo *repository.OrderRepository, productRepo *repository.ProductRepository) *OrderService {
	return &OrderService{
		orderRepo:   orderRepo,
		productRepo: productRepo,
	}
}

type CreateOrderRequest struct {
	Items  []OrderItemRequest `json:"items"`
	Remark string             `json:"remark"`
}

type OrderItemRequest struct {
	ProductID uuid.UUID `json:"product_id"`
	Quantity  int       `json:"quantity"`
}

func (s *OrderService) Create(userID uuid.UUID, req CreateOrderRequest) (*model.Order, error) {
	if len(req.Items) == 0 {
		return nil, errors.New("order items cannot be empty")
	}

	var orderItems []model.OrderItem
	var totalAmount float64

	for _, item := range req.Items {
		if item.Quantity <= 0 {
			return nil, errors.New("quantity must be positive")
		}

		product, err := s.productRepo.FindByID(item.ProductID)
		if err != nil {
			return nil, fmt.Errorf("product not found: %s", item.ProductID)
		}

		if product.Status != "active" {
			return nil, fmt.Errorf("product is not available: %s", product.Name)
		}

		orderItem := model.OrderItem{
			ProductID: product.ID,
			Quantity:  item.Quantity,
			Price:     product.Price,
		}

		orderItems = append(orderItems, orderItem)
		totalAmount += product.Price * float64(item.Quantity)
	}

	orderNo := fmt.Sprintf("ORD%s%s", time.Now().Format("20060102150405"), uuid.New().String()[:8])

	order := &model.Order{
		UserID:      userID,
		OrderNo:     orderNo,
		TotalAmount: totalAmount,
		Status:      "pending",
		Remark:      req.Remark,
		OrderItems:  orderItems,
	}

	if err := s.orderRepo.Create(order); err != nil {
		return nil, err
	}

	return order, nil
}

func (s *OrderService) Get(id uuid.UUID) (*model.Order, error) {
	return s.orderRepo.FindByID(id)
}

func (s *OrderService) ListByUser(userID uuid.UUID, limit, offset int) ([]model.Order, error) {
	return s.orderRepo.FindByUserID(userID, limit, offset)
}

func (s *OrderService) UpdateStatus(id uuid.UUID, status string) error {
	validStatuses := map[string]bool{
		"pending":   true,
		"confirmed": true,
		"preparing": true,
		"ready":     true,
		"completed": true,
		"cancelled": true,
	}

	if !validStatuses[status] {
		return errors.New("invalid status")
	}

	return s.orderRepo.UpdateStatus(id, status)
}
