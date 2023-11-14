// Package mock provides a store to aid testing of mserv.
package mock

import (
	"context"
	"errors"

	"github.com/TykTechnologies/mserv/storage"
)

// ErrEmptyUID is returned when middleware UID is empty.
var ErrEmptyUID = errors.New("UID cannot be empty")

// Storage is a mock store for testing mserv.
type Storage struct{}

// GetMWByID is a test mock.
func (s *Storage) GetMWByID(_ context.Context, id string) (*storage.MW, error) {
	if id == "" {
		return nil, ErrEmptyUID
	}

	return &storage.MW{UID: id}, nil
}

// GetMWByAPIID is a test mock.
func (s *Storage) GetMWByAPIID(_ context.Context, _ string) (*storage.MW, error) {
	panic("TODO: Implement")
}

// GetAllActive is a test mock.
func (s *Storage) GetAllActive(_ context.Context) ([]*storage.MW, error) {
	panic("TODO: Implement")
}

// CreateMW is a test mock.
func (s *Storage) CreateMW(_ context.Context, mw *storage.MW) (string, error) {
	if mw.UID == "" {
		return "", ErrEmptyUID
	}

	return mw.UID, nil
}

// UpdateMW is a test mock.
func (s *Storage) UpdateMW(_ context.Context, mw *storage.MW) (string, error) {
	panic("TODO: Implement")
}

// DeleteMW is a test mock.
func (s *Storage) DeleteMW(_ context.Context, id string) error {
	return nil
}

// InitMservStore is a test mock.
func (s *Storage) InitMservStore(_ context.Context, tag string) error {
	panic("TODO: Implement")
}
