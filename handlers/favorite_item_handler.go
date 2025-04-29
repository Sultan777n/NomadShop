package handlers

import (
	"NomadShop/models"
	"github.com/gin-gonic/gin"
	_ "github.com/golang-migrate/migrate/v4"
	_ "github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	"gorm.io/gorm"
	"net/http"
	"strconv"
)

type FavoriteItemHandler struct {
	DB *gorm.DB
}

func NewFavoriteItemHandler(db *gorm.DB) *FavoriteItemHandler {
	return &FavoriteItemHandler{DB: db}
}

func (fh *FavoriteItemHandler) GetAllFavoriteItems(c *gin.Context) {
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

	var favoriteItems []models.FavoriteItem
	query := fh.DB.Preload("Product").Preload("Product.Category")

	if filter != "" {
		query = query.Joins("JOIN products ON products.id = favorite_items.product_id").
			Where("LOWER(products.name) LIKE ?", "%"+filter+"%")
	}

	var total int64
	if err := query.Model(&models.FavoriteItem{}).Count(&total).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"message": "Error counting favorite items"})
		return
	}

	if err := query.Offset(offset).Limit(limit).Find(&favoriteItems).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"message": "Error fetching favorite items"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"data":       favoriteItems,
		"total":      total,
		"page":       page,
		"limit":      limit,
		"totalPages": (total + int64(limit) - 1) / int64(limit),
	})
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

	var favoriteItems []models.FavoriteItem
	query := fh.DB.Preload("Product").Preload("Product.Category").Where("user_id = ?", userID)

	if filter != "" {
		query = query.Joins("JOIN products ON products.id = favorite_items.product_id").
			Where("LOWER(products.name) LIKE ?", "%"+filter+"%").
			Where("favorite_items.user_id = ?", userID)
	}

	var total int64
	if err := query.Model(&models.FavoriteItem{}).Count(&total).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"message": "Error counting favorite items"})
		return
	}

	if err := query.Offset(offset).Limit(limit).Find(&favoriteItems).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"message": "Error fetching favorite items"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"data":       favoriteItems,
		"total":      total,
		"page":       page,
		"limit":      limit,
		"totalPages": (total + int64(limit) - 1) / int64(limit),
	})
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

	pageStr := c.DefaultQuery("page", "1")
	limitStr := c.DefaultQuery("limit", "10")

	page, err := strconv.Atoi(pageStr)
	if err != nil || page < 1 {
		page = 1
	}

	limit, err := strconv.Atoi(limitStr)
	if err != nil || limit < 1 {
		limit = 10
	}

	offset := (page - 1) * limit

	var favoriteItems []models.FavoriteItem
	query := fh.DB.Preload("Product").Preload("Product.Category").Where("product_id = ?", productID)

	var total int64
	if err := query.Model(&models.FavoriteItem{}).Count(&total).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"message": "Error counting favorite items"})
		return
	}

	if err := query.Offset(offset).Limit(limit).Find(&favoriteItems).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"message": "Error fetching favorite items"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"data":       favoriteItems,
		"total":      total,
		"page":       page,
		"limit":      limit,
		"totalPages": (total + int64(limit) - 1) / int64(limit),
	})
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

// Белгілі бір пайдаланушының сүйікті өнімі бар ма?
func (fh *FavoriteItemHandler) IsProductFavoritedByUser(c *gin.Context) {
	userIDStr := c.Query("user_id")
	productIDStr := c.Query("product_id")

	if userIDStr == "" || productIDStr == "" {
		c.JSON(http.StatusBadRequest, gin.H{"message": "User ID and Product ID are required"})
		return
	}

	userID, err := strconv.Atoi(userIDStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": "Invalid user ID"})
		return
	}

	productID, err := strconv.Atoi(productIDStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": "Invalid product ID"})
		return
	}

	var favoriteItem models.FavoriteItem
	err = fh.DB.Where("user_id = ? AND product_id = ?", userID, productID).First(&favoriteItem).Error
	if err != nil {
		c.JSON(http.StatusOK, gin.H{"favorited": false})
		return
	}

	c.JSON(http.StatusOK, gin.H{"favorited": true, "favorite_item_id": favoriteItem.ID})
}

// Пайдаланушы мен өнім бойынша сүйікті жазбаны өшіру
func (fh *FavoriteItemHandler) DeleteFavoriteItemByUserAndProduct(c *gin.Context) {
	userIDStr := c.Query("user_id")
	productIDStr := c.Query("product_id")

	if userIDStr == "" || productIDStr == "" {
		c.JSON(http.StatusBadRequest, gin.H{"message": "User ID and Product ID are required"})
		return
	}

	userID, err := strconv.Atoi(userIDStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": "Invalid user ID"})
		return
	}

	productID, err := strconv.Atoi(productIDStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": "Invalid product ID"})
		return
	}

	var favoriteItem models.FavoriteItem
	err = fh.DB.Where("user_id = ? AND product_id = ?", userID, productID).First(&favoriteItem).Error
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"message": "Favorite item not found"})
		return
	}

	if err := fh.DB.Delete(&favoriteItem).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"message": "Error deleting favorite item"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Favorite item deleted"})
}
