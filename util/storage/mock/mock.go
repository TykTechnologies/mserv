// Package mock provides a store to aid testing of mserv.
package mock

import (
	"github.com/TykTechnologies/mserv/storage"
)

// Storage is a mock store for testing mserv.
type Storage struct{}

// GetMWByID is a test mock.
func (s *Storage) GetMWByID(id string) (*storage.MW, error) {
	panic("TODO: Implement")
}

// GetMWByAPIID is a test mock.
func (s *Storage) GetMWByAPIID(apiID string) (*storage.MW, error) {
	panic("TODO: Implement")
}

// GetAllActive is a test mock.
func (s *Storage) GetAllActive() ([]*storage.MW, error) {
	panic("TODO: Implement")
}

// CreateMW is a test mock.
func (s *Storage) CreateMW(mw *storage.MW) (string, error) {
	if mw.UID == "" {
		return "", storage.ErrEmptyUID
	}

	return mw.UID, nil
}

// UpdateMW is a test mock.
func (s *Storage) UpdateMW(mw *storage.MW) (string, error) {
	panic("TODO: Implement")
}

// DeleteMW is a test mock.
func (s *Storage) DeleteMW(id string) error {
	panic("TODO: Implement")
}

// InitMservStore is a test mock.
func (s *Storage) InitMservStore(tag string) error {
	panic("TODO: Implement")
}
