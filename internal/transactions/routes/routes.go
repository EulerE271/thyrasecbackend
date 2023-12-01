// /internal/transactions/api/routes.go

package transactions

import (
	middleware "thyra/internal/common/middleware"      // Middleware imports
	api "thyra/internal/transactions/api/transactions" // Import other necessary packages

	"github.com/gin-gonic/gin"
)

func SetupRoutes(router *gin.RouterGroup) {
	router.GET("/user/:userId/transactions", middleware.TokenMiddleware, api.GetTransactionByUserHandler)
	router.GET("/transactions", middleware.TokenMiddleware, api.GetAllTransactionsHandler)
	router.GET("/transaction/types", middleware.TokenMiddleware, api.GetTransactionTypesHandler)
	router.POST("/transaction/create/deposit", middleware.TokenMiddleware, api.CreateDeposit)
	router.POST("/transaction/create/withdrawal", middleware.TokenMiddleware, api.CreateWithdrawal)
	router.GET("/assets/id", api.GetAssetID)
}
