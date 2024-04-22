package rest

import (
	"fmt"
	"net/http"

	"github.com/alimasry/gopad/internal/database"
	"github.com/alimasry/gopad/internal/models"
	"github.com/alimasry/gopad/internal/services/editor"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

func handleGetDocuments(c *gin.Context) {
	db := database.GetDb()

	var response []models.GetDocumentResponse

	if err := db.Model(&models.Document{}).Find(&response).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to retrieve documents"})
	}

	c.JSON(http.StatusOK, response)
}

func handleViewDocument(c *gin.Context) {
	document_uuid := c.Param("document_uuid")

	document, err := editor.GetDocument(document_uuid)

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err})
	}

	c.HTML(http.StatusOK, "editor.html", document)
}

func handleGetDocument(c *gin.Context) {

	document_uuid := c.Param("document_uuid")

	document := editor.GetDocumentFromMap(document_uuid)

	c.JSON(http.StatusOK, document)
}
func handleGetDocumentHistory(c *gin.Context) {
	db := database.GetDb()

	document_uuid := c.Param("document_uuid")

	var response []models.GetDocumentResponse

	if err := db.Model(&models.DocumentVersion{}).Where("uuid = ?", document_uuid).Find(&response).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to retrieve document"})
	}

	c.JSON(http.StatusOK, response)
}

func handleCreateDocument(c *gin.Context) {
	db := database.GetDb()

	var body models.CreateDocumentRequest

	if err := c.ShouldBindJSON(&body); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	document := models.Document{
		UUID:    uuid.New().String(),
		Title:   body.Title,
		Content: body.Content,
		Size:    int(len(body.Content)),
		Version: 1,
	}

	document_version := models.DocumentVersion{
		Document: document,
	}

	err := db.Transaction(func(tx *gorm.DB) error {
		if err := tx.Create(&document).Error; err != nil {
			return err
		}
		if err := tx.Create(&document_version).Error; err != nil {
			return err
		}
		return nil
	})

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err})
		return
	}

	c.JSON(http.StatusOK, gin.H{"uuid": document.UUID})
}

func handleUpdateDocument(c *gin.Context) {
	db := database.GetDb()

	var body models.UpdateDocumentRequest

	if err := c.BindJSON(&body); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	var document models.Document
	if err := db.Where("uuid = ?", body.UUID).First(&document).Error; err != nil {
		c.JSON(
			http.StatusInternalServerError,
			gin.H{"message": fmt.Sprintf("failed to update the document : %v", err)},
		)
	}

	document.Version++

	if body.Title != nil {
		document.Title = *body.Title
	}

	if body.Content != nil {
		document.Content = *body.Content
		document.Size = int(len(document.Content))
	}

	if err := editor.SaveDocument(document); err != nil {
		c.JSON(
			http.StatusInternalServerError,
			gin.H{"message": fmt.Sprintf("failed to update the document : %v", err)},
		)
	}

	c.JSON(http.StatusOK, gin.H{"version": document.Version})
}
