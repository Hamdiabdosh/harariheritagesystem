package export

import (
	"errors"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"

	"github.com/qirs-mezgeb/api/internal/dashboard"
	"github.com/qirs-mezgeb/api/internal/middleware"
	"github.com/qirs-mezgeb/api/internal/models"
)

type Handler struct {
	service *Service
}

func NewHandler(service *Service) *Handler {
	return &Handler{service: service}
}

func (h *Handler) ExportCSV(c *gin.Context) {
	actor, ok := middleware.GetAuthUser(c)
	if !ok {
		middleware.RespondError(c, http.StatusUnauthorized, "Authentication required")
		return
	}

	filters, err := dashboard.ParseExportFilters(c)
	if err != nil {
		dashboard.RespondFilterError(c, err)
		return
	}

	result, err := h.service.ExportCSV(c.Request.Context(), filters, actor.ID, actor.Role)
	if err != nil {
		respondServiceError(c, err)
		return
	}

	c.Header("Content-Disposition", "attachment; filename="+result.Filename)
	c.Data(http.StatusOK, "text/csv; charset=utf-8", result.Content)
}

func (h *Handler) ExportPDF(c *gin.Context) {
	actor, ok := middleware.GetAuthUser(c)
	if !ok {
		middleware.RespondError(c, http.StatusUnauthorized, "Authentication required")
		return
	}

	recordType := models.RecordType(c.Param("type"))
	if !recordType.IsValid() {
		middleware.RespondError(c, http.StatusBadRequest, "Record type must be immovable or movable")
		return
	}

	recordID, err := uuid.Parse(c.Param("id"))
	if err != nil {
		middleware.RespondError(c, http.StatusBadRequest, "Invalid record ID")
		return
	}

	result, err := h.service.ExportPDF(c.Request.Context(), recordType, recordID, actor.ID, actor.Role)
	if err != nil {
		respondServiceError(c, err)
		return
	}

	c.Header("Content-Disposition", "attachment; filename="+result.Filename)
	c.Data(http.StatusOK, "application/pdf", result.Content)
}

func respondServiceError(c *gin.Context, err error) {
	switch {
	case errors.Is(err, ErrRecordNotFound):
		middleware.RespondError(c, http.StatusNotFound, "Record not found")
	case errors.Is(err, ErrForbidden):
		middleware.RespondError(c, http.StatusForbidden, "Forbidden")
	case errors.Is(err, ErrInvalidRecordType):
		middleware.RespondError(c, http.StatusBadRequest, "Record type must be immovable or movable")
	case errors.Is(err, ErrDraftNotPrintable):
		middleware.RespondError(c, http.StatusConflict, "Draft records cannot be printed")
	case errors.Is(err, ErrNotApproved):
		middleware.RespondError(c, http.StatusForbidden, "Only approved records can be printed")
	default:
		middleware.RespondError(c, http.StatusInternalServerError, "Export failed")
	}
}
