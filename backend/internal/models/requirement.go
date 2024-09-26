package models

import "gorm.io/gorm"

type Status string

const (
	Compliant    Status = "compliant"
	NonCompliant Status = "non-compliant"
	Pending      Status = "pending"
)

// Requirement defines the CSR requirement model
type Requirement struct {
	gorm.Model
	Name        string `gorm:"not null"`
	Description string `gorm:"type:text"`
	Status      Status `gorm:"default:'non-compliant'"` // 'compliant', 'non-compliant', 'pending'
	// Documents   []Document
}
