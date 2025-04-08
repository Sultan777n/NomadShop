package main

import (
	"NomadShop/handlers"
	"NomadShop/models"
	"fmt"
	"github.com/gin-gonic/gin"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"log"
)

var db *gorm.DB
var err error

func setupDatabase() *gorm.DB {
	dsn := "user=postgres password=asd12345 dbname=nomadshop port=5432 sslmode=disable"
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Fatal("Could not connect to the database:", err)
	}

	err = db.AutoMigrate(&models.User{}, &models.Role{}, &models.UserRole{}, &models.Product{}, &models.Category{},
		&models.CartItem{}, &models.FavoriteItem{}, &models.Order{}, &models.OrderItem{})
	if err != nil {
		log.Fatal("Error during migration:", err)
	}

	if err := resetAutoIncrement(db, "products"); err != nil {
		log.Println("Error resetting auto increment for products:", err)
	}
	if err := resetAutoIncrement(db, "cart_items"); err != nil {
		log.Println("Error resetting auto increment for cart_items:", err)
	}

	return db
}

// Автоинкрементті 1-ден бастап орнату
func resetAutoIncrement(db *gorm.DB, tableName string) error {
	query := fmt.Sprintf("ALTER SEQUENCE %s_id_seq RESTART WITH 1;", tableName)
	return db.Exec(query).Error
}

func main() {
	db = setupDatabase()

	r := gin.Default()

	handler := handlers.Handler{DB: db}
	r.GET("/products_all", handler.GetProducts)
	r.GET("/products/:id", handler.GetProductByID)
	r.GET("/products", handler.GetProductsByCategory)
	r.POST("/products/create", handler.CreateProduct)
	r.PUT("/products/:id", handler.UpdateProduct)
	r.DELETE("/products/:id", handler.DeleteProduct)

	categoryHandler := handlers.NewCategoryHandler(db)
	r.GET("/categories", categoryHandler.GetAllCategories)
	r.POST("/categories", categoryHandler.CreateCategory)
	r.GET("/categories/:id", categoryHandler.GetCategoryByID)

	userHandler := handlers.NewUserHandler(db)
	r.POST("/users", userHandler.CreateUser)
	r.GET("/users", userHandler.GetUsers)
	r.GET("/users/:id", userHandler.GetUserByID)
	r.PUT("/users/:id", userHandler.UpdateUser)
	r.DELETE("/users/:id", userHandler.DeleteUser)

	roleHandler := handlers.NewRoleHandler(db)
	r.GET("/roles", roleHandler.GetAllRoles)
	r.GET("/roles/:id", roleHandler.GetRoleByID)
	r.POST("/roles", roleHandler.CreateRole)
	r.PUT("/roles/:id", roleHandler.UpdateRole)
	r.DELETE("/roles/:id", roleHandler.DeleteRole)

	userRoleHandler := handlers.NewUserRoleHandler(db)
	r.GET("/user_roles/all", userRoleHandler.GetAllUserRoles)
	r.POST("/user_roles", userRoleHandler.AddUserRole)
	r.GET("/user_roles/", userRoleHandler.GetUserRoles)
	r.GET("/user-roles", userRoleHandler.GetUserRolesByRole)
	r.DELETE("/user_roles/:user_id/:role_id", userRoleHandler.DeleteUserRole)

	cartItemHandler := handlers.NewCartItemHandler(db)
	r.GET("/cart_items/:user_id", cartItemHandler.GetCartItems)
	r.POST("/cart_items", cartItemHandler.CreateCartItem)
	r.GET("/cart_items", cartItemHandler.GetCartItemsByUser)
	r.GET("/cart-items", cartItemHandler.GetCartItemsByProduct)
	r.PUT("/cart_items/:id", cartItemHandler.UpdateCartItem)
	r.DELETE("/cart_items/:id", cartItemHandler.DeleteCartItem)
	r.GET("/cart_items_all", cartItemHandler.GetAllCartItems)

	favoriteItemHandler := handlers.NewFavoriteItemHandler(db)
	r.GET("/favorite_items_all", favoriteItemHandler.GetAllFavoriteItems)
	r.GET("/favorite_items/:id", favoriteItemHandler.GetFavoriteItemByID)
	r.GET("/favorite-items", favoriteItemHandler.GetFavoriteItemsByUser)
	r.GET("/favorite_items", favoriteItemHandler.GetFavoriteItemsByProduct)
	r.POST("/favorite_items", favoriteItemHandler.CreateFavoriteItem)
	r.DELETE("/favorite_items/:id", favoriteItemHandler.DeleteFavoriteItem)

	orderHandler := handlers.NewOrderHandler(db)
	r.POST("/orders", orderHandler.CreateOrder)
	r.GET("/orders/", orderHandler.GetOrdersByUser)
	r.GET("/orders/by_id/", orderHandler.GetOrderByID)
	r.GET("/orders/all", orderHandler.GetAllOrders)
	r.PUT("/orders/:order_id", orderHandler.UpdateOrder)
	r.DELETE("/orders/:order_id", orderHandler.DeleteOrder)

	orderItemHandler := handlers.NewOrderItemHandler(db)
	r.GET("/order_items_all", orderItemHandler.GetAllOrderItems)
	r.POST("/order_items", orderItemHandler.CreateOrderItem)
	r.GET("/order_items", orderItemHandler.GetOrderItemsByOrderID)
	r.GET("/order_items/by_product_id/", orderItemHandler.GetOrderItemsByProductID)
	r.PUT("/order_items/:id", orderItemHandler.UpdateOrderItem)
	r.DELETE("/order_items/:id", orderItemHandler.DeleteOrderItem)

	err := r.Run(":8080")
	if err != nil {
		log.Fatal("Server run error:", err)
	}
}
