package storage

import (
	"errors"
	"nowyouseeme/models"
	"sync"
)

// MemoryStore is an in-memory storage implementation
type MemoryStore struct {
	visualizations map[string]*models.Visualization
	mu             sync.RWMutex
}

// NewMemoryStore creates a new in-memory store
func NewMemoryStore() *MemoryStore {
	return &MemoryStore{
		visualizations: make(map[string]*models.Visualization),
	}
}

// CreateVisualization stores a new visualization
func (s *MemoryStore) CreateVisualization(v *models.Visualization) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	if _, exists := s.visualizations[v.ID]; exists {
		return errors.New("visualization with this ID already exists")
	}

	s.visualizations[v.ID] = v
	return nil
}

// GetVisualization retrieves a visualization by ID
func (s *MemoryStore) GetVisualization(id string) (*models.Visualization, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	v, exists := s.visualizations[id]
	if !exists {
		return nil, errors.New("visualization not found")
	}

	return v, nil
}

// GetAllVisualizations retrieves all visualizations
func (s *MemoryStore) GetAllVisualizations() []*models.Visualization {
	s.mu.RLock()
	defer s.mu.RUnlock()

	result := make([]*models.Visualization, 0, len(s.visualizations))
	for _, v := range s.visualizations {
		result = append(result, v)
	}

	return result
}
