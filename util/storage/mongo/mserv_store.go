package mongo

import (
	"errors"
	mservStorage "github.com/TykTechnologies/mserv/storage"
	"github.com/TykTechnologies/mserv/util/storage/mongo/mgo_models"
	"gopkg.in/mgo.v2/bson"
)

func (m *Store) GetMWByID(id string) (*mservStorage.MW, error) {
	s := m.ms.Copy()
	defer s.Close()

	mm := &mgo_models.MgoMW{}
	if err := s.DB("").C(mservCol).Find(bson.M{"mw.uid": id}).One(mm); err != nil {
		return nil, err
	}

	return mm.MW, nil
}

func (m *Store) GetMWByApiID(ApiID string) (*mservStorage.MW, error) {
	s := m.ms.Copy()
	defer s.Close()

	mm := &mgo_models.MgoMW{}
	if err := s.DB("").C(mservCol).Find(bson.M{"mw.apiid": ApiID}).One(mm); err != nil {
		return nil, err
	}

	return mm.MW, nil
}
func (m *Store) GetAllActive() ([]*mservStorage.MW, error) {
	s := m.ms.Copy()
	defer s.Close()

	mm := make([]mgo_models.MgoMW, 0)
	if err := s.DB("").C(mservCol).Find(bson.M{"mw.active": true}).All(&mm); err != nil {
		return nil, err
	}

	mws := make([]*mservStorage.MW, len(mm))
	for i, mmw := range mm {
		mws[i] = mmw.MW
	}

	return mws, nil
}

func (m *Store) UpdateMW(mw *mservStorage.MW) (string, error) {
	s := m.ms.Copy()
	defer s.Close()

	if mw.UID == "" {
		return "", errors.New("UID cannot be empty")
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

func (m *Store) CreateMW(mw *mservStorage.MW) (string, error) {
	s := m.ms.Copy()
	defer s.Close()

	if mw.UID == "" {
		return "", errors.New("UID cannot be empty")
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
		return errors.New("UID cannot be empty")
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
