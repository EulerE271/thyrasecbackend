// handlers/account_handler.go
package handlers

import (
	"net/http"
	"thyra/internal/accounts/models" // Import your Account model
	"thyra/internal/accounts/services"

	"github.com/gin-gonic/gin"
)

type AccountHandler struct {
	service *services.AccountService
}

func NewAccountHandler(service *services.AccountService) *AccountHandler {
	return &AccountHandler{service: service}
}

func (h *AccountHandler) CreateAccount(c *gin.Context) {
	var account models.Account
	if err := c.ShouldBindJSON(&account); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	authUserID, exists := c.Get("userID")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "User not authenticated"})
		return
	}

	ctx := c.Request.Context()
	if err := h.service.CreateAccount(ctx, account, authUserID.(string)); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create account", "details": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, account)
}

func (h *AccountHandler) GetAccountsByUser(c *gin.Context) {
	targetUserID := c.Param("userId")
	authUserID, _ := c.Get("userID")
	authUserRole, _ := c.Get("userType")

	ctx := c.Request.Context()
	accounts, err := h.service.GetAccountsByUser(ctx, targetUserID, authUserID.(string), authUserRole.(string))
	if err != nil {
		c.JSON(http.StatusForbidden, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, accounts)
}

func (h *AccountHandler) GetAllAccounts(c *gin.Context) {
	authUserRole, _ := c.Get("userType")

	ctx := c.Request.Context()
	accounts, err := h.service.GetAllAccounts(ctx, authUserRole.(string))
	if err != nil {
		c.JSON(http.StatusForbidden, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, accounts)
}

func (h *AccountHandler) GetAccountTypes(c *gin.Context) {
	authUserRole, _ := c.Get("userType")

	ctx := c.Request.Context()
	accountTypes, err := h.service.GetAccountTypes(ctx, authUserRole.(string))
	if err != nil {
		c.JSON(http.StatusForbidden, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, accountTypes)
}

func (h *AccountHandler) GetHouseAccount(c *gin.Context) {
	authUserRole, _ := c.Get("userType")

	ctx := c.Request.Context()
	accountID, err := h.service.GetHouseAccount(ctx, authUserRole.(string))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"house_account_id": accountID})
}
