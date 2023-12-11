package middleware

import (
	"net/http"
	"trackprosto/delivery/utils"

	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
)

func JWTAuthMiddleware(allowedRoles ...string) gin.HandlerFunc {
	return func(c *gin.Context) {
		// Mendapatkan token dari header Authorization
		authHeader := c.GetHeader("Authorization")
		if authHeader == "" {
			logrus.Error("Authorization header is required")
			utils.SendResponse(c, http.StatusUnauthorized, "Authorization header is required", nil)
			c.Abort()
			return
		}

		token, err := utils.ExtractTokenFromAuthHeader(authHeader)
		if err != nil {
			logrus.Error("Invalid authorization token")
			utils.SendResponse(c, http.StatusBadRequest, "Invalid authorization token", nil)
			c.Abort()
			return
		}

		claims, err := utils.VerifyJWTToken(token)
		if err != nil {
			logrus.Error("Invalid token or expired")
			utils.SendResponse(c, http.StatusUnauthorized, "Invalid token or expired", nil)
			c.Abort()
			return
		}

		// Memeriksa peran pengguna
		role, ok := claims["role"].(string)
		if !ok || !contains(allowedRoles, role) {
			logrus.Error("Access denied. Invalid role")
			utils.SendResponse(c, http.StatusForbidden, "Access denied. Invalid role", nil)
			c.Abort()
			return
		}
		c.Set("claims", claims)
		c.Set("claims", claims)

		c.Next()
	}
}

func contains(slice []string, value string) bool {
	for _, v := range slice {
		if v == value {
			return true
		}
	}
	return false
}

func JSONMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Writer.Header().Set("Content-Type", "application/json")
		c.Next()
	}
}
