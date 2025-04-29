package handlers

import (
	"NomadShop/middlewares"
	"NomadShop/models"
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/golang-migrate/migrate/v4"
	_ "github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	"github.com/stretchr/testify/assert"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"log"
	"net/http"
	"net/http/httptest"
	"os"
	"path/filepath"
	"testing"
)

func setupUserTestDB() *gorm.DB {
	dsn := "user=postgres password=asd12345 dbname=nomadshop port=5432 sslmode=disable"
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Fatalf("Could not connect to test DB: %v", err)
	}

	err = applyMigrationsForTest()
	if err != nil {
		log.Fatalf("Migration failed: %v", err)
	}

	models.FixPasswords(db)

	return db
}

func applyMigrationsForTest() error {
	url := "postgres://postgres:asd12345@localhost:5432/nomadshop?sslmode=disable"

	// Windows жолын URI-ге сәйкестендіру
	basePath, err := os.Getwd()
	if err != nil {
		return fmt.Errorf("could not get working dir: %w", err)
	}

	// Толық жол (UNIX style path for migrate)
	absPath := filepath.Join(basePath, "..", "db", "migrations")
	absPath = filepath.ToSlash(absPath) // <-- басты қадам
	migrationPath := "file://" + absPath

	m, err := migrate.New(migrationPath, url)
	if err != nil {
		return fmt.Errorf("migration init error: %w", err)
	}

	if err := m.Up(); err != nil && err != migrate.ErrNoChange {
		return fmt.Errorf("migration apply error: %w", err)
	}

	return nil
}

func setupUserRouter() *gin.Engine {
	r := gin.Default()
	db := setupUserTestDB()
	h := NewUserHandler(db)

	r.POST("/users", middlewares.AuthMiddleware(), middlewares.RoleAuthorization("Admin", "User"), h.CreateUser)
	r.GET("/users", middlewares.AuthMiddleware(), middlewares.RoleAuthorization("Admin", "Support"), h.GetUsers)
	r.GET("/users/:id", middlewares.AuthMiddleware(), middlewares.RoleAuthorization("Admin", "Support", "Seller"), h.GetUserByID)
	r.PUT("/users/:id", middlewares.AuthMiddleware(), middlewares.RoleAuthorization("Admin", "User"), h.UpdateUser)
	r.DELETE("/users/:id", middlewares.AuthMiddleware(), middlewares.RoleAuthorization("Admin", "User"), h.DeleteUser)

	return r
}

func TestCreateUser_PasswordHashed(t *testing.T) {
	db := setupUserTestDB()
	handler := NewUserHandler(db)
	router := gin.Default()

	router.POST("/users", handler.CreateUser)

	user := models.User{
		Username: "Secureuser4",
		Email:    "secure4@example.com",
		Password: "rawpassword",
	}
	body, _ := json.Marshal(user)

	req, _ := http.NewRequest("POST", "/users", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusCreated, w.Code)

	var created models.User
	json.Unmarshal(w.Body.Bytes(), &created)

	// 🔐 Пароль хэштелген бе?
	assert.NotEqual(t, "rawpassword", created.Password)
	err := bcrypt.CompareHashAndPassword([]byte(created.Password), []byte("rawpassword"))
	assert.NoError(t, err)
}

func TestGetUsersUnauthorized(t *testing.T) {
	r := setupUserRouter()
	req, _ := http.NewRequest("GET", "/users", nil)
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusUnauthorized, w.Code)
}

func TestGetUserByIDUnauthorized(t *testing.T) {
	r := setupUserRouter()
	req, _ := http.NewRequest("GET", "/users/1", nil)
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusUnauthorized, w.Code)
}

func TestUpdateUserUnauthorized(t *testing.T) {
	r := setupUserRouter()

	update := models.User{
		Username: "Secure3",
		Email:    "secure3@example.com",
	}
	body, _ := json.Marshal(update)

	req, _ := http.NewRequest("PUT", "/users/1", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()

	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusUnauthorized, w.Code)
}

func TestDeleteUserUnauthorized(t *testing.T) {
	r := setupUserRouter()
	req, _ := http.NewRequest("DELETE", "/users/1", nil)
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusUnauthorized, w.Code)
}
