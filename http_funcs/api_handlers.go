package http_funcs

import (
	"fmt"
	"github.com/gorilla/mux"
	"github.com/jpillora/overseer"
	"github.com/satori/go.uuid"
	"io"
	"net/http"
	"os"
	"path"
)

func (h *HttpServ) ExtractBundleFromPost(w http.ResponseWriter, r *http.Request) (string, string) {
	// Save the file to disk
	r.ParseMultipartForm(32 << 20)
	file, handler, err := r.FormFile("uploadfile")
	if err != nil {
		h.HandleError(err, w, r)
		return "", ""
	}

	defer file.Close()

	bundleName := uuid.NewV4().String()

	dir := path.Join("./tmp", bundleName)
	os.Mkdir(dir, os.ModePerm)

	tmpFileLoc := path.Join(dir, handler.Filename)
	f, err := os.OpenFile(tmpFileLoc, os.O_WRONLY|os.O_CREATE, 0666)
	if err != nil {
		h.HandleError(err, w, r)
		return "", ""
	}

	defer f.Close()

	_, err = io.Copy(f, file)
	if err != nil {
		h.HandleError(err, w, r)
		return "", ""
	}

	return tmpFileLoc, bundleName
}

// API endpoints
func (h *HttpServ) AddMW(w http.ResponseWriter, r *http.Request) {
	apiID := r.FormValue("api_id")
	if apiID == "" {
		h.HandleError(fmt.Errorf("api_id must be specified"), w, r)
		return
	}

	tmpFileLoc, bundleName := h.ExtractBundleFromPost(w, r)
	if tmpFileLoc == "" || bundleName == "" {
		return
	}
	log.Info("saved bundle to ", tmpFileLoc)

	storeOnly := r.FormValue("store_only")

	if storeOnly != "" {
		// This is a python or JS bundle, just proxy it to a store
		mw, err := h.api.StoreBundleOnly(tmpFileLoc, apiID, bundleName)
		if err != nil {
			h.HandleError(err, w, r)
			return
		}

		ret := map[string]interface{}{"BundleID": mw.UID}
		h.HandleOK(ret, w, r)
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
	tmpFileLoc, _ := h.ExtractBundleFromPost(w, r)
	if tmpFileLoc == "" {
		return
	}

	mw, err := h.api.HandleUpdateBundle(tmpFileLoc, id)
	if err != nil {
		h.HandleError(err, w, r)
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
