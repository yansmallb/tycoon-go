package api

import (
	"net/http"

	"github.com/gorilla/mux"
)

type handler func(w http.ResponseWriter, r *http.Request)

var routes = map[string]map[string]handler{
	"GET": {
		"/services": getServices,
	},
}

// NewPrimary creates a new API router.
func NewPrimary() *mux.Router {
	// Register the API events handler in the cluster.
	r := mux.NewRouter()
	for method, mappings := range routes {
		for route, fct := range mappings {
			//log.WithFields(log.Fields{"method": method, "route": route}).Debug("Registering HTTP route")

			localRoute := route
			localFct := fct
			wrap := func(w http.ResponseWriter, r *http.Request) {
				//log.WithFields(log.Fields{"method": r.Method, "uri": r.RequestURI}).Debug("HTTP request received")
				localFct(w, r)
			}
			localMethod := method

			r.Path(localRoute).Methods(localMethod).HandlerFunc(wrap)
		}
	}

	return r
}
