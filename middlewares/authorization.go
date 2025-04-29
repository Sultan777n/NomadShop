package middlewares

import (
	"github.com/gin-gonic/gin"
	"net/http"
)

func RoleAuthorization(allowedRoles ...string) gin.HandlerFunc {
	return func(c *gin.Context) {
		userRoles, exists := c.Get("roles")
		if !exists {
			c.JSON(http.StatusUnauthorized, gin.H{"message": "Role not found in token"})
			c.Abort()
			return
		}

		// Рөлдермен салыстыру
		roleFound := false
		for _, userRole := range userRoles.([]string) {
			for _, allowedRole := range allowedRoles {
				if userRole == allowedRole {
					roleFound = true
					break
				}
			}
		}

		if !roleFound {
			c.JSON(http.StatusForbidden, gin.H{"message": "Access denied"})
			c.Abort()
			return
		}

		c.Next()
	}
}
