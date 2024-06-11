package editor

import (
	"encoding/json"
	"errors"
	"log"
	"time"

	"github.com/alimasry/gopad/internal/database"
	"github.com/alimasry/gopad/internal/models"
	"github.com/patrickmn/go-cache"
)

var (
	ErrDocumentNotFound           = errors.New("document not found")
	ErrFailedToAddDocumentToCache = errors.New("failed to add document to cache")
	ErrFailedToRetriefDocument    = errors.New("failed to retrieve document")
	ErrFailedToSaveDocument       = errors.New("failed to save document")
)

const (
	// time for before each document expires from the cache
	defaultExpiration = 10 * time.Minute

	// time for cache cleanup
	cleanupInterval = 20 * time.Minute
)

// cache for the loaded documents so that they could be retrieved quickly
var DocumentCache *cache.Cache = cache.New(defaultExpiration, cleanupInterval)

// gets document from database
func GetDocument(documentUUID string) (models.Document, error) {
	db := database.GetDb()

	var response models.Document

	if err := db.Model(&models.Document{}).Where("uuid = ?", documentUUID).Find(&response).Error; err != nil {
		log.Println("Error occured: ", err)
		return models.Document{}, ErrFailedToRetriefDocument
	}

	return response, nil
}

// save document to database
func SaveDocument(document models.Document) error {
	db := database.GetDb()

	tx := db.Begin()

	if err := tx.Save(&document).Error; err != nil {
		tx.Rollback()
		return ErrFailedToSaveDocument
	}

	document_version := models.DocumentVersion{
		Document: document,
	}

	if err := tx.Create(&document_version).Error; err != nil {
		tx.Rollback()
		return ErrFailedToSaveDocument
	}

	if err := addDocumentToCache(document); err != nil {
		return err
	}

	return tx.Commit().Error
}

// get document from cache and add it if it isn't there
// if document is in cache, reset expiration timer
func GetDocumentFromCache(uuid string) (models.Document, error) {
	var document models.Document
	if documentJSON, ok := DocumentCache.Get(uuid); ok {
		DocumentCache.Set(uuid, documentJSON, defaultExpiration)
		if err := json.Unmarshal(documentJSON.([]byte), &document); err != nil {
			return models.Document{}, err
		}
	} else {
		var err error
		document, err = GetDocument(uuid)
		if err != nil {
			return models.Document{}, ErrDocumentNotFound
		}
		if err := addDocumentToCache(document); err != nil {
			return models.Document{}, err
		}
	}
	return document, nil
}

func addDocumentToCache(document models.Document) error {
	documentJSON, err := json.Marshal(document)
	if err != nil {
		return ErrFailedToAddDocumentToCache
	}
	DocumentCache.Set(document.UUID, documentJSON, defaultExpiration)
	return nil
}
