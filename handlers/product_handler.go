package handlers

import (
	"errors"
	"net/http"
	"strconv"

	"NomadShop/models"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type Handler struct {
	DB *gorm.DB
}

func (h *Handler) GetProducts(c *gin.Context) {
	// Query параметрлерін оқу
	limitStr := c.DefaultQuery("limit", "10")
	pageStr := c.DefaultQuery("page", "1")
	categoryIDStr := c.Query("category_id")
	nameFilter := c.Query("name")

	limit, err := strconv.Atoi(limitStr)
	if err != nil || limit <= 0 {
		limit = 10
	}

	page, err := strconv.Atoi(pageStr)
	if err != nil || page <= 0 {
		page = 1
	}

	offset := (page - 1) * limit

	// Сұранысты дайындау
	query := h.DB.Model(&models.Product{}).Preload("Category")

	if categoryIDStr != "" {
		categoryID, err := strconv.Atoi(categoryIDStr)
		if err == nil {
			query = query.Where("category_id = ?", categoryID)
		}
	}

	if nameFilter != "" {
		query = query.Where("name ILIKE ?", "%"+nameFilter+"%")
	}

	// Жалпы өнім санын есептеу
	var total int64
	if err := query.Count(&total).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"message": "Failed to count products"})
		return
	}

	// Өнімдерді алып келу
	var products []models.Product
	if err := query.Limit(limit).Offset(offset).Find(&products).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"message": "Failed to get products"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"products": products,
		"pagination": gin.H{
			"page":       page,
			"limit":      limit,
			"total":      total,
			"totalPages": (total + int64(limit) - 1) / int64(limit),
		},
	})
}

func (h *Handler) GetProductByID(c *gin.Context) {
	idParam := c.Param("id")
	id, err := strconv.Atoi(idParam)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": "Invalid product ID"})
		return
	}

	// Өнімді ID бойынша алу және категорияны алдын ала жүктеу
	var product models.Product
	err = h.DB.Preload("Category").First(&product, id).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			c.JSON(http.StatusNotFound, gin.H{"message": "Product not found"})
		} else {
			c.JSON(http.StatusInternalServerError, gin.H{"message": "Failed to retrieve product"})
		}
		return
	}

	c.JSON(http.StatusOK, product)
}

func (h *Handler) GetProductsByCategory(c *gin.Context) {
	categoryIDStr := c.DefaultQuery("category_id", "")
	if categoryIDStr == "" {
		c.JSON(http.StatusBadRequest, gin.H{"message": "Category ID is required"})
		return
	}

	categoryID, err := strconv.Atoi(categoryIDStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": "Invalid category ID"})
		return
	}

	var category models.Category
	err = h.DB.First(&category, categoryID).Error
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"message": "Category not found"})
		return
	}

	var products []models.Product
	err = h.DB.Preload("Category").Where("category_id = ?", categoryID).Find(&products).Error
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"message": "Error fetching products"})
		return
	}

	// Егер сол категорияда ешқандай продукт жоқ болса
	if len(products) == 0 {
		c.JSON(http.StatusNotFound, gin.H{"message": "No products found in this category"})
		return
	}

	c.JSON(http.StatusOK, products)
}

func (h *Handler) CreateProduct(c *gin.Context) {
	var product models.Product

	// JSON деректерін байланыстыру
	if err := c.ShouldBindJSON(&product); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": "Invalid input"})
		return
	}

	if _, err := models.GetCategoryByID(h.DB, product.CategoryID); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": "Category not found"})
		return
	}

	createdProduct, err := models.CreateProduct(h.DB, &product)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"message": "Error creating product"})
		return
	}

	err = h.DB.Preload("Category").First(&createdProduct, createdProduct.ID).Error
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"message": "Failed to load category"})
		return
	}

	c.JSON(http.StatusOK, createdProduct)
}

func (h *Handler) UpdateProduct(c *gin.Context) {
	idParam := c.Param("id")
	id, err := strconv.Atoi(idParam)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": "Invalid product ID"})
		return
	}

	var updatedData models.Product
	if err := c.ShouldBindJSON(&updatedData); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": "Invalid input"})
		return
	}

	// Категория бар-жоғын тексеру
	if _, err := models.GetCategoryByID(h.DB, updatedData.CategoryID); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": "Category not found"})
		return
	}

	product, err := models.UpdateProduct(h.DB, uint(id), &updatedData)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"message": "Failed to update product"})
		return
	}

	c.JSON(http.StatusOK, product)
}

func (h *Handler) DeleteProduct(c *gin.Context) {
	idParam := c.Param("id")
	id, err := strconv.Atoi(idParam)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": "Invalid product ID"})
		return
	}

	if err := models.DeleteProduct(h.DB, uint(id)); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"message": "Failed to delete product"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Product deleted successfully"})
}
