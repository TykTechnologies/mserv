package storage

import (
	"context"
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
	GetMWByID(ctx context.Context, id string) (*MW, error)
	GetMWByAPIID(ctx context.Context, APIID string) (*MW, error)
	GetAllActive(ctx context.Context) ([]*MW, error)
	CreateMW(ctx context.Context, mw *MW) (string, error)
	UpdateMW(ctx context.Context, mw *MW) (string, error)
	DeleteMW(ctx context.Context, id string) error
	InitMservStore(ctx context.Context, tag string) error
}
