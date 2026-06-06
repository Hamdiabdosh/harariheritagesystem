package users

import (
	"errors"
	"net/http"
	"strconv"

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

type createUserRequest struct {
	FullName string          `json:"full_name" binding:"required,max=100"`
	Email    string          `json:"email" binding:"required,email"`
	Password string          `json:"password" binding:"required,min=8"`
	Role     models.Role     `json:"role" binding:"required"`
	Language models.Language `json:"language"`
}

type updateUserRequest struct {
	FullName *string          `json:"full_name" binding:"omitempty,max=100"`
	Role     *models.Role     `json:"role"`
	Language *models.Language `json:"language"`
	IsActive *bool            `json:"is_active"`
}

type updateLanguageRequest struct {
	Language models.Language `json:"language" binding:"required"`
}

func (h *Handler) List(c *gin.Context) {
	filters := ListFilters{
		Page:  queryInt(c, "page", 1),
		Limit: queryInt(c, "limit", 20),
		Role:  models.Role(c.Query("role")),
	}

	if isActiveStr := c.Query("is_active"); isActiveStr != "" {
		isActive := isActiveStr == "true"
		filters.IsActive = &isActive
	}

	result, err := h.service.List(c.Request.Context(), filters)
	if err != nil {
		middleware.RespondError(c, http.StatusInternalServerError, "Failed to list users")
		return
	}

	middleware.RespondSuccess(c, result, "Users retrieved successfully")
}

func (h *Handler) Create(c *gin.Context) {
	var req createUserRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		middleware.RespondError(c, http.StatusUnprocessableEntity, "Validation failed")
		return
	}

	item, err := h.service.Create(c.Request.Context(), CreateInput{
		FullName: req.FullName,
		Email:    req.Email,
		Password: req.Password,
		Role:     req.Role,
		Language: req.Language,
	})
	if err != nil {
		respondServiceError(c, err)
		return
	}

	c.JSON(http.StatusCreated, gin.H{
		"success": true,
		"data":    gin.H{"user": item},
		"message": "User created successfully",
	})
}

func (h *Handler) Update(c *gin.Context) {
	targetID, err := uuid.Parse(c.Param("id"))
	if err != nil {
		middleware.RespondError(c, http.StatusBadRequest, "Invalid user ID")
		return
	}

	actor, ok := middleware.GetAuthUser(c)
	if !ok {
		middleware.RespondError(c, http.StatusUnauthorized, "Authentication required")
		return
	}

	var req updateUserRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		middleware.RespondError(c, http.StatusUnprocessableEntity, "Validation failed")
		return
	}

	item, err := h.service.Update(c.Request.Context(), actor.ID, targetID, UpdateInput{
		FullName: req.FullName,
		Role:     req.Role,
		Language: req.Language,
		IsActive: req.IsActive,
	})
	if err != nil {
		respondServiceError(c, err)
		return
	}

	middleware.RespondSuccess(c, gin.H{"user": item}, "User updated successfully")
}

func (h *Handler) Delete(c *gin.Context) {
	targetID, err := uuid.Parse(c.Param("id"))
	if err != nil {
		middleware.RespondError(c, http.StatusBadRequest, "Invalid user ID")
		return
	}

	actor, ok := middleware.GetAuthUser(c)
	if !ok {
		middleware.RespondError(c, http.StatusUnauthorized, "Authentication required")
		return
	}

	if err := h.service.Deactivate(c.Request.Context(), actor.ID, targetID); err != nil {
		respondServiceError(c, err)
		return
	}

	middleware.RespondSuccess(c, gin.H{}, "User deactivated successfully")
}

func (h *Handler) GetMe(c *gin.Context) {
	actor, ok := middleware.GetAuthUser(c)
	if !ok {
		middleware.RespondError(c, http.StatusUnauthorized, "Authentication required")
		return
	}

	item, err := h.service.GetMe(c.Request.Context(), actor.ID)
	if err != nil {
		respondServiceError(c, err)
		return
	}

	middleware.RespondSuccess(c, gin.H{"user": item}, "Profile retrieved successfully")
}

func (h *Handler) UpdateMyLanguage(c *gin.Context) {
	actor, ok := middleware.GetAuthUser(c)
	if !ok {
		middleware.RespondError(c, http.StatusUnauthorized, "Authentication required")
		return
	}

	var req updateLanguageRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		middleware.RespondError(c, http.StatusUnprocessableEntity, "Validation failed")
		return
	}

	language, err := h.service.UpdateMyLanguage(c.Request.Context(), actor.ID, req.Language)
	if err != nil {
		respondServiceError(c, err)
		return
	}

	middleware.RespondSuccess(c, gin.H{"language": language}, "Language updated successfully")
}

func respondServiceError(c *gin.Context, err error) {
	switch {
	case errors.Is(err, ErrUserNotFound):
		middleware.RespondError(c, http.StatusNotFound, "User not found")
	case errors.Is(err, ErrEmailAlreadyExists):
		middleware.RespondError(c, http.StatusConflict, "Email already exists")
	case errors.Is(err, ErrCannotDeactivateSelf):
		middleware.RespondError(c, http.StatusForbidden, "Cannot deactivate your own account")
	case errors.Is(err, ErrCannotChangeOwnRole):
		middleware.RespondError(c, http.StatusForbidden, "Cannot change your own role")
	case errors.Is(err, ErrInvalidRole):
		middleware.RespondError(c, http.StatusUnprocessableEntity, "Invalid role")
	case errors.Is(err, ErrInvalidLanguage):
		middleware.RespondError(c, http.StatusUnprocessableEntity, "Invalid language")
	case errors.Is(err, ErrInvalidPassword):
		middleware.RespondError(c, http.StatusUnprocessableEntity, err.Error())
	case errors.Is(err, ErrInvalidFullName):
		middleware.RespondError(c, http.StatusUnprocessableEntity, err.Error())
	default:
		middleware.RespondError(c, http.StatusInternalServerError, "Request failed")
	}
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
