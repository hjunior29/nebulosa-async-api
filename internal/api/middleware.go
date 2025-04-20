package api

import (
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/hjunior29/nebulosa-async-api/internal/utils"
)

func AuthMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		authHeader := c.GetHeader("Authorization")
		if authHeader == "" {
			utils.NotAuthorized(c, "Authorization header is missing")
			return
		}

		tokenParts := strings.Split(authHeader, " ")
		if len(tokenParts) != 2 || strings.ToLower(tokenParts[0]) != "bearer" {
			utils.NotAuthorized(c, "Authorization header format must be Bearer {token}")
			return
		}

		token := tokenParts[1]
		isValid, err := utils.VerifyJWT(token)
		if err != nil || !isValid {
			utils.NotAuthorized(c, "Invalid or expired token")
			return
		}

		c.Next()
	}
}
