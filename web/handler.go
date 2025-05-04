package web

import (
	"net/http"
)

func HealthcheckHandler(w http.ResponseWriter, r *http.Request) {
	Respond(w, http.StatusOK, map[string]string{"status": "ok"})
}
