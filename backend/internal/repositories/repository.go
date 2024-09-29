package repositories

import (
	"gorm.io/gorm"
)

// Repository is a generic struct that holds a gorm.DB instance andd simplifies queries to the database using generics
type Repository[T any] struct {
	db *gorm.DB
}

func NewRepository[T any](db *gorm.DB) *Repository[T] {
	return &Repository[T]{db: db}
}

// FIndByFilter finds a single entity by a filter and preloads the specified relationships
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

// FindAllByFilter finds multiple entities by a filter and preloads the specified relationships
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

// Save saves an entity to the database
func (r *Repository[T]) Save(entity *T) error {
	return r.db.Save(entity).Error
}
