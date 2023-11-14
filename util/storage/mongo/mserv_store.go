// Package mongo implements Mserv MongoDB storage.
package mongo

import (
	"context"
	"fmt"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"

	"github.com/TykTechnologies/mserv/storage"
)

// InitMservStore initializes Mserv storage.
func (m *Store) InitMservStore(_ context.Context, tag string) error {
	m.tag = tag

	return m.Init()
}

// getByKey func fetches middleware from database for give key/value pair.
func (m *Store) getByKey(ctx context.Context, key, value string) (*storage.MW, error) {
	mm := mgoMW{}

	f := bson.M{key: value}

	if err := m.db.Collection(mservCol).FindOne(ctx, f).Decode(&mm); err != nil {
		return nil, err
	}

	return mm.MW, nil
}

// GetMWByID gets middleware from the store based on its UID.
func (m *Store) GetMWByID(ctx context.Context, id string) (*storage.MW, error) {
	return m.getByKey(ctx, "mw.uid", id)
}

// GetMWByAPIID gets middleware from the store based on its API ID.
func (m *Store) GetMWByAPIID(ctx context.Context, apiID string) (*storage.MW, error) {
	return m.getByKey(ctx, "mw.apiid", apiID)
}

// GetAllActive returns all active middleware from the store.
func (m *Store) GetAllActive(ctx context.Context) ([]*storage.MW, error) {
	mm := make([]mgoMW, 0)

	f := bson.M{"mw.active": true}

	cur, err := m.db.Collection(mservCol).Find(ctx, f)
	if err != nil {
		return nil, fmt.Errorf("collect error: %w", err)
	}

	if err := cur.All(ctx, &mm); err != nil {
		return nil, fmt.Errorf("cursor fetch error: %w", err)
	}

	mws := make([]*storage.MW, len(mm))
	for i, mmw := range mm {
		mws[i] = mmw.MW
	}

	return mws, nil
}

// CreateMW stores the given middleware.
func (m *Store) CreateMW(ctx context.Context, mw *storage.MW) (string, error) {
	if mw.UID == "" {
		return "", ErrEmptyUID
	}

	mMw := mgoMW{
		MID: primitive.NewObjectID(),
		MW:  mw,
	}

	if _, err := m.db.Collection(mservCol).InsertOne(ctx, mMw); err != nil {
		return "", fmt.Errorf("insert error: %w", err)
	}

	return mw.UID, nil
}

// UpdateMW will update the given middleware in-place in storage.
func (m *Store) UpdateMW(ctx context.Context, mw *storage.MW) (string, error) {
	if mw.UID == "" {
		return "", ErrEmptyUID
	}

	mMw := mgoMW{}

	f := bson.M{"mw.uid": mw.UID}

	if err := m.db.Collection(mservCol).FindOne(ctx, f).Decode(&mMw); err != nil {
		return "", fmt.Errorf("find error: %w", err)
	}

	mMw.MW = mw

	update := bson.M{
		"$set": mMw,
	}

	res, err := m.db.Collection(mservCol).UpdateOne(ctx, f, update)
	if err != nil {
		return "", fmt.Errorf("update error: %w", err)
	}

	if res.MatchedCount == 0 {
		return "", ErrNotFound
	}

	return mw.UID, nil
}

// DeleteMW removes given middleware.
func (m *Store) DeleteMW(ctx context.Context, id string) error {
	if id == "" {
		return ErrEmptyUID
	}

	mMw := mgoMW{}

	f := bson.M{"mw.uid": id}

	if err := m.db.Collection(mservCol).FindOne(ctx, f).Decode(&mMw); err != nil {
		return fmt.Errorf("find error: %w", err)
	}

	if _, err := m.db.Collection(mservCol).DeleteOne(ctx, f); err != nil {
		return fmt.Errorf("delete error: %w", err)
	}

	return nil
}
