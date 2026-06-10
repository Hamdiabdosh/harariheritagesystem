package middleware

import (
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"

	"github.com/qirs-mezgeb/api/internal/models"
)

type contextKey string

const AuthUserKey contextKey = "auth_user"

type AuthUser struct {
	ID    uuid.UUID
	Email string
	Role  models.Role
}

type AccessTokenParser func(tokenString string) (AuthUser, error)

func AuthRequired(parse AccessTokenParser) gin.HandlerFunc {
	return func(c *gin.Context) {
		header := c.GetHeader("Authorization")
		if header == "" {
			RespondError(c, http.StatusUnauthorized, "Authorization header required")
			return
		}

		parts := strings.SplitN(header, " ", 2)
		if len(parts) != 2 || !strings.EqualFold(parts[0], "Bearer") {
			RespondError(c, http.StatusUnauthorized, "Invalid authorization header")
			return
		}

		user, err := parse(parts[1])
		if err != nil {
			RespondError(c, http.StatusUnauthorized, "Invalid or expired access token")
			return
		}

		c.Set(string(AuthUserKey), user)
		c.Next()
	}
}

func GetAuthUser(c *gin.Context) (AuthUser, bool) {
	value, ok := c.Get(string(AuthUserKey))
	if !ok {
		return AuthUser{}, false
	}

	user, ok := value.(AuthUser)
	return user, ok
}

func RequireRole(roles ...models.Role) gin.HandlerFunc {
	allowed := make(map[models.Role]struct{}, len(roles))
	for _, role := range roles {
		allowed[role] = struct{}{}
	}

	return func(c *gin.Context) {
		user, ok := GetAuthUser(c)
		if !ok {
			RespondError(c, http.StatusUnauthorized, "Authentication required")
			return
		}

		if _, ok := allowed[user.Role]; !ok {
			RespondError(c, http.StatusForbidden, "Insufficient permissions")
			return
		}

		c.Next()
	}
}
