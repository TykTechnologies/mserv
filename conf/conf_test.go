package config_test

import (
	"os"
	"testing"

	"github.com/matryer/is"

	config "github.com/TykTechnologies/mserv/conf"
)

// TestGetConf makes sure that envconfig defaults get overridden by values set (in order of preference) first in the
// environment, and then second in the config file.
func TestGetConf(t *testing.T) {
	is := is.New(t)
	is.NoErr(os.Setenv("TYK_MSERV_CONFIG", "testdata/mserv-conf.json"))
	is.NoErr(os.Setenv("MS_MSERV_RETAINUPLOADS", "true"))

	ms := config.GetConf()
	is.Equal(ms.Mserv.FileStore.Kind, "s3")                   // cfg.Mserv.FileStore.Kind mismatch
	is.Equal(ms.Mserv.FileStore.S3.ConfigRegion, "eu-west-2") // cfg.Mserv.FileStore.S3.ConfigRegion mismatch
	is.Equal(ms.Mserv.HTTPAddr, ":8989")                      // cfg.Mserv.HTTPAddr mismatch
	is.Equal(ms.Mserv.MiddlewarePath, "/mnt/mserv/plugins")   // cfg.Mserv.MiddlewarePath mismatch
	is.Equal(ms.Mserv.PluginDir, "/tmp/mserv/plugins")        // cfg.Mserv.PluginDir mismatch
	is.Equal(ms.Mserv.RetainUploads, true)                    // cfg.Mserv.RetainUploads mismatch
}
