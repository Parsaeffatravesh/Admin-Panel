package services

import (
        "context"

        "admin-panel/internal/models"
        "admin-panel/internal/repository"
)

type AuditService struct {
        auditRepo *repository.AuditLogRepository
}

func NewAuditService(auditRepo *repository.AuditLogRepository) *AuditService {
        return &AuditService{
                auditRepo: auditRepo,
        }
}

func (s *AuditService) List(ctx context.Context, params *models.ListParams) (*models.PaginatedResponse, error) {
        logs, total, err := s.auditRepo.List(ctx, params)
        if err != nil {
                return nil, err
        }

        totalPages := int(total) / params.PerPage
        if int(total)%params.PerPage > 0 {
                totalPages++
        }

        return &models.PaginatedResponse{
                Data:       logs,
                Total:      total,
                Page:       params.Page,
                PerPage:    params.PerPage,
                TotalPages: totalPages,
        }, nil
}

func (s *AuditService) ListAll(ctx context.Context, params *models.ListParams) ([]*models.AuditLog, error) {
        return s.auditRepo.ListAllForExport(ctx, params)
}
