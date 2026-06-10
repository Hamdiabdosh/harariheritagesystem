package dashboard

import (
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

func (h *Handler) GetStats(c *gin.Context) {
	actor, ok := middleware.GetAuthUser(c)
	if !ok {
		middleware.RespondError(c, http.StatusUnauthorized, "Authentication required")
		return
	}

	stats, err := h.service.GetStats(c.Request.Context(), actor.ID, actor.Role)
	if err != nil {
		middleware.RespondError(c, http.StatusInternalServerError, "Failed to load dashboard stats")
		return
	}

	middleware.RespondSuccess(c, stats, "Dashboard stats retrieved successfully")
}

func (h *Handler) ListRecords(c *gin.Context) {
	actor, ok := middleware.GetAuthUser(c)
	if !ok {
		middleware.RespondError(c, http.StatusUnauthorized, "Authentication required")
		return
	}

	filters, err := ParseListFilters(c)
	if err != nil {
		RespondFilterError(c, err)
		return
	}

	result, err := h.service.ListRecords(c.Request.Context(), filters, actor.ID, actor.Role)
	if err != nil {
		middleware.RespondError(c, http.StatusInternalServerError, "Failed to load records")
		return
	}

	middleware.RespondSuccess(c, result, "Records retrieved successfully")
}
