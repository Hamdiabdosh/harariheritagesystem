package photos

import (
	"errors"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"

	"github.com/qirs-mezgeb/api/internal/middleware"
	"github.com/qirs-mezgeb/api/internal/models"
)

type Handler struct {
	service *Service
}

func NewHandler(service *Service) *Handler {
	return &Handler{service: service}
}

func (h *Handler) Upload(c *gin.Context) {
	actor, ok := middleware.GetAuthUser(c)
	if !ok {
		middleware.RespondError(c, http.StatusUnauthorized, "Authentication required")
		return
	}

	recordType, recordID, err := parseRecordParams(c)
	if err != nil {
		middleware.RespondError(c, http.StatusBadRequest, err.Error())
		return
	}

	fileHeader, err := c.FormFile("photo")
	if err != nil {
		middleware.RespondError(c, http.StatusBadRequest, "Photo file is required")
		return
	}

	result, err := h.service.Upload(c.Request.Context(), recordType, recordID, actor.ID, actor.Role, fileHeader)
	if err != nil {
		respondServiceError(c, err)
		return
	}

	c.JSON(http.StatusCreated, gin.H{
		"success": true,
		"data":    result,
		"message": "Photo uploaded successfully",
	})
}

func (h *Handler) Delete(c *gin.Context) {
	actor, ok := middleware.GetAuthUser(c)
	if !ok {
		middleware.RespondError(c, http.StatusUnauthorized, "Authentication required")
		return
	}

	recordType, recordID, err := parseRecordParams(c)
	if err != nil {
		middleware.RespondError(c, http.StatusBadRequest, err.Error())
		return
	}

	photoID, err := uuid.Parse(c.Param("photo_id"))
	if err != nil {
		middleware.RespondError(c, http.StatusBadRequest, "Invalid photo ID")
		return
	}

	if err := h.service.Delete(c.Request.Context(), recordType, recordID, photoID, actor.ID, actor.Role); err != nil {
		respondServiceError(c, err)
		return
	}

	middleware.RespondSuccess(c, gin.H{}, "Photo deleted")
}

func parseRecordParams(c *gin.Context) (models.RecordType, uuid.UUID, error) {
	recordType := models.RecordType(c.Param("type"))
	if !recordType.IsValid() {
		return "", uuid.Nil, ErrInvalidRecordType
	}

	recordID, err := uuid.Parse(c.Param("id"))
	if err != nil {
		return "", uuid.Nil, errors.New("invalid record ID")
	}

	return recordType, recordID, nil
}

func respondServiceError(c *gin.Context, err error) {
	switch {
	case errors.Is(err, ErrPhotoNotFound), errors.Is(err, ErrRecordNotFound):
		middleware.RespondError(c, http.StatusNotFound, err.Error())
	case errors.Is(err, ErrForbidden):
		middleware.RespondError(c, http.StatusForbidden, "Forbidden")
	case errors.Is(err, ErrNotEditable):
		middleware.RespondError(c, http.StatusConflict, "Record is not in an editable status")
	case errors.Is(err, ErrInvalidRecordType):
		middleware.RespondError(c, http.StatusBadRequest, "Record type must be immovable or movable")
	case errors.Is(err, ErrMaxPhotosReached):
		middleware.RespondError(c, http.StatusBadRequest, "Maximum of 10 photos per record")
	case errors.Is(err, ErrFileTooLarge):
		middleware.RespondError(c, http.StatusRequestEntityTooLarge, err.Error())
	case errors.Is(err, ErrUnsupportedType):
		middleware.RespondError(c, http.StatusUnsupportedMediaType, err.Error())
	default:
		middleware.RespondError(c, http.StatusInternalServerError, "Request failed")
	}
}
