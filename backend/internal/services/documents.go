package services

import (
	"fmt"
	"mime/multipart"
	"path/filepath"
	"strconv"

	"github.com/VincentBaron/beavr_technical_test/backend/internal/models"
	"github.com/VincentBaron/beavr_technical_test/backend/internal/repositories"

	"github.com/gin-gonic/gin"
)

type DocumentsService struct {
	documentsRepo *repositories.Repository[models.Document]
	versionRepo   *repositories.Repository[models.DocumentVersions] // Ensure you have a repository for history
}

func NewDocumentsService(documentsRepo *repositories.Repository[models.Document], versionRepo *repositories.Repository[models.DocumentVersions]) *DocumentsService {
	return &DocumentsService{
		documentsRepo: documentsRepo,
		versionRepo:   versionRepo,
	}
}

// GetDocuments returns a list of documents
func (s *DocumentsService) GetDocuments(c *gin.Context) ([]models.Document, error) {
	documents, err := s.documentsRepo.FindAllByFilter(map[string]interface{}{}, "Versions")
	if err != nil {
		return nil, err
	}
	return documents, nil
}

// UpdateDocument updates the current document and creates a new DocumentHistory entry
func (s *DocumentsService) UpdateVersion(c *gin.Context, ID string, version models.DocumentVersions) error {

	versionID, err := strconv.Atoi(ID)
	if err != nil {
		return err // Return error if the ID is not a valid integer
	}
	version.ID = uint(versionID)
	// Step 3: Insert the history record
	if err := s.versionRepo.Save(&version); err != nil {
		return err // Return error if saving history fails
	}

	return nil
}

// Upload Document uploads a file and updates the document's path
func (s *DocumentsService) UploadDocument(c *gin.Context, versionID string, file *multipart.FileHeader) error {
	// Get the document by ID
	version, err := s.versionRepo.FindByFilter(map[string]interface{}{"id": versionID})
	if err != nil {
		return err // Return error if the document is not found
	}

	extension := filepath.Ext(file.Filename)

	// Save the file to a local folder
	savePath := filepath.Join("uploads", fmt.Sprint(version.DocumentID)+"_"+fmt.Sprint(version.Version+1)+extension)
	if err := c.SaveUploadedFile(file, savePath); err != nil {
		return err // Return error if saving the file fails
	}

	// Step 2: Create a new DocumentHistory entry
	newVersion := models.DocumentVersions{
		DocumentID: version.DocumentID,  // Use the updated document ID
		Version:    version.Version + 1, // Function to determine the next version number
		Path:       savePath,
		Archived:   version.Archived,
	}

	// Step 3: Insert the history record
	if err := s.versionRepo.Save(&newVersion); err != nil {
		return err // Return error if saving history fails
	}

	return nil
}
