package main

import (
	"encoding/csv"
	"fmt"
	"log"
	"os"
	"strings"

	"github.com/VincentBaron/beavr_technical_test/backend/internal/config"
	"github.com/VincentBaron/beavr_technical_test/backend/internal/models"
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

// Import CSV to the DB with associated documents
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
			Name:        row[0],         // First column: Name
			Description: row[1],         // Second column: Description
			Status:      models.Pending, // Default status (you can change this to parse it from the row[2] if required)
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
				Description:   fmt.Sprintf("Auto-generated document for %s", req.Name),
				RequirementID: req.ID, // Foreign key association
			}

			// Insert the document into the database
			if err := db.Create(&doc).Error; err != nil {
				log.Printf("Failed to insert document for row %d: %+v", i, err)
			}
		}
	}

	return nil
}

func main() {

	// Path to your CSV file
	filename := "/Users/vincentbaron/Downloads/beavr.csv"

	// Import the CSV data into the database
	if err := importCSVToDB(filename, config.DB); err != nil {
		log.Fatalf("Failed to import CSV data: %v", err)
	}

	fmt.Println("CSV data with documents successfully imported into the database!")
}
