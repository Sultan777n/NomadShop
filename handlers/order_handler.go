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

type OrderHandler struct {
	DB *gorm.DB
}

func NewOrderHandler(db *gorm.DB) *OrderHandler {
	return &OrderHandler{DB: db}
}

func (h *OrderHandler) CreateOrder(c *gin.Context) {
	var order models.Order
	if err := c.ShouldBindJSON(&order); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": "Invalid input"})
		return
	}

	// Тапсырысты базада сақтау
	if err := h.DB.Create(&order).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"message": "Error creating order"})
		return
	}

	// Егер OrderItems болса — оларды OrderID мен байланыстыру
	for i := range order.OrderItems {
		order.OrderItems[i].OrderID = order.ID
		if err := h.DB.Create(&order.OrderItems[i]).Error; err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"message": "Error saving order items"})
			return
		}

		// Product және Category ақпаратын жүктеу
		if err := h.DB.Preload("Product").Preload("Product.Category").First(&order.OrderItems[i], order.OrderItems[i].ID).Error; err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"message": "Error loading product data"})
			return
		}
	}

	// Дайын тапсырысты толық мәліметімен қайтадан жүктеу (User, OrderItems)
	var fullOrder models.Order
	if err := h.DB.Preload("User").Preload("OrderItems.Product").Preload("OrderItems.Product.Category").First(&fullOrder, order.ID).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"message": "Error loading full order data"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Order created successfully", "order": fullOrder})
}

func (h *OrderHandler) GetOrdersByUser(c *gin.Context) {
	userIDStr := c.DefaultQuery("user_id", "")
	if userIDStr == "" {
		c.JSON(http.StatusBadRequest, gin.H{"message": "Missing user_id parameter"})
		return
	}

	userID, err := strconv.Atoi(userIDStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": "Invalid user ID"})
		return
	}

	var orders []models.Order
	if err := h.DB.
		Preload("User").
		Preload("User.Roles").
		Preload("OrderItems.Product").
		Preload("OrderItems.Product.Category").
		Where("user_id = ?", userID).Find(&orders).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"message": "Error fetching orders"})
		return
	}

	c.JSON(http.StatusOK, orders)
}

func (h *OrderHandler) GetOrderByID(c *gin.Context) {
	// Параметрді алу
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

	var order models.Order
	// Тапсырысты ID бойынша алу
	if err := h.DB.
		Preload("User").
		Preload("User.Roles").
		Preload("OrderItems.Product").
		Preload("OrderItems.Product.Category").
		First(&order, orderID).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"message": "Error fetching order"})
		return
	}

	c.JSON(http.StatusOK, order)
}

func (h *OrderHandler) GetAllOrders(c *gin.Context) {
	var orders []models.Order
	if err := h.DB.Preload("User").Preload("OrderItems.Product").Preload("OrderItems.Product.Category").Find(&orders).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"message": "Error fetching orders"})
		return
	}
	c.JSON(http.StatusOK, orders)
}

func (h *OrderHandler) UpdateOrder(c *gin.Context) {
	orderIDStr := c.Param("order_id")
	orderID, err := strconv.Atoi(orderIDStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": "Invalid order ID"})
		return
	}

	var updatedOrder models.Order
	if err := c.ShouldBindJSON(&updatedOrder); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": "Invalid input"})
		return
	}

	var existingOrder models.Order
	if err := h.DB.First(&existingOrder, orderID).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"message": "Order not found"})
		return
	}

	// Қолмен жаңарту
	existingOrder.Status = updatedOrder.Status
	existingOrder.Total = updatedOrder.Total

	if err := h.DB.Save(&existingOrder).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"message": "Failed to update order"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Order updated successfully", "order": existingOrder})
}

func (h *OrderHandler) DeleteOrder(c *gin.Context) {
	orderIDStr := c.Param("order_id")
	orderID, err := strconv.Atoi(orderIDStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": "Invalid order ID"})
		return
	}

	// Алдымен OrderItem-дерді жою
	if err := h.DB.Where("order_id = ?", orderID).Delete(&models.OrderItem{}).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"message": "Failed to delete order items"})
		return
	}

	// Сосын тапсырыстың өзін жою
	if err := h.DB.Delete(&models.Order{}, orderID).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"message": "Failed to delete order"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Order deleted successfully"})
}
