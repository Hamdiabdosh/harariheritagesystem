package export

import (
	"testing"

	"github.com/google/uuid"

	"github.com/qirs-mezgeb/api/internal/models"
)

func TestAuthorizePDFExport(t *testing.T) {
	registrarID := uuid.New()
	otherID := uuid.New()

	tests := []struct {
		name    string
		status  models.RecordStatus
		ownerID uuid.UUID
		userID  uuid.UUID
		role    models.Role
		wantErr error
	}{
		{
			name:    "draft blocked",
			status:  models.StatusDraft,
			ownerID: registrarID,
			userID:  registrarID,
			role:    models.RoleRegistrar,
			wantErr: ErrDraftNotPrintable,
		},
		{
			name:    "registrar own approved",
			status:  models.StatusApproved,
			ownerID: registrarID,
			userID:  registrarID,
			role:    models.RoleRegistrar,
			wantErr: nil,
		},
		{
			name:    "registrar other record",
			status:  models.StatusApproved,
			ownerID: otherID,
			userID:  registrarID,
			role:    models.RoleRegistrar,
			wantErr: ErrForbidden,
		},
		{
			name:    "registrar pending review",
			status:  models.StatusPendingReview,
			ownerID: registrarID,
			userID:  registrarID,
			role:    models.RoleRegistrar,
			wantErr: ErrNotApproved,
		},
		{
			name:    "manager approved",
			status:  models.StatusApproved,
			ownerID: registrarID,
			userID:  uuid.New(),
			role:    models.RoleManager,
			wantErr: nil,
		},
		{
			name:    "supervisor pending review",
			status:  models.StatusPendingReview,
			ownerID: registrarID,
			userID:  uuid.New(),
			role:    models.RoleSupervisor,
			wantErr: nil,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := authorizePDFExport(tt.status, tt.ownerID, tt.userID, tt.role)
			if tt.wantErr == nil && err != nil {
				t.Fatalf("expected nil, got %v", err)
			}
			if tt.wantErr != nil && err != tt.wantErr {
				t.Fatalf("expected %v, got %v", tt.wantErr, err)
			}
		})
	}
}

func TestBuildCSV(t *testing.T) {
	content, err := buildCSV(nil)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(content) == 0 {
		t.Fatal("expected csv header bytes")
	}
}

func TestBuildRecordPDF(t *testing.T) {
	content, err := buildRecordPDF(pdfRecord{
		Title:    "Test Record",
		RecordID: "ET-HR-AN-I-2026-0001",
		Status:   string(models.StatusApproved),
		Fields: [][2]string{
			{"Woreda", "Amir-Nur"},
			{"Kebele", "01"},
		},
	})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(content) < 100 {
		t.Fatal("expected non-trivial pdf output")
	}
}
