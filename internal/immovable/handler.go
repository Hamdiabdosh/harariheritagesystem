package immovable

import (
	"errors"
	"net/http"
	"strconv"
	"time"

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

func (h *Handler) Create(c *gin.Context) {
	actor, ok := middleware.GetAuthUser(c)
	if !ok {
		middleware.RespondError(c, http.StatusUnauthorized, "Authentication required")
		return
	}

	var input models.ImmovableRecordInput
	if err := c.ShouldBindJSON(&input); err != nil {
		middleware.RespondError(c, http.StatusUnprocessableEntity, "Validation failed")
		return
	}

	result, err := h.service.Create(c.Request.Context(), actor.ID, input)
	if err != nil {
		respondServiceError(c, err)
		return
	}

	c.JSON(http.StatusCreated, gin.H{
		"success": true,
		"data":    result,
		"message": "Record created successfully",
	})
}

func (h *Handler) List(c *gin.Context) {
	actor, ok := middleware.GetAuthUser(c)
	if !ok {
		middleware.RespondError(c, http.StatusUnauthorized, "Authentication required")
		return
	}

	filters, err := parseListFilters(c)
	if err != nil {
		middleware.RespondError(c, http.StatusBadRequest, "Invalid query parameters")
		return
	}

	result, err := h.service.List(c.Request.Context(), filters, actor.ID, actor.Role)
	if err != nil {
		respondServiceError(c, err)
		return
	}

	middleware.RespondSuccess(c, result, "Records retrieved successfully")
}

func (h *Handler) GetByID(c *gin.Context) {
	actor, ok := middleware.GetAuthUser(c)
	if !ok {
		middleware.RespondError(c, http.StatusUnauthorized, "Authentication required")
		return
	}

	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		middleware.RespondError(c, http.StatusBadRequest, "Invalid record ID")
		return
	}

	result, err := h.service.GetByID(c.Request.Context(), id, actor.ID, actor.Role)
	if err != nil {
		respondServiceError(c, err)
		return
	}

	middleware.RespondSuccess(c, result, "Record retrieved successfully")
}

func (h *Handler) Update(c *gin.Context) {
	actor, ok := middleware.GetAuthUser(c)
	if !ok {
		middleware.RespondError(c, http.StatusUnauthorized, "Authentication required")
		return
	}

	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		middleware.RespondError(c, http.StatusBadRequest, "Invalid record ID")
		return
	}

	var input models.ImmovableRecordInput
	if err := c.ShouldBindJSON(&input); err != nil {
		middleware.RespondError(c, http.StatusUnprocessableEntity, "Validation failed")
		return
	}

	record, err := h.service.Update(c.Request.Context(), id, actor.ID, actor.Role, input)
	if err != nil {
		respondServiceError(c, err)
		return
	}

	middleware.RespondSuccess(c, gin.H{"record": record}, "Record updated successfully")
}

func (h *Handler) Submit(c *gin.Context) {
	actor, ok := middleware.GetAuthUser(c)
	if !ok {
		middleware.RespondError(c, http.StatusUnauthorized, "Authentication required")
		return
	}

	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		middleware.RespondError(c, http.StatusBadRequest, "Invalid record ID")
		return
	}

	result, err := h.service.Submit(c.Request.Context(), id, actor.ID, actor.Role)
	if err != nil {
		respondServiceError(c, err)
		return
	}

	middleware.RespondSuccess(c, result, "Record submitted for review")
}

func parseListFilters(c *gin.Context) (ListFilters, error) {
	filters := ListFilters{
		Page:   queryInt(c, "page", 1),
		Limit:  queryInt(c, "limit", 20),
		Status: models.RecordStatus(c.Query("status")),
		Woreda: c.Query("woreda"),
		Search: c.Query("search"),
	}

	if dateFrom := c.Query("date_from"); dateFrom != "" {
		parsed, err := time.Parse("2006-01-02", dateFrom)
		if err != nil {
			return filters, err
		}
		filters.DateFrom = &parsed
	}

	if dateTo := c.Query("date_to"); dateTo != "" {
		parsed, err := time.Parse("2006-01-02", dateTo)
		if err != nil {
			return filters, err
		}
		end := parsed.Add(24*time.Hour - time.Nanosecond)
		filters.DateTo = &end
	}

	return filters, nil
}

func queryInt(c *gin.Context, key string, fallback int) int {
	value := c.Query(key)
	if value == "" {
		return fallback
	}
	parsed, err := strconv.Atoi(value)
	if err != nil || parsed < 1 {
		return fallback
	}
	return parsed
}

func respondServiceError(c *gin.Context, err error) {
	var validationErr *ValidationError
	switch {
	case errors.As(err, &validationErr):
		middleware.RespondError(c, http.StatusUnprocessableEntity, "Validation failed")
	case errors.Is(err, ErrRecordNotFound):
		middleware.RespondError(c, http.StatusNotFound, "Record not found")
	case errors.Is(err, ErrForbidden):
		middleware.RespondError(c, http.StatusForbidden, "Forbidden")
	case errors.Is(err, ErrNotEditable):
		middleware.RespondError(c, http.StatusConflict, "Record is not in an editable status")
	case errors.Is(err, ErrInvalidStatusTransition):
		middleware.RespondError(c, http.StatusConflict, "Invalid status transition")
	default:
		middleware.RespondError(c, http.StatusInternalServerError, "Request failed")
	}
}
