package routes

import (
	handlers "thyra/internal/accounts/api/accounts" // Import the handlers package

	"github.com/gin-gonic/gin"
)

func SetupRoutes(router *gin.RouterGroup, accountValueHandler *handlers.AccountBalanceHandler) {
	//router.GET("/account/:accountId/performance-change", accountPerformanceHandler.GetAccountPerformanceChange)
	//router.GET("/user/:userId/performance-change", accountPerformanceHandler.GetUserPerformanceChange)
}
