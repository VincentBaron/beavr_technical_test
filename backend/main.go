package main

import (
	"log"
	"net/http"
	"os"

	"github.com/VincentBaron/beavr_technical_test/backend/internal/config"
	"github.com/VincentBaron/beavr_technical_test/backend/internal/handlers"
	"github.com/VincentBaron/beavr_technical_test/backend/internal/models"
	"github.com/VincentBaron/beavr_technical_test/backend/internal/repositories"
	"github.com/VincentBaron/beavr_technical_test/backend/internal/services"
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"gopkg.in/yaml.v2"
)

func init() {
	config.LoadEnvVariables()
	config.ConnectToDb()
	config.SyncDatabase()
}

func LoadConfig(file string) (models.Config, error) {
	var config models.Config
	data, err := os.ReadFile(file)
	if err != nil {
		return config, err
	}
	err = yaml.Unmarshal(data, &config)
	return config, err
}

func main() {
	// Set up the Gin router
	r := gin.New()
	corsConfig := cors.DefaultConfig()
	corsConfig.AllowOrigins = []string{"http://localhost:5173"}

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
	r.POST("/documents/:id/versions", documentsHandler.CreateVersion)
	r.PATCH("/documents/versions/:id", documentsHandler.UpdateVersion)
	r.PATCH("/documents/versions/:id/upload-file", documentsHandler.UploadFile)

	// Start the server
	log.Printf("Server started at http://localhost:8080...")
	log.Fatal(http.ListenAndServe("localhost:8080", r))
}
