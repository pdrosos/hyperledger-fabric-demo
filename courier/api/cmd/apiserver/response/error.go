package response

import (
	"encoding/json"
	"net/http"

	"github.com/pdrosos/hyperledger-fabric-demo/courier/api/logger"
	"github.com/pdrosos/hyperledger-fabric-demo/courier/api/viewmodel"
)

func ResponseError(
	rw http.ResponseWriter,
	status int,
	code string,
	title string,
	errors []viewmodel.ErrorDetails,
	err error,
) {
	if err != nil {
		logger.Log.WithField("err", err).Error(title)
	}

	errorViewModel := viewmodel.Error{}
	errorViewModel.Status = status
	errorViewModel.Code = code
	errorViewModel.Title = title
	errorViewModel.Errors = errors
	errorResponse, encodeError := json.Marshal(errorViewModel)
	if encodeError != nil {
		logger.Log.WithField("err", encodeError).Error("Unable to encode json response error")

		http.Error(rw, title, status)

		return
	}

	rw.Header().Add("Content-Type", "application/json; charset=utf-8")
	rw.WriteHeader(status)
	rw.Write(errorResponse)
}
