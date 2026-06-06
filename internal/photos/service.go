package photos

import (
	"bytes"
	"context"
	"fmt"
	"io"
	"mime/multipart"
	"os"
	"path/filepath"
	"strings"

	"github.com/google/uuid"

	"github.com/qirs-mezgeb/api/internal/models"
)

const (
	maxPhotosPerRecord = 10
	maxFileSizeBytes   = 5 * 1024 * 1024
)

type ImmovableReader interface {
	GetByID(ctx context.Context, id uuid.UUID, userID uuid.UUID, role models.Role) (*models.ImmovableRecord, error)
}

type MovableReader interface {
	GetByID(ctx context.Context, id uuid.UUID, userID uuid.UUID, role models.Role) (*models.MovableRecord, error)
}

type RepositoryInterface interface {
	Create(ctx context.Context, photo *models.RecordPhoto) error
	GetByID(ctx context.Context, id uuid.UUID) (*models.RecordPhoto, error)
	ListByRecord(ctx context.Context, recordType models.RecordType, recordID uuid.UUID) ([]models.RecordPhoto, error)
	CountByRecord(ctx context.Context, recordType models.RecordType, recordID uuid.UUID) (int, error)
	Delete(ctx context.Context, id uuid.UUID) error
}

type Service struct {
	repo       RepositoryInterface
	immovable  ImmovableReader
	movable    MovableReader
	mediaPath  string
}

func NewService(repo RepositoryInterface, immovable ImmovableReader, movable MovableReader, mediaPath string) *Service {
	return &Service{
		repo:      repo,
		immovable: immovable,
		movable:   movable,
		mediaPath: mediaPath,
	}
}

type UploadResult struct {
	PhotoID  uuid.UUID `json:"photo_id"`
	FilePath string    `json:"file_path"`
}

func (s *Service) Upload(ctx context.Context, recordType models.RecordType, recordID, userID uuid.UUID, role models.Role, fileHeader *multipart.FileHeader) (*UploadResult, error) {
	if !recordType.IsValid() {
		return nil, ErrInvalidRecordType
	}
	if role != models.RoleRegistrar {
		return nil, ErrForbidden
	}

	if err := s.verifyEditableRecord(ctx, recordType, recordID, userID); err != nil {
		return nil, err
	}

	count, err := s.repo.CountByRecord(ctx, recordType, recordID)
	if err != nil {
		return nil, err
	}
	if count >= maxPhotosPerRecord {
		return nil, ErrMaxPhotosReached
	}

	if fileHeader.Size > maxFileSizeBytes {
		return nil, ErrFileTooLarge
	}

	file, err := fileHeader.Open()
	if err != nil {
		return nil, fmt.Errorf("open upload file: %w", err)
	}
	defer file.Close()

	content, err := io.ReadAll(io.LimitReader(file, maxFileSizeBytes+1))
	if err != nil {
		return nil, fmt.Errorf("read upload file: %w", err)
	}
	if len(content) > maxFileSizeBytes {
		return nil, ErrFileTooLarge
	}

	ext, err := detectImageExtension(content)
	if err != nil {
		return nil, err
	}

	photoID := uuid.New()
	relativePath := filepath.ToSlash(filepath.Join(string(recordType), recordID.String(), photoID.String()+ext))
	absolutePath := filepath.Join(s.mediaPath, relativePath)

	if err := os.MkdirAll(filepath.Dir(absolutePath), 0o755); err != nil {
		return nil, fmt.Errorf("create photo directory: %w", err)
	}

	if err := os.WriteFile(absolutePath, content, 0o644); err != nil {
		return nil, fmt.Errorf("write photo file: %w", err)
	}

	fileName := fileHeader.Filename
	size := len(content)
	photo := &models.RecordPhoto{
		ID:            photoID,
		RecordType:    recordType,
		RecordID:      recordID,
		FilePath:      relativePath,
		FileName:      &fileName,
		FileSizeBytes: &size,
		UploadedBy:    userID,
	}

	if err := s.repo.Create(ctx, photo); err != nil {
		_ = os.Remove(absolutePath)
		return nil, err
	}

	return &UploadResult{
		PhotoID:  photo.ID,
		FilePath: photo.FilePath,
	}, nil
}

func (s *Service) Delete(ctx context.Context, recordType models.RecordType, recordID, photoID, userID uuid.UUID, role models.Role) error {
	if !recordType.IsValid() {
		return ErrInvalidRecordType
	}
	if role != models.RoleRegistrar {
		return ErrForbidden
	}

	if err := s.verifyEditableRecord(ctx, recordType, recordID, userID); err != nil {
		return err
	}

	photo, err := s.repo.GetByID(ctx, photoID)
	if err != nil {
		return err
	}
	if photo == nil {
		return ErrPhotoNotFound
	}
	if photo.RecordType != recordType || photo.RecordID != recordID {
		return ErrPhotoNotFound
	}

	if err := s.repo.Delete(ctx, photoID); err != nil {
		return err
	}

	absolutePath := filepath.Join(s.mediaPath, filepath.FromSlash(photo.FilePath))
	_ = os.Remove(absolutePath)

	return nil
}

func (s *Service) ListByRecord(ctx context.Context, recordType models.RecordType, recordID uuid.UUID) ([]models.RecordPhoto, error) {
	return s.repo.ListByRecord(ctx, recordType, recordID)
}

func (s *Service) verifyEditableRecord(ctx context.Context, recordType models.RecordType, recordID, userID uuid.UUID) error {
	switch recordType {
	case models.RecordTypeImmovable:
		record, err := s.immovable.GetByID(ctx, recordID, userID, models.RoleRegistrar)
		if err != nil {
			return err
		}
		if record == nil {
			return ErrRecordNotFound
		}
		if record.RegistrarID != userID {
			return ErrForbidden
		}
		if !record.Status.IsEditable() {
			return ErrNotEditable
		}
		return nil
	case models.RecordTypeMovable:
		record, err := s.movable.GetByID(ctx, recordID, userID, models.RoleRegistrar)
		if err != nil {
			return err
		}
		if record == nil {
			return ErrRecordNotFound
		}
		if record.RegistrarID != userID {
			return ErrForbidden
		}
		if !record.Status.IsEditable() {
			return ErrNotEditable
		}
		return nil
	default:
		return ErrInvalidRecordType
	}
}

func detectImageExtension(content []byte) (string, error) {
	if len(content) < 3 {
		return "", ErrUnsupportedType
	}

	switch {
	case bytes.HasPrefix(content, []byte{0xFF, 0xD8, 0xFF}):
		return ".jpg", nil
	case bytes.HasPrefix(content, []byte{0x89, 0x50, 0x4E, 0x47}):
		return ".png", nil
	default:
		contentType := httpDetectContentType(content)
		switch contentType {
		case "image/jpeg":
			return ".jpg", nil
		case "image/png":
			return ".png", nil
		default:
			return "", ErrUnsupportedType
		}
	}
}

func httpDetectContentType(content []byte) string {
	sample := content
	if len(sample) > 512 {
		sample = sample[:512]
	}
	return strings.SplitN(detectContentType(sample), ";", 2)[0]
}

func detectContentType(data []byte) string {
	if len(data) >= 3 && data[0] == 0xFF && data[1] == 0xD8 && data[2] == 0xFF {
		return "image/jpeg"
	}
	if len(data) >= 8 && string(data[:8]) == "\x89PNG\r\n\x1a\n" {
		return "image/png"
	}
	return "application/octet-stream"
}
