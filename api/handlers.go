package api

import (
	"encoding/json"
	"net/http"
)

func getServices(w http.ResponseWriter, r *http.Request) {
	json.NewEncoder(w).Encode("testGetServices")
}
