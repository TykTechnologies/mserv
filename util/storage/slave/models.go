package slave

import (
	"github.com/TykTechnologies/mserv/models"
	"github.com/TykTechnologies/mserv/storage"
)

type MWPayload struct {
	models.BasePayload
	Payload *storage.MW
}

type AllActiveMWPayload struct {
	models.BasePayload
	Payload []*storage.MW
}
