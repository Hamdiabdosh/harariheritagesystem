package dashboard

import (
	"context"

	"github.com/google/uuid"

	"github.com/qirs-mezgeb/api/internal/models"
)

type RepositoryInterface interface {
	GetStats(ctx context.Context, userID uuid.UUID, role models.Role) (*Stats, error)
	ListRecords(ctx context.Context, filters ListFilters, userID uuid.UUID, role models.Role) ([]RecordSummary, int, error)
}

type Service struct {
	repo RepositoryInterface
}

func NewService(repo RepositoryInterface) *Service {
	return &Service{repo: repo}
}

func (s *Service) GetStats(ctx context.Context, userID uuid.UUID, role models.Role) (*Stats, error) {
	return s.repo.GetStats(ctx, userID, role)
}

func (s *Service) ListRecords(ctx context.Context, filters ListFilters, userID uuid.UUID, role models.Role) (*PaginatedRecords, error) {
	items, total, err := s.repo.ListRecords(ctx, filters, userID, role)
	if err != nil {
		return nil, err
	}

	limit := filters.Limit
	if limit < 1 {
		limit = 20
	}
	page := filters.Page
	if page < 1 {
		page = 1
	}

	totalPages := 0
	if total > 0 {
		totalPages = (total + limit - 1) / limit
	}

	return &PaginatedRecords{
		Items:      items,
		Total:      total,
		Page:       page,
		Limit:      limit,
		TotalPages: totalPages,
	}, nil
}
