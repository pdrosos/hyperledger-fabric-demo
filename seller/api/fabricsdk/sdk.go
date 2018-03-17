package fabricsdk

import (
	"github.com/hyperledger/fabric-sdk-go/pkg/client/channel"
	"github.com/hyperledger/fabric-sdk-go/pkg/core/config"
	"github.com/hyperledger/fabric-sdk-go/pkg/fabsdk"

	"github.com/sirupsen/logrus"
	"github.com/spf13/viper"

	"github.com/pdrosos/hyperledger-fabric-demo/seller/api/logger"
)

func GetFabricSdk() (*fabsdk.FabricSDK, error) {
	// load config
	configOpt := config.FromFile("./config/network_config.yaml")

	sdk, err := fabsdk.New(configOpt)
	if err != nil {
		logger.Log.WithError(err).Error("Unable to create Fabric SDK")

		return nil, err
	}

	return sdk, nil
}

func GetChannelClient(fabricSDK *fabsdk.FabricSDK) (*channel.Client, error) {
	organization := viper.GetString("app.fabric.organization")
	username := viper.GetString("app.fabric.username")
	channelID := viper.GetString("app.fabric.channelID")

	logger.Log.WithFields(logrus.Fields{
		"organization": organization,
		"username":     username,
		"channelID":    channelID,
	}).Debug("Creating Channel Client")

	//prepare channel client context using client context
	clientChannelContext := fabricSDK.ChannelContext(
		channelID,
		fabsdk.WithUser(username),
		fabsdk.WithOrg(organization),
	)

	client, err := channel.New(clientChannelContext)
	if err != nil {
		logger.Log.WithError(err).Error("Unable to create Fabric Channel Client")

		return nil, err
	}

	return client, nil
}
