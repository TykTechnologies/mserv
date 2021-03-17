package http_funcs_test

import (
	"encoding/json"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"path/filepath"
	"testing"

	"github.com/matryer/is"

	config "github.com/TykTechnologies/mserv/conf"
	"github.com/TykTechnologies/mserv/models"
)

func TestDeleteMW(t *testing.T) {
	is := is.New(t)
	srv := setupServerAndTempDir(t)

	t.Run("Deleted middleware does not leave stow.Container behind", func(t *testing.T) {
		// Paths to check for things to make sure the handlers are behaving and cleaning up properly
		fileCountPath := filepath.Join(config.GetConf().Mserv.MiddlewarePath, "plugins")
		localContainerPath := config.GetConf().Mserv.FileStore.Local.ConfigKeyPath

		startFileCount, err := ioutil.ReadDir(fileCountPath)
		is.NoErr(err) // could not read 'config.Mserv.MiddlewarePath+"/plugins"' directory

		addReq := prepareAddRequest(t, compressedTestData)
		addReq.Form.Add("store_only", "false") // target the 'HandleNewBundle' code path

		// Handle the AddMW request
		w := httptest.NewRecorder()
		srv.AddMW(w, addReq)

		// Parse the AddMW response
		addResp := w.Result()
		defer is.NoErr(addResp.Body.Close())        // could not close AddMW response body cleanly
		is.Equal(http.StatusOK, addResp.StatusCode) // expected response status does not equal actual from AddMW

		addBod, errRead := ioutil.ReadAll(addResp.Body)
		is.NoErr(errRead) // could not read response body
		t.Logf("response from %s %s: %s %s", addReq.Method, addReq.URL, addResp.Status, addBod)

		// Confirm that stow created one new container in this test's temp directory
		startContainerCount, err := ioutil.ReadDir(localContainerPath)
		is.NoErr(err)                         // could not read directory
		is.Equal(1, len(startContainerCount)) // should have one stow.Container after calling AddMW

		// Get ID of added middleware
		payload := &models.Payload{}
		is.NoErr(json.Unmarshal(addBod, payload)) // could not unmarshal response payload

		internalPayload, ok := payload.Payload.(map[string]interface{})
		is.True(ok) // could not assert type on internal payload

		addedMWID, hasID := internalPayload["BundleID"]
		is.True(hasID) // internal payload does not contain "BundleID"

		sAddedMWID, ok := addedMWID.(string)
		is.True(ok) // could not assert type on added middleware ID
		t.Logf("ID of newly-added middleware: '%s'", sAddedMWID)

		// Do the DeleteMW things
		deleteReq := prepareDeleteRequest(t, sAddedMWID)
		w = httptest.NewRecorder() // reinitialise
		srv.DeleteMW(w, deleteReq)

		// Parse the DeleteMW response
		deleteResp := w.Result()
		defer is.NoErr(deleteResp.Body.Close())        // could not close DeleteMW response body cleanly
		is.Equal(http.StatusOK, deleteResp.StatusCode) // expected response status does not equal actual from DeleteMW

		deleteBod, errRead := ioutil.ReadAll(deleteResp.Body)
		is.NoErr(errRead) // could not read response body
		t.Logf("response from %s %s: %s %s", deleteReq.Method, deleteReq.URL, deleteResp.Status, deleteBod)

		finishFileCount, err := ioutil.ReadDir(fileCountPath)
		is.NoErr(err) // could not read 'config.Mserv.MiddlewarePath+"/plugins"' directory

		is.Equal(len(startFileCount), len(finishFileCount)) // should not leave uploads behind unless configured to do so

		// Confirm that DeleteMW directed stow to clean up after itself
		finishContainerCount, err := ioutil.ReadDir(localContainerPath)
		is.NoErr(err)                          // could not read directory
		is.Equal(0, len(finishContainerCount)) // should have zero stow.Container after calling DeleteMW
	})
}
