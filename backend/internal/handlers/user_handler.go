package handlers

import (
	"encoding/json"
	"net/http"
	"strconv"

	"admin-panel/internal/middleware"
	"admin-panel/internal/models"
	"admin-panel/internal/services"
	"admin-panel/internal/utils"

	"github.com/go-chi/chi/v5"
	"github.com/go-playground/validator/v10"
	"github.com/google/uuid"
)

type UserHandler struct {
	userService *services.UserService
	validate    *validator.Validate
}

func NewUserHandler(userService *services.UserService, validate *validator.Validate) *UserHandler {
	return &UserHandler{
		userService: userService,
		validate:    validate,
	}
}

func (h *UserHandler) Create(w http.ResponseWriter, r *http.Request) {
	claims := middleware.GetUserFromContext(r.Context())
	if claims == nil {
		utils.Unauthorized(w, "Not authenticated")
		return
	}

	var req services.CreateUserRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		utils.BadRequest(w, "Invalid request body", nil)
		return
	}

	if err := h.validate.Struct(&req); err != nil {
		details := make(map[string]string)
		for _, e := range err.(validator.ValidationErrors) {
			details[e.Field()] = e.Tag()
		}
		utils.BadRequest(w, "Validation failed", details)
		return
	}

	user, err := h.userService.Create(r.Context(), &req, claims.TenantID)
	if err != nil {
		if err == services.ErrEmailExists {
			utils.Conflict(w, "Email already exists")
			return
		}
		utils.InternalError(w, "Failed to create user")
		return
	}

	utils.JSON(w, http.StatusCreated, user)
}

func (h *UserHandler) Get(w http.ResponseWriter, r *http.Request) {
	idStr := chi.URLParam(r, "id")
	id, err := uuid.Parse(idStr)
	if err != nil {
		utils.BadRequest(w, "Invalid user ID", nil)
		return
	}

	user, err := h.userService.GetByID(r.Context(), id)
	if err != nil {
		utils.NotFound(w, "User not found")
		return
	}

	utils.JSON(w, http.StatusOK, user)
}

func (h *UserHandler) List(w http.ResponseWriter, r *http.Request) {
	claims := middleware.GetUserFromContext(r.Context())
	if claims == nil {
		utils.Unauthorized(w, "Not authenticated")
		return
	}

	page, _ := strconv.Atoi(r.URL.Query().Get("page"))
	if page < 1 {
		page = 1
	}

	perPage, _ := strconv.Atoi(r.URL.Query().Get("per_page"))
	if perPage < 1 || perPage > 100 {
		perPage = 20
	}

	params := &models.ListParams{
		Page:     page,
		PerPage:  perPage,
		Sort:     r.URL.Query().Get("sort"),
		Order:    r.URL.Query().Get("order"),
		Search:   r.URL.Query().Get("search"),
		Filters:  make(map[string]interface{}),
		TenantID: claims.TenantID,
	}

	if status := r.URL.Query().Get("status"); status != "" {
		params.Filters["status"] = status
	}

	result, err := h.userService.List(r.Context(), params)
	if err != nil {
		utils.InternalError(w, "Failed to list users")
		return
	}

	utils.JSON(w, http.StatusOK, result)
}

func (h *UserHandler) Update(w http.ResponseWriter, r *http.Request) {
	idStr := chi.URLParam(r, "id")
	id, err := uuid.Parse(idStr)
	if err != nil {
		utils.BadRequest(w, "Invalid user ID", nil)
		return
	}

	var req services.UpdateUserRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		utils.BadRequest(w, "Invalid request body", nil)
		return
	}

	if err := h.validate.Struct(&req); err != nil {
		details := make(map[string]string)
		for _, e := range err.(validator.ValidationErrors) {
			details[e.Field()] = e.Tag()
		}
		utils.BadRequest(w, "Validation failed", details)
		return
	}

	user, err := h.userService.Update(r.Context(), id, &req)
	if err != nil {
		if err == services.ErrEmailExists {
			utils.Conflict(w, "Email already exists")
			return
		}
		utils.InternalError(w, "Failed to update user")
		return
	}

	utils.JSON(w, http.StatusOK, user)
}

func (h *UserHandler) Delete(w http.ResponseWriter, r *http.Request) {
	idStr := chi.URLParam(r, "id")
	id, err := uuid.Parse(idStr)
	if err != nil {
		utils.BadRequest(w, "Invalid user ID", nil)
		return
	}

	if err := h.userService.Delete(r.Context(), id); err != nil {
		utils.InternalError(w, "Failed to delete user")
		return
	}

	utils.JSON(w, http.StatusOK, map[string]string{"message": "User deleted successfully"})
}

func (h *UserHandler) ResetPassword(w http.ResponseWriter, r *http.Request) {
	idStr := chi.URLParam(r, "id")
	id, err := uuid.Parse(idStr)
	if err != nil {
		utils.BadRequest(w, "Invalid user ID", nil)
		return
	}

	var req services.ResetPasswordRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		utils.BadRequest(w, "Invalid request body", nil)
		return
	}

	if err := h.validate.Struct(&req); err != nil {
		utils.BadRequest(w, "Password must be at least 8 characters", nil)
		return
	}

	if err := h.userService.ResetPassword(r.Context(), id, &req); err != nil {
		utils.InternalError(w, "Failed to reset password")
		return
	}

	utils.JSON(w, http.StatusOK, map[string]string{"message": "Password reset successfully"})
}

func (h *UserHandler) GetRoles(w http.ResponseWriter, r *http.Request) {
	idStr := chi.URLParam(r, "id")
	id, err := uuid.Parse(idStr)
	if err != nil {
		utils.BadRequest(w, "Invalid user ID", nil)
		return
	}

	roles, err := h.userService.GetUserRoles(r.Context(), id)
	if err != nil {
		utils.InternalError(w, "Failed to get user roles")
		return
	}

	utils.JSON(w, http.StatusOK, roles)
}
