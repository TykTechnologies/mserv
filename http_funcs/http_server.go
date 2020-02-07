package http_funcs

import (
	"encoding/json"
	"github.com/TykTechnologies/mserv/api"
	"github.com/TykTechnologies/mserv/health"
	"github.com/TykTechnologies/mserv/models"
	"github.com/TykTechnologies/mserv/storage"
	"github.com/TykTechnologies/mserv/util/logger"
	"github.com/gorilla/mux"
	"net"
	"net/http"
	"time"
)

var moduleName = "mserv.http"
var log = logger.GetAndExcludeLoggerFromTrace(moduleName)

func NewServer(listenOn string, store storage.MservStore) *HttpServ {
	return &HttpServ{
		addr: listenOn,
		api:  api.NewAPI(store),
	}
}

type HttpServ struct {
	addr string
	api  *api.API
}

func (h *HttpServ) Listen(m *mux.Router, l net.Listener) error {
	srv := &http.Server{
		Handler: m,
		Addr:    h.addr,

		WriteTimeout: 1 * time.Minute,
		ReadTimeout:  1 * time.Minute,
	}

	err := srv.Serve(l)
	if err != nil {
		return err
	}

	return nil
}

// Health endpoint
func (h *HttpServ) HealthHandler(w http.ResponseWriter, r *http.Request) {
	if health.Report.HTTPStarted || health.Report.GRPCStarted {
		h.HandleOK(health.Report, w, r)
		return
	}

	h.writeToClient(w, r, models.NewPayload("error", health.Report, ""), 500)
}

func (h *HttpServ) HandleError(err error, w http.ResponseWriter, r *http.Request) {
	log.Error("api error: ", err)
	h.writeToClient(w, r, models.NewPayload("error", nil, err.Error()), 500)
}

func (h *HttpServ) HandleOK(payload interface{}, w http.ResponseWriter, r *http.Request) {
	h.writeToClient(w, r, models.NewPayload("ok", payload, ""), 200)
}

func (h *HttpServ) writeToClient(w http.ResponseWriter, r *http.Request, payload models.Payload, status int) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)

	// marshall the payload and handle encoding errors
	js, err := json.Marshal(payload)
	if err != nil {
		// Write big error
		es, err := json.Marshal(models.NewPayload("error", nil, err.Error()))
		if err != nil {
			log.Fatal("This is a terrible place to be: ", err)
		}

		w.Write(es)
		return
	}

	w.Write(js)
}
