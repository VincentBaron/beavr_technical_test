package models

import (
	"gorm.io/gorm"
)

// Document defines the structure for the current version of documents
type Document struct {
	gorm.Model
	Name          string `gorm:"not null"`
	Description   string `gorm:"type:text"`
	Path          string `gorm:"not null"` // Path to the stored document
	Status        Status // 'compliant', 'non-compliant', 'pending'
	Version       int    `gorm:"default:1"` // Version number
	RequirementID uint   // Foreign key for Requirement
}

// DocumentHistory defines the structure for document version history
type DocumentHistory struct {
	gorm.Model
	DocumentID  uint   // Foreign key for the current Document
	Version     int    // Version number
	Name        string `gorm:"not null"`
	Description string `gorm:"type:text"`
	Path        string `gorm:"not null"` // Path to the stored document
}
