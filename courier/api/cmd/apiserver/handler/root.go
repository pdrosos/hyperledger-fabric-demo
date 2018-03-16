package handler

import (
	"encoding/json"
	"net/http"
	"os"

	"github.com/pdrosos/hyperledger-fabric-demo/courier/api/logger"
)

func rootHandler(responseWriter http.ResponseWriter, request *http.Request) {
	hostname, err := os.Hostname()
	if err != nil {
		logger.Log.Error("Unable to get hostname")

		return
	}

	view := struct {
		Hostname string `json:"hostname"`
		Revision string `json:"revision"`
	}{
		hostname,
		logger.Revision,
	}

	responseWriter.Header().Set("Content-Type", "application/json")
	response, cerr := json.Marshal(view)
	if cerr != nil {
		logger.Log.Error("Unable to encode json")

		return
	}

	responseWriter.Write(response)
}
