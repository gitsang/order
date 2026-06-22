package service

import (
	"errors"
	"strings"
	"testing"

	"github.com/google/uuid"

	"github.com/gitsang/order/internal/model"
)

func TestOrderServiceCreate(t *testing.T) {
	userID := uuid.New()
	activeProductID := uuid.New()
	inactiveProductID := uuid.New()
	missingProductID := uuid.New()
	categoryID := uuid.New()

	tests := []struct {
		name       string
		req        CreateOrderRequest
		setup      func(*serviceTestStore)
		wantErr    string
		wantAmount float64
		wantItems  int
	}{
		{
			name: "creates order with totals",
			req: CreateOrderRequest{
				Remark: "less sugar",
				Items:  []OrderItemRequest{{ProductID: activeProductID, Quantity: 2}},
			},
			wantAmount: 57,
			wantItems:  1,
		},
		{
			name:    "rejects empty items boundary",
			req:     CreateOrderRequest{},
			wantErr: "order items cannot be empty",
		},
		{
			name:    "rejects non-positive quantity",
			req:     CreateOrderRequest{Items: []OrderItemRequest{{ProductID: activeProductID, Quantity: 0}}},
			wantErr: "quantity must be positive",
		},
		{
			name:    "rejects missing product",
			req:     CreateOrderRequest{Items: []OrderItemRequest{{ProductID: missingProductID, Quantity: 1}}},
			wantErr: "product not found",
		},
		{
			name:    "rejects inactive product",
			req:     CreateOrderRequest{Items: []OrderItemRequest{{ProductID: inactiveProductID, Quantity: 1}}},
			wantErr: "product is not available",
		},
		{
			name: "propagates repository create error",
			req:  CreateOrderRequest{Items: []OrderItemRequest{{ProductID: activeProductID, Quantity: 1}}},
			setup: func(store *serviceTestStore) {
				store.errors["order.create"] = errors.New("create order failed")
			},
			wantErr: "create order failed",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			store := newServiceTestStore()
			store.products[activeProductID] = model.Product{ID: activeProductID, CategoryID: categoryID, Name: "Latte", Price: 28.5, Status: "active"}
			store.products[inactiveProductID] = model.Product{ID: inactiveProductID, CategoryID: categoryID, Name: "Old Latte", Price: 24, Status: "inactive"}
			if tt.setup != nil {
				tt.setup(store)
			}
			_, productRepo, orderRepo := newTestRepositories(t, store)
			service := NewOrderService(orderRepo, productRepo)

			order, err := service.Create(userID, tt.req)
			if tt.wantErr != "" {
				if err == nil || !strings.Contains(err.Error(), tt.wantErr) {
					t.Fatalf("Create() error = %v, want containing %q", err, tt.wantErr)
				}
				return
			}
			if err != nil {
				t.Fatalf("Create() error = %v", err)
			}
			if order == nil {
				t.Fatal("Create() returned nil order")
			}
			assertEqual(t, "order user", order.UserID, userID)
			assertEqual(t, "order amount", order.TotalAmount, tt.wantAmount)
			assertEqual(t, "order status", order.Status, "pending")
			assertEqual(t, "order item count", len(order.OrderItems), tt.wantItems)
			if !strings.HasPrefix(order.OrderNo, "ORD") {
				t.Fatalf("Create() order number = %q", order.OrderNo)
			}
			if len(store.orders) != 1 {
				t.Fatalf("Create() stored %d orders, want 1", len(store.orders))
			}
		})
	}
}

func TestOrderServiceGetByID(t *testing.T) {
	orderID := uuid.New()
	userID := uuid.New()

	tests := []struct {
		name    string
		id      uuid.UUID
		seed    []model.Order
		wantErr bool
	}{
		{
			name: "returns order",
			id:   orderID,
			seed: []model.Order{{ID: orderID, UserID: userID, OrderNo: "ORD001", TotalAmount: 57, Status: "pending"}},
		},
		{
			name:    "returns error for missing order",
			id:      uuid.New(),
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			store := newServiceTestStore()
			for _, order := range tt.seed {
				store.orders[order.ID] = order
			}
			_, productRepo, orderRepo := newTestRepositories(t, store)
			service := NewOrderService(orderRepo, productRepo)

			got, err := service.Get(tt.id)
			if tt.wantErr {
				if err == nil {
					t.Fatal("Get() expected error")
				}
				return
			}
			if err != nil {
				t.Fatalf("Get() error = %v", err)
			}
			assertEqual(t, "order id", got.ID, orderID)
			assertEqual(t, "order no", got.OrderNo, "ORD001")
		})
	}
}

func TestOrderServiceListByUser(t *testing.T) {
	userID := uuid.New()
	otherUserID := uuid.New()
	firstOrder := model.Order{ID: uuid.New(), UserID: userID, OrderNo: "ORD001", TotalAmount: 20, Status: "pending"}
	secondOrder := model.Order{ID: uuid.New(), UserID: userID, OrderNo: "ORD002", TotalAmount: 30, Status: "completed"}
	otherOrder := model.Order{ID: uuid.New(), UserID: otherUserID, OrderNo: "ORD003", TotalAmount: 40, Status: "pending"}

	tests := []struct {
		name    string
		userID  uuid.UUID
		limit   int
		offset  int
		wantIDs []uuid.UUID
	}{
		{
			name:    "lists current user orders",
			userID:  userID,
			limit:   10,
			offset:  0,
			wantIDs: []uuid.UUID{firstOrder.ID, secondOrder.ID},
		},
		{
			name:    "empty result boundary",
			userID:  uuid.New(),
			limit:   10,
			offset:  0,
			wantIDs: []uuid.UUID{},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			store := newServiceTestStore()
			for _, order := range []model.Order{firstOrder, secondOrder, otherOrder} {
				store.orders[order.ID] = order
			}
			_, productRepo, orderRepo := newTestRepositories(t, store)
			service := NewOrderService(orderRepo, productRepo)

			got, err := service.ListByUser(tt.userID, tt.limit, tt.offset)
			if err != nil {
				t.Fatalf("ListByUser() error = %v", err)
			}
			gotIDs := make([]uuid.UUID, 0, len(got))
			for _, order := range got {
				gotIDs = append(gotIDs, order.ID)
			}
			assertDeepEqual(t, "order ids", gotIDs, tt.wantIDs)
		})
	}
}

func TestOrderServiceUpdateStatus(t *testing.T) {
	orderID := uuid.New()
	userID := uuid.New()

	tests := []struct {
		name       string
		status     string
		setup      func(*serviceTestStore)
		wantErr    string
		wantStatus string
	}{
		{
			name:       "updates to confirmed",
			status:     "confirmed",
			wantStatus: "confirmed",
		},
		{
			name:       "allows cancelled boundary status",
			status:     "cancelled",
			wantStatus: "cancelled",
		},
		{
			name:    "rejects invalid status",
			status:  "unknown",
			wantErr: "invalid status",
		},
		{
			name:   "propagates repository error",
			status: "ready",
			setup: func(store *serviceTestStore) {
				store.errors["order.update_status"] = errors.New("update status failed")
			},
			wantErr: "update status failed",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			store := newServiceTestStore()
			store.orders[orderID] = model.Order{ID: orderID, UserID: userID, OrderNo: "ORD001", TotalAmount: 57, Status: "pending"}
			if tt.setup != nil {
				tt.setup(store)
			}
			_, productRepo, orderRepo := newTestRepositories(t, store)
			service := NewOrderService(orderRepo, productRepo)

			err := service.UpdateStatus(orderID, tt.status)
			if tt.wantErr != "" {
				if err == nil || err.Error() != tt.wantErr {
					t.Fatalf("UpdateStatus() error = %v, want %q", err, tt.wantErr)
				}
				return
			}
			if err != nil {
				t.Fatalf("UpdateStatus() error = %v", err)
			}
			assertEqual(t, "order status", store.orders[orderID].Status, tt.wantStatus)
		})
	}
}
