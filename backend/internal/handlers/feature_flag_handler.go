package handlers

import (
        "encoding/json"
        "net/http"
        "time"

        "admin-panel/internal/middleware"
        "admin-panel/internal/models"
        "admin-panel/internal/repository"
        "admin-panel/internal/utils"

        "github.com/go-chi/chi/v5"
        "github.com/go-playground/validator/v10"
        "github.com/google/uuid"
)

type FeatureFlagHandler struct {
        flagRepo  *repository.FeatureFlagRepository
        auditRepo *repository.AuditLogRepository
        validate  *validator.Validate
}

func NewFeatureFlagHandler(flagRepo *repository.FeatureFlagRepository, auditRepo *repository.AuditLogRepository, validate *validator.Validate) *FeatureFlagHandler {
        return &FeatureFlagHandler{
                flagRepo:  flagRepo,
                auditRepo: auditRepo,
                validate:  validate,
        }
}

func (h *FeatureFlagHandler) List(w http.ResponseWriter, r *http.Request) {
        claims := middleware.GetUserFromContext(r.Context())
        if claims == nil {
                utils.Unauthorized(w, "Not authenticated")
                return
        }

        flags, err := h.flagRepo.List(r.Context(), claims.TenantID)
        if err != nil {
                utils.InternalError(w, "Failed to list feature flags")
                return
        }

        if flags == nil {
                flags = []*models.FeatureFlag{}
        }

        utils.JSON(w, http.StatusOK, flags)
}

func (h *FeatureFlagHandler) Get(w http.ResponseWriter, r *http.Request) {
        idStr := chi.URLParam(r, "id")
        id, err := uuid.Parse(idStr)
        if err != nil {
                utils.BadRequest(w, "Invalid flag ID", nil)
                return
        }

        flag, err := h.flagRepo.GetByID(r.Context(), id)
        if err != nil {
                utils.NotFound(w, "Feature flag not found")
                return
        }

        utils.JSON(w, http.StatusOK, flag)
}

type CreateFeatureFlagRequest struct {
        Key         string  `json:"key" validate:"required,min=1,max=100"`
        Name        string  `json:"name" validate:"required,min=1,max=255"`
        Description string  `json:"description"`
        Enabled     bool    `json:"enabled"`
        Metadata    *string `json:"metadata"`
}

func (h *FeatureFlagHandler) Create(w http.ResponseWriter, r *http.Request) {
        claims := middleware.GetUserFromContext(r.Context())
        if claims == nil {
                utils.Unauthorized(w, "Not authenticated")
                return
        }

        var req CreateFeatureFlagRequest
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

        flag := &models.FeatureFlag{
                TenantID:    claims.TenantID,
                Key:         req.Key,
                Name:        req.Name,
                Description: req.Description,
                Enabled:     req.Enabled,
                Metadata:    req.Metadata,
        }

        if err := h.flagRepo.Create(r.Context(), flag); err != nil {
                utils.InternalError(w, "Failed to create feature flag")
                return
        }

        h.auditRepo.Log(r.Context(), &models.AuditLog{
                ID:         uuid.New(),
                TenantID:   claims.TenantID,
                UserID:     &claims.UserID,
                Action:     "create",
                Resource:   "feature_flag",
                ResourceID: &flag.ID,
                IPAddress:  r.RemoteAddr,
                UserAgent:  r.UserAgent(),
                CreatedAt:  time.Now(),
        })

        utils.JSON(w, http.StatusCreated, flag)
}

type UpdateFeatureFlagRequest struct {
        Name        string  `json:"name" validate:"required,min=1,max=255"`
        Description string  `json:"description"`
        Enabled     bool    `json:"enabled"`
        Metadata    *string `json:"metadata"`
}

func (h *FeatureFlagHandler) Update(w http.ResponseWriter, r *http.Request) {
        claims := middleware.GetUserFromContext(r.Context())
        if claims == nil {
                utils.Unauthorized(w, "Not authenticated")
                return
        }

        idStr := chi.URLParam(r, "id")
        id, err := uuid.Parse(idStr)
        if err != nil {
                utils.BadRequest(w, "Invalid flag ID", nil)
                return
        }

        flag, err := h.flagRepo.GetByID(r.Context(), id)
        if err != nil {
                utils.NotFound(w, "Feature flag not found")
                return
        }

        var req UpdateFeatureFlagRequest
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

        flag.Name = req.Name
        flag.Description = req.Description
        flag.Enabled = req.Enabled
        flag.Metadata = req.Metadata

        if err := h.flagRepo.Update(r.Context(), flag); err != nil {
                utils.InternalError(w, "Failed to update feature flag")
                return
        }

        h.auditRepo.Log(r.Context(), &models.AuditLog{
                ID:         uuid.New(),
                TenantID:   claims.TenantID,
                UserID:     &claims.UserID,
                Action:     "update",
                Resource:   "feature_flag",
                ResourceID: &flag.ID,
                IPAddress:  r.RemoteAddr,
                UserAgent:  r.UserAgent(),
                CreatedAt:  time.Now(),
        })

        utils.JSON(w, http.StatusOK, flag)
}

func (h *FeatureFlagHandler) Delete(w http.ResponseWriter, r *http.Request) {
        claims := middleware.GetUserFromContext(r.Context())
        if claims == nil {
                utils.Unauthorized(w, "Not authenticated")
                return
        }

        idStr := chi.URLParam(r, "id")
        id, err := uuid.Parse(idStr)
        if err != nil {
                utils.BadRequest(w, "Invalid flag ID", nil)
                return
        }

        if err := h.flagRepo.Delete(r.Context(), id); err != nil {
                utils.InternalError(w, "Failed to delete feature flag")
                return
        }

        h.auditRepo.Log(r.Context(), &models.AuditLog{
                ID:         uuid.New(),
                TenantID:   claims.TenantID,
                UserID:     &claims.UserID,
                Action:     "delete",
                Resource:   "feature_flag",
                ResourceID: &id,
                IPAddress:  r.RemoteAddr,
                UserAgent:  r.UserAgent(),
                CreatedAt:  time.Now(),
        })

        utils.JSON(w, http.StatusOK, map[string]string{"message": "Feature flag deleted"})
}

func (h *FeatureFlagHandler) Toggle(w http.ResponseWriter, r *http.Request) {
        claims := middleware.GetUserFromContext(r.Context())
        if claims == nil {
                utils.Unauthorized(w, "Not authenticated")
                return
        }

        idStr := chi.URLParam(r, "id")
        id, err := uuid.Parse(idStr)
        if err != nil {
                utils.BadRequest(w, "Invalid flag ID", nil)
                return
        }

        flag, err := h.flagRepo.GetByID(r.Context(), id)
        if err != nil {
                utils.NotFound(w, "Feature flag not found")
                return
        }

        flag.Enabled = !flag.Enabled

        if err := h.flagRepo.Update(r.Context(), flag); err != nil {
                utils.InternalError(w, "Failed to toggle feature flag")
                return
        }

        h.auditRepo.Log(r.Context(), &models.AuditLog{
                ID:         uuid.New(),
                TenantID:   claims.TenantID,
                UserID:     &claims.UserID,
                Action:     "toggle",
                Resource:   "feature_flag",
                ResourceID: &flag.ID,
                IPAddress:  r.RemoteAddr,
                UserAgent:  r.UserAgent(),
                CreatedAt:  time.Now(),
        })

        utils.JSON(w, http.StatusOK, flag)
}
