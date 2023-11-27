// /internal/accounts/routes.go

package routes

import (
	handlers "thyra/internal/accounts/api/accounts" // Import the handlers package

	"github.com/gin-gonic/gin"
)

func SetupRoutes(router *gin.RouterGroup, accountValueHandler *handlers.AccountValueHandler, accountPerformanceHandler *handlers.AccountPerformanceHandler) {
	router.POST("/create/account", handlers.CreateAccountHandler)
	router.GET("/user/:userId/accounts", handlers.GetAccountsByUser)
	router.GET("/accounts", handlers.GetAllAccounts)
	router.GET("/account-types", handlers.GetAccountTypes)
	router.GET("/account/house", handlers.GetHouseAccount)

	router.GET("/user/:userId/aggregated-values", accountValueHandler.GetAggregatedValues)
	router.GET("/account/:accountId/values", accountValueHandler.GetSpecificAccountValue)

	// Correctly use the method from accountPerformanceHandler
	router.GET("/account/:accountId/performance-change", accountPerformanceHandler.GetAccountPerformanceChange)
	router.GET("/user/:userId/performance-change", accountPerformanceHandler.GetUserPerformanceChange)
}
