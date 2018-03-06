package main

import (
	"flag"
	"fmt"
	"net/http"

	"github.com/sirupsen/logrus"

	"github.com/pdrosos/hyperledger-fabric-demo/customer/cmd/apiserver/handler"
	"github.com/pdrosos/hyperledger-fabric-demo/customer/logger"
	"github.com/pdrosos/hyperledger-fabric-demo/customer/customer"
)

func main() {
	customer.LoadConfig()
	logger.Log.Info("* Web server")

	var port uint
	var verbose bool = false

	flag.UintVar(&port, "port", 9999, "Port to run web server on")
	flag.BoolVar(&verbose, "verbose", false, "Run in verbose mode. Prints debug statement to stdout")
	flag.Parse()
	if verbose {
		logger.Log.Level = logrus.DebugLevel
	}

	hostAndPort := fmt.Sprintf(":%d", port)

	handler.Register()

	logger.Log.Infof("Starting the web server on %s. Revision: %s", hostAndPort, logger.Revision)

	if err := http.ListenAndServe(hostAndPort, nil); err != nil {
		logger.Log.WithFields(logrus.Fields{
			"error": err,
			"Host":  hostAndPort,
		}).Fatal("Unable to start web server")
	}
}
