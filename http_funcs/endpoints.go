package http_funcs

import "github.com/gorilla/mux"

var rt *mux.Router

func GetRouter() *mux.Router {
	rt = mux.NewRouter()
	rt.Use(setupGlobalMiddleware)
	return rt
}

func InitEndpoints(r *mux.Router, serv *HttpServ) {
	// Health endpoint
	r.HandleFunc("/health", serv.HealthHandler).Methods("GET")
}

func InitHttpInvocationServer(r *mux.Router, serv *HttpServ) {
	// Invocation endpoints
	r.HandleFunc("/execute/{name}", serv.Execute).Methods("POST")
}

func InitAPI(r *mux.Router, serv *HttpServ) {
	r.HandleFunc("/api/mw/master/all", serv.FetchAllActiveMW).Methods("GET")
	r.HandleFunc("/api/mw/bundle/{id}", serv.FetchBundleFile).Methods("GET")
	r.HandleFunc("/api/mw/{id}", serv.UpdateMW).Methods("PUT")
	r.HandleFunc("/api/mw/{id}", serv.FetchMW).Methods("GET")
	r.HandleFunc("/api/mw/{id}", serv.DeleteMW).Methods("DELETE")
	r.HandleFunc("/api/mw", serv.AddMW).Methods("POST")
}
