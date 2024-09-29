package handlers

import (
	"net/http"

	"github.com/VincentBaron/beavr_technical_test/backend/internal/models"
	"github.com/VincentBaron/beavr_technical_test/backend/internal/services"
	"github.com/gin-gonic/gin"
)

type DocumentsHandler struct {
	documentsService *services.DocumentsService
}

func NewDocumentsHandler(documentsService *services.DocumentsService) *DocumentsHandler {
	return &DocumentsHandler{
		documentsService: documentsService,
	}
}

// List returns a list of documents
func (h *DocumentsHandler) List(c *gin.Context) {
	var params models.GetDocumentsParams

	// Bind query parameters to the GetDocumentsParams struct
	if err := c.ShouldBindQuery(&params); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid query parameters"})
		return
	}

	documents, err := h.documentsService.GetDocuments(c, &params)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	// Return the list of documents
	c.JSON(http.StatusOK, gin.H{"documents": documents})
}

// CreateVersion creates a new version of a document
func (h *DocumentsHandler) CreateVersion(c *gin.Context) {
	documentID := c.Param("id")
	err := h.documentsService.CreateVersion(c, documentID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "Version created successfully"})
}

// Update updates general infors of a document
func (h *DocumentsHandler) Update(c *gin.Context) {
	documentID := c.Param("id")
	var document models.Document
	if err := c.ShouldBindJSON(&document); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	err := h.documentsService.UpdateDocument(c, documentID, document)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "Document updated successfully"})
}

// UpdateVersion updates a document version
func (h *DocumentsHandler) UpdateVersion(c *gin.Context) {
	documentID := c.Param("id")
	var document models.DocumentVersions
	if err := c.ShouldBindJSON(&document); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	err := h.documentsService.UpdateVersion(c, documentID, document)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Document updated successfully"})
}

// UploadFile handles file uploads to a document version and updates the document's path
func (h *DocumentsHandler) UploadFile(c *gin.Context) {
	documentID := c.Param("id")

	// Parse the file from the request
	file, err := c.FormFile("file")
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "File is required"})
		return
	}

	err = h.documentsService.UploadDocument(c, documentID, file)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "File uploaded and document updated successfully"})
}
