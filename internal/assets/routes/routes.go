package routes

import (
	handlers "thyra/internal/assets/api/assets"

	"github.com/gin-gonic/gin"
)

func SetupRoutes(router *gin.RouterGroup) {

	router.GET("/instruments", handlers.GetAllInstruments)
	router.GET("/types/asset", handlers.GetAllAssetTypes)
	router.POST("/create/instruments", handlers.CreateInstrument)

}
