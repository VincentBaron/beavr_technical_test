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

func (h *DocumentsHandler) List(c *gin.Context) {

	// Create a new slice to store the playlist names
	var documents []models.Document

	documents, err := h.documentsService.GetDocuments(c)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	// Return the list of playlist names
	c.JSON(http.StatusOK, gin.H{"documents": documents})
}

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

// UploadFile handles file uploads and updates the document's path
func (h *DocumentsHandler) UploadFile(c *gin.Context) {
	// Get the document ID from the endpoint path
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
