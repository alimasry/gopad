package editor

import (
	"fmt"

	"github.com/alimasry/gopad/internal/database"
	"github.com/alimasry/gopad/internal/models"
)

// Cache for the loaded documents so that they could be retrieved quickly
// TODO: find better way to make this caching
var DocumentCache = make(map[string]models.Document)

// gets document from database
func GetDocument(documentUUID string) (*models.Document, error) {
	db := database.GetDb()

	var response models.Document

	if err := db.Model(&models.Document{}).Where("uuid = ?", documentUUID).Find(&response).Error; err != nil {
		return nil, fmt.Errorf("failed to retrieve document : %v", err)
	}

	return &response, nil
}

// save document to database
func SaveDocument(document models.Document) error {
	db := database.GetDb()

	tx := db.Begin()

	if err := tx.Save(&document).Error; err != nil {
		tx.Rollback()
		return err
	}

	document_version := models.DocumentVersion{
		Document: document,
	}

	if err := tx.Create(&document_version).Error; err != nil {
		tx.Rollback()
		return err
	}

	DocumentCache[document.UUID] = document

	return tx.Commit().Error
}

// get document from cache and add it if it isn't there
func GetDocumentFromCache(uuid string) models.Document {
	document, ok := DocumentCache[uuid]
	if !ok {
		document, _ := GetDocument(uuid)
		DocumentCache[uuid] = *document
	}
	return document
}
