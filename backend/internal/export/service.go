package export

import (
	"context"
	"errors"

	"github.com/google/uuid"

	"github.com/qirs-mezgeb/api/internal/dashboard"
	"github.com/qirs-mezgeb/api/internal/immovable"
	"github.com/qirs-mezgeb/api/internal/models"
	"github.com/qirs-mezgeb/api/internal/movable"
)

type RecordLister interface {
	ListRecords(ctx context.Context, filters dashboard.ListFilters, userID uuid.UUID, role models.Role) (*dashboard.PaginatedRecords, error)
}

type PhotoLister interface {
	ListByRecord(ctx context.Context, recordType models.RecordType, recordID uuid.UUID) ([]models.RecordPhoto, error)
}

type Service struct {
	records   RecordLister
	immovable immovable.RepositoryInterface
	movable   movable.RepositoryInterface
	photos    PhotoLister
	mediaPath string
}

func NewService(
	records RecordLister,
	immovableRepo immovable.RepositoryInterface,
	movableRepo movable.RepositoryInterface,
	photos PhotoLister,
	mediaPath string,
) *Service {
	return &Service{
		records:   records,
		immovable: immovableRepo,
		movable:   movableRepo,
		photos:    photos,
		mediaPath: mediaPath,
	}
}

type CSVExport struct {
	Filename string
	Content  []byte
}

type PDFExport struct {
	Filename string
	Content  []byte
}

func (s *Service) ExportCSV(ctx context.Context, filters dashboard.ListFilters, userID uuid.UUID, role models.Role) (*CSVExport, error) {
	if role != models.RoleSupervisor && role != models.RoleManager {
		return nil, ErrForbidden
	}

	result, err := s.records.ListRecords(ctx, filters, userID, role)
	if err != nil {
		return nil, err
	}

	content, err := buildCSV(result.Items)
	if err != nil {
		return nil, err
	}

	return &CSVExport{
		Filename: csvFilename(),
		Content:  content,
	}, nil
}

func (s *Service) ExportPDF(
	ctx context.Context,
	recordType models.RecordType,
	recordID uuid.UUID,
	userID uuid.UUID,
	role models.Role,
) (*PDFExport, error) {
	if !recordType.IsValid() {
		return nil, ErrInvalidRecordType
	}

	switch recordType {
	case models.RecordTypeImmovable:
		return s.exportImmovablePDF(ctx, recordID, userID, role)
	case models.RecordTypeMovable:
		return s.exportMovablePDF(ctx, recordID, userID, role)
	default:
		return nil, ErrInvalidRecordType
	}
}

func (s *Service) exportImmovablePDF(ctx context.Context, recordID, userID uuid.UUID, role models.Role) (*PDFExport, error) {
	record, err := s.immovable.GetByID(ctx, recordID, userID, role)
	if err != nil {
		return nil, mapRecordError(err)
	}
	if record == nil {
		return nil, ErrRecordNotFound
	}

	if err := authorizePDFExport(record.Status, record.RegistrarID, userID, role); err != nil {
		return nil, err
	}

	photos, err := s.photos.ListByRecord(ctx, models.RecordTypeImmovable, recordID)
	if err != nil {
		return nil, err
	}

	content, err := buildImmovablePDF(record, photos, s.mediaPath)
	if err != nil {
		return nil, err
	}

	return &PDFExport{
		Filename: pdfFilename(record.RecordID),
		Content:  content,
	}, nil
}

func (s *Service) exportMovablePDF(ctx context.Context, recordID, userID uuid.UUID, role models.Role) (*PDFExport, error) {
	record, err := s.movable.GetByID(ctx, recordID, userID, role)
	if err != nil {
		return nil, mapRecordError(err)
	}
	if record == nil {
		return nil, ErrRecordNotFound
	}

	if err := authorizePDFExport(record.Status, record.RegistrarID, userID, role); err != nil {
		return nil, err
	}

	photos, err := s.photos.ListByRecord(ctx, models.RecordTypeMovable, recordID)
	if err != nil {
		return nil, err
	}

	content, err := buildMovablePDF(record, photos, s.mediaPath)
	if err != nil {
		return nil, err
	}

	return &PDFExport{
		Filename: pdfFilename(record.RecordID),
		Content:  content,
	}, nil
}

func authorizePDFExport(status models.RecordStatus, registrarID, userID uuid.UUID, role models.Role) error {
	if status == models.StatusDraft {
		return ErrDraftNotPrintable
	}

	switch role {
	case models.RoleRegistrar:
		if registrarID != userID {
			return ErrForbidden
		}
		if status != models.StatusApproved {
			return ErrNotApproved
		}
	case models.RoleSupervisor, models.RoleManager:
		return nil
	default:
		return ErrForbidden
	}

	return nil
}

func mapRecordError(err error) error {
	switch {
	case errors.Is(err, immovable.ErrRecordNotFound), errors.Is(err, movable.ErrRecordNotFound):
		return ErrRecordNotFound
	case errors.Is(err, immovable.ErrForbidden), errors.Is(err, movable.ErrForbidden):
		return ErrForbidden
	default:
		return err
	}
}
