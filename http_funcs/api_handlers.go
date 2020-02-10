package http_funcs

import (
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"os"
	"path"

	"github.com/gorilla/mux"
	"github.com/jpillora/overseer"
	"github.com/satori/go.uuid"
)

func (h *HttpServ) ExtractBundleFromPost(r *http.Request) (string, error) {
	// Save the file to disk
	if err := r.ParseMultipartForm(32 << 20); err != nil {
		return "", err
	}

	uploadedFile, _, err := r.FormFile("uploadfile")
	if err != nil {
		return "", err
	}

	defer uploadedFile.Close()

	tmpDir := path.Join(os.TempDir(), "mserv-bundles")
	if err := os.Mkdir(tmpDir, 0700); err != nil {
		if !os.IsExist(err) {
			return "", err
		}
	}

	tmpFile, err := ioutil.TempFile(tmpDir, "bundle-*.zip")
	if err != nil {
		return "", err
	}
	defer tmpFile.Close()

	_, err = io.Copy(tmpFile, uploadedFile)
	if err != nil {
		return "", err
	}

	return tmpFile.Name(), nil
}

// swagger:route POST /api/mw mw mwAdd
// Adds a new middleware. If `store_only` field is true it will only be available for download.
// Expects a file bundle in `uploadfile` form field.
//
// Security:
//   api_key:
//
// Responses:
//   200: mwIDResponse
//   500: genericErrorResponse
func (h *HttpServ) AddMW(w http.ResponseWriter, r *http.Request) {
	apiID := r.FormValue("api_id")

	tmpFileLoc, err := h.ExtractBundleFromPost(r)
	if err != nil {
		h.HandleError(err, w, r)
		return
	}
	log.Info("saved bundle to ", tmpFileLoc)

	storeOnly := r.FormValue("store_only")
	bundleName := uuid.NewV4().String()

	if storeOnly == "true" {
		// This is a python or JS bundle, just proxy it to a store
		mw, err := h.api.StoreBundleOnly(tmpFileLoc, apiID, bundleName)
		if err != nil {
			h.HandleError(err, w, r)
			return
		}

		ret := map[string]interface{}{"BundleID": mw.UID}
		h.HandleOK(ret, w, r)
		return
	}

	// This is a plugin bundle (.so) so we should process it differently
	mw, err := h.api.HandleNewBundle(tmpFileLoc, apiID, bundleName)
	if err != nil {
		h.HandleError(err, w, r)
		return
	}

	ret := map[string]interface{}{"BundleID": mw.UID}
	h.HandleOK(ret, w, r)
}

// swagger:route PUT /api/mw/{id} mw mwUpdate
// Updates a middleware specified by {id}.
// Expects a file bundle in `uploadfile` form field.
//
// Security:
//   api_key:
//
// Responses:
//   200: mwIDResponse
//   500: genericErrorResponse
func (h *HttpServ) UpdateMW(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id, ok := vars["id"]
	if !ok {
		h.HandleError(fmt.Errorf("bundle_id must be specified"), w, r)
		return
	}

	// We do not want the generated bundle name since we already have a reference ID
	tmpFileLoc, err := h.ExtractBundleFromPost(r)
	if err != nil {
		h.HandleError(err, w, r)
		return
	}

	mw, err := h.api.HandleUpdateBundle(tmpFileLoc, id)
	if err != nil {
		h.HandleError(err, w, r)
		return
	}

	h.HandleOK(map[string]interface{}{"BundleID": mw.UID}, w, r)
}

// swagger:route DELETE /api/mw/{id} mw mwDelete
// Deletes a middleware specified by {id}.
//
// Security:
//   api_key:
//
// Responses:
//   200: mwIDResponse
//   500: genericErrorResponse
func (h *HttpServ) DeleteMW(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id, ok := vars["id"]
	if !ok {
		h.HandleError(fmt.Errorf("bundle id must be specified"), w, r)
		return
	}

	err := h.api.HandleDeleteBundle(id)
	if err != nil {
		h.HandleError(err, w, r)
		return
	}

	h.HandleOK(map[string]interface{}{"BundleID": id}, w, r)

	overseer.Restart()
}

// swagger:route GET /api/mw/{id} mw mwFetch
// Fetches a middleware specified by {id}.
//
// Security:
//   api_key:
//
// Responses:
//   200: mwResponse
//   500: genericErrorResponse
func (h *HttpServ) FetchMW(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id, ok := vars["id"]
	if !ok {
		h.HandleError(fmt.Errorf("id must be specified"), w, r)
		return
	}

	dat, err := h.api.GetMWByID(id)
	if err != nil {
		h.HandleError(err, w, r)
		return
	}

	h.HandleOK(dat, w, r)
}

// swagger:route GET /api/mw/bundle/{id} mw mwFetchBundle
// Fetches a middleware bundle file specified by {id}.
//
// Produces:
// - application/octet-stream
//
// Security:
//   api_key:
//
// Responses:
//   200: mwBundleResponse
//   500: genericErrorResponse
func (h *HttpServ) FetchBundleFile(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id, ok := vars["id"]
	if !ok {
		h.HandleError(fmt.Errorf("id must be specified"), w, r)
		return
	}

	dat, err := h.api.GetMWByID(id)
	if err != nil {
		h.HandleError(err, w, r)
		return
	}

	bf, err := h.api.FetchAndServeBundleFile(dat)
	if err != nil {
		h.HandleError(err, w, r)
		return
	}

	// Serve the file only
	http.ServeFile(w, r, bf)
}

// swagger:route GET /api/mw/master/all mw mwListAll
// Lists all middleware.
//
// Security:
//   api_key:
//
// Responses:
//   200: mwListResponse
//   500: genericErrorResponse
func (h *HttpServ) FetchAllActiveMW(w http.ResponseWriter, r *http.Request) {
	mws, err := h.api.GetAllActiveMW()
	if err != nil {
		h.HandleError(err, w, r)
		return
	}

	h.HandleOK(mws, w, r)
}
