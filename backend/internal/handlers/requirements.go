package handlers

import (
	"net/http"

	"github.com/VincentBaron/beavr_technical_test/backend/internal/models"
	"github.com/VincentBaron/beavr_technical_test/backend/internal/services"
	"github.com/gin-gonic/gin"
)

type RequirementsHandler struct {
	requirementsService *services.RequirementsService
}

func NewRequirementsHandler(requirementsService *services.RequirementsService) *RequirementsHandler {
	return &RequirementsHandler{
		requirementsService: requirementsService,
	}
}

// List returns a list of requirements
func (h *RequirementsHandler) List(c *gin.Context) {
	var requirements []models.Requirement

	requirements, err := h.requirementsService.GetRequirements(c)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	// Return the list of playlist names
	c.JSON(http.StatusOK, gin.H{"requirements": requirements})
}
