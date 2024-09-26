package config

import "github.com/VincentBaron/beavr_technical_test/backend/internal/models"

func SyncDatabase() {
	DB.AutoMigrate(&models.Requirement{}, &models.Document{}, &models.DocumentHistory{})
}
