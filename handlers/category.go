package handlers

import (
	"NomadShop/models"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
	"net/http"
	"strconv"
)

type CategoryHandler struct {
	DB *gorm.DB
}

func NewCategoryHandler(db *gorm.DB) *CategoryHandler {
	return &CategoryHandler{DB: db}
}


func (h *CategoryHandler) GetAllCategories(c *gin.Context) {
	categories, err := models.GetAllCategories(h.DB)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"message": "Failed to get categories"})
		return
	}
	c.JSON(http.StatusOK, categories)
}


func (h *CategoryHandler) CreateCategory(c *gin.Context) {
	var category models.Category
	if err := c.ShouldBindJSON(&category); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": "Invalid input"})
		return
	}

	if _, err := models.CreateCategory(h.DB, &category); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"message": "Failed to create category"})
		return
	}

	c.JSON(http.StatusOK, category)
}


func (h *CategoryHandler) GetCategoryByID(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": "Invalid category ID"})
		return
	}

	category, err := models.GetCategoryByID(h.DB, uint(id))
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"message": "Category not found"})
		return
	}

	c.JSON(http.StatusOK, category)
}
