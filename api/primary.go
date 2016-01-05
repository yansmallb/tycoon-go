package api

import (
	log "github.com/Sirupsen/logrus"
	"github.com/gorilla/mux"
	"net/http"
)

type handler func(w http.ResponseWriter, r *http.Request)

var routes = map[string]map[string]handler{
	"GET": {
		"/services":              getServices,
		"/servicesInfo":          getServicesInfo,
		"/service/{name:.*}/get": getService,
	},
	"POST": {
		"/service/{name:.*}/stop":    stopService,
		"/service/{name:.*}/start":   startService,
		"/service/{name:.*}/restart": restartService,
		"/service/{name:.*}/delete":  deleteService,
		"/service/create":            createService,
	},
}

// NewPrimary creates a new API router.
func NewPrimary() *mux.Router {
	// Register the API events handler in the cluster.
	r := mux.NewRouter()
	for method, mappings := range routes {
		for route, fct := range mappings {
			//log.WithFields(log.Fields{"method": method, "route": route}).Debug("api.NewPrimary():Registering HTTP route")

			localRoute := route
			localFct := fct
			wrap := func(w http.ResponseWriter, r *http.Request) {
				log.WithFields(log.Fields{"method": r.Method, "uri": r.RequestURI}).Debug("api.NewPrimary():HTTP request received")
				localFct(w, r)
			}
			localMethod := method

			r.Path(localRoute).Methods(localMethod).HandlerFunc(wrap)
		}
	}
	return r
}
