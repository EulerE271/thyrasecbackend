// /internal/accounts/routes.go

package accounts

import (
	accounts "thyra/internal/accounts/api/accounts" // Adjust the import path as necessary

	"github.com/gin-gonic/gin"
)

func SetupRoutes(router *gin.RouterGroup) {
	router.POST("/create/account", accounts.CreateAccountHandler)
	router.GET("/user/:userId/accounts", accounts.GetAccountsByUser)
	router.GET("/accounts", accounts.GetAllAccounts)
	router.GET("/account-types", accounts.GetAccountTypes)
	router.GET("/account/house", accounts.GetHouseAccount)
}
