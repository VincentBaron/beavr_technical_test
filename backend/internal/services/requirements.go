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
	requirements, err := s.requirementsRepository.FindAllByFilter(map[string]interface{}{}, "Documents")
	if err != nil {
		return nil, err
	}

	for i, requirement := range requirements {
		allCompliant := true
		for _, document := range requirement.Documents {
			if document.Status != "compliant" {
				allCompliant = false
				break
			}
		}
		if allCompliant {
			requirements[i].Status = "compliant"
		} else {
			requirements[i].Status = "non-compliant"
		}
	}

	return requirements, nil
}
