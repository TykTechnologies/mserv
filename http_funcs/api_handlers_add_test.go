package http_funcs_test

import (
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"path/filepath"
	"testing"

	"github.com/matryer/is"

	config "github.com/TykTechnologies/mserv/conf"
)

func TestAddMWStoreBundleOnly(t *testing.T) {
	is := is.New(t)
	srv := setupServerAndTempDir(t)

	for name, tc := range addMWTestCases {
		name, tc := name, tc
		t.Run(name, func(t *testing.T) {
			req := prepareAddRequest(t, tc.testBodyBytes)
			req.Form.Add("store_only", "true") // target the 'StoreBundleOnly' code path

			// Handle the request
			w := httptest.NewRecorder()
			srv.AddMW(w, req)

			// Parse the response
			resp := w.Result()
			defer is.NoErr(resp.Body.Close()) // could not close response body cleanly

			is.Equal(tc.expectedStatus, resp.StatusCode) // expected response status does not equal actual

			bod, errRead := ioutil.ReadAll(resp.Body)
			is.NoErr(errRead) // could not read response body
			t.Logf("response: %s %s", resp.Status, bod)
		})
	}
}

func TestAddMWHandleNewBundle(t *testing.T) {
	is := is.New(t)
	srv := setupServerAndTempDir(t)

	t.Run("Compressed (ZIP) upload is OK", func(t *testing.T) {
		fileCountPath := filepath.Join(config.GetConf().Mserv.MiddlewarePath, "plugins")
		startCount, err := ioutil.ReadDir(fileCountPath)
		is.NoErr(err) // could not read 'config.Mserv.MiddlewarePath+"/plugins"' directory

		req := prepareAddRequest(t, compressedTestData)
		req.Form.Add("store_only", "false") // target the 'HandleNewBundle' code path

		// Handle the request
		w := httptest.NewRecorder()
		srv.AddMW(w, req)

		// Parse the response
		resp := w.Result()
		defer is.NoErr(resp.Body.Close()) // could not close response body cleanly

		is.Equal(http.StatusOK, resp.StatusCode) // expected response status does not equal actual

		bod, errRead := ioutil.ReadAll(resp.Body)
		is.NoErr(errRead) // could not read response body
		t.Logf("response: %s %s", resp.Status, bod)

		finishCount, err := ioutil.ReadDir(fileCountPath)
		is.NoErr(err)                               // could not read 'config.Mserv.MiddlewarePath+"/plugins"' directory
		is.Equal(len(startCount), len(finishCount)) // should not leave uploads behind unless configured to do so
	})
}
