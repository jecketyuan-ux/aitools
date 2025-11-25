package middleware

import (
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/playedu/playedu-go/internal/pkg/jwt"
	"github.com/playedu/playedu-go/internal/pkg/response"
	"github.com/playedu/playedu-go/pkg/constants"
	"github.com/redis/go-redis/v9"
)

type AuthMiddleware struct {
	jwtManager *jwt.JWTManager
	rdb        *redis.Client
	roleType   string
}

func NewAuthMiddleware(jwtManager *jwt.JWTManager, rdb *redis.Client, roleType string) *AuthMiddleware {
	return &AuthMiddleware{
		jwtManager: jwtManager,
		rdb:        rdb,
		roleType:   roleType,
	}
}

func (m *AuthMiddleware) Authenticate() gin.HandlerFunc {
	return func(c *gin.Context) {
		authHeader := c.GetHeader("Authorization")
		if authHeader == "" {
			response.Unauthorized(c)
			c.Abort()
			return
		}

		parts := strings.SplitN(authHeader, " ", 2)
		if len(parts) != 2 || parts[0] != "Bearer" {
			response.Unauthorized(c)
			c.Abort()
			return
		}

		token := parts[1]
		claims, err := m.jwtManager.Verify(token)
		if err != nil {
			response.Unauthorized(c)
			c.Abort()
			return
		}

		if claims.Role != m.roleType {
			response.Unauthorized(c)
			c.Abort()
			return
		}

		c.Set(constants.ContextKeyUserID, claims.UserID)
		c.Set(constants.ContextKeyEmail, claims.Email)
		c.Set(constants.ContextKeyRole, claims.Role)
		c.Set(constants.ContextKeyJTI, claims.JTI)

		c.Next()
	}
}
