package editor

import (
	"fmt"

	"github.com/alimasry/gopad/internal/database"
	"github.com/alimasry/gopad/internal/models"
)

var DocumentMap = make(map[string]models.Document)

func GetDocument(documentUuid string) (*models.Document, error) {
	db := database.GetDb()

	var response models.Document

	if err := db.Model(&models.Document{}).Where("uuid = ?", documentUuid).Find(&response).Error; err != nil {
		return nil, fmt.Errorf("failed to retrieve document : %v", err)
	}

	return &response, nil
}

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

	DocumentMap[document.UUID] = document

	return tx.Commit().Error
}

func GetDocumentFromMap(uuid string) models.Document {
	document, ok := DocumentMap[uuid]
	if !ok {
		document, _ := GetDocument(uuid)
		DocumentMap[uuid] = *document
	}
	return document
}
