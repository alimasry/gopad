package editor

import (
	"errors"
	"log"

	"github.com/alimasry/gopad/internal/database"
	"github.com/alimasry/gopad/internal/models"
)

var (
	ErrDocumentNotFound        = errors.New("document not found")
	ErrFailedToRetriefDocument = errors.New("failed to retrieve document")
	ErrFailedToSaveDocument    = errors.New("failed to save document")
)

// Cache for the loaded documents so that they could be retrieved quickly
// TODO: find better way to make this caching
var DocumentCache = make(map[string]models.Document)

// gets document from database
func GetDocument(documentUUID string) (*models.Document, error) {
	db := database.GetDb()

	var response models.Document

	if err := db.Model(&models.Document{}).Where("uuid = ?", documentUUID).Find(&response).Error; err != nil {
		log.Println("Error occured: ", err)
		return nil, ErrFailedToRetriefDocument
	}

	return &response, nil
}

// save document to database
func SaveDocument(document models.Document) error {
	db := database.GetDb()

	tx := db.Begin()

	if err := tx.Save(&document).Error; err != nil {
		log.Println("Error occured: ", err.Error())
		tx.Rollback()
		return ErrFailedToSaveDocument
	}

	document_version := models.DocumentVersion{
		Document: document,
	}

	if err := tx.Create(&document_version).Error; err != nil {
		log.Println("Error occured: ", err.Error())
		tx.Rollback()
		return ErrFailedToSaveDocument
	}

	DocumentCache[document.UUID] = document

	return tx.Commit().Error
}

// get document from cache and add it if it isn't there
func GetDocumentFromCache(uuid string) (models.Document, error) {
	document, ok := DocumentCache[uuid]
	if !ok {
		document, err := GetDocument(uuid)
		if err != nil {
			return models.Document{}, ErrDocumentNotFound
		}
		DocumentCache[uuid] = *document
	}
	return document, nil
}
