package handlers

import (
	"fmt"
	"net/http"
	"strconv"

	"NomadShop/models"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type CartItemHandler struct {
	DB *gorm.DB
}

func NewCartItemHandler(db *gorm.DB) *CartItemHandler {
	return &CartItemHandler{DB: db}
}

func (ch *CartItemHandler) GetAllCartItems(c *gin.Context) {
	var cartItems []models.CartItem

	err := ch.DB.Preload("Product").Preload("Product.Category").Find(&cartItems).Error

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"message": "Error fetching all cart items"})
		return
	}

	for _, cartItem := range cartItems {
		fmt.Printf("CartItem: %v\n", cartItem)
	}

	c.JSON(http.StatusOK, cartItems)
}

func (ch *CartItemHandler) GetCartItems(c *gin.Context) {
	userIDStr := c.Param("user_id")
	userID, err := strconv.Atoi(userIDStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": "Invalid user ID"})
		return
	}

	cartItems, err := models.GetCartItems(ch.DB, uint(userID))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"message": "Error fetching cart items"})
		return
	}

	c.JSON(http.StatusOK, cartItems)
}

func (ch *CartItemHandler) GetCartItemsByUser(c *gin.Context) {
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

	var cartItems []models.CartItem
	// Продукция мен оның категориясын алдын ала жүктеу
	err = ch.DB.Preload("Product").Preload("Product.Category").Where("user_id = ?", userID).Find(&cartItems).Error
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"message": "Error fetching cart items"})
		return
	}

	c.JSON(http.StatusOK, cartItems)
}

func (ch *CartItemHandler) GetCartItemsByProduct(c *gin.Context) {
	// product_id сұрау параметрін алу
	productIDStr := c.DefaultQuery("product_id", "") // product_id query параметрі
	if productIDStr == "" {
		c.JSON(http.StatusBadRequest, gin.H{"message": "Product ID is required"})
		return
	}

	// product_id санға түрлендіріледі
	productID, err := strconv.Atoi(productIDStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": "Invalid product ID"})
		return
	}

	// Preload арқылы өнім мен оның категориясын жүктейміз
	var cartItems []models.CartItem
	err = ch.DB.Preload("Product").Preload("Product.Category").Where("product_id = ?", productID).Find(&cartItems).Error
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"message": "Error fetching cart items"})
		return
	}

	// Себеттегі өнімдермен бірге өнімнің толық мәліметтері қайтарылады
	c.JSON(http.StatusOK, cartItems)
}

func (ch *CartItemHandler) CreateCartItem(c *gin.Context) {
	var cartItem models.CartItem
	if err := c.ShouldBindJSON(&cartItem); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": "Invalid input"})
		return
	}

	fmt.Printf("Received ProductID: %d\n", cartItem.ProductID)

	var product models.Product
	if err := ch.DB.First(&product, cartItem.ProductID).Error; err != nil {
		fmt.Printf("Error fetching product: %v\n", err)
		c.JSON(http.StatusBadRequest, gin.H{"message": "Product not found"})
		return
	}

	if cartItem.Quantity > uint(product.Stock) {
		c.JSON(http.StatusBadRequest, gin.H{"message": "Not enough stock"})
		return
	}

	var existingCartItem models.CartItem
	if err := ch.DB.Where("user_id = ? AND product_id = ?", cartItem.UserID, cartItem.ProductID).First(&existingCartItem).Error; err == nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": "Product already in cart"})
		return
	}

	addedCartItem, err := models.AddToCart(ch.DB, &cartItem)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"message": "Error creating cart item"})
		return
	}

	c.JSON(http.StatusOK, addedCartItem)
}

func (ch *CartItemHandler) UpdateCartItem(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": "Invalid ID format"})
		return
	}

	var updatedCartItem models.CartItem
	if err := c.ShouldBindJSON(&updatedCartItem); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": "Invalid input"})
		return
	}

	cartItem, err := models.GetCartItemByID(ch.DB, uint(id))
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"message": "Cart item not found"})
		return
	}

	// Жаңарту: тек санды өзгерту
	if updatedCartItem.Quantity > uint(cartItem.Product.Stock) {
		c.JSON(http.StatusBadRequest, gin.H{"message": "Not enough stock to update quantity"})
		return
	}
	cartItem.Quantity = updatedCartItem.Quantity

	if err := ch.DB.Save(&cartItem).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"message": "Error updating cart item"})
		return
	}

	c.JSON(http.StatusOK, cartItem)
}

func (ch *CartItemHandler) DeleteCartItem(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": "Invalid ID format"})
		return
	}

	// Өнімді табу және өшіру
	if err := ch.DB.Where("id = ?", id).Delete(&models.CartItem{}).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"message": "Error deleting cart item"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Cart item deleted"})
}
