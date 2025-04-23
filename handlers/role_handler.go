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

type RoleHandler struct {
	DB *gorm.DB
}

func NewRoleHandler(db *gorm.DB) *RoleHandler {
	return &RoleHandler{DB: db}
}

func (h *RoleHandler) GetAllRoles(c *gin.Context) {
	roles, err := models.GetRoles(h.DB)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"message": "Failed to get roles"})
		return
	}
	c.JSON(http.StatusOK, roles)
}

func (h *RoleHandler) GetRoleByID(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": "Invalid role ID"})
		return
	}

	role, err := models.GetRoleByID(h.DB, uint(id))
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"message": "Role not found"})
		return
	}

	c.JSON(http.StatusOK, role)
}

func (h *RoleHandler) CreateRole(c *gin.Context) {
	var role models.Role
	if err := c.ShouldBindJSON(&role); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": "Invalid input"})
		return
	}

	if _, err := models.CreateRole(h.DB, &role); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"message": "Failed to create role"})
		return
	}

	c.JSON(http.StatusOK, role)
}

func (h *RoleHandler) UpdateRole(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": "Invalid role ID"})
		return
	}

	var role models.Role
	if err := c.ShouldBindJSON(&role); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": "Invalid input"})
		return
	}

	updatedRole, err := models.UpdateRole(h.DB, uint(id), &role)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"message": "Failed to update role"})
		return
	}

	c.JSON(http.StatusOK, updatedRole)
}

func (h *RoleHandler) DeleteRole(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": "Invalid role ID"})
		return
	}

	if err := models.DeleteRole(h.DB, uint(id)); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"message": "Failed to delete role"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Role deleted"})
}
