package mongo

import (
	"gopkg.in/mgo.v2/bson"

	"github.com/TykTechnologies/mserv/storage"
	"github.com/TykTechnologies/mserv/util/storage/mongo/mgo_models"
)

// GetMWByID gets middleware from the store based on its UID.
func (m *Store) GetMWByID(id string) (*storage.MW, error) {
	return m.getByID("mw.uid", id)
}

// GetMWByAPIID gets middleware from the store based on its API ID.
func (m *Store) GetMWByAPIID(apiID string) (*storage.MW, error) {
	return m.getByID("mw.apiid", apiID)
}

func (m *Store) getByID(key, value string) (*storage.MW, error) {
	s := m.ms.Copy()
	defer s.Close()

	mm := &mgo_models.MgoMW{}
	if err := s.DB("").C(mservCol).Find(bson.M{key: value}).One(mm); err != nil {
		return nil, err
	}

	return mm.MW, nil
}

// GetAllActive returns all active middleware from the store.
func (m *Store) GetAllActive() ([]*storage.MW, error) {
	s := m.ms.Copy()
	defer s.Close()

	mm := make([]mgo_models.MgoMW, 0)
	if err := s.DB("").C(mservCol).Find(bson.M{"mw.active": true}).All(&mm); err != nil {
		return nil, err
	}

	mws := make([]*storage.MW, len(mm))
	for i, mmw := range mm {
		mws[i] = mmw.MW
	}

	return mws, nil
}

// UpdateMW will update the given middleware in-place in storage.
func (m *Store) UpdateMW(mw *storage.MW) (string, error) {
	s := m.ms.Copy()
	defer s.Close()

	if mw.UID == "" {
		return "", storage.ErrEmptyUID
	}

	mMw := &mgo_models.MgoMW{}
	if err := s.DB("").C(mservCol).Find(bson.M{"mw.uid": mw.UID}).One(mMw); err != nil {
		return "", err
	}

	mMw.MW = mw

	if err := s.DB("").C(mservCol).Update(bson.M{"mw.uid": mw.UID}, mMw); err != nil {
		return "", err
	}

	return mw.UID, nil
}

// CreateMW stores the given middleware.
func (m *Store) CreateMW(mw *storage.MW) (string, error) {
	s := m.ms.Copy()
	defer s.Close()

	if mw.UID == "" {
		return "", storage.ErrEmptyUID
	}

	mMw := &mgo_models.MgoMW{
		MID: bson.NewObjectId(),
		MW:  mw,
	}

	if err := s.DB("").C(mservCol).Insert(mMw); err != nil {
		return "", err
	}

	return mw.UID, nil
}

func (m *Store) DeleteMW(id string) error {
	s := m.ms.Copy()
	defer s.Close()

	if id == "" {
		return storage.ErrEmptyUID
	}

	mMw := &mgo_models.MgoMW{}
	if err := s.DB("").C(mservCol).Find(bson.M{"mw.uid": id}).One(mMw); err != nil {
		return err
	}

	if err := s.DB("").C(mservCol).Remove(bson.M{"mw.uid": id}); err != nil {
		return err
	}

	return nil
}

func (m *Store) InitMservStore(tag string) error {
	m.tag = tag
	return m.Init()
}
