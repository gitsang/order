package handler

import (
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/google/uuid"

	"github.com/gitsang/order/internal/service"
	"github.com/gitsang/order/pkg/response"
)

type ProductHandler struct {
	productService *service.ProductService
}

func NewProductHandler(productService *service.ProductService) *ProductHandler {
	return &ProductHandler{productService: productService}
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
