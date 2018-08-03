package storage

import (
	"github.com/TykTechnologies/mserv/coprocess/bindings/go"
	"github.com/TykTechnologies/tyk/apidef"
	"time"
)

type Plugin struct {
	UID      string
	Name     string
	FileName string
	FileRef  string
	Type     coprocess.HookType
}

type MW struct {
	UID          string
	APIID        string
	OrgID        string
	Manifest     *apidef.BundleManifest
	Plugins      []*Plugin
	Active       bool
	Added        time.Time
	BundleRef    string
	DownloadOnly bool
}

type MservStore interface {
	GetMWByID(id string) (*MW, error)
	GetMWByApiID(ApiID string) (*MW, error)
	GetAllActive() ([]*MW, error)
	CreateMW(mw *MW) (string, error)
	UpdateMW(mw *MW) (string, error)
	DeleteMW(id string) error
	InitMservStore(tag string) error
}
