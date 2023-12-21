package routes

import (
	handlers "thyra/internal/orders/api" // Import the handlers package

	"github.com/gin-gonic/gin"
)

func SetupRoutes(router *gin.RouterGroup, settlementHandler *handlers.SettlementHandler, orderHandler *handlers.OrderHandler) {
	router.GET("/orders", handlers.GetAllOrdersHandler(orderHandler.Service))
	router.POST("/orders/create/sell", handlers.CreateSellOrderHandler(*orderHandler.Service))
	router.POST("/orders/create/buy", handlers.CreateBuyOrderHandler(*orderHandler.Service))
	router.PUT("/orders/:orderId/confirm", handlers.ConfirmOrderHandler(*orderHandler.Service))
	router.PUT("/orders/:orderId/execute", handlers.ExecuteOrderHandler(*orderHandler.Service))

	// Updated to use the method of the settlementHandler instance
	router.PUT("/orders/:orderId/settle/sell", settlementHandler.SettlementSellHandler)

	//router.PUT("/orders/:orderId/settle/buy", settlementHandler.SettlementBuyHandler) // Check if you need to update this line too
	router.GET("/orders/type/name", handlers.GetOrderTypeByName(*orderHandler.Service))
	router.GET("/orders/type/id", handlers.GetOrderTypeByID(*orderHandler.Service))
}
