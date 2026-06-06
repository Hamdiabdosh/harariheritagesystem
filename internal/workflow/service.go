package workflow

import (
	"context"
	"errors"
	"strings"
	"time"

	"github.com/google/uuid"

	"github.com/qirs-mezgeb/api/internal/immovable"
	"github.com/qirs-mezgeb/api/internal/models"
	"github.com/qirs-mezgeb/api/internal/movable"
)

type recordMeta struct {
	Status      models.RecordStatus
	RegistrarID uuid.UUID
}

type ImmovableStore interface {
	GetByID(ctx context.Context, id uuid.UUID, userID uuid.UUID, role models.Role) (*models.ImmovableRecord, error)
	UpdateStatus(ctx context.Context, id uuid.UUID, fromStatus, toStatus models.RecordStatus, changedBy uuid.UUID, note *string) error
	FinalApprove(ctx context.Context, id uuid.UUID, fromStatus models.RecordStatus, approvedBy uuid.UUID, note *string) (time.Time, error)
}

type MovableStore interface {
	GetByID(ctx context.Context, id uuid.UUID, userID uuid.UUID, role models.Role) (*models.MovableRecord, error)
	UpdateStatus(ctx context.Context, id uuid.UUID, fromStatus, toStatus models.RecordStatus, changedBy uuid.UUID, note *string) error
	FinalApprove(ctx context.Context, id uuid.UUID, fromStatus models.RecordStatus, approvedBy uuid.UUID, note *string) (time.Time, error)
}

type RepositoryInterface interface {
	CreateComment(ctx context.Context, comment *models.RecordComment) error
	ListComments(ctx context.Context, recordType models.RecordType, recordID uuid.UUID) ([]models.RecordComment, error)
}

type AuditReader interface {
	ListByRecord(ctx context.Context, recordType models.RecordType, recordID uuid.UUID) ([]models.StatusHistoryEntry, error)
}

type Service struct {
	repo      RepositoryInterface
	audit     AuditReader
	immovable ImmovableStore
	movable   MovableStore
}

func NewService(repo RepositoryInterface, audit AuditReader, immovable ImmovableStore, movable MovableStore) *Service {
	return &Service{repo: repo, audit: audit, immovable: immovable, movable: movable}
}

type StatusResult struct {
	Status     models.RecordStatus `json:"status"`
	ApprovedAt *time.Time          `json:"approved_at,omitempty"`
}

func (s *Service) ReviewApprove(ctx context.Context, recordType models.RecordType, recordID, supervisorID uuid.UUID, comment *string) (*StatusResult, error) {
	meta, err := s.getRecordMeta(ctx, recordType, recordID, supervisorID, models.RoleSupervisor)
	if err != nil {
		return nil, err
	}
	if meta.Status != models.StatusPendingReview {
		return nil, ErrInvalidStatusTransition
	}

	note := optionalCommentNote(comment)
	if err := s.updateStatus(ctx, recordType, recordID, meta.Status, models.StatusUnderReview, supervisorID, note); err != nil {
		return nil, mapRecordError(err)
	}

	if trimmed := trimComment(comment); trimmed != "" {
		if err := s.addComment(ctx, recordType, recordID, supervisorID, trimmed); err != nil {
			return nil, err
		}
	}

	return &StatusResult{Status: models.StatusUnderReview}, nil
}

func (s *Service) ReviewReturn(ctx context.Context, recordType models.RecordType, recordID, supervisorID uuid.UUID, comment string) (*StatusResult, error) {
	if trimComment(&comment) == "" {
		return nil, ErrCommentRequired
	}

	meta, err := s.getRecordMeta(ctx, recordType, recordID, supervisorID, models.RoleSupervisor)
	if err != nil {
		return nil, err
	}
	if meta.Status != models.StatusPendingReview {
		return nil, ErrInvalidStatusTransition
	}

	text := trimComment(&comment)
	note := &text
	if err := s.updateStatus(ctx, recordType, recordID, meta.Status, models.StatusReturned, supervisorID, note); err != nil {
		return nil, mapRecordError(err)
	}

	if err := s.addComment(ctx, recordType, recordID, supervisorID, text); err != nil {
		return nil, err
	}

	return &StatusResult{Status: models.StatusReturned}, nil
}

func (s *Service) FinalApprove(ctx context.Context, recordType models.RecordType, recordID, managerID uuid.UUID, comment *string) (*StatusResult, error) {
	meta, err := s.getRecordMeta(ctx, recordType, recordID, managerID, models.RoleManager)
	if err != nil {
		return nil, err
	}
	if meta.Status != models.StatusUnderReview {
		return nil, ErrInvalidStatusTransition
	}

	note := optionalCommentNote(comment)
	approvedAt, err := s.finalApprove(ctx, recordType, recordID, meta.Status, managerID, note)
	if err != nil {
		return nil, mapRecordError(err)
	}

	if trimmed := trimComment(comment); trimmed != "" {
		if err := s.addComment(ctx, recordType, recordID, managerID, trimmed); err != nil {
			return nil, err
		}
	}

	return &StatusResult{Status: models.StatusApproved, ApprovedAt: &approvedAt}, nil
}

func (s *Service) FinalReturn(ctx context.Context, recordType models.RecordType, recordID, managerID uuid.UUID, comment string) (*StatusResult, error) {
	if trimComment(&comment) == "" {
		return nil, ErrCommentRequired
	}

	meta, err := s.getRecordMeta(ctx, recordType, recordID, managerID, models.RoleManager)
	if err != nil {
		return nil, err
	}
	if meta.Status != models.StatusUnderReview {
		return nil, ErrInvalidStatusTransition
	}

	text := trimComment(&comment)
	note := &text
	if err := s.updateStatus(ctx, recordType, recordID, meta.Status, models.StatusPendingReview, managerID, note); err != nil {
		return nil, mapRecordError(err)
	}

	if err := s.addComment(ctx, recordType, recordID, managerID, text); err != nil {
		return nil, err
	}

	return &StatusResult{Status: models.StatusPendingReview}, nil
}

func (s *Service) AddComment(ctx context.Context, recordType models.RecordType, recordID, authorID uuid.UUID, role models.Role, text string) (*models.RecordComment, error) {
	if role != models.RoleSupervisor && role != models.RoleManager {
		return nil, ErrForbidden
	}
	if strings.TrimSpace(text) == "" {
		return nil, ErrCommentRequired
	}

	if _, err := s.getRecordMeta(ctx, recordType, recordID, authorID, role); err != nil {
		return nil, err
	}

	comment := &models.RecordComment{
		RecordType:  recordType,
		RecordID:    recordID,
		AuthorID:    authorID,
		CommentText: strings.TrimSpace(text),
	}
	if err := s.repo.CreateComment(ctx, comment); err != nil {
		return nil, err
	}

	return comment, nil
}

func (s *Service) GetComments(ctx context.Context, recordType models.RecordType, recordID, requesterID uuid.UUID, role models.Role) ([]models.RecordComment, error) {
	if _, err := s.getRecordMeta(ctx, recordType, recordID, requesterID, role); err != nil {
		return nil, err
	}
	return s.repo.ListComments(ctx, recordType, recordID)
}

func (s *Service) GetHistory(ctx context.Context, recordType models.RecordType, recordID, requesterID uuid.UUID, role models.Role) ([]models.StatusHistoryEntry, error) {
	if _, err := s.getRecordMeta(ctx, recordType, recordID, requesterID, role); err != nil {
		return nil, err
	}
	return s.audit.ListByRecord(ctx, recordType, recordID)
}

func mapRecordError(err error) error {
	if err == nil {
		return nil
	}
	switch {
	case errors.Is(err, immovable.ErrRecordNotFound), errors.Is(err, movable.ErrRecordNotFound):
		return ErrRecordNotFound
	case errors.Is(err, immovable.ErrForbidden), errors.Is(err, movable.ErrForbidden):
		return ErrForbidden
	case errors.Is(err, immovable.ErrInvalidStatusTransition), errors.Is(err, movable.ErrInvalidStatusTransition):
		return ErrInvalidStatusTransition
	default:
		return err
	}
}

func (s *Service) getRecordMeta(ctx context.Context, recordType models.RecordType, recordID, userID uuid.UUID, role models.Role) (*recordMeta, error) {
	if !recordType.IsValid() {
		return nil, ErrInvalidRecordType
	}

	switch recordType {
	case models.RecordTypeImmovable:
		record, err := s.immovable.GetByID(ctx, recordID, userID, role)
		if err != nil {
			return nil, mapRecordError(err)
		}
		if record == nil {
			return nil, ErrRecordNotFound
		}
		if role == models.RoleRegistrar && record.RegistrarID != userID {
			return nil, ErrForbidden
		}
		return &recordMeta{Status: record.Status, RegistrarID: record.RegistrarID}, nil
	case models.RecordTypeMovable:
		record, err := s.movable.GetByID(ctx, recordID, userID, role)
		if err != nil {
			return nil, mapRecordError(err)
		}
		if record == nil {
			return nil, ErrRecordNotFound
		}
		if role == models.RoleRegistrar && record.RegistrarID != userID {
			return nil, ErrForbidden
		}
		return &recordMeta{Status: record.Status, RegistrarID: record.RegistrarID}, nil
	default:
		return nil, ErrInvalidRecordType
	}
}

func (s *Service) updateStatus(ctx context.Context, recordType models.RecordType, recordID uuid.UUID, from, to models.RecordStatus, changedBy uuid.UUID, note *string) error {
	switch recordType {
	case models.RecordTypeImmovable:
		return s.immovable.UpdateStatus(ctx, recordID, from, to, changedBy, note)
	case models.RecordTypeMovable:
		return s.movable.UpdateStatus(ctx, recordID, from, to, changedBy, note)
	default:
		return ErrInvalidRecordType
	}
}

func (s *Service) finalApprove(ctx context.Context, recordType models.RecordType, recordID uuid.UUID, from models.RecordStatus, approvedBy uuid.UUID, note *string) (time.Time, error) {
	switch recordType {
	case models.RecordTypeImmovable:
		return s.immovable.FinalApprove(ctx, recordID, from, approvedBy, note)
	case models.RecordTypeMovable:
		return s.movable.FinalApprove(ctx, recordID, from, approvedBy, note)
	default:
		return time.Time{}, ErrInvalidRecordType
	}
}

func (s *Service) addComment(ctx context.Context, recordType models.RecordType, recordID, authorID uuid.UUID, text string) error {
	comment := &models.RecordComment{
		RecordType:  recordType,
		RecordID:    recordID,
		AuthorID:    authorID,
		CommentText: text,
	}
	return s.repo.CreateComment(ctx, comment)
}

func trimComment(comment *string) string {
	if comment == nil {
		return ""
	}
	return strings.TrimSpace(*comment)
}

func optionalCommentNote(comment *string) *string {
	text := trimComment(comment)
	if text == "" {
		return nil
	}
	return &text
}
