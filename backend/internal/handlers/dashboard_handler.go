package handlers

import (
	"net/http"

	"admin-panel/internal/middleware"
	"admin-panel/internal/services"
	"admin-panel/internal/utils"
)

type DashboardHandler struct {
	dashboardService *services.DashboardService
}

func NewDashboardHandler(dashboardService *services.DashboardService) *DashboardHandler {
	return &DashboardHandler{
		dashboardService: dashboardService,
	}
}

func (h *DashboardHandler) GetStats(w http.ResponseWriter, r *http.Request) {
	claims := middleware.GetUserFromContext(r.Context())
	if claims == nil {
		utils.Unauthorized(w, "Not authenticated")
		return
	}

	stats, err := h.dashboardService.GetStats(r.Context(), claims.TenantID)
	if err != nil {
		utils.InternalError(w, "Failed to get dashboard stats")
		return
	}

	utils.JSON(w, http.StatusOK, stats)
}
