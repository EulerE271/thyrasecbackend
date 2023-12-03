package routes

import (
	middleware "thyra/internal/common/middleware"
	handlers "thyra/internal/users/api/users"

	"github.com/gin-gonic/gin"
)

func SetupRoutes(router *gin.Engine, userHandler *handlers.UserHandler, authHandler *handlers.AuthHandler) {
	// Public route
	router.POST("/v1/login", authHandler.LoginHandler) // This route is public and outside the protected group
	// Group for version 1 APIs with Token Middleware
	v1 := router.Group("/v1")
	v1.Use(middleware.TokenMiddleware) // Apply token middleware to all routes in this group

	// Protected routes
	v1.POST("/register/admin", userHandler.RegisterAdminHandler)
	v1.POST("/register/partner", userHandler.RegisterPartnerAdvisorHandler)
	v1.POST("/register/customer", userHandler.RegisterCustomerHandler)
	v1.GET("/fetch/users", userHandler.GetAllUsersHandler)
	v1.GET("/fetch/username", userHandler.GetUserNameByUuid) // Protected route
}
