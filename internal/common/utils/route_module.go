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

	analyticshandler "thyra/internal/analytics/api/performance"
	analyticsrepo "thyra/internal/analytics/repositories/performance"
	analyticsroutes "thyra/internal/analytics/routes"
	analyticsservice "thyra/internal/analytics/services/performance"

	positionshandlers "thyra/internal/positions/api"
	positionsrepo "thyra/internal/positions/repositories"
	positionsroutes "thyra/internal/positions/routes"
	positionsservices "thyra/internal/positions/services"

	userhandlers "thyra/internal/users/api/users"
	userrepo "thyra/internal/users/repositories"
	usersroutes "thyra/internal/users/routes"
	userservices "thyra/internal/users/services"

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

func InitializeAnalyticsModule(db *sql.DB, router *gin.RouterGroup) {
	// Initialize repositories
	accountPerformanceRepo := analyticsrepo.NewAccountPerformanceRepository(db)

	// Initialize services
	accountPerformanceService := analyticsservice.NewAccountPerformanceService(accountPerformanceRepo)

	// Initialize handlers
	accountPerformanceHandler := analyticshandler.NewAccountPerformanceHandler(accountPerformanceService)

	// Setup routes specific to the Analytics module
	analyticsroutes.SetupRoutes(router, accountPerformanceHandler)
}

func InitializePositionsModule(dbx *sqlx.DB, router *gin.RouterGroup) {
	// Initialize repositories
	holdingRepo := positionsrepo.NewHoldingsRepository(dbx)

	// Initialize services
	holdingService := positionsservices.NewHoldingsService(holdingRepo)

	// Initialize handlers
	holdingHandler := positionshandlers.NewHoldingsHandler(holdingService)

	// Setup routes specific to the Positions module
	positionsroutes.SetupRoutes(router, holdingHandler)
}

func InitializeUsersModule(dbx *sqlx.DB, router *gin.RouterGroup) {
	// Initialize repositories
	userRepo := userrepo.NewUserRepository(dbx)
	// Initialize services
	userService := userservices.NewUserService(userRepo)
	// Initialize handlers
	userHandler := userhandlers.NewUserHandler(userService)

	authRepo := userrepo.NewAuthRepository(dbx)
	authService := userservices.NewAuthService(authRepo)
	authHandler := userhandlers.NewAuthHandler(authService)
	// Setup routes specific to the Users module
	usersroutes.SetupRoutes(router, userHandler, authHandler)
}
