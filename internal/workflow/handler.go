package workflow

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

type addCommentRequest struct {
	CommentText string `json:"comment_text" binding:"required"`
}

func (h *Handler) ReviewApprove(c *gin.Context) {
	actor, recordType, recordID, ok := parseWorkflowRequest(c)
	if !ok {
		return
	}

	var req optionalCommentRequest
	_ = c.ShouldBindJSON(&req)

	result, err := h.service.ReviewApprove(c.Request.Context(), recordType, recordID, actor.ID, req.value())
	if err != nil {
		respondServiceError(c, err)
		return
	}

	middleware.RespondSuccess(c, result, "Record reviewed successfully")
}

func (h *Handler) ReviewReturn(c *gin.Context) {
	actor, recordType, recordID, ok := parseWorkflowRequest(c)
	if !ok {
		return
	}

	var req requiredCommentRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		middleware.RespondError(c, http.StatusUnprocessableEntity, "Comment is required")
		return
	}

	comment, ok := req.value()
	if !ok {
		middleware.RespondError(c, http.StatusUnprocessableEntity, "Comment is required")
		return
	}

	result, err := h.service.ReviewReturn(c.Request.Context(), recordType, recordID, actor.ID, comment)
	if err != nil {
		respondServiceError(c, err)
		return
	}

	middleware.RespondSuccess(c, result, "Record returned to registrar")
}

func (h *Handler) FinalApprove(c *gin.Context) {
	actor, recordType, recordID, ok := parseWorkflowRequest(c)
	if !ok {
		return
	}

	var req optionalCommentRequest
	_ = c.ShouldBindJSON(&req)

	result, err := h.service.FinalApprove(c.Request.Context(), recordType, recordID, actor.ID, req.value())
	if err != nil {
		respondServiceError(c, err)
		return
	}

	middleware.RespondSuccess(c, result, "Record approved")
}

func (h *Handler) FinalReturn(c *gin.Context) {
	actor, recordType, recordID, ok := parseWorkflowRequest(c)
	if !ok {
		return
	}

	var req requiredCommentRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		middleware.RespondError(c, http.StatusUnprocessableEntity, "Comment is required")
		return
	}

	comment, ok := req.value()
	if !ok {
		middleware.RespondError(c, http.StatusUnprocessableEntity, "Comment is required")
		return
	}

	result, err := h.service.FinalReturn(c.Request.Context(), recordType, recordID, actor.ID, comment)
	if err != nil {
		respondServiceError(c, err)
		return
	}

	middleware.RespondSuccess(c, result, "Record returned to supervisor")
}

func (h *Handler) AddComment(c *gin.Context) {
	actor, recordType, recordID, ok := parseWorkflowRequest(c)
	if !ok {
		return
	}

	var req addCommentRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		middleware.RespondError(c, http.StatusUnprocessableEntity, "Validation failed")
		return
	}

	comment, err := h.service.AddComment(c.Request.Context(), recordType, recordID, actor.ID, actor.Role, req.CommentText)
	if err != nil {
		respondServiceError(c, err)
		return
	}

	c.JSON(http.StatusCreated, gin.H{
		"success": true,
		"data":    gin.H{"comment": comment},
		"message": "Comment added",
	})
}

func (h *Handler) GetComments(c *gin.Context) {
	actor, recordType, recordID, ok := parseWorkflowRequest(c)
	if !ok {
		return
	}

	comments, err := h.service.GetComments(c.Request.Context(), recordType, recordID, actor.ID, actor.Role)
	if err != nil {
		respondServiceError(c, err)
		return
	}

	middleware.RespondSuccess(c, gin.H{"comments": comments}, "Comments retrieved successfully")
}

func (h *Handler) GetHistory(c *gin.Context) {
	actor, recordType, recordID, ok := parseWorkflowRequest(c)
	if !ok {
		return
	}

	history, err := h.service.GetHistory(c.Request.Context(), recordType, recordID, actor.ID, actor.Role)
	if err != nil {
		respondServiceError(c, err)
		return
	}

	middleware.RespondSuccess(c, gin.H{"history": history}, "History retrieved successfully")
}

func parseWorkflowRequest(c *gin.Context) (middleware.AuthUser, models.RecordType, uuid.UUID, bool) {
	actor, ok := middleware.GetAuthUser(c)
	if !ok {
		middleware.RespondError(c, http.StatusUnauthorized, "Authentication required")
		return middleware.AuthUser{}, "", uuid.Nil, false
	}

	recordType := models.RecordType(c.Param("type"))
	if !recordType.IsValid() {
		middleware.RespondError(c, http.StatusBadRequest, "Record type must be immovable or movable")
		return middleware.AuthUser{}, "", uuid.Nil, false
	}

	recordID, err := uuid.Parse(c.Param("id"))
	if err != nil {
		middleware.RespondError(c, http.StatusBadRequest, "Invalid record ID")
		return middleware.AuthUser{}, "", uuid.Nil, false
	}

	return actor, recordType, recordID, true
}

func respondServiceError(c *gin.Context, err error) {
	switch {
	case errors.Is(err, ErrRecordNotFound):
		middleware.RespondError(c, http.StatusNotFound, "Record not found")
	case errors.Is(err, ErrForbidden):
		middleware.RespondError(c, http.StatusForbidden, "Forbidden")
	case errors.Is(err, ErrInvalidRecordType):
		middleware.RespondError(c, http.StatusBadRequest, "Record type must be immovable or movable")
	case errors.Is(err, ErrCommentRequired):
		middleware.RespondError(c, http.StatusUnprocessableEntity, "Comment is required")
	case errors.Is(err, ErrInvalidStatusTransition):
		middleware.RespondError(c, http.StatusConflict, "Record is not in the correct status for this action")
	default:
		middleware.RespondError(c, http.StatusInternalServerError, "Request failed")
	}
}
