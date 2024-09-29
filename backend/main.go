package main

import (
	"log"
	"net/http"

	"github.com/VincentBaron/beavr_technical_test/backend/internal/config"
	"github.com/VincentBaron/beavr_technical_test/backend/internal/handlers"
	"github.com/VincentBaron/beavr_technical_test/backend/internal/models"
	"github.com/VincentBaron/beavr_technical_test/backend/internal/repositories"
	"github.com/VincentBaron/beavr_technical_test/backend/internal/services"
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

func init() {
	config.ConnectToDb()
	config.SyncDatabase()
}

func main() {
	// Set up the Gin router
	r := gin.New()
	corsConfig := cors.DefaultConfig()
	corsConfig.AllowOrigins = []string{"https://frontend-blue-silence-594.fly.dev"}
	corsConfig.AllowCredentials = true
	corsConfig.AllowMethods = []string{"GET", "POST", "PUT", "PATCH", "DELETE", "OPTIONS"}
	corsConfig.AllowHeaders = []string{"Origin", "Content-Type", "Accept", "Authorization"}

	r.Use(cors.New(corsConfig))

	// Initialize repositories
	requirementsRepo := repositories.NewRepository[models.Requirement](config.DB)
	documentsRepo := repositories.NewRepository[models.Document](config.DB)
	documentsHistoryRepo := repositories.NewRepository[models.DocumentVersions](config.DB)

	// Initialize services
	requirementsService := services.NewRequirementsService(requirementsRepo)
	documentsService := services.NewDocumentsService(documentsRepo, documentsHistoryRepo)

	// Initialize handlers
	requirementsHandler := handlers.NewRequirementsHandler(requirementsService)
	documentsHandler := handlers.NewDocumentsHandler(documentsService)

	// Set up routes with handlers
	r.GET("/requirements", requirementsHandler.List)
	r.GET("/documents", documentsHandler.List)
	r.PATCH("documents/:id", documentsHandler.Update)
	r.POST("/documents/:id/versions", documentsHandler.CreateVersion)
	r.PATCH("/documents/versions/:id", documentsHandler.UpdateVersion)
	r.PATCH("/documents/versions/:id/upload-file", documentsHandler.UploadFile)

	// Start the server
	log.Printf("Server started at http://localhost:8080...")
	log.Fatal(http.ListenAndServe("0.0.0.0:8080", r))
}
