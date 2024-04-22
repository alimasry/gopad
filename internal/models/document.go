package models

import (
	"time"

	"gorm.io/gorm"
)

type Document struct {
	ID        int
	UUID      string `gorm:"type:uuid;index"`
	Title     string `gorm:"index"`
	Content   string
	Size      int
	Version   int
	CreatedAt time.Time
	UpdatedAt time.Time      `gorm:"index"`
	DeletedAt gorm.DeletedAt `gorm:"index"`
}

type DocumentVersion struct {
	ID       int
	Document `gorm:"embedded"`
}
