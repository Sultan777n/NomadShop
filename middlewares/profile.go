package middlewares

import (
	"NomadShop/models"
	"errors"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
	"net/http"
)

func ProfileHandler(db *gorm.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		// user_id контекстен алу
		userIDRaw, exists := c.Get("user_id")
		if !exists {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
			return
		}

		// user_id түрін uint-қа түрлендіру
		userID, ok := userIDRaw.(uint)
		if !ok {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid user ID"})
			return
		}

		var user models.User
		// user-ді db-ден алу
		if err := db.First(&user, userID).Error; err != nil {
			// Егер пайдаланушы табылмаса, 404 қатесі қайтару
			if errors.Is(err, gorm.ErrRecordNotFound) {
				c.JSON(http.StatusNotFound, gin.H{"error": "User not found"})
			} else {
				c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch user"})
			}
			return
		}

		// Пайдаланушы табылған жағдайда оны JSON ретінде қайтару
		c.JSON(http.StatusOK, gin.H{
			"id":    user.ID,
			"name":  user.Username,
			"email": user.Email,
		})
	}
}
