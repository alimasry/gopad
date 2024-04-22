package rest

import (
	"github.com/alimasry/gopad/internal/pkg/ws"
	"github.com/gin-gonic/gin"
)

func AddRoutes(incommingRoutes *gin.Engine) {
	incommingRoutes.GET("/documents", handleGetDocuments)
	incommingRoutes.GET("/documents/:document_uuid", handleViewDocument)
	incommingRoutes.GET("/documents/:document_uuid/json", handleGetDocument)
	incommingRoutes.GET("/documents/:document_uuid/history", handleGetDocumentHistory)
	incommingRoutes.POST("/documents", handleCreateDocument)
	incommingRoutes.PUT("/documents", handleUpdateDocument)
	incommingRoutes.GET("/documents/:document_uuid/ws", ws.ServeWs)
}
