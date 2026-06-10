package dashboard

import (
	"errors"
	"net/http"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"

	"github.com/qirs-mezgeb/api/internal/middleware"
	"github.com/qirs-mezgeb/api/internal/models"
)

const MaxExportRows = 10000

func ParseListFilters(c *gin.Context) (ListFilters, error) {
	recordType := models.RecordType(c.Query("type"))
	if recordType != "" && !recordType.IsValid() {
		return ListFilters{}, errInvalidRecordType
	}

	status := models.RecordStatus(c.Query("status"))
	if status != "" && !status.IsValid() {
		return ListFilters{}, errInvalidStatus
	}

	filters := ListFilters{
		Page:   queryInt(c, "page", 1),
		Limit:  queryInt(c, "limit", 20),
		Type:   recordType,
		Status: status,
		Woreda: c.Query("woreda"),
		Kebele: c.Query("kebele"),
		Search: c.Query("search"),
	}

	if dateFrom := c.Query("date_from"); dateFrom != "" {
		parsed, err := time.Parse("2006-01-02", dateFrom)
		if err != nil {
			return filters, errInvalidDate
		}
		filters.DateFrom = &parsed
	}

	if dateTo := c.Query("date_to"); dateTo != "" {
		parsed, err := time.Parse("2006-01-02", dateTo)
		if err != nil {
			return filters, errInvalidDate
		}
		end := parsed.Add(24*time.Hour - time.Nanosecond)
		filters.DateTo = &end
	}

	return filters, nil
}

func ParseExportFilters(c *gin.Context) (ListFilters, error) {
	filters, err := ParseListFilters(c)
	if err != nil {
		return filters, err
	}
	filters.Page = 1
	filters.Limit = MaxExportRows
	return filters, nil
}

func RespondFilterError(c *gin.Context, err error) {
	switch {
	case errors.Is(err, errInvalidRecordType):
		respondFilterErrorMessage(c, "Record type must be immovable or movable")
	case errors.Is(err, errInvalidStatus):
		respondFilterErrorMessage(c, "Invalid status filter")
	case errors.Is(err, errInvalidDate):
		respondFilterErrorMessage(c, "Invalid date filter; use YYYY-MM-DD")
	default:
		respondFilterErrorMessage(c, "Invalid query parameters")
	}
}

func respondFilterErrorMessage(c *gin.Context, message string) {
	middleware.RespondError(c, http.StatusBadRequest, message)
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
