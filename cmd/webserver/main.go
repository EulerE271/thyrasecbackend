package main

import (
	"log"
	"os"
	accountHandlers "thyra/internal/accounts/api/accounts"
	repository "thyra/internal/accounts/repositories"
	accounts "thyra/internal/accounts/routes"
	account "thyra/internal/accounts/services"
	assetHandlers "thyra/internal/assets/api/assets"
	"thyra/internal/assets/repositories" // Replace with the actual path to your repositories
	assets "thyra/internal/assets/routes"
	"thyra/internal/assets/services" // Replace with the actual path to your services
	"thyra/internal/common/db"
	middle "thyra/internal/common/middleware"
	transactions "thyra/internal/transactions/routes"
	routes "thyra/internal/users/routes"
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
		AllowOrigins:     []string{"https://dev.thyrasolutions.se"}, // Replace with your frontend's URL
		AllowMethods:     []string{"POST", "OPTIONS", "GET", "PUT", "DELETE"},
		AllowHeaders:     []string{"Accept", "Authorization", "Content-Type", "Content-Length", "X-CSRF-Token", "Token", "session", "Origin", "Host", "Connection", "Accept-Encoding", "Accept-Language", "X-Requested-With"},
		ExposeHeaders:    []string{"Content-Length"},
		AllowCredentials: true,
	}

	r.Use(cors.New(config))

	v1 := r.Group("/v1")

	repo := repositories.NewHoldingsRepository(db.GetDB())
	service := services.NewHoldingsService(repo)
	holdingsHandler := assetHandlers.NewHoldingsHandler(service)

	accountValueRepo := repository.NewAccountValueRepository(db.GetDB().DB)            // If renamed
	accountValueService := account.NewAccountValueService(accountValueRepo)            // If renamed
	accountValueHandler := accountHandlers.NewAccountValueHandler(accountValueService) // If renamed

	routes.SetupRoutes(r)
	v1.Use(middle.DBContext())
	v1.Use(middle.TokenMiddleware)
	// Setup module-specific routes
	transactions.SetupRoutes(v1) // Setup rout
	accounts.SetupRoutes(v1, accountValueHandler)
	assets.SetupRoutes(v1, holdingsHandler)

	// Set up your routes by calling the SetupRoutes function from the "routes" package

	// Start the server
	r.Run(":8082")
}
