package routes

import (
	handlers "thyra/internal/accounts/api/accounts" // Import the handlers package

	"github.com/gin-gonic/gin"
)

func SetupRoutes(router *gin.RouterGroup, accountValueHandler *handlers.AccountBalanceHandler) {
	/*router.GET("/orders", handlers.GetAllOrdersHandler)
	  router.POST("/orders/create/sell", handlers.CreateSellOrderHandler)
	  router.POST("/orders/create/buy", handlers.CreateBuyOrderHandler)
	  router.PUT("/orders/:orderId/confirm", handlers.ConfirmOrderHandler)
	  router.PUT("/orders/:orderId/execute", handlers.ExecuteOrderHandler)
	  router.PUT("/orders/:orderId/settle", handlers.SettlementHandler) */
}
