package handlers

import (
	"bytes"
	"encoding/json"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"log"
	"net/http"
	"net/http/httptest"
	"testing"

	"NomadShop/middlewares"
	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
)

func setupTestDB() *gorm.DB {
	dsn := "user=postgres password=asd12345 dbname=nomadshop port=5432 sslmode=disable"
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Fatalf("Failed to connect to test database: %v", err)
	}
	return db
}

func setupRouter() *gin.Engine {
	r := gin.Default()
	db := setupTestDB() // << нақты базаға қосыламыз

	h := &Handler{DB: db}

	r.GET("/products_all", h.GetProducts)
	r.GET("/products/:id", h.GetProductByID)
	r.GET("/products", h.GetProductsByCategory)

	auth := r.Group("/products", middlewares.AuthMiddleware())
	{
		auth.POST("/create", middlewares.RoleAuthorization("Admin", "Seller"), h.CreateProduct)
		auth.PUT("/:id", middlewares.RoleAuthorization("Admin", "Seller"), h.UpdateProduct)
		auth.DELETE("/:id", middlewares.RoleAuthorization("Admin", "Seller"), h.DeleteProduct)
	}

	return r
}

// Тест: Барлық өнімдерді алу
func TestGetProducts(t *testing.T) {
	r := setupRouter()

	req, _ := http.NewRequest("GET", "/products_all", nil)
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
}

// Тест: Өнімді ID бойынша алу
func TestGetProductByID(t *testing.T) {
	r := setupRouter()

	req, _ := http.NewRequest("GET", "/products/1", nil)
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
}

// Тест: Категория бойынша өнімдерді алу
func TestGetProductsByCategory(t *testing.T) {
	r := setupRouter()

	req, _ := http.NewRequest("GET", "/products?category_id=1", nil)
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
}

// Тест: Авторизациясыз өнім құру (күтілетіні: Unauthorized)
func TestCreateProductUnauthorized(t *testing.T) {
	r := setupRouter()

	product := map[string]interface{}{
		"name":        "Test Product",
		"price":       100,
		"description": "Test product description",
		"image":       "image.jpg",
		"color":       "Red",
		"size":        "M",
		"category_id": 1,
		"stock":       10,
	}

	body, _ := json.Marshal(product)
	req, _ := http.NewRequest("POST", "/products/create", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")

	// Авторизация хедері жоқ
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusUnauthorized, w.Code)
}

// Тест: Авторизациясыз өнімді жаңарту (күтілетіні: Unauthorized)
func TestUpdateProductUnauthorized(t *testing.T) {
	r := setupRouter()

	update := map[string]interface{}{
		"name":        "Updated Product",
		"price":       120,
		"description": "Updated description",
		"image":       "newimage.jpg",
		"color":       "Blue",
		"size":        "L",
		"category_id": 1,
		"stock":       15,
	}

	body, _ := json.Marshal(update)
	req, _ := http.NewRequest("PUT", "/products/1", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")

	// Авторизация хедері жоқ
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusUnauthorized, w.Code)
}

// Тест: Авторизациясыз өнімді өшіру (күтілетіні: Unauthorized)
func TestDeleteProductUnauthorized(t *testing.T) {
	r := setupRouter()

	req, _ := http.NewRequest("DELETE", "/products/1", nil)
	// Авторизация хедері жоқ
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusUnauthorized, w.Code)
}
