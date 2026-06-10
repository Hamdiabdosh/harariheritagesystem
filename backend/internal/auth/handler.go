package auth

import (
	"errors"
	"net/http"

	"github.com/gin-gonic/gin"

	"github.com/qirs-mezgeb/api/internal/middleware"
)

type Handler struct {
	service *Service
}

func NewHandler(service *Service) *Handler {
	return &Handler{service: service}
}

type loginRequest struct {
	Email    string `json:"email" binding:"required,email"`
	Password string `json:"password" binding:"required,min=8"`
}

type refreshRequest struct {
	RefreshToken string `json:"refresh_token" binding:"required"`
}

type logoutRequest struct {
	RefreshToken string `json:"refresh_token" binding:"required"`
}

func (h *Handler) Login(c *gin.Context) {
	var req loginRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		middleware.RespondError(c, http.StatusUnprocessableEntity, "Validation failed")
		return
	}

	pair, err := h.service.Login(c.Request.Context(), req.Email, req.Password)
	if err != nil {
		switch {
		case errors.Is(err, ErrInvalidCredentials):
			middleware.RespondError(c, http.StatusUnauthorized, "Invalid email or password")
		case errors.Is(err, ErrAccountDeactivated):
			middleware.RespondError(c, http.StatusForbidden, "Account is deactivated")
		default:
			middleware.RespondError(c, http.StatusInternalServerError, "Login failed")
		}
		return
	}

	middleware.RespondSuccess(c, pair, "Logged in successfully")
}

func (h *Handler) Refresh(c *gin.Context) {
	var req refreshRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		middleware.RespondError(c, http.StatusUnprocessableEntity, "Validation failed")
		return
	}

	accessToken, err := h.service.Refresh(c.Request.Context(), req.RefreshToken)
	if err != nil {
		switch {
		case errors.Is(err, ErrInvalidRefreshToken):
			middleware.RespondError(c, http.StatusUnauthorized, "Invalid or expired refresh token")
		case errors.Is(err, ErrAccountDeactivated):
			middleware.RespondError(c, http.StatusForbidden, "Account is deactivated")
		default:
			middleware.RespondError(c, http.StatusInternalServerError, "Token refresh failed")
		}
		return
	}

	middleware.RespondSuccess(c, gin.H{"access_token": accessToken}, "Token refreshed successfully")
}

func (h *Handler) Logout(c *gin.Context) {
	var req logoutRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		middleware.RespondError(c, http.StatusUnprocessableEntity, "Validation failed")
		return
	}

	if err := h.service.Logout(c.Request.Context(), req.RefreshToken); err != nil {
		if errors.Is(err, ErrInvalidRefreshToken) {
			middleware.RespondError(c, http.StatusUnauthorized, "Invalid or expired refresh token")
			return
		}
		middleware.RespondError(c, http.StatusInternalServerError, "Logout failed")
		return
	}

	middleware.RespondSuccess(c, gin.H{}, "Logged out")
}
