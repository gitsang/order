package handler

import (
	"encoding/json"
	"net/http"

	"github.com/gitsang/order/internal/service"
	"github.com/gitsang/order/pkg/response"
)

type AuthHandler struct {
	authService *service.AuthService
}

func NewAuthHandler(authService *service.AuthService) *AuthHandler {
	return &AuthHandler{authService: authService}
}

type LoginRequest struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

type RegisterRequest struct {
	Username string `json:"username"`
	Password string `json:"password"`
	Name     string `json:"name"`
	Phone    string `json:"phone"`
}

func (h *AuthHandler) Login(w http.ResponseWriter, r *http.Request) {
	var req LoginRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		response.Error(w, http.StatusBadRequest, "Invalid request body")
		return
	}

	token, err := h.authService.Login(req.Username, req.Password)
	if err != nil {
		response.Error(w, http.StatusUnauthorized, "Invalid credentials")
		return
	}

	response.Success(w, map[string]string{"token": token})
}

func (h *AuthHandler) Register(w http.ResponseWriter, r *http.Request) {
	var req RegisterRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		response.Error(w, http.StatusBadRequest, "Invalid request body")
		return
	}

	user, err := h.authService.Register(req.Username, req.Password, req.Name, req.Phone)
	if err != nil {
		response.Error(w, http.StatusBadRequest, err.Error())
		return
	}

	response.Success(w, map[string]string{
		"id":       user.ID.String(),
		"username": user.Username,
	})
}

func (h *AuthHandler) GetMe(w http.ResponseWriter, r *http.Request) {
	userID := GetUserID(r)
	username := GetUsername(r)
	role := GetRole(r)

	response.Success(w, map[string]string{
		"id":       userID.String(),
		"user_id":  userID.String(),
		"username": username,
		"role":     role,
	})
}
