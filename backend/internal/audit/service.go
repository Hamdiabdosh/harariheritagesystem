package audit

import (
	"context"

	"github.com/google/uuid"

	"github.com/qirs-mezgeb/api/internal/models"
)

type RepositoryInterface interface {
	ListByRecord(ctx context.Context, recordType models.RecordType, recordID uuid.UUID) ([]models.StatusHistoryEntry, error)
}

type Service struct {
	repo RepositoryInterface
}

func NewService(repo RepositoryInterface) *Service {
	return &Service{repo: repo}
}

func (s *Service) ListByRecord(ctx context.Context, recordType models.RecordType, recordID uuid.UUID) ([]models.StatusHistoryEntry, error) {
	return s.repo.ListByRecord(ctx, recordType, recordID)
}
