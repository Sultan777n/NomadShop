package main

import (
	"NomadShop/handlers"
	"NomadShop/middlewares"
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

	if err := FixUserSequence(db); err != nil {
		log.Println("Could not fix user ID sequence:", err)
	}

	if err := resetAutoIncrement(db, "products"); err != nil {
		log.Println("Error resetting auto increment for products:", err)
	}
	if err := resetAutoIncrement(db, "cart_items"); err != nil {
		log.Println("Error resetting auto increment for cart_items:", err)
	}

	return db
}

func FixUserSequence(db *gorm.DB) error {
	var maxID int64
	if err := db.Table("users").Select("MAX(id)").Scan(&maxID).Error; err != nil {
		return err
	}
	// sequence-ті жаңарту
	return db.Exec(fmt.Sprintf("ALTER SEQUENCE users_id_seq RESTART WITH %d;", maxID+1)).Error
}

// Автоинкрементті 1-ден бастап орнату
func resetAutoIncrement(db *gorm.DB, tableName string) error {
	query := fmt.Sprintf("ALTER SEQUENCE %s_id_seq RESTART WITH 1;", tableName)
	return db.Exec(query).Error
}

func main() {
	db = setupDatabase()

	models.FixPasswords(db)

	r := gin.Default()

	handler := handlers.Handler{DB: db}
	r.GET("/products_all", handler.GetProducts)
	r.GET("/products/:id", handler.GetProductByID)
	r.GET("/products", handler.GetProductsByCategory)

	categoryHandler := handlers.NewCategoryHandler(db)
	categoryGroup := r.Group("/categories")
	{
		categoryGroup.GET("/", categoryHandler.GetAllCategories)
		categoryGroup.GET("/:id", categoryHandler.GetCategoryByID)
		categoryGroup.POST("/", middlewares.AuthMiddleware(), middlewares.RoleAuthorization("Admin"), categoryHandler.CreateCategory)
	}

	userHandler := handlers.NewUserHandler(db)
	r.POST("/users", middlewares.AuthMiddleware(), middlewares.RoleAuthorization("Admin", "User"), userHandler.CreateUser)
	r.GET("/users", middlewares.AuthMiddleware(), middlewares.RoleAuthorization("Admin", "Support"), userHandler.GetUsers)
	r.GET("/users/:id", middlewares.AuthMiddleware(), middlewares.RoleAuthorization("Admin", "Support", "Seller"), userHandler.GetUserByID)
	r.PUT("/users/:id", middlewares.AuthMiddleware(), middlewares.AuthMiddleware(), middlewares.RoleAuthorization("Admin", "User"), userHandler.UpdateUser)
	r.DELETE("/users/:id", middlewares.AuthMiddleware(), middlewares.RoleAuthorization("Admin", "User"), userHandler.DeleteUser)

	roleHandler := handlers.NewRoleHandler(db)
	r.GET("/roles", middlewares.AuthMiddleware(), middlewares.RoleAuthorization("Admin"), roleHandler.GetAllRoles)
	r.GET("/roles/:id", middlewares.AuthMiddleware(), middlewares.RoleAuthorization("Admin"), roleHandler.GetRoleByID)
	r.POST("/roles", middlewares.AuthMiddleware(), middlewares.RoleAuthorization("Admin"), roleHandler.CreateRole)
	r.PUT("/roles/:id", middlewares.AuthMiddleware(), middlewares.RoleAuthorization("Admin"), roleHandler.UpdateRole)
	r.DELETE("/roles/:id", middlewares.AuthMiddleware(), middlewares.RoleAuthorization("Admin"), roleHandler.DeleteRole)

	userRoleHandler := handlers.NewUserRoleHandler(db)
	r.GET("/user_roles/all", middlewares.AuthMiddleware(), middlewares.RoleAuthorization("Admin"), userRoleHandler.GetAllUserRoles)
	r.POST("/user_roles", middlewares.AuthMiddleware(), middlewares.RoleAuthorization("Admin"), userRoleHandler.AddUserRole)
	r.GET("/user_roles/", middlewares.AuthMiddleware(), middlewares.RoleAuthorization("Admin"), userRoleHandler.GetUserRoles)
	r.GET("/user-roles", middlewares.AuthMiddleware(), middlewares.RoleAuthorization("Admin"), userRoleHandler.GetUserRolesByRole)
	r.DELETE("/user_roles/:user_id/:role_id", middlewares.AuthMiddleware(), middlewares.RoleAuthorization("Admin"), userRoleHandler.DeleteUserRole)

	cartItemHandler := handlers.NewCartItemHandler(db)
	r.GET("/cart_items/:user_id", middlewares.AuthMiddleware(), middlewares.RoleAuthorization("Admin", "User"), cartItemHandler.GetCartItems)
	r.POST("/cart_items", middlewares.AuthMiddleware(), middlewares.RoleAuthorization("User"), cartItemHandler.CreateCartItem)
	r.GET("/cart_items", middlewares.AuthMiddleware(), middlewares.RoleAuthorization("Admin", "User"), cartItemHandler.GetCartItemsByUser)
	r.GET("/cart-items", middlewares.AuthMiddleware(), middlewares.RoleAuthorization("Admin", "Support"), cartItemHandler.GetCartItemsByProduct)
	r.PUT("/cart_items/:id", middlewares.AuthMiddleware(), middlewares.RoleAuthorization("User"), cartItemHandler.UpdateCartItem)
	r.DELETE("/cart_items/:id", middlewares.AuthMiddleware(), middlewares.RoleAuthorization("User"), cartItemHandler.DeleteCartItem)
	r.GET("/cart_items_all", middlewares.AuthMiddleware(), middlewares.RoleAuthorization("Admin", "Support"), cartItemHandler.GetAllCartItems)

	favoriteItemHandler := handlers.NewFavoriteItemHandler(db)
	r.GET("/favorite_items_all", middlewares.AuthMiddleware(), middlewares.RoleAuthorization("Admin", "Support"), favoriteItemHandler.GetAllFavoriteItems)
	r.GET("/favorite_items/:id", middlewares.AuthMiddleware(), middlewares.RoleAuthorization("Admin", "User"), favoriteItemHandler.GetFavoriteItemByID)
	r.GET("/favorite-items", middlewares.AuthMiddleware(), middlewares.RoleAuthorization("User"), favoriteItemHandler.GetFavoriteItemsByUser)
	r.GET("/favorite_items", middlewares.AuthMiddleware(), middlewares.RoleAuthorization("Admin", "Support"), favoriteItemHandler.GetFavoriteItemsByProduct)
	r.POST("/favorite_items", middlewares.AuthMiddleware(), middlewares.RoleAuthorization("User"), favoriteItemHandler.CreateFavoriteItem)
	r.DELETE("/favorite_items/:id", middlewares.AuthMiddleware(), middlewares.RoleAuthorization("User"), favoriteItemHandler.DeleteFavoriteItem)
	r.GET("/favorite-items/check", middlewares.AuthMiddleware(), middlewares.RoleAuthorization("Admin", "Support", "User", "Seller"), favoriteItemHandler.IsProductFavoritedByUser)
	r.DELETE("/favorite-items/delete-by-user-product", middlewares.AuthMiddleware(), middlewares.RoleAuthorization("Admin", "User"), favoriteItemHandler.DeleteFavoriteItemByUserAndProduct)

	orderHandler := handlers.NewOrderHandler(db)
	r.POST("/orders", middlewares.AuthMiddleware(), middlewares.RoleAuthorization("User"), orderHandler.CreateOrder)
	r.GET("/orders/", middlewares.AuthMiddleware(), middlewares.RoleAuthorization("User"), orderHandler.GetOrdersByUser)
	r.GET("/orders/by_id/", middlewares.AuthMiddleware(), middlewares.RoleAuthorization("Admin", "Support"), orderHandler.GetOrderByID)
	r.GET("/orders/all", middlewares.AuthMiddleware(), middlewares.RoleAuthorization("Admin", "Support"), orderHandler.GetAllOrders)
	r.PUT("/orders/:order_id", middlewares.AuthMiddleware(), middlewares.RoleAuthorization("Admin", "User"), orderHandler.UpdateOrder)
	r.DELETE("/orders/:order_id", middlewares.AuthMiddleware(), middlewares.RoleAuthorization("Admin", "User"), orderHandler.DeleteOrder)

	orderItemHandler := handlers.NewOrderItemHandler(db)
	r.GET("/order_items_all", middlewares.AuthMiddleware(), middlewares.RoleAuthorization("Admin", "Support"), orderItemHandler.GetAllOrderItems)
	r.POST("/order_items", middlewares.AuthMiddleware(), middlewares.RoleAuthorization("Admin", "Seller"), orderItemHandler.CreateOrderItem)
	r.GET("/order_items", middlewares.AuthMiddleware(), middlewares.RoleAuthorization("Admin", "User", "Support"), orderItemHandler.GetOrderItemsByOrderID)
	r.GET("/order_items/by_product_id/", middlewares.AuthMiddleware(), middlewares.RoleAuthorization("Admin", "Support", "Seller"), orderItemHandler.GetOrderItemsByProductID)
	r.PUT("/order_items/:id", middlewares.AuthMiddleware(), middlewares.RoleAuthorization("Admin", "Seller"), orderItemHandler.UpdateOrderItem)
	r.DELETE("/order_items/:id", middlewares.AuthMiddleware(), middlewares.RoleAuthorization("Admin", "Seller"), orderItemHandler.DeleteOrderItem)

	r.POST("/register", middlewares.RegisterHandler(db))
	r.POST("/login", middlewares.LoginHandler(db))
	r.GET("/profile", middlewares.AuthMiddleware(), middlewares.ProfileHandler(db))

	auth := r.Group("/products", middlewares.AuthMiddleware())
	{
		auth.POST("/create", middlewares.RoleAuthorization("Admin", "Seller"), handler.CreateProduct)
		auth.PUT("/:id", middlewares.RoleAuthorization("Admin", "Seller"), handler.UpdateProduct)
		auth.DELETE("/:id", middlewares.RoleAuthorization("Admin", "Seller"), handler.DeleteProduct)
	}

	err := r.Run(":8080")
	if err != nil {
		log.Fatal("Server run error:", err)
	}
}
