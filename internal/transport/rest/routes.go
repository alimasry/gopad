package rest

import (
	"github.com/alimasry/gopad/internal/pkg/ws"
	"github.com/gin-gonic/gin"
)

// add routes for handler functions
func AddRoutes(incommingRoutes *gin.Engine) {
	incommingRoutes.GET("/documents/:document_uuid", handleViewDocument)
	incommingRoutes.POST("/documents", handleCreateDocument)

	incommingRoutes.GET("/documents/:document_uuid/ws", ws.ServeWs)
}
