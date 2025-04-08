package handlers

import (
	"log"
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
	var products []models.Product
	err := h.DB.Preload("Category").Find(&products).Error
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"message": "Failed to get products"})
		return
	}

	for _, product := range products {
		if product.Category.ID == 0 {
			log.Println("No category for product:", product.Name)
		}
	}

	c.JSON(http.StatusOK, products)
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
		if err == gorm.ErrRecordNotFound {
			c.JSON(http.StatusNotFound, gin.H{"message": "Product not found"})
		} else {
			c.JSON(http.StatusInternalServerError, gin.H{"message": "Failed to retrieve product"})
		}
		return
	}

	c.JSON(http.StatusOK, product)
}

func (h *Handler) GetProductsByCategory(c *gin.Context) {
	categoryIDStr := c.DefaultQuery("category_id", "") // category_id параметрін query параметр ретінде алу
	if categoryIDStr == "" {
		c.JSON(http.StatusBadRequest, gin.H{"message": "Category ID is required"})
		return
	}

	categoryID, err := strconv.Atoi(categoryIDStr) // category_id санға түрлендіріледі
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": "Invalid category ID"})
		return
	}

	// category_id бойынша өнімдерді алу
	var products []models.Product
	err = h.DB.Where("category_id = ?", categoryID).Find(&products).Error
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"message": "Error fetching products"})
		return
	}

	// Продуктілерді қайтару
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
