package rest

import (
	"github.com/alimasry/gopad/internal/handlers"
	"github.com/alimasry/gopad/internal/pkg/ws"
	"github.com/gin-gonic/gin"

	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"

	_ "github.com/alimasry/gopad/docs"
)

// add routes for handler functions
func AddRoutes(incommingRoutes *gin.Engine) {
	incommingRoutes.GET("/documents/:document_uuid", handlers.HandleViewDocument)
	incommingRoutes.POST("/documents", handlers.HandleCreateDocument)

	incommingRoutes.GET("/documents/:document_uuid/ws", ws.ServeWs)

	// setup swagger
	incommingRoutes.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))
}
