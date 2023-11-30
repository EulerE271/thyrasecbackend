package routes

import (
	handlers "thyra/internal/assets/api/assets"

	"github.com/gin-gonic/gin"
)

func SetupRoutes(router *gin.RouterGroup, holdingsHandler *handlers.HoldingsHandler) {

	router.GET("/instruments", handlers.GetAllInstruments)
	router.GET("/types/asset", handlers.GetAllAssetTypes)
	router.POST("/create/instruments", handlers.CreateInstrument)
	router.GET("/orders", handlers.GetAllOrdersHandler)
	router.POST("/orders/create/sell", handlers.CreateSellOrderHandler)
	router.POST("/orders/create/buy", handlers.CreateBuyOrderHandler)
	router.PUT("/orders/:orderId/confirm", handlers.ConfirmOrderHandler)
	router.PUT("/orders/:orderId/execute", handlers.ExecuteOrderHandler)
	router.PUT("/orders/:orderId/settle", handlers.SettlementHandler)
	router.GET("/currency", holdingsHandler.GetCurrencyID)
	router.GET("/account/:accountId/holdings", holdingsHandler.GetAccountHoldingsWithDetails)

}
