package utils

import (
	"database/sql"
	accounthandler "thyra/internal/accounts/api/accounts"
	accountrepo "thyra/internal/accounts/repositories"
	accountroutes "thyra/internal/accounts/routes"
	accountservices "thyra/internal/accounts/services"

	assethandlers "thyra/internal/assets/api/assets"
	assetrepo "thyra/internal/assets/repositories"
	assetroutes "thyra/internal/assets/routes"
	assetservice "thyra/internal/assets/services"

	analyticsrepo "thyra/internal/analytics/repositories"
	analyticsroutes "thyra/internal/analytics/routes"
	analyticsservice "thyra/internal/anayltics/service"
	analyticshandler "thyra/internal/anyltics/api"

	"github.com/gin-gonic/gin"
	"github.com/jmoiron/sqlx"
)

func InitializeAccountModule(dbx *sqlx.DB, db *sql.DB, router *gin.RouterGroup) {
	// Initialize repositories
	accountValueRepo := accountrepo.NewAccountBalanceRepository(db)
	// Initialize services
	accountValueService := accountservices.NewAccountBalanceService(accountValueRepo)
	// Initialize handlers
	accountValueHandler := accounthandler.NewAccountBalanceHandler(accountValueService)

	accountRepo := accountrepo.NewAccountRepository(dbx)
	accountService := accountservices.NewAccountService(accountRepo)
	accountHandler := accounthandler.NewAccountHandler(accountService)

	// Setup routes specific to the Account module
	accountroutes.SetupRoutes(router, accountValueHandler, accountHandler)
}

func InitializeAssetModule(dbx *sqlx.DB, router *gin.RouterGroup) {
	// Initialize repositories
	assetRepo := assetrepo.NewAssetsRepository(dbx)

	// Initialize services
	assetService := assetservice.NewAssetsService(assetRepo)

	// Initialize handlers
	assetHandler := assethandlers.NewAssetsHandler(assetService)

	// Setup routes specific to the Asset module
	assetroutes.SetupRoutes(router, assetHandler)
}

func InitializeAnalyticsModule(dbx *sqlx.DB, router *gin.RouterGroup) {
	// Initialize repositories
	accountPerformanceRepo := analyticsrepo.NewAccountPerformanceRepository(dbx)

	// Initialize services
	accountPerformanceService := analyticsservice.NewAccountPerformanceService(accountPerformanceRepo)

	// Initialize handlers
	accountPerformanceHandler := analyticshandler.NewAccountPerformanceHandler(accountPerformanceService)

	// Setup routes specific to the Analytics module
	analyticsroutes.SetupRoutes(router, accountPerformanceHandler)
}
