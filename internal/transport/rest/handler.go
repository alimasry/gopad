package rest

import (
	"net/http"

	"github.com/alimasry/gopad/internal/models"
	"github.com/alimasry/gopad/internal/services/editor"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

// renders the editor
func handleViewDocument(c *gin.Context) {
	document_uuid := c.Param("document_uuid")

	document := editor.GetDocumentFromCache(document_uuid)

	c.HTML(http.StatusOK, "editor.html", document)
}

// creates a new document
func handleCreateDocument(c *gin.Context) {
	var body models.CreateDocumentRequest

	if err := c.ShouldBindJSON(&body); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	document := models.Document{
		UUID:    uuid.NewString(),
		Title:   body.Title,
		Content: body.Content,
		Size:    len(body.Content),
		Version: 1,
	}

	err := editor.SaveDocument(document)

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err})
		return
	}

	c.JSON(http.StatusOK, gin.H{"uuid": document.UUID})
}
