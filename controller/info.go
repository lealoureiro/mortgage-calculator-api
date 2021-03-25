package controller

import (
	"encoding/json"
	"net/http"

	"github.com/lealoureiro/mortgage-calculator-api/config"
	"github.com/lealoureiro/mortgage-calculator-api/model"
	"github.com/lealoureiro/mortgage-calculator-api/utils"
	log "github.com/sirupsen/logrus"
)

// ShowInfo : create server info response
func ShowInfo(w http.ResponseWriter, r *http.Request) {

	// swagger:route GET /info API info
	// ---
	// description: Get application name and version information
	// responses:
	//   200: infoResponse

	log.Printf("Show application information to: %s", r.RemoteAddr)

	if !utils.IsClientAcceptingJSON(r) {
		utils.RespondHTTPError(406, "Unsupported media!", w)
		return
	}

	info := model.Info{ApplicationName: "MortgageCalculatorAPI", ApplicationVersion: config.Version}

	j, _ := json.Marshal(info)

	w.Header().Set("Content-Type", "application/json")

	w.Write([]byte(j))
}
