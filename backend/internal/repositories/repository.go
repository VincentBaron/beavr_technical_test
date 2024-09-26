package repositories

import (
	"github.com/VincentBaron/beavr_technical_test/backend/internal/models"
	"gorm.io/gorm"
)

type Repository[T any] struct {
	db *gorm.DB
}

func NewRepository[T any](db *gorm.DB) *Repository[T] {
	return &Repository[T]{db: db}
}

func (r *Repository[T]) FindByFilter(filter map[string]interface{}, preload ...string) (*T, error) {
	var entity T
	query := r.db.Where(filter)
	for _, p := range preload {
		query = query.Preload(p)
	}
	if err := query.First(&entity).Error; err != nil {
		return nil, err
	}
	return &entity, nil
}

func (r *Repository[T]) FindAllByFilter(filter map[string]interface{}, preload ...string) ([]T, error) {
	var entities []T
	query := r.db.Where(filter)
	for _, p := range preload {
		query = query.Preload(p)
	}
	if err := query.Find(&entities).Error; err != nil {
		return nil, err
	}
	return entities, nil
}

// getNextVersion retrieves the next version number for a given DocumentID
func getNextVersion(documentID uint, db *gorm.DB) int {
	var count int64
	db.Model(&models.DocumentHistory{}).Where("document_id = ?", documentID).Count(&count)
	return int(count) + 1
}

func (r *Repository[T]) Save(entity *T) error {
	return r.db.Save(entity).Error
}
