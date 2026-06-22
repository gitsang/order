package handler

import (
	"encoding/json"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/google/uuid"

	"github.com/gitsang/order/internal/model"
	"github.com/gitsang/order/internal/service"
	"github.com/gitsang/order/pkg/response"
)

type ProductHandler struct {
	productService *service.ProductService
}

func NewProductHandler(productService *service.ProductService) *ProductHandler {
	return &ProductHandler{productService: productService}
}

type CreateProductRequest struct {
	CategoryID  uuid.UUID `json:"category_id"`
	Name        string    `json:"name"`
	Description string    `json:"description"`
	Price       float64   `json:"price"`
	Image       string    `json:"image"`
	Status      string    `json:"status"`
	SortOrder   int       `json:"sort_order"`
}

type UpdateProductRequest struct {
	CategoryID  *uuid.UUID `json:"category_id"`
	Name        *string    `json:"name"`
	Description *string    `json:"description"`
	Price       *float64   `json:"price"`
	Image       *string    `json:"image"`
	Status      *string    `json:"status"`
	SortOrder   *int       `json:"sort_order"`
}

func (h *ProductHandler) Get(w http.ResponseWriter, r *http.Request) {
	idStr := chi.URLParam(r, "id")
	id, err := uuid.Parse(idStr)
	if err != nil {
		response.Error(w, http.StatusBadRequest, "Invalid product ID")
		return
	}

	product, err := h.productService.Get(id)
	if err != nil {
		response.Error(w, http.StatusNotFound, "Product not found")
		return
	}

	response.Success(w, product)
}

func (h *ProductHandler) List(w http.ResponseWriter, r *http.Request) {
	var categoryID *uuid.UUID
	categoryIDStr := r.URL.Query().Get("category_id")
	if categoryIDStr != "" {
		id, err := uuid.Parse(categoryIDStr)
		if err != nil {
			response.Error(w, http.StatusBadRequest, "Invalid category ID")
			return
		}
		categoryID = &id
	}

	products, err := h.productService.List(categoryID)
	if err != nil {
		response.Error(w, http.StatusInternalServerError, "Failed to list products")
		return
	}

	response.Success(w, products)
}

func (h *ProductHandler) Create(w http.ResponseWriter, r *http.Request) {
	var req CreateProductRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		response.Error(w, http.StatusBadRequest, "Invalid request body")
		return
	}

	if req.CategoryID == uuid.Nil {
		response.Error(w, http.StatusBadRequest, "category_id is required")
		return
	}
	if req.Name == "" {
		response.Error(w, http.StatusBadRequest, "name is required")
		return
	}
	if req.Price <= 0 {
		response.Error(w, http.StatusBadRequest, "price must be positive")
		return
	}
	if req.Status == "" {
		req.Status = "active"
	}

	product := &model.Product{
		CategoryID:  req.CategoryID,
		Name:        req.Name,
		Description: req.Description,
		Price:       req.Price,
		Image:       req.Image,
		Status:      req.Status,
		SortOrder:   req.SortOrder,
	}

	if err := h.productService.Create(product); err != nil {
		response.Error(w, http.StatusInternalServerError, "Failed to create product")
		return
	}

	response.Success(w, product)
}

func (h *ProductHandler) Update(w http.ResponseWriter, r *http.Request) {
	idStr := chi.URLParam(r, "id")
	id, err := uuid.Parse(idStr)
	if err != nil {
		response.Error(w, http.StatusBadRequest, "Invalid product ID")
		return
	}

	var req UpdateProductRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		response.Error(w, http.StatusBadRequest, "Invalid request body")
		return
	}

	product, err := h.productService.Get(id)
	if err != nil {
		response.Error(w, http.StatusNotFound, "Product not found")
		return
	}

	if req.CategoryID != nil {
		if *req.CategoryID == uuid.Nil {
			response.Error(w, http.StatusBadRequest, "category_id is required")
			return
		}
		product.CategoryID = *req.CategoryID
	}
	if req.Name != nil {
		if *req.Name == "" {
			response.Error(w, http.StatusBadRequest, "name is required")
			return
		}
		product.Name = *req.Name
	}
	if req.Description != nil {
		product.Description = *req.Description
	}
	if req.Price != nil {
		if *req.Price <= 0 {
			response.Error(w, http.StatusBadRequest, "price must be positive")
			return
		}
		product.Price = *req.Price
	}
	if req.Image != nil {
		product.Image = *req.Image
	}
	if req.Status != nil {
		if *req.Status == "" {
			response.Error(w, http.StatusBadRequest, "status is required")
			return
		}
		product.Status = *req.Status
	}
	if req.SortOrder != nil {
		product.SortOrder = *req.SortOrder
	}

	if err := h.productService.Update(product); err != nil {
		response.Error(w, http.StatusInternalServerError, "Failed to update product")
		return
	}

	response.Success(w, product)
}

func (h *ProductHandler) Delete(w http.ResponseWriter, r *http.Request) {
	idStr := chi.URLParam(r, "id")
	id, err := uuid.Parse(idStr)
	if err != nil {
		response.Error(w, http.StatusBadRequest, "Invalid product ID")
		return
	}

	if _, err := h.productService.Get(id); err != nil {
		response.Error(w, http.StatusNotFound, "Product not found")
		return
	}

	if err := h.productService.Delete(id); err != nil {
		response.Error(w, http.StatusInternalServerError, "Failed to delete product")
		return
	}

	response.Success(w, nil)
}
