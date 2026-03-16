package middleware

import (
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"

	"gra/pkg/config"
	"gra/pkg/response"
)

type Claims struct {
	UserID   uint   `json:"user_id"`
	Username string `json:"username"`
	jwt.RegisteredClaims
}

func JWTAuth() gin.HandlerFunc {
	return func(c *gin.Context) {
		token := c.GetHeader("Authorization")
		if token == "" || !strings.HasPrefix(token, "Bearer ") {
			response.Error(c, http.StatusUnauthorized, "未登录或Token已过期")
			c.Abort()
			return
		}
		token = strings.TrimPrefix(token, "Bearer ")

		claims := &Claims{}
		t, err := jwt.ParseWithClaims(token, claims, func(t *jwt.Token) (interface{}, error) {
			return []byte(config.Cfg.JWT.Secret), nil
		})
		if err != nil || !t.Valid {
			response.Error(c, http.StatusUnauthorized, "Token无效")
			c.Abort()
			return
		}

		c.Set("user_id", claims.UserID)
		c.Set("username", claims.Username)
		c.Next()
	}
}
