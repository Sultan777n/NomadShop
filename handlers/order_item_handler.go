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

type OrderItemHandler struct {
	DB *gorm.DB
}

func NewOrderItemHandler(db *gorm.DB) *OrderItemHandler {
	return &OrderItemHandler{DB: db}
}

func (h *OrderItemHandler) GetAllOrderItems(c *gin.Context) {
	var orderItems []models.OrderItem

	if err := h.DB.Preload("Product.Category").Find(&orderItems).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"message": "Failed to fetch order items"})
		return
	}

	c.JSON(http.StatusOK, orderItems)
}

func (h *OrderItemHandler) CreateOrderItem(c *gin.Context) {
	var orderItem models.OrderItem
	if err := c.ShouldBindJSON(&orderItem); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": "Invalid input"})
		return
	}

	// OrderItem-ді базада сақтау
	if err := h.DB.Create(&orderItem).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"message": "Error creating order item"})
		return
	}

	// Product және Category ақпаратын жүктеу
	if err := h.DB.
		Preload("Product").
		Preload("Product.Category").
		First(&orderItem, orderItem.ID).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"message": "Error loading product data"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message":   "Order item added successfully",
		"orderItem": orderItem,
	})
}

func (h *OrderItemHandler) GetOrderItemsByOrderID(c *gin.Context) {
	orderIDStr := c.DefaultQuery("order_id", "")
	if orderIDStr == "" {
		c.JSON(http.StatusBadRequest, gin.H{"message": "Missing order_id parameter"})
		return
	}

	orderID, err := strconv.Atoi(orderIDStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": "Invalid order ID"})
		return
	}

	var orderItems []models.OrderItem
	if err := h.DB.
		Preload("Product").
		Preload("Product.Category").
		Where("order_id = ?", orderID).
		Find(&orderItems).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"message": "Error fetching order items"})
		return
	}

	if len(orderItems) == 0 {
		c.JSON(http.StatusNotFound, gin.H{"message": "No order items found for the provided order ID"})
		return
	}

	c.JSON(http.StatusOK, orderItems)
}

func (h *OrderItemHandler) GetOrderItemsByProductID(c *gin.Context) {
	productIDStr := c.DefaultQuery("product_id", "")
	if productIDStr == "" {
		c.JSON(http.StatusBadRequest, gin.H{"message": "Missing product_id parameter"})
		return
	}

	productID, err := strconv.Atoi(productIDStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": "Invalid product ID"})
		return
	}

	var orderItems []models.OrderItem
	if err := h.DB.
		Preload("Product").
		Preload("Product.Category").
		Where("product_id = ?", productID).
		Find(&orderItems).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"message": "Error fetching order items"})
		return
	}

	if len(orderItems) == 0 {
		c.JSON(http.StatusNotFound, gin.H{"message": "No order items found for the provided product ID"})
		return
	}

	c.JSON(http.StatusOK, orderItems)
}

func (h *OrderItemHandler) UpdateOrderItem(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": "Invalid order item ID"})
		return
	}

	var updatedData models.OrderItem
	if err := c.ShouldBindJSON(&updatedData); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": "Invalid input"})
		return
	}

	var existingOrderItem models.OrderItem
	if err := h.DB.First(&existingOrderItem, id).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"message": "Order item not found"})
		return
	}

	existingOrderItem.ProductID = updatedData.ProductID
	existingOrderItem.Quantity = updatedData.Quantity
	existingOrderItem.Price = updatedData.Price

	if err := h.DB.Save(&existingOrderItem).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"message": "Failed to update order item"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message":   "Order item updated successfully",
		"orderItem": existingOrderItem,
	})
}

func (h *OrderItemHandler) DeleteOrderItem(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": "Invalid order item ID"})
		return
	}

	if err := h.DB.Delete(&models.OrderItem{}, id).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"message": "Failed to delete order item"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Order item deleted successfully"})
}
