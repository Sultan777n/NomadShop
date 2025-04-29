package handlers

import (
	"NomadShop/models"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	_ "github.com/golang-migrate/migrate/v4"
	_ "github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	"gorm.io/gorm"
)

type UserHandler struct {
	DB *gorm.DB
}

func NewUserHandler(db *gorm.DB) *UserHandler {
	return &UserHandler{DB: db}
}

func (uh *UserHandler) CreateUser(c *gin.Context) {
	var user models.User
	if err := c.ShouldBindJSON(&user); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": "Invalid input"})
		return
	}

	// Email немесе username қайталанбасын
	var existing models.User
	if err := uh.DB.Where("email = ? OR username = ?", user.Email, user.Username).First(&existing).Error; err == nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": "User with this email or username already exists"})
		return
	}

	// БД-ға сақтау (пароль BeforeCreate арқылы хэштеледі)
	newUser, err := models.CreateUser(uh.DB, &user)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"message": "Error creating user"})
		return
	}

	c.JSON(http.StatusCreated, newUser)
}

func (uh *UserHandler) GetUserByID(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": "Invalid ID format"})
		return
	}

	var user models.User
	err = uh.DB.Preload("Roles").First(&user, id).Error
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"message": "User not found"})
		return
	}

	// Егер roles [] болса, "User" рөлін қосу
	if len(user.Roles) == 0 {
		var defaultRole models.Role
		if err := uh.DB.First(&defaultRole, "name = ?", "User").Error; err == nil {
			// Қолданушыға "User" рөлін қосамыз
			uh.DB.Model(&user).Association("Roles").Append(&defaultRole)

			// Қайта жүктеп аламыз жаңа рөлмен
			uh.DB.Preload("Roles").First(&user, id)
		}
	}

	c.JSON(http.StatusOK, user)
}

func (uh *UserHandler) GetUsers(c *gin.Context) {
	// URL параметрлерін оқу
	pageStr := c.DefaultQuery("page", "1")
	limitStr := c.DefaultQuery("limit", "10")
	filter := c.DefaultQuery("filter", "")

	page, err := strconv.Atoi(pageStr)
	if err != nil || page < 1 {
		page = 1
	}

	limit, err := strconv.Atoi(limitStr)
	if err != nil || limit < 1 {
		limit = 10
	}

	offset := (page - 1) * limit

	var users []models.User
	query := uh.DB.Model(&models.User{}).Preload("Roles")

	// Егер filter берілсе, username немесе email бойынша іздейміз
	if filter != "" {
		query = query.Where("username ILIKE ? OR email ILIKE ?", "%"+filter+"%", "%"+filter+"%")
	}

	// Жалпы қолданушылар санын есептеу (фильтрмен)
	var total int64
	if err := query.Count(&total).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"message": "Error counting users"})
		return
	}

	// Қолданушыларды шектеу мен бет бойынша шығару
	if err := query.Limit(limit).Offset(offset).Find(&users).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"message": "Error fetching users"})
		return
	}

	var defaultRole models.Role
	err = uh.DB.First(&defaultRole, "name = ?", "User").Error
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"message": "Default role not found"})
		return
	}

	// Егер кейбір қолданушыларда рөлдер жоқ болса, "User" рөлін тағайындау
	for i := range users {
		if len(users[i].Roles) == 0 {
			_ = uh.DB.Model(&users[i]).Association("Roles").Append(&defaultRole)
			_ = uh.DB.Preload("Roles").First(&users[i], users[i].ID).Error
		}
	}

	c.JSON(http.StatusOK, gin.H{
		"users": users,
		"pagination": gin.H{
			"page":       page,
			"limit":      limit,
			"total":      total,
			"totalPages": (total + int64(limit) - 1) / int64(limit),
		},
	})
}

func (uh *UserHandler) UpdateUser(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": "Invalid ID format"})
		return
	}

	var updatedUser models.User
	if err := c.ShouldBindJSON(&updatedUser); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": "Invalid input"})
		return
	}

	user, err := models.UpdateUser(uh.DB, uint(id), &updatedUser)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"message": "Error updating user"})
		return
	}

	if err := uh.DB.Preload("Roles").First(user, id).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"message": "Failed to load roles"})
		return
	}

	c.JSON(http.StatusOK, user)
}

func (uh *UserHandler) DeleteUser(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": "Invalid ID format"})
		return
	}

	err = models.DeleteUser(uh.DB, uint(id))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"message": "Error deleting user"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "User deleted successfully"})
}
