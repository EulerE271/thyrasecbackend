package routes

import (
	handlers "thyra/internal/analytics/api/performance" // Import the handlers package

	"github.com/gin-gonic/gin"
)

func SetupRoutes(router *gin.RouterGroup, accountPerformanceHandler *handlers.AccountPerformanceHandler) {
	//router.GET("/account/:accountId/performance-change", accountPerformanceHandler.GetAccountPerformanceChange)
	//router.GET("/user/:userId/performance-change", accountPerformanceHandler.GetUserPerformanceChange)
}
