package http_funcs

import (
	"encoding/json"
	"errors"
	"net"
	"net/http"
	"time"

	"github.com/gorilla/mux"
	"github.com/rs/cors"

	"github.com/TykTechnologies/mserv/api"
	"github.com/TykTechnologies/mserv/health"
	"github.com/TykTechnologies/mserv/models"
	"github.com/TykTechnologies/mserv/storage"
	"github.com/TykTechnologies/mserv/util/logger"
)

var (
	moduleName = "mserv.http"
	log        = logger.GetLogger(moduleName)
)

func NewServer(listenOn string, store storage.MservStore) *HttpServ {
	return &HttpServ{
		addr: listenOn,
		api:  api.NewAPI(store),
	}
}

type HttpServ struct {
	api  *api.API
	addr string
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

// swagger:route GET /health system health
// Query health status of Mserv service.
//
// Responses:
//   200: healthResponse
//   500: healthResponse
func (h *HttpServ) HealthHandler(w http.ResponseWriter, r *http.Request) {
	if health.Report.HTTPStarted || health.Report.GRPCStarted {
		h.HandleOK(health.Report, w, r)
		return
	}

	h.writeToClient(w, r, models.NewPayload("error", health.Report, ""), http.StatusInternalServerError) // 500
}

func (h *HttpServ) HandleError(err error, w http.ResponseWriter, r *http.Request) {
	var (
		message string
		status  int
	)

	switch {
	case errors.Is(err, ErrGenericMimeDetected):
		message = "unsupported media type"
		status = http.StatusUnsupportedMediaType // 415

	case errors.Is(err, ErrUploadNotZip):
		message = "unprocessable entity"
		status = http.StatusUnprocessableEntity // 422

	default:
		message = "internal server error"
		status = http.StatusInternalServerError // 500
	}

	log.WithError(err).Error(message)
	h.writeToClient(w, r, models.NewPayload("error", nil, err.Error()), status)
}

func (h *HttpServ) HandleOK(payload interface{}, w http.ResponseWriter, r *http.Request) {
	h.writeToClient(w, r, models.NewPayload("ok", payload, ""), http.StatusOK) // 200
}

func (h *HttpServ) writeToClient(w http.ResponseWriter, r *http.Request, payload models.Payload, status int) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)

	// marshall the payload and handle encoding errors
	js, err := json.Marshal(payload)
	if err != nil {
		// Write big error
		es, errMarshal := json.Marshal(models.NewPayload("error", nil, err.Error()))
		if errMarshal != nil {
			log.WithError(errMarshal).Fatal("This is a terrible place to be")
		}

		w.Write(es)
		return
	}

	w.Write(js)
}

func setupGlobalMiddleware(handler http.Handler) http.Handler {
	handleCORS := cors.Default().Handler

	return handleCORS(handler)
}
