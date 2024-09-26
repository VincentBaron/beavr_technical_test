package models

import "gorm.io/gorm"

type Document struct {
	gorm.Model
	Name          string `gorm:"not null"`
	Description   string `gorm:"type:text"`
	RequirementID uint
}
