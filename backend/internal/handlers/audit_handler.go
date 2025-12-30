package handlers

import (
	"net/http"
	"strconv"

	"admin-panel/internal/middleware"
	"admin-panel/internal/models"
	"admin-panel/internal/services"
	"admin-panel/internal/utils"
)

type AuditHandler struct {
	auditService *services.AuditService
}

func NewAuditHandler(auditService *services.AuditService) *AuditHandler {
	return &AuditHandler{
		auditService: auditService,
	}
}

func (h *AuditHandler) List(w http.ResponseWriter, r *http.Request) {
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

	if action := r.URL.Query().Get("action"); action != "" {
		params.Filters["action"] = action
	}
	if resource := r.URL.Query().Get("resource"); resource != "" {
		params.Filters["resource"] = resource
	}

	result, err := h.auditService.List(r.Context(), params)
	if err != nil {
		utils.InternalError(w, "Failed to list audit logs")
		return
	}

	utils.JSON(w, http.StatusOK, result)
}
