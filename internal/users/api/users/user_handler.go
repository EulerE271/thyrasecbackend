// user_handler.go

package handlers

import (
	"database/sql"
	"net/http"
	"thyra/internal/users/models"
	"thyra/internal/users/services"

	"github.com/gin-gonic/gin"
)

type UserHandler struct {
	service *services.UserService
}

func NewUserHandler(service *services.UserService) *UserHandler {
	return &UserHandler{service: service}
}

func (h *UserHandler) GetAllUsersHandler(c *gin.Context) {
	role := c.DefaultQuery("role", "customer")
	users, err := h.service.GetAllUsers(c.Request.Context(), role)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"users": users})
}

func (h *UserHandler) GetUserNameByUuid(c *gin.Context) {
	uuid := c.Query("uuid")
	if uuid == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "UUID parameter is required"})
		return
	}

	username, err := h.service.GetUsernameByUUID(c.Request.Context(), uuid)
	if err != nil {
		if err == sql.ErrNoRows {
			c.JSON(http.StatusNotFound, gin.H{"error": "User not found"})
		} else {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch username", "details": err.Error()})
		}
		return
	}

	c.JSON(http.StatusOK, gin.H{"username": username})
}

func (h *UserHandler) RegisterAdminHandler(c *gin.Context) {
	var admin models.AdminRegistrationRequest
	if err := c.BindJSON(&admin); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Failed to parse JSON data"})
		return
	}

	if err := h.service.RegisterAdmin(c.Request.Context(), admin); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Error registering admin", "detail": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Admin registration successful"})
}

func (h *UserHandler) RegisterPartnerAdvisorHandler(c *gin.Context) {
	var advisor models.PartnerAdvisorRegistrationRequest
	if err := c.BindJSON(&advisor); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Failed to parse JSON data"})
		return
	}

	if err := h.service.RegisterPartnerAdvisor(c.Request.Context(), advisor); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Error registering partner/advisor", "detail": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Partner/Advisor registration successful"})
}

func (h *UserHandler) RegisterCustomerHandler(c *gin.Context) {
	var customer models.CustomerRegistrationRequest
	if err := c.BindJSON(&customer); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Failed to parse JSON data"})
		return
	}

	if err := h.service.RegisterCustomer(c.Request.Context(), customer); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Error registering customer", "detail": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Customer registration successful"})
}
