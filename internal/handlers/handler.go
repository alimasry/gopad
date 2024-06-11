package handlers

import (
	"net/http"

	"github.com/alimasry/gopad/internal/models"
	"github.com/alimasry/gopad/internal/services/editor"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

// render HTML page for the document
func HandleViewDocument(c *gin.Context) {
	document_uuid := c.Param("document_uuid")

	document, err := editor.GetDocumentFromCache(document_uuid)

	if err == editor.ErrDocumentNotFound {
		c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
	}

	c.HTML(http.StatusOK, "editor.html", document)
}

// @Summary Create a new document
// @Description takes title and content and create a new document
// @Tags 	document
// @Accept  json
// @Produce json
// @Param   input   body models.CreateDocumentRequest true "title"
// @Success 200 {object} models.CreateDocumentResponse "JSON Content"
// @Router /documents [post]
func HandleCreateDocument(c *gin.Context) {
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

	if err == editor.ErrFailedToRetriefDocument {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"uuid": document.UUID})
}
