// /internal/accounts/routes.go

package routes

import (
	handlers "thyra/internal/accounts/api/accounts" // Import the handlers package

	"github.com/gin-gonic/gin"
)

func SetupRoutes(router *gin.RouterGroup, accountBalanceHandler *handlers.AccountBalanceHandler, accountHandler *handlers.AccountHandler) {
	router.POST("/create/account", accountHandler.CreateAccount)
	router.GET("/user/:userId/accounts", accountHandler.GetAccountsByUser)
	router.GET("/accounts", accountHandler.GetAllAccounts)
	router.GET("/account-types", accountHandler.GetAccountTypes)
	router.GET("/account/house", accountHandler.GetHouseAccount)

	router.GET("/user/:userId/aggregated-values", accountBalanceHandler.GetAggregatedValues)
	router.GET("/account/:accountId/values", accountBalanceHandler.GetSpecificAccountValue)
}
