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

type CategoryHandler struct {
	categoryService *service.CategoryService
}

func NewCategoryHandler(categoryService *service.CategoryService) *CategoryHandler {
	return &CategoryHandler{categoryService: categoryService}
}

type CreateCategoryRequest struct {
	Name      string `json:"name"`
	SortOrder int    `json:"sort_order"`
}

type UpdateCategoryRequest struct {
	Name      *string `json:"name"`
	SortOrder *int    `json:"sort_order"`
}

func (h *CategoryHandler) List(w http.ResponseWriter, r *http.Request) {
	categories, err := h.categoryService.List()
	if err != nil {
		response.Error(w, http.StatusInternalServerError, "Failed to list categories")
		return
	}

	response.Success(w, categories)
}

func (h *CategoryHandler) Create(w http.ResponseWriter, r *http.Request) {
	var req CreateCategoryRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		response.Error(w, http.StatusBadRequest, "Invalid request body")
		return
	}

	if req.Name == "" {
		response.Error(w, http.StatusBadRequest, "name is required")
		return
	}

	category := &model.Category{
		Name:      req.Name,
		SortOrder: req.SortOrder,
	}

	if err := h.categoryService.Create(category); err != nil {
		response.Error(w, http.StatusInternalServerError, "Failed to create category")
		return
	}

	response.Success(w, category)
}

func (h *CategoryHandler) Update(w http.ResponseWriter, r *http.Request) {
	idStr := chi.URLParam(r, "id")
	id, err := uuid.Parse(idStr)
	if err != nil {
		response.Error(w, http.StatusBadRequest, "Invalid category ID")
		return
	}

	var req UpdateCategoryRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		response.Error(w, http.StatusBadRequest, "Invalid request body")
		return
	}

	category, err := h.categoryService.Get(id)
	if err != nil {
		response.Error(w, http.StatusNotFound, "Category not found")
		return
	}

	if req.Name != nil {
		if *req.Name == "" {
			response.Error(w, http.StatusBadRequest, "name is required")
			return
		}
		category.Name = *req.Name
	}
	if req.SortOrder != nil {
		category.SortOrder = *req.SortOrder
	}

	if err := h.categoryService.Update(category); err != nil {
		response.Error(w, http.StatusInternalServerError, "Failed to update category")
		return
	}

	response.Success(w, category)
}

func (h *CategoryHandler) Delete(w http.ResponseWriter, r *http.Request) {
	idStr := chi.URLParam(r, "id")
	id, err := uuid.Parse(idStr)
	if err != nil {
		response.Error(w, http.StatusBadRequest, "Invalid category ID")
		return
	}

	if _, err := h.categoryService.Get(id); err != nil {
		response.Error(w, http.StatusNotFound, "Category not found")
		return
	}

	if err := h.categoryService.Delete(id); err != nil {
		response.Error(w, http.StatusInternalServerError, "Failed to delete category")
		return
	}

	response.Success(w, nil)
}
