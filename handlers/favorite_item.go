package handlers

import (
	"net/http"
	"strconv"

	"NomadShop/models"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type FavoriteItemHandler struct {
	DB *gorm.DB
}

func NewFavoriteItemHandler(db *gorm.DB) *FavoriteItemHandler {
	return &FavoriteItemHandler{DB: db}
}

func (fh *FavoriteItemHandler) GetAllFavoriteItems(c *gin.Context) {
	var favoriteItems []models.FavoriteItem

	err := fh.DB.Preload("Product").Preload("Product.Category").Find(&favoriteItems).Error
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"message": "Error fetching all favorite items"})
		return
	}

	c.JSON(http.StatusOK, favoriteItems)
}

func (fh *FavoriteItemHandler) GetFavoriteItemByID(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": "Invalid ID format"})
		return
	}

	var favoriteItem models.FavoriteItem
	err = fh.DB.Preload("Product").Preload("Product.Category").First(&favoriteItem, id).Error
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"message": "Favorite item not found"})
		return
	}

	c.JSON(http.StatusOK, favoriteItem)
}

func (fh *FavoriteItemHandler) GetFavoriteItemsByUser(c *gin.Context) {
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

	var favoriteItems []models.FavoriteItem
	err = fh.DB.Preload("Product").Preload("Product.Category").Where("user_id = ?", userID).Find(&favoriteItems).Error
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"message": "Error fetching favorite items"})
		return
	}

	c.JSON(http.StatusOK, favoriteItems)
}

func (fh *FavoriteItemHandler) GetFavoriteItemsByProduct(c *gin.Context) {
	productIDStr := c.DefaultQuery("product_id", "")
	if productIDStr == "" {
		c.JSON(http.StatusBadRequest, gin.H{"message": "Product ID is required"})
		return
	}

	productID, err := strconv.Atoi(productIDStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": "Invalid product ID"})
		return
	}

	var favoriteItems []models.FavoriteItem
	err = fh.DB.Preload("Product").Preload("Product.Category").Where("product_id = ?", productID).Find(&favoriteItems).Error
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"message": "Error fetching favorite items"})
		return
	}

	c.JSON(http.StatusOK, favoriteItems)
}

func (fh *FavoriteItemHandler) CreateFavoriteItem(c *gin.Context) {
	var favoriteItem models.FavoriteItem
	if err := c.ShouldBindJSON(&favoriteItem); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": "Invalid input"})
		return
	}

	var product models.Product
	if err := fh.DB.First(&product, favoriteItem.ProductID).Error; err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": "Product not found"})
		return
	}

	var category models.Category
	if err := fh.DB.First(&category, product.CategoryID).Error; err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": "Category not found"})
		return
	}

	createdItem, err := models.AddToFavorites(fh.DB, &favoriteItem)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"message": "Error creating favorite item"})
		return
	}

	var fullFavoriteItem models.FavoriteItem
	err = fh.DB.Preload("Product").Preload("Product.Category").First(&fullFavoriteItem, createdItem.ID).Error
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"message": "Error fetching favorite item with product and category"})
		return
	}

	c.JSON(http.StatusOK, fullFavoriteItem)
}

func (fh *FavoriteItemHandler) DeleteFavoriteItem(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": "Invalid ID format"})
		return
	}

	// ID бойынша сүйікті өнімді алу
	var favoriteItem models.FavoriteItem
	err = fh.DB.Preload("Product").Preload("Product.Category").First(&favoriteItem, id).Error
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"message": "Favorite item not found"})
		return
	}

	// Сүйікті өнімді өшіру
	if err := fh.DB.Delete(&favoriteItem).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"message": "Error deleting favorite item"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Favorite item deleted"})
}
