package services

import (
	"fmt"
	"mime/multipart"
	"path/filepath"

	"github.com/VincentBaron/beavr_technical_test/backend/internal/models"
	"github.com/VincentBaron/beavr_technical_test/backend/internal/repositories"

	"github.com/gin-gonic/gin"
)

type DocumentsService struct {
	documentsRepo *repositories.Repository[models.Document]
	historyRepo   *repositories.Repository[models.DocumentHistory] // Ensure you have a repository for history
}

func NewDocumentsService(documentsRepo *repositories.Repository[models.Document], historyRepo *repositories.Repository[models.DocumentHistory]) *DocumentsService {
	return &DocumentsService{
		documentsRepo: documentsRepo,
		historyRepo:   historyRepo,
	}
}

// GetDocuments returns a list of documents
func (s *DocumentsService) GetDocuments(c *gin.Context) ([]models.Document, error) {
	documents, err := s.documentsRepo.FindAllByFilter(map[string]interface{}{})
	if err != nil {
		return nil, err
	}
	return documents, nil
}

// UpdateDocument updates the current document and creates a new DocumentHistory entry
func (s *DocumentsService) UpdateDocument(c *gin.Context, ID string, document models.Document) error {
	// Check if a file is provided in the request
	file, err := c.FormFile("file")
	if err == nil {
		// Save the file to a local folder
		savePath := filepath.Join("uploads", file.Filename)
		if err := c.SaveUploadedFile(file, savePath); err != nil {
			return err // Return error if saving the file fails
		}
		// Update the document's path with the saved file's path
		document.Path = savePath
	}

	currentDocument, err := s.documentsRepo.FindByFilter(map[string]interface{}{"id": ID})
	if err != nil {
		return err // Return error if the document is not found
	}
	document.Version = currentDocument.Version + 1

	// Step 1: Update the document
	err = s.documentsRepo.Save(&document)
	if err != nil {
		return err // Return error if the update fails
	}

	// Step 2: Create a new DocumentHistory entry
	historyEntry := models.DocumentHistory{
		DocumentID:  document.ID,                 // Use the updated document ID
		Version:     currentDocument.Version + 1, // Function to determine the next version number
		Name:        document.Name,
		Description: document.Description,
		Path:        document.Path,
		Archived:    document.Archived,
	}

	// Step 3: Insert the history record
	if err := s.historyRepo.Save(&historyEntry); err != nil {
		return err // Return error if saving history fails
	}

	return nil
}

// Upload Document uploads a file and updates the document's path
func (s *DocumentsService) UploadDocument(c *gin.Context, documentID string, file *multipart.FileHeader) error {
	// Get the document by ID
	document, err := s.documentsRepo.FindByFilter(map[string]interface{}{"id": documentID})
	if err != nil {
		return err // Return error if the document is not found
	}

	extension := filepath.Ext(file.Filename)

	// Save the file to a local folder
	savePath := filepath.Join("uploads", document.Name+"_"+fmt.Sprint(document.Version+1)+extension)
	if err := c.SaveUploadedFile(file, savePath); err != nil {
		return err // Return error if saving the file fails
	}

	// Update the document's path with the saved file's path
	document.Path = savePath
	document.Version = document.Version + 1
	document.Status = models.Pending

	// Update the document in the database
	if err := s.documentsRepo.Save(document); err != nil {
		return err // Return error if updating the document fails
	}

	// Step 2: Create a new DocumentHistory entry
	historyEntry := models.DocumentHistory{
		DocumentID:  document.ID,          // Use the updated document ID
		Version:     document.Version + 1, // Function to determine the next version number
		Name:        document.Name,
		Description: document.Description,
		Path:        document.Path,
		Archived:    document.Archived,
	}

	// Step 3: Insert the history record
	if err := s.historyRepo.Save(&historyEntry); err != nil {
		return err // Return error if saving history fails
	}

	return nil
}
