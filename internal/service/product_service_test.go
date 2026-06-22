package service

import (
	"errors"
	"testing"

	"github.com/google/uuid"

	"github.com/gitsang/order/internal/model"
)

func TestProductServiceCreate(t *testing.T) {
	categoryID := uuid.New()

	tests := []struct {
		name    string
		product *model.Product
		setup   func(*serviceTestStore)
		wantErr string
	}{
		{
			name:    "creates product",
			product: &model.Product{ID: uuid.New(), CategoryID: categoryID, Name: "Latte", Description: "Milk coffee", Price: 28.5, Status: "active", SortOrder: 1},
		},
		{
			name:    "propagates repository error",
			product: &model.Product{ID: uuid.New(), CategoryID: categoryID, Name: "Latte", Price: 28.5, Status: "active"},
			setup: func(store *serviceTestStore) {
				store.errors["product.create"] = errors.New("create failed")
			},
			wantErr: "create failed",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			store := newServiceTestStore()
			if tt.setup != nil {
				tt.setup(store)
			}
			_, productRepo, _ := newTestRepositories(t, store)
			service := NewProductService(productRepo)

			err := service.Create(tt.product)
			if tt.wantErr != "" {
				if err == nil || err.Error() != tt.wantErr {
					t.Fatalf("Create() error = %v, want %q", err, tt.wantErr)
				}
				return
			}
			if err != nil {
				t.Fatalf("Create() error = %v", err)
			}
			got, ok := store.products[tt.product.ID]
			if !ok {
				t.Fatalf("Create() did not persist product %s", tt.product.ID)
			}
			assertEqual(t, "stored name", got.Name, tt.product.Name)
			assertEqual(t, "stored price", got.Price, tt.product.Price)
		})
	}
}

func TestProductServiceGetByID(t *testing.T) {
	productID := uuid.New()
	categoryID := uuid.New()

	tests := []struct {
		name    string
		id      uuid.UUID
		seed    []model.Product
		want    *model.Product
		wantErr bool
	}{
		{
			name: "returns product",
			id:   productID,
			seed: []model.Product{{ID: productID, CategoryID: categoryID, Name: "Latte", Price: 28.5, Status: "active"}},
			want: &model.Product{ID: productID, CategoryID: categoryID, Name: "Latte", Price: 28.5, Status: "active"},
		},
		{
			name:    "returns error for missing product",
			id:      uuid.New(),
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			store := newServiceTestStore()
			for _, product := range tt.seed {
				store.products[product.ID] = product
			}
			_, productRepo, _ := newTestRepositories(t, store)
			service := NewProductService(productRepo)

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
			assertEqual(t, "product id", got.ID, tt.want.ID)
			assertEqual(t, "product name", got.Name, tt.want.Name)
			assertEqual(t, "product price", got.Price, tt.want.Price)
		})
	}
}

func TestProductServiceList(t *testing.T) {
	categoryID := uuid.New()
	otherCategoryID := uuid.New()
	activeProduct := model.Product{ID: uuid.New(), CategoryID: categoryID, Name: "Americano", Price: 20, Status: "active"}
	inactiveProduct := model.Product{ID: uuid.New(), CategoryID: categoryID, Name: "Old Latte", Price: 24, Status: "inactive"}
	otherProduct := model.Product{ID: uuid.New(), CategoryID: otherCategoryID, Name: "Mocha", Price: 30, Status: "active"}

	tests := []struct {
		name       string
		categoryID *uuid.UUID
		wantIDs    []uuid.UUID
	}{
		{
			name:    "lists active products only",
			wantIDs: []uuid.UUID{activeProduct.ID, otherProduct.ID},
		},
		{
			name:       "filters by category",
			categoryID: &categoryID,
			wantIDs:    []uuid.UUID{activeProduct.ID},
		},
		{
			name:       "empty result boundary",
			categoryID: ptrUUID(uuid.New()),
			wantIDs:    []uuid.UUID{},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			store := newServiceTestStore()
			for _, product := range []model.Product{activeProduct, inactiveProduct, otherProduct} {
				store.products[product.ID] = product
			}
			_, productRepo, _ := newTestRepositories(t, store)
			service := NewProductService(productRepo)

			got, err := service.List(tt.categoryID)
			if err != nil {
				t.Fatalf("List() error = %v", err)
			}
			gotIDs := make([]uuid.UUID, 0, len(got))
			for _, product := range got {
				if product.Status != "active" {
					t.Fatalf("List() returned inactive product %#v", product)
				}
				gotIDs = append(gotIDs, product.ID)
			}
			assertDeepEqual(t, "product ids", gotIDs, tt.wantIDs)
		})
	}
}

func TestProductServiceUpdate(t *testing.T) {
	productID := uuid.New()
	categoryID := uuid.New()

	tests := []struct {
		name    string
		input   *model.Product
		setup   func(*serviceTestStore)
		wantErr string
	}{
		{
			name:  "updates existing product",
			input: &model.Product{ID: productID, CategoryID: categoryID, Name: "New Latte", Description: "updated", Price: 31, Image: "latte.png", Status: "active", SortOrder: 3},
		},
		{
			name:  "propagates repository error",
			input: &model.Product{ID: productID, CategoryID: categoryID, Name: "New Latte", Price: 31, Status: "active"},
			setup: func(store *serviceTestStore) {
				store.errors["product.update"] = errors.New("update failed")
			},
			wantErr: "update failed",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			store := newServiceTestStore()
			store.products[productID] = model.Product{ID: productID, CategoryID: categoryID, Name: "Latte", Price: 28.5, Status: "active"}
			if tt.setup != nil {
				tt.setup(store)
			}
			_, productRepo, _ := newTestRepositories(t, store)
			service := NewProductService(productRepo)

			err := service.Update(tt.input)
			if tt.wantErr != "" {
				if err == nil || err.Error() != tt.wantErr {
					t.Fatalf("Update() error = %v, want %q", err, tt.wantErr)
				}
				return
			}
			if err != nil {
				t.Fatalf("Update() error = %v", err)
			}
			got := store.products[productID]
			assertEqual(t, "updated name", got.Name, tt.input.Name)
			assertEqual(t, "updated price", got.Price, tt.input.Price)
		})
	}
}

func TestProductServiceDelete(t *testing.T) {
	productID := uuid.New()
	categoryID := uuid.New()

	tests := []struct {
		name    string
		id      uuid.UUID
		setup   func(*serviceTestStore)
		wantErr string
	}{
		{
			name: "deletes existing product",
			id:   productID,
		},
		{
			name: "missing product is idempotent boundary",
			id:   uuid.New(),
		},
		{
			name: "propagates repository error",
			id:   productID,
			setup: func(store *serviceTestStore) {
				store.errors["product.update"] = errors.New("delete failed")
			},
			wantErr: "delete failed",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			store := newServiceTestStore()
			store.products[productID] = model.Product{ID: productID, CategoryID: categoryID, Name: "Latte", Price: 28.5, Status: "active"}
			if tt.setup != nil {
				tt.setup(store)
			}
			_, productRepo, _ := newTestRepositories(t, store)
			service := NewProductService(productRepo)

			err := service.Delete(tt.id)
			if tt.wantErr != "" {
				if err == nil || err.Error() != tt.wantErr {
					t.Fatalf("Delete() error = %v, want %q", err, tt.wantErr)
				}
				return
			}
			if err != nil {
				t.Fatalf("Delete() error = %v", err)
			}
			if tt.id == productID {
				if _, ok := store.products[productID]; ok {
					t.Fatal("Delete() did not remove product")
				}
			}
		})
	}
}

func ptrUUID(id uuid.UUID) *uuid.UUID { return &id }
