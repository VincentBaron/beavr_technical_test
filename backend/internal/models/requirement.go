package models

import "gorm.io/gorm"

type Status string

const (
	Compliant    Status = "compliant"
	NonCompliant Status = "non-compliant"
)

// Requirement defines the CSR requirement model
type Requirement struct {
	gorm.Model
	Name        string `gorm:"not null"`
	Description string `gorm:"type:text"`
	Status      Status `gorm:"-"`
	Documents   []Document
}
