package main

import (
	"encoding/csv"
	"fmt"
	"log"
	"net/http"
	"os"
	"strings"

	"github.com/VincentBaron/beavr_technical_test/backend/internal/config"
	"github.com/VincentBaron/beavr_technical_test/backend/internal/handlers"
	"github.com/VincentBaron/beavr_technical_test/backend/internal/models"
	"github.com/VincentBaron/beavr_technical_test/backend/internal/repositories"
	"github.com/VincentBaron/beavr_technical_test/backend/internal/services"
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"gopkg.in/yaml.v2"
	"gorm.io/gorm"
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

func importCSVToDB(filename string, db *gorm.DB) error {
	// Open the CSV file
	file, err := os.Open(filename)
	if err != nil {
		return err
	}
	defer file.Close()

	// Create a new CSV reader
	reader := csv.NewReader(file)

	// Read all rows from the CSV
	rows, err := reader.ReadAll()
	if err != nil {
		return err
	}

	// Loop through each row and insert data into the database
	for i, row := range rows {
		// Skip the first row (headers)
		if i == 0 {
			continue
		}

		// Assume columns are in the order: Name, Description, Status, Documents (comma-separated)
		if len(row) < 4 {
			log.Printf("Skipping row %d due to missing data: %+v", i, row)
			continue
		}

		// Create the Requirement object
		req := models.Requirement{
			Name:        row[0], // First column: Name
			Description: row[1], // Second column: Description
		}

		// Insert the requirement into the database
		if err := db.Create(&req).Error; err != nil {
			log.Printf("Failed to insert row %d: %+v", i, err)
			continue
		}

		// Parse the document list (fourth column) and create associated Document records
		documentNames := strings.Split(row[2], ",") // Assuming the documents are comma-separated
		for _, docName := range documentNames {
			// Trim space and skip empty document names
			docName = strings.TrimSpace(docName)
			if docName == "" {
				continue
			}

			// Create the Document object and associate it with the requirement
			doc := models.Document{
				Name:          docName,
				Description:   fmt.Sprintf("Document for %s", req.Name),
				RequirementID: req.ID, // Foreign key association
			}

			// Insert the document into the database
			if err := db.Create(&doc).Error; err != nil {
				log.Printf("Failed to insert document for row %d: %+v", i, err)
			}

			// Create a new DocumentHistory entry for this version
			history := models.DocumentVersions{
				DocumentID: doc.ID, // Associate with the current document
				Version:    1,      // First version for this document
			}

			// Insert the history record into the database
			if err := db.Create(&history).Error; err != nil {
				log.Printf("Failed to insert document history for row %d: %+v", i, err)
			}
		}
	}

	return nil
}

// Handler function to upload and process CSV file
func uploadCSVHandler(c *gin.Context) {
	// Get the uploaded file
	file, err := c.FormFile("file")
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "No file is received"})
		return
	}

	// Save the uploaded file to a temporary location
	tempFile := "/tmp/" + file.Filename
	if err := c.SaveUploadedFile(file, tempFile); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to save file"})
		return
	}

	// Import the CSV data into the database
	if err := importCSVToDB(tempFile, config.DB); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to import CSV data"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "CSV data with documents successfully imported into the database!"})
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
	r.POST("/upload-csv", uploadCSVHandler)

	// Start the server
	log.Printf("Server started at http://localhost:8080...")
	log.Fatal(http.ListenAndServe("0.0.0.0:8080", r))
}
