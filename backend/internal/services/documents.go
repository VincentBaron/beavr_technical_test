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
	versionRepo   *repositories.Repository[models.DocumentVersions]
}

func NewDocumentsService(documentsRepo *repositories.Repository[models.Document], versionRepo *repositories.Repository[models.DocumentVersions]) *DocumentsService {
	return &DocumentsService{
		documentsRepo: documentsRepo,
		versionRepo:   versionRepo,
	}
}

// GetDocuments returns a list of documents
func (s *DocumentsService) GetDocuments(c *gin.Context, params *models.GetDocumentsParams) ([]models.Document, error) {
	var filter map[string]interface{}
	if params != nil && params.RequirementID != nil {
		filter = map[string]interface{}{
			"requirement_id": *params.RequirementID,
		}
	} else {
		filter = map[string]interface{}{}
	}
	documents, err := s.documentsRepo.FindAllByFilter(filter, "Versions")
	if err != nil {
		return nil, err
	}
	return documents, nil
}

// UpdateDocument updates a document
func (s *DocumentsService) UpdateDocument(c *gin.Context, ID string, document models.Document) error {
	documentID, err := strconv.Atoi(ID)
	if err != nil {
		return err
	}

	document.ID = uint(documentID)
	if err := s.documentsRepo.Save(&document); err != nil {
		return err
	}
	return nil
}

// CreateVersion creates a new version of the document
func (s *DocumentsService) CreateVersion(c *gin.Context, ID string) error {
	document, err := s.documentsRepo.FindByFilter(map[string]interface{}{"id": ID}, "Versions")
	if err != nil {
		return err
	}

	version := models.DocumentVersions{
		DocumentID: document.ID,
		Version:    len(document.Versions) + 1,
	}
	if err := s.versionRepo.Save(&version); err != nil {
		return err
	}
	return nil
}

// UpdateVersion updates a document version
func (s *DocumentsService) UpdateVersion(c *gin.Context, ID string, version models.DocumentVersions) error {
	versionID, err := strconv.Atoi(ID)
	if err != nil {
		return err
	}

	version.ID = uint(versionID)
	if err := s.versionRepo.Save(&version); err != nil {
		return err
	}

	return nil
}

// Upload Document uploads a file to a document version and updates the document's path
func (s *DocumentsService) UploadDocument(c *gin.Context, versionID string, file *multipart.FileHeader) error {
	version, err := s.versionRepo.FindByFilter(map[string]interface{}{"id": versionID})
	if err != nil {
		return err
	}
	document, err := s.documentsRepo.FindByFilter(map[string]interface{}{"id": version.DocumentID})
	if err != nil {
		return err
	}

	//  Save the file to the uploads directory
	extension := filepath.Ext(file.Filename)
	savePath := filepath.Join("uploads", fmt.Sprint(document.Name)+"_"+fmt.Sprint(version.Version+1)+extension)
	if err := c.SaveUploadedFile(file, savePath); err != nil {
		return err
	}

	newVersion := models.DocumentVersions{
		DocumentID: version.DocumentID,
		Version:    version.Version + 1,
		Path:       savePath,
		Archived:   version.Archived,
	}

	if err := s.versionRepo.Save(&newVersion); err != nil {
		return err
	}

	return nil
}
