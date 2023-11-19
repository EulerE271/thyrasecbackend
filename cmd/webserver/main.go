package main

import (
	"log"
	accounts "thyra/internal/accounts/routes"
	assets "thyra/internal/assets/routes"
	"thyra/internal/common/db"
	middle "thyra/internal/common/middleware"
	transactions "thyra/internal/transactions/routes"
	routes "thyra/internal/users/routes"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

// In main.go of the transaction service

func main() {

	// Initialize the database connection
	err := db.Initialize()
	if err != nil {
		log.Fatalf("Failed to initialize the database: %v", err)
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
		AllowOrigins:     []string{"*"}, // Replace with your frontend's URL
		AllowMethods:     []string{"GET", "POST", "PUT", "DELETE"},
		AllowHeaders:     []string{"Origin", "Content-Type"},
		AllowCredentials: true,
	}

	r.Use(cors.New(config))

	v1 := r.Group("/v1")

	routes.SetupRoutes(r)
	v1.Use(middle.DBContext())
	v1.Use(middle.TokenMiddleware)
	// Setup module-specific routes
	transactions.SetupRoutes(v1) // Setup rout
	accounts.SetupRoutes(v1)
	assets.SetupRoutes(v1)

	// Set up your routes by calling the SetupRoutes function from the "routes" package

	// Start the server
	r.Run(":8082")
}
