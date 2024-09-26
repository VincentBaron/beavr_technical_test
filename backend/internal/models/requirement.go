package models

import "gorm.io/gorm"

type Status string

const (
	Compliant    Status = "compliant"
	NonCompliant Status = "non-compliant"
	Pending      Status = "pending"
	Waiting      Status = "waiting"
)

// Requirement defines the CSR requirement model
type Requirement struct {
	gorm.Model
	Name        string `gorm:"not null"`
	Description string `gorm:"type:text"`
	Status      Status
	Documents   []Document
}
