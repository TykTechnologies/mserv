package storage

import (
	"time"

	"github.com/TykTechnologies/tyk/apidef"

	coprocess "github.com/TykTechnologies/mserv/coprocess/bindings/go"
)

type Plugin struct {
	UID      string
	Name     string
	FileName string
	FileRef  string
	Type     coprocess.HookType
}

type MW struct {
	Added        time.Time
	Manifest     *apidef.BundleManifest
	APIID        string
	OrgID        string
	UID          string
	BundleRef    string
	Plugins      []*Plugin
	Active       bool
	DownloadOnly bool
}

type MservStore interface {
	GetMWByID(id string) (*MW, error)
	GetMWByAPIID(APIID string) (*MW, error)
	GetAllActive() ([]*MW, error)
	CreateMW(mw *MW) (string, error)
	UpdateMW(mw *MW) (string, error)
	DeleteMW(id string) error
	InitMservStore(tag string) error
}
