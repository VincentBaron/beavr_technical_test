package services

import (
	"github.com/VincentBaron/beavr_technical_test/backend/internal/models"
	"github.com/VincentBaron/beavr_technical_test/backend/internal/repositories"
	"github.com/gin-gonic/gin"
)

type RequirementsService struct {
	requirementsRepository *repositories.Repository[models.Requirement]
}

func NewRequirementsService(requirementsRepo *repositories.Repository[models.Requirement]) *RequirementsService {
	return &RequirementsService{
		requirementsRepository: requirementsRepo,
	}
}

func (s *RequirementsService) GetRequirements(c *gin.Context) ([]models.Requirement, error) {
	requirements, err := s.requirementsRepository.FindAllByFilter(map[string]interface{}{})
	if err != nil {
		return nil, err
	}
	return requirements, nil
}
