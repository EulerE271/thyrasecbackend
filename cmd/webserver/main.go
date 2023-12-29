package main

import (
	"log"
	"os"
	"thyra/internal/common/db"
	helpers "thyra/internal/common/middleware"
	"thyra/internal/common/utils"
	transactionroutes "thyra/internal/transactions/routes"
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

	dbConn := db.GetDB()
	dbxConn := db.GetDB() // Assuming you have a function to get sqlx DB connection

	utils.InitializeUsersModule(dbxConn, v1)
	v1.Use(helpers.DBContext(), helpers.TokenMiddleware)

	// Initialize modules
	utils.InitializeAccountModule(dbxConn, dbConn.DB, v1)
	utils.InitializeAssetModule(dbxConn, v1)
	utils.InitializeAnalyticsModule(dbConn.DB, v1)
	utils.InitializePositionsModule(dbxConn, v1)

	// Setup routes for other modules if needed
	transactionroutes.SetupRoutes(v1)
	//orderroutes.SetupRoutes(v1)
	// Set up your routes by calling the SetupRoutes function from the "routes" package

	// Start the server
	r.Run(":8082")
}
