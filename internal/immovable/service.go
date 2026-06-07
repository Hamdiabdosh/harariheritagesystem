package immovable

import (
	"context"
	"strings"
	"time"

	"github.com/google/uuid"

	"github.com/qirs-mezgeb/api/internal/models"
)

var validCategories = map[string]struct{}{
	"I": {}, "II": {}, "VII": {}, "VIII": {},
}

type RepositoryInterface interface {
	Create(ctx context.Context, record *models.ImmovableRecord) error
	GetByID(ctx context.Context, id uuid.UUID, userID uuid.UUID, role models.Role) (*models.ImmovableRecord, error)
	List(ctx context.Context, filters ListFilters, userID uuid.UUID, role models.Role) ([]models.ImmovableRecord, int, error)
	Update(ctx context.Context, record *models.ImmovableRecord) error
	UpdateStatus(ctx context.Context, id uuid.UUID, fromStatus, toStatus models.RecordStatus, changedBy uuid.UUID, note *string) error
	FinalApprove(ctx context.Context, id uuid.UUID, fromStatus models.RecordStatus, approvedBy uuid.UUID, note *string) (time.Time, error)
}

type PhotoLister interface {
	ListByRecord(ctx context.Context, recordType models.RecordType, recordID uuid.UUID) ([]models.RecordPhoto, error)
}

type HistoryLister interface {
	ListByRecord(ctx context.Context, recordType models.RecordType, recordID uuid.UUID) ([]models.StatusHistoryEntry, error)
}

type CommentLister interface {
	ListComments(ctx context.Context, recordType models.RecordType, recordID uuid.UUID) ([]models.RecordComment, error)
}

type Service struct {
	repo          RepositoryInterface
	photoLister   PhotoLister
	historyLister HistoryLister
	commentLister CommentLister
}

func NewService(repo RepositoryInterface, photoLister PhotoLister, historyLister HistoryLister, commentLister CommentLister) *Service {
	return &Service{repo: repo, photoLister: photoLister, historyLister: historyLister, commentLister: commentLister}
}

func (s *Service) Create(ctx context.Context, registrarID uuid.UUID, input models.ImmovableRecordInput) (*CreateResult, error) {
	record := models.NewDraftImmovableRecord(registrarID)
	models.ApplyImmovableInput(&record, input)

	if err := s.repo.Create(ctx, &record); err != nil {
		return nil, err
	}

	return &CreateResult{
		ID:       record.ID,
		RecordID: record.RecordID,
		Status:   record.Status,
	}, nil
}

func (s *Service) GetByID(ctx context.Context, id, userID uuid.UUID, role models.Role) (*RecordDetail, error) {
	record, err := s.repo.GetByID(ctx, id, userID, role)
	if err != nil {
		return nil, err
	}
	if record == nil {
		return nil, ErrRecordNotFound
	}

	return &RecordDetail{
		Record:   *record,
		Photos:   listPhotos(ctx, s.photoLister, models.RecordTypeImmovable, id),
		Comments: listComments(ctx, s.commentLister, models.RecordTypeImmovable, id),
		History:  listHistory(ctx, s.historyLister, models.RecordTypeImmovable, id),
	}, nil
}

func listPhotos(ctx context.Context, lister PhotoLister, recordType models.RecordType, recordID uuid.UUID) []models.RecordPhoto {
	if lister == nil {
		return []models.RecordPhoto{}
	}
	photos, err := lister.ListByRecord(ctx, recordType, recordID)
	if err != nil {
		return []models.RecordPhoto{}
	}
	return photos
}

func listHistory(ctx context.Context, lister HistoryLister, recordType models.RecordType, recordID uuid.UUID) []models.StatusHistoryEntry {
	if lister == nil {
		return []models.StatusHistoryEntry{}
	}
	history, err := lister.ListByRecord(ctx, recordType, recordID)
	if err != nil {
		return []models.StatusHistoryEntry{}
	}
	return history
}

func listComments(ctx context.Context, lister CommentLister, recordType models.RecordType, recordID uuid.UUID) []models.RecordComment {
	if lister == nil {
		return []models.RecordComment{}
	}
	comments, err := lister.ListComments(ctx, recordType, recordID)
	if err != nil {
		return []models.RecordComment{}
	}
	return comments
}

func (s *Service) List(ctx context.Context, filters ListFilters, userID uuid.UUID, role models.Role) (*PaginatedRecords, error) {
	records, total, err := s.repo.List(ctx, filters, userID, role)
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
		Items:      records,
		Total:      total,
		Page:       page,
		Limit:      limit,
		TotalPages: totalPages,
	}, nil
}

func (s *Service) Update(ctx context.Context, id, userID uuid.UUID, role models.Role, input models.ImmovableRecordInput) (*models.ImmovableRecord, error) {
	if role != models.RoleRegistrar {
		return nil, ErrForbidden
	}

	record, err := s.repo.GetByID(ctx, id, userID, role)
	if err != nil {
		return nil, err
	}
	if record == nil {
		return nil, ErrRecordNotFound
	}
	if record.RegistrarID != userID {
		return nil, ErrForbidden
	}
	if !record.Status.IsEditable() {
		return nil, ErrNotEditable
	}

	models.ApplyImmovableInput(record, input)

	if err := s.repo.Update(ctx, record); err != nil {
		return nil, err
	}

	return record, nil
}

func (s *Service) Submit(ctx context.Context, id, userID uuid.UUID, role models.Role) (*SubmitResult, error) {
	if role != models.RoleRegistrar {
		return nil, ErrForbidden
	}

	record, err := s.repo.GetByID(ctx, id, userID, role)
	if err != nil {
		return nil, err
	}
	if record == nil {
		return nil, ErrRecordNotFound
	}
	if record.RegistrarID != userID {
		return nil, ErrForbidden
	}
	if !record.Status.IsEditable() {
		return nil, ErrNotEditable
	}

	if err := validateForSubmit(record); err != nil {
		return nil, err
	}

	fromStatus := record.Status
	toStatus := models.StatusPendingReview
	if !canTransition(fromStatus, toStatus) {
		return nil, ErrInvalidStatusTransition
	}

	if err := s.repo.UpdateStatus(ctx, id, fromStatus, toStatus, userID, nil); err != nil {
		return nil, err
	}

	return &SubmitResult{Status: toStatus}, nil
}

func (s *Service) UpdateStatus(ctx context.Context, id uuid.UUID, fromStatus, toStatus models.RecordStatus, changedBy uuid.UUID, note *string) error {
	if !canTransition(fromStatus, toStatus) {
		return ErrInvalidStatusTransition
	}
	return s.repo.UpdateStatus(ctx, id, fromStatus, toStatus, changedBy, note)
}

func canTransition(from, to models.RecordStatus) bool {
	switch from {
	case models.StatusDraft:
		return to == models.StatusPendingReview
	case models.StatusReturned:
		return to == models.StatusPendingReview
	case models.StatusPendingReview:
		return to == models.StatusUnderReview || to == models.StatusReturned
	case models.StatusUnderReview:
		return to == models.StatusApproved || to == models.StatusPendingReview
	case models.StatusApproved:
		return false
	default:
		return false
	}
}

func validateForSubmit(record *models.ImmovableRecord) error {
	fields := map[string]string{}

	if strings.TrimSpace(record.NameAmharic) == "" {
		fields["name_amharic"] = "required"
	}
	if len(record.Category) == 0 || !hasValidCategory(record.Category) {
		fields["category"] = "required"
	}
	if strings.TrimSpace(record.Woreda) == "" {
		fields["woreda"] = "required"
	}
	if strings.TrimSpace(record.Kebele) == "" {
		fields["kebele"] = "required"
	}
	if len(record.CurrentUse) == 0 {
		fields["current_use"] = "required"
	}
	if record.OwnerType == nil {
		fields["owner_type"] = "required"
	}
	if record.ConstructionPeriod == nil || strings.TrimSpace(*record.ConstructionPeriod) == "" {
		fields["construction_period"] = "required"
	}
	if record.AgeMethod == nil {
		fields["age_method"] = "required"
	}

	if len(fields) > 0 {
		return NewValidationError(fields)
	}

	return nil
}

func hasValidCategory(categories []string) bool {
	for _, category := range categories {
		if _, ok := validCategories[strings.ToUpper(strings.TrimSpace(category))]; ok {
			return true
		}
	}
	return false
}
