package handlers

import (
	"NomadShop/models"
	"github.com/gin-gonic/gin"
	_ "github.com/golang-migrate/migrate/v4"
	_ "github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	"gorm.io/gorm"
	"log"
	"net/http"
	"strconv"
)

type UserRoleHandler struct {
	DB *gorm.DB
}

func NewUserRoleHandler(db *gorm.DB) *UserRoleHandler {
	return &UserRoleHandler{DB: db}
}

func (h *UserRoleHandler) GetAllUserRoles(c *gin.Context) {
	var userRoles []models.UserRole
	err := h.DB.Preload("User").Preload("Role").Find(&userRoles).Error
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"message": "Error fetching all user roles"})
		return
	}

	c.JSON(http.StatusOK, userRoles)
}

func (h *UserRoleHandler) AddUserRole(c *gin.Context) {
	var userRole models.UserRole
	if err := c.ShouldBindJSON(&userRole); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": "Invalid input"})
		return
	}

	// userRole.UserID мәнін журналға шығару
	log.Printf("Received user_id: %d", userRole.UserID)

	// Пайдаланушы мен рөлдің бар екеніне көз жеткізу
	var user models.User
	if err := h.DB.First(&user, userRole.UserID).Error; err != nil {
		log.Printf("Error finding user: %v", err)
		c.JSON(http.StatusNotFound, gin.H{"message": "User not found"})
		return
	}

	log.Printf("User found: %v", user)

	var role models.Role
	if err := h.DB.First(&role, userRole.RoleID).Error; err != nil {
		log.Printf("Error finding role: %v", err)
		c.JSON(http.StatusNotFound, gin.H{"message": "Role not found"})
		return
	}

	log.Printf("Role found: %v", role)

	// Рөлді тексеру: Пайдаланушының осы рөлі бар ма?
	existingUserRole, err := models.GetRoleByUserAndRoleID(h.DB, userRole.UserID, userRole.RoleID)
	if err == nil && existingUserRole != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": "User already has this role"})
		return
	}

	// Рөлді қосу
	newUserRole, err := models.AddUserRole(h.DB, &userRole)
	if err != nil {
		log.Printf("Error adding user role: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"message": "Error adding user role"})
		return
	}

	// Рөлді сәтті қосқаннан кейін толық деректермен қайта сұрау (Preload пайдаланылғандықтан, қайта сұрау қажет емес)
	c.JSON(http.StatusOK, gin.H{
		"message":  "User role added successfully",
		"userRole": newUserRole,
	})
}

func (h *UserRoleHandler) GetUserRoles(c *gin.Context) {
	userIDStr := c.DefaultQuery("user_id", "")
	if userIDStr == "" {
		c.JSON(http.StatusBadRequest, gin.H{"message": "User ID is required"})
		return
	}

	userID, err := strconv.Atoi(userIDStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": "Invalid user ID"})
		return
	}

	// Пайдаланушының рөлдерін алу
	roles, err := models.GetUserRoles(h.DB, uint(userID))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"message": "Error fetching user roles"})
		return
	}

	c.JSON(http.StatusOK, roles)
}

func (h *UserRoleHandler) GetUserRolesByRole(c *gin.Context) {
	roleIDStr := c.DefaultQuery("role_id", "")
	if roleIDStr == "" {
		c.JSON(http.StatusBadRequest, gin.H{"message": "Role ID is required"})
		return
	}

	roleID, err := strconv.Atoi(roleIDStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": "Invalid role ID"})
		return
	}

	// Рөлге байланысты пайдаланушылардың рөлдерін алу
	var userRoles []models.UserRole
	err = h.DB.Where("role_id = ?", uint(roleID)).Preload("User").Preload("Role").Find(&userRoles).Error
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"message": "Error fetching user roles for the given role"})
		return
	}

	// Рөлдер бойынша пайдаланушыларды қайтару
	c.JSON(http.StatusOK, userRoles)
}

func (h *UserRoleHandler) DeleteUserRole(c *gin.Context) {
	userIDStr := c.Param("user_id")
	roleIDStr := c.Param("role_id")

	userID, err1 := strconv.Atoi(userIDStr)
	roleID, err2 := strconv.Atoi(roleIDStr)

	if err1 != nil || err2 != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": "Invalid ID format"})
		return
	}

	// Пайдаланушыда осы рөл бар ма? Рөлді өшіру алдында тексеру
	existingUserRole, err := models.GetRoleByUserAndRoleID(h.DB, uint(userID), uint(roleID))
	if err != nil || existingUserRole == nil {
		c.JSON(http.StatusNotFound, gin.H{"message": "User does not have this role"})
		return
	}

	// Рөлді өшіру
	if err := models.DeleteUserRole(h.DB, uint(userID), uint(roleID)); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"message": "Error deleting user role"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "User role deleted successfully"})
}
