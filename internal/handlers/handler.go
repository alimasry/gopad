package handlers

import (
	"net/http"

	"github.com/alimasry/gopad/internal/models"
	"github.com/alimasry/gopad/internal/services/document"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

// render HTML page for the document
func HandleViewDocument(c *gin.Context) {
	document_uuid := c.Param("document_uuid")

	doc, err := document.GetDocumentFromCache(document_uuid)

	if err == document.ErrDocumentNotFound {
		c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		return
	} else if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.HTML(http.StatusOK, "editor.html", doc)
}

// creates and open a new document
func HandleNewDocument(c *gin.Context) {
	uuid := uuid.NewString()
	doc := models.Document{
		UUID:    uuid,
		Title:   "Untitled",
		Version: 1,
	}

	err := document.SaveDocument(doc)

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.Redirect(http.StatusFound, "/documents/"+uuid)
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

	doc := models.Document{
		UUID:    uuid.NewString(),
		Title:   body.Title,
		Version: 1,
	}

	err := document.SaveDocument(doc)

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"uuid": doc.UUID})
}
