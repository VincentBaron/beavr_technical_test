package models

import (
	"gorm.io/gorm"
)

// Document defines the structure for a document's general infos
type Document struct {
	gorm.Model
	Name          string `gorm:"not null"`
	Description   string `gorm:"type:text"`
	RequirementID uint   // Foreign key for Requirement
	Versions      []DocumentVersions
	Status        Status `gorm:"default:'non-compliant'"`
}

type GetDocumentsParams struct {
	RequirementID *uint
}

// DocumentHistory defines the structure for document version history
type DocumentVersions struct {
	gorm.Model
	DocumentID uint // Foreign key for the current Document
	Version    int
	Path       string `gorm:"not null"`
	Archived   bool   `gorm:"default:false"`
}
