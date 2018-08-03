package storage

import (
	"fmt"
	"github.com/TykTechnologies/mserv/coprocess/bindings/go"
	"github.com/TykTechnologies/mserv/util/storage/errors"
	"github.com/TykTechnologies/tyk/apidef"
	"sync"
)

var GlobalRtStore *RuntimeStore

type RuntimeStore struct {
	manifests sync.Map
	functions sync.Map
}

func NewRuntimeStore() *RuntimeStore {
	return &RuntimeStore{
		manifests: sync.Map{},
		functions: sync.Map{},
	}
}

func GenerateStoreKey(org, api, hType, name string) string {
	return fmt.Sprintf("%s.%s", hType, name)
}

func (s *RuntimeStore) GetManifest(apiID string) (*apidef.BundleManifest, error) {
	bm, ok := s.manifests.Load(apiID)
	if !ok {
		return nil, errors.New("not found")
	}

	manifest, ok := bm.(*apidef.BundleManifest)
	if !ok {
		return nil, errors.New("data is not a bundle manifest")
	}

	return manifest, nil
}

func (s *RuntimeStore) GetHookFunc(name string) (func(*coprocess.Object) (*coprocess.Object, error), error) {

	hf, ok := s.functions.Load(name)
	if !ok {
		return nil, errors.New("not found")
	}

	hook, ok := hf.(func(*coprocess.Object) (*coprocess.Object, error))
	if !ok {
		return nil, errors.New("data is not a hook function")
	}

	return hook, nil
}

func (s *RuntimeStore) UpdateOrStoreManifest(apiID string, manifest *apidef.BundleManifest) (bool, error) {
	_, updated := s.manifests.Load(apiID)

	s.manifests.Store(apiID, manifest)

	return updated, nil
}

func (s *RuntimeStore) UpdateOrStoreHook(name string, hook func(*coprocess.Object) (*coprocess.Object, error)) (bool, error) {
	_, updated := s.functions.Load(name)

	s.functions.Store(name, hook)

	return updated, nil
}

type DiffReport struct {
	Added   []*MW
	Removed []*MW
}

func (s *RuntimeStore) FilterNewMW(fetched []*MW) (*DiffReport, error) {
	diff := &DiffReport{
		Added:   make([]*MW, 0),
		Removed: make([]*MW, 0),
	}

	s.manifests.Range(func(key, value interface{}) bool {
		found := false
		for _, fo := range fetched {
			if key == fo.UID {
				found = true
			}
		}

		// In loaded list, but not fetched list, so must be deleted
		if !found {
			diff.Removed = append(diff.Removed, value.(*MW))
		}

		return true
	})

	for _, o := range fetched {
		//fmt.Printf("checking %s\n", o.UID)
		iMw, exists := s.manifests.Load(o.UID)
		if exists {
			mw, ok := iMw.(*MW)
			if !ok {
				return nil, fmt.Errorf("mw not the correct type: %v", o.UID)
			}

			//fmt.Printf("comparing %v to %v\n", mw.Added, o.Added)
			if mw.Added == o.Added {
				// no change, skip
				continue
			}
		}

		// doesn't exist, and has not been updated
		//fmt.Printf("%s does not exist, adding to diff\n", o.UID)
		diff.Added = append(diff.Added, o)
	}

	return diff, nil
}

func (s *RuntimeStore) AddMW(name string, mw *MW) {
	s.manifests.Store(name, mw)
}
