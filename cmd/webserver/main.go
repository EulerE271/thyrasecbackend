package main

import (
	"log"
	"os"
	accounthandler "thyra/internal/accounts/api/accounts"
	accountrepository "thyra/internal/accounts/repositories"
	accountroutes "thyra/internal/accounts/routes"
	accountservice "thyra/internal/accounts/services"
	assetroutes "thyra/internal/assets/routes"
	"thyra/internal/common/db"
	middle "thyra/internal/common/middleware"
	transactionroutes "thyra/internal/transactions/routes"
	"thyra/internal/users/routes"
	"time"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

// In main.go of the transaction service

func main() {

	cwd, _ := os.Getwd()
	log.Println("Current working directory:", cwd)

	const maxRetries = 50
	const retryInterval = 10 * time.Second

	var err error
	for i := 0; i < maxRetries; i++ {
		err = db.Initialize()
		if err == nil {
			break
		}

		log.Printf("Failed to initialize the database (attempt %d/%d): %v", i+1, maxRetries, err)
		time.Sleep(retryInterval)
	}

	if err != nil {
		log.Fatalf("Failed to initialize the database after %d attempts: %v", maxRetries, err)
	}

	// Test database connection
	testDB := db.GetDB()
	if testDB == nil {
		log.Fatal("Database connection is nil after initialization")
	}

	var testVar int
	err = testDB.Get(&testVar, "SELECT 1")
	if err != nil {
		log.Fatalf("Failed to query database: %v", err)
	}

	// Initialize the Gin router
	r := gin.Default()

	// Configure CORS middleware
	config := cors.Config{
		AllowOrigins: []string{
			"https://dev.thyrasolutions.se",
			"http://localhost:5173",
		}, // Replace with your frontend's URL
		AllowMethods:     []string{"POST", "OPTIONS", "GET", "PUT", "DELETE"},
		AllowHeaders:     []string{"Accept", "Authorization", "Content-Type", "Content-Length", "X-CSRF-Token", "Token", "session", "Origin", "Host", "Connection", "Accept-Encoding", "Accept-Language", "X-Requested-With"},
		ExposeHeaders:    []string{"Content-Length"},
		AllowCredentials: true,
	}

	r.Use(cors.New(config))

	v1 := r.Group("/v1")

	//repo := positionrepository.NewHoldingsRepository(db.GetDB())
	//service := positionservice.NewHoldingsService(repo)
	//holdingsHandler := positionshandler.NewHoldingsHandler(service) */

	accountValueRepo := accountrepository.NewAccountBalanceRepository(db.GetDB().DB)    // If renamed
	accountValueService := accountservice.NewAccountBalanceService(accountValueRepo)    // If renamed
	accountValueHandler := accounthandler.NewAccountBalanceHandler(accountValueService) // If renamed

	//accountPerformanceRepo := accountperformancerepository.NewAccountPerformanceRepository(db.GetDB().DB)
	//accountPerformanceService := accountperformanceservice.NewAccountPerformanceService(accountPerformanceRepo)
	//accountPerformanceHandler := accountperformancehandler.NewAccountPerformanceHandler(accountPerformanceService)

	// Setup routes
	routes.SetupRoutes(r)
	v1.Use(middle.DBContext())
	v1.Use(middle.TokenMiddleware)
	// Setup module-specific routes
	transactionroutes.SetupRoutes(v1) // Setup rout
	accountroutes.SetupRoutes(v1, accountValueHandler)
	assetroutes.SetupRoutes(v1)

	// Set up your routes by calling the SetupRoutes function from the "routes" package

	// Start the server
	r.Run(":8082")
}
