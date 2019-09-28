package controller

import (
	"encoding/json"
	"net/http"

	"github.com/lealoureiro/mortgage-calculator-api/utils"
	log "github.com/sirupsen/logrus"
)

func ShowInfo(w http.ResponseWriter, r *http.Request) {

	log.Printf("Show application information to: %s", r.RemoteAddr)

	if !utils.IsClientAcceptingJSON(r) {
		utils.RespondHTTPError(406, "Unsupported media!", w)
		return
	}

	m := map[string]string{"applicationName": "Mortgage Calculator", "applicationVersion": "0.0.1"}
	j, _ := json.Marshal(m)

	w.Header().Set("Content-Type", "application/json")

	w.Write([]byte(j))
}
