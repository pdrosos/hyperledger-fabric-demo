package main

import (
	"flag"
	"fmt"
	"net/http"

	"github.com/sirupsen/logrus"

	"github.com/pdrosos/hyperledger-fabric-demo/seller/api/cmd/apiserver/handler"
	"github.com/pdrosos/hyperledger-fabric-demo/seller/api/common"
	"github.com/pdrosos/hyperledger-fabric-demo/seller/api/fabricsdk"
	"github.com/pdrosos/hyperledger-fabric-demo/seller/api/logger"
)

func main() {
	common.LoadConfig()
	logger.Log.Info("* Web server")

	// connect to Fabric SDK
	fabricSDK, err := fabricsdk.GetFabricSdk()
	if err != nil {
		logger.Log.WithFields(logrus.Fields{
			"error": err,
		}).Fatal("Unable to connect to Fabric SDK")
	}

	defer fabricSDK.Close()

	// get Fabric channel client
	channelClient, err := fabricsdk.GetChannelClient()
	if err != nil {
		logger.Log.WithFields(logrus.Fields{
			"error": err,
		}).Fatal("Unable to create Fabric channel client")
	}

	handler.Register(channelClient)

	var port uint
	var verbose bool = false

	flag.UintVar(&port, "port", 7777, "Port to run web server on")
	flag.BoolVar(&verbose, "verbose", false, "Run in verbose mode. Prints debug statement to stdout")
	flag.Parse()
	if verbose {
		logger.Log.Level = logrus.DebugLevel
	}

	hostAndPort := fmt.Sprintf(":%d", port)

	logger.Log.Infof("Starting the web server on %s. Revision: %s", hostAndPort, logger.Revision)

	if err := http.ListenAndServe(hostAndPort, nil); err != nil {
		logger.Log.WithFields(logrus.Fields{
			"error": err,
			"Host":  hostAndPort,
		}).Fatal("Unable to start web server")
	}
}
