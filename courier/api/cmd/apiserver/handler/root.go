package handler

import (
	"encoding/json"
	"net/http"
	"os"

	"github.com/pdrosos/hyperledger-fabric-demo/courier/api/logger"
)

type RootHandler struct {
}

func NewRootHandler() *RootHandler {
	return &RootHandler{}
}

func (this *RootHandler) ServeHTTP(rw http.ResponseWriter, r *http.Request) {
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

	rw.Header().Set("Content-Type", "application/json")
	response, cerr := json.Marshal(view)
	if cerr != nil {
		logger.Log.Error("Unable to encode json")

		return
	}

	rw.Write(response)
}
