package routes

import (
	handlers "thyra/internal/accounts/api/accounts" // Import the handlers package

	"github.com/gin-gonic/gin"
)

func SetupRoutes(router *gin.RouterGroup, accountValueHandler *handlers.AccountBalanceHandler) {
	/*router.GET("/currency", holdingsHandler.GetCurrencyID)
	  router.GET("/account/:accountId/holdings", holdingsHandler.GetAccountHoldingsWithDetails)
	*/
}
