package models

import (
	"gorm.io/gorm"
)

// Document defines the structure for the current version of documents
type Document struct {
	gorm.Model
	Name          string `gorm:"not null"`
	Description   string `gorm:"type:text"`
	RequirementID uint   // Foreign key for Requirement
	Versions      []DocumentVersions
}

// DocumentHistory defines the structure for document version history
type DocumentVersions struct {
	gorm.Model
	DocumentID uint   // Foreign key for the current Document
	Version    int    // Version number
	Status     Status `gorm:"default:'non-compliant'"` // Status of the document
	Path       string `gorm:"not null"`                // Path to the stored document
	Archived   bool   `gorm:"default:false"`
}
