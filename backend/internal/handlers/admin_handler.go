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

type AdminHandler struct {
        adminAuthRepo *repository.AdminAuthRepository
        auditRepo     *repository.AuditLogRepository
        validate      *validator.Validate
}

func NewAdminHandler(adminAuthRepo *repository.AdminAuthRepository, auditRepo *repository.AuditLogRepository, validate *validator.Validate) *AdminHandler {
        return &AdminHandler{
                adminAuthRepo: adminAuthRepo,
                auditRepo:     auditRepo,
                validate:      validate,
        }
}

type SetAdminRequest struct {
        Enabled  bool   `json:"enabled"`
        Password string `json:"password" validate:"required_if=Enabled true,min=8"`
}

func (h *AdminHandler) SetAdmin(w http.ResponseWriter, r *http.Request) {
        userIDStr := chi.URLParam(r, "id")
        userID, err := uuid.Parse(userIDStr)
        if err != nil {
                utils.BadRequest(w, "Invalid user ID", nil)
                return
        }

        var req SetAdminRequest
        if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
                utils.BadRequest(w, "Invalid request body", nil)
                return
        }

        if req.Enabled {
                if len(req.Password) < 8 {
                        utils.BadRequest(w, "Password must be at least 8 characters", map[string]string{"password": "min=8"})
                        return
                }

                passwordHash, err := utils.HashPasswordArgon2id(req.Password)
                if err != nil {
                        utils.InternalError(w, "Failed to hash password")
                        return
                }

                if err := h.adminAuthRepo.SetAdmin(r.Context(), userID, passwordHash); err != nil {
                        utils.InternalError(w, "Failed to set admin")
                        return
                }

                claims := middleware.GetUserFromContext(r.Context())
                if claims != nil {
                        h.auditRepo.Log(r.Context(), &models.AuditLog{
                                ID:         uuid.New(),
                                TenantID:   claims.TenantID,
                                UserID:     &claims.UserID,
                                Action:     "set_admin",
                                Resource:   "user",
                                ResourceID: &userID,
                                IPAddress:  r.RemoteAddr,
                                UserAgent:  r.UserAgent(),
                                CreatedAt:  time.Now(),
                        })
                }

                utils.JSON(w, http.StatusOK, map[string]interface{}{
                        "message": "Admin access granted",
                        "user_id": userID,
                })
        } else {
                if err := h.adminAuthRepo.UnsetAdmin(r.Context(), userID); err != nil {
                        utils.InternalError(w, "Failed to revoke admin")
                        return
                }

                claims := middleware.GetUserFromContext(r.Context())
                if claims != nil {
                        h.auditRepo.Log(r.Context(), &models.AuditLog{
                                ID:         uuid.New(),
                                TenantID:   claims.TenantID,
                                UserID:     &claims.UserID,
                                Action:     "unset_admin",
                                Resource:   "user",
                                ResourceID: &userID,
                                IPAddress:  r.RemoteAddr,
                                UserAgent:  r.UserAgent(),
                                CreatedAt:  time.Now(),
                        })
                }

                utils.JSON(w, http.StatusOK, map[string]interface{}{
                        "message": "Admin access revoked",
                        "user_id": userID,
                })
        }
}

func (h *AdminHandler) GetAdminStatus(w http.ResponseWriter, r *http.Request) {
        userIDStr := chi.URLParam(r, "id")
        userID, err := uuid.Parse(userIDStr)
        if err != nil {
                utils.BadRequest(w, "Invalid user ID", nil)
                return
        }

        adminAuth, err := h.adminAuthRepo.GetByUserID(r.Context(), userID)
        if err != nil {
                utils.JSON(w, http.StatusOK, map[string]interface{}{
                        "user_id":  userID,
                        "is_admin": false,
                })
                return
        }

        utils.JSON(w, http.StatusOK, map[string]interface{}{
                "user_id":    userID,
                "is_admin":   adminAuth.IsAdmin,
                "enabled_at": adminAuth.EnabledAt,
        })
}
