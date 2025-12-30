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

type RoleHandler struct {
	roleService *services.RoleService
	validate    *validator.Validate
}

func NewRoleHandler(roleService *services.RoleService, validate *validator.Validate) *RoleHandler {
	return &RoleHandler{
		roleService: roleService,
		validate:    validate,
	}
}

func (h *RoleHandler) Create(w http.ResponseWriter, r *http.Request) {
	claims := middleware.GetUserFromContext(r.Context())
	if claims == nil {
		utils.Unauthorized(w, "Not authenticated")
		return
	}

	var req services.CreateRoleRequest
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

	role, err := h.roleService.Create(r.Context(), &req, claims.TenantID)
	if err != nil {
		if err == services.ErrRoleNameExists {
			utils.Conflict(w, "Role name already exists")
			return
		}
		utils.InternalError(w, "Failed to create role")
		return
	}

	utils.JSON(w, http.StatusCreated, role)
}

func (h *RoleHandler) Get(w http.ResponseWriter, r *http.Request) {
	idStr := chi.URLParam(r, "id")
	id, err := uuid.Parse(idStr)
	if err != nil {
		utils.BadRequest(w, "Invalid role ID", nil)
		return
	}

	role, err := h.roleService.GetByID(r.Context(), id)
	if err != nil {
		utils.NotFound(w, "Role not found")
		return
	}

	utils.JSON(w, http.StatusOK, role)
}

func (h *RoleHandler) List(w http.ResponseWriter, r *http.Request) {
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

	result, err := h.roleService.List(r.Context(), params)
	if err != nil {
		utils.InternalError(w, "Failed to list roles")
		return
	}

	utils.JSON(w, http.StatusOK, result)
}

func (h *RoleHandler) Update(w http.ResponseWriter, r *http.Request) {
	idStr := chi.URLParam(r, "id")
	id, err := uuid.Parse(idStr)
	if err != nil {
		utils.BadRequest(w, "Invalid role ID", nil)
		return
	}

	var req services.UpdateRoleRequest
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

	role, err := h.roleService.Update(r.Context(), id, &req)
	if err != nil {
		switch err {
		case services.ErrRoleNameExists:
			utils.Conflict(w, "Role name already exists")
		case services.ErrSystemRole:
			utils.Forbidden(w, "Cannot modify system role")
		default:
			utils.InternalError(w, "Failed to update role")
		}
		return
	}

	utils.JSON(w, http.StatusOK, role)
}

func (h *RoleHandler) Delete(w http.ResponseWriter, r *http.Request) {
	idStr := chi.URLParam(r, "id")
	id, err := uuid.Parse(idStr)
	if err != nil {
		utils.BadRequest(w, "Invalid role ID", nil)
		return
	}

	if err := h.roleService.Delete(r.Context(), id); err != nil {
		if err == services.ErrSystemRole {
			utils.Forbidden(w, "Cannot delete system role")
			return
		}
		utils.InternalError(w, "Failed to delete role")
		return
	}

	utils.JSON(w, http.StatusOK, map[string]string{"message": "Role deleted successfully"})
}

func (h *RoleHandler) GetPermissions(w http.ResponseWriter, r *http.Request) {
	idStr := chi.URLParam(r, "id")
	id, err := uuid.Parse(idStr)
	if err != nil {
		utils.BadRequest(w, "Invalid role ID", nil)
		return
	}

	permissions, err := h.roleService.GetRolePermissions(r.Context(), id)
	if err != nil {
		utils.InternalError(w, "Failed to get role permissions")
		return
	}

	utils.JSON(w, http.StatusOK, permissions)
}

func (h *RoleHandler) GetAllPermissions(w http.ResponseWriter, r *http.Request) {
	permissions, err := h.roleService.GetAllPermissions(r.Context())
	if err != nil {
		utils.InternalError(w, "Failed to get permissions")
		return
	}

	utils.JSON(w, http.StatusOK, permissions)
}
