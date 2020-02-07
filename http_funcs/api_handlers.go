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

// API endpoints
func (h *HttpServ) AddMW(w http.ResponseWriter, r *http.Request) {
	apiID := r.FormValue("api_id")
	if apiID == "" {
		h.HandleError(fmt.Errorf("api_id must be specified"), w, r)
		return
	}

	tmpFileLoc, err := h.ExtractBundleFromPost(r)
	if err != nil {
		h.HandleError(err, w, r)
		return
	}
	log.Info("saved bundle to ", tmpFileLoc)

	storeOnly := r.FormValue("store_only")
	bundleName := uuid.NewV4().String()

	if storeOnly != "" {
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

func (h *HttpServ) FetchAllActiveMW(w http.ResponseWriter, r *http.Request) {
	mws, err := h.api.GetAllActiveMW()
	if err != nil {
		h.HandleError(err, w, r)
		return
	}

	h.HandleOK(mws, w, r)
}
