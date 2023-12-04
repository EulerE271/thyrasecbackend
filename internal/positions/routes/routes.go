package routes

import (
	handlers "thyra/internal/positions/api" // Import the handlers package

	"github.com/gin-gonic/gin"
)

func SetupRoutes(router *gin.RouterGroup, holdingHandler *handlers.HoldingsHandler) {
	router.GET("/currency", holdingHandler.GetCurrencyID)
	router.GET("/account/:accountId/holdings", holdingHandler.GetAccountHoldingsWithDetails)

}
