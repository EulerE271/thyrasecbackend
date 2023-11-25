package routes

import (
	middleware "thyra/internal/common/middleware"
	handlers "thyra/internal/users/api/users"

	"github.com/gin-gonic/gin"
)

func SetupRoutes(router *gin.Engine) {
	// Public route
	router.POST("/login", handlers.LoginHandler) // This route is public and outside the protected group
	router.GET("/cookie-test", handlers.CookieTestHandler)
	// Group for version 1 APIs with Token Middleware
	v1 := router.Group("/v1")
	v1.Use(middleware.TokenMiddleware) // Apply token middleware to all routes in this group

	// Protected routes
	v1.POST("/register/admin", handlers.RegisterAdminHandler)
	v1.POST("/register/partner", handlers.RegisterPartnerAdvisorHandler)
	v1.POST("/register/customer", handlers.RegisterCustomerHandler)
	v1.GET("/fetch/users", handlers.GetAllUsersHandler)
	v1.GET("/fetch/username", handlers.GetUserNameByUuid) // Protected route
}
