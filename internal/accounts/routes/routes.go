// /internal/accounts/routes.go

package accounts

import (
	accounts "thyra/internal/accounts/api/accounts" // Adjust the import path as necessary
	handlers "thyra/internal/accounts/api/accounts"

	"github.com/gin-gonic/gin"
)

func SetupRoutes(router *gin.RouterGroup, accountValueHandler *handlers.AccountValueHandler) {

	router.POST("/create/account", accounts.CreateAccountHandler)
	router.GET("/user/:userId/accounts", accounts.GetAccountsByUser)
	router.GET("/accounts", accounts.GetAllAccounts)
	router.GET("/account-types", accounts.GetAccountTypes)
	router.GET("/account/house", accounts.GetHouseAccount)

	router.GET("/user/:userId/aggregated-values", accountValueHandler.GetAggregatedValues)
	router.GET("/account/:accountId/values", accountValueHandler.GetSpecificAccountValue)
}
