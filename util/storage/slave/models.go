package slave

import (
	"encoding/json"

	clientmodels "github.com/TykTechnologies/mserv/mservclient/models"
	"github.com/TykTechnologies/mserv/storage"
)

func clientToStorageMW(clientMw *clientmodels.MW) (*storage.MW, error) {
	marshalled, err := clientMw.MarshalBinary()
	if err != nil {
		return nil, err
	}

	mw := &storage.MW{}
	if err := json.Unmarshal(marshalled, mw); err != nil {
		return nil, err
	}

	return mw, nil
}
