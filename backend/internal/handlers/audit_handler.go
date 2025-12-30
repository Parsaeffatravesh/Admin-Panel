package handlers

import (
        "encoding/csv"
        "fmt"
        "net/http"
        "strconv"
        "time"

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

func (h *AuditHandler) Export(w http.ResponseWriter, r *http.Request) {
        claims := middleware.GetUserFromContext(r.Context())
        if claims == nil {
                utils.Unauthorized(w, "Not authenticated")
                return
        }

        params := &models.ListParams{
                Sort:     "created_at",
                Order:    "desc",
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

        logs, err := h.auditService.ListAll(r.Context(), params)
        if err != nil {
                utils.InternalError(w, "Failed to export audit logs")
                return
        }

        filename := fmt.Sprintf("audit_logs_%s.csv", time.Now().Format("2006-01-02"))
        w.Header().Set("Content-Type", "text/csv")
        w.Header().Set("Content-Disposition", fmt.Sprintf("attachment; filename=%s", filename))

        writer := csv.NewWriter(w)
        defer writer.Flush()

        writer.Write([]string{"ID", "User ID", "Action", "Resource", "Resource ID", "IP Address", "User Agent", "Created At"})

        for _, log := range logs {
                userID := ""
                if log.UserID != nil {
                        userID = log.UserID.String()
                }
                resourceID := ""
                if log.ResourceID != nil {
                        resourceID = log.ResourceID.String()
                }

                writer.Write([]string{
                        log.ID.String(),
                        userID,
                        log.Action,
                        log.Resource,
                        resourceID,
                        log.IPAddress,
                        log.UserAgent,
                        log.CreatedAt.Format(time.RFC3339),
                })
        }
}
