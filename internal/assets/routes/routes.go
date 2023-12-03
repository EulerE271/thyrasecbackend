package routes

import (
	handlers "thyra/internal/assets/api/assets"

	"github.com/gin-gonic/gin"
)

func SetupRoutes(router *gin.RouterGroup, assetHandler *handlers.AssetsHandler) {

	router.GET("/instruments", assetHandler.GetAllInstruments)
	router.GET("/types/asset", assetHandler.GetAllAssetTypes)
	router.POST("/create/instruments", assetHandler.CreateInstrument)

}
