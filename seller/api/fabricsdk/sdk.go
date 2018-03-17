package fabricsdk

import (
	"errors"

	"github.com/hyperledger/fabric-sdk-go/pkg/client/channel"
	"github.com/hyperledger/fabric-sdk-go/pkg/common/providers/core"
	"github.com/hyperledger/fabric-sdk-go/pkg/common/providers/fab"
	"github.com/hyperledger/fabric-sdk-go/pkg/core/config"
	"github.com/hyperledger/fabric-sdk-go/pkg/fab/peer"
	"github.com/hyperledger/fabric-sdk-go/pkg/fabsdk"

	"github.com/sirupsen/logrus"
	"github.com/spf13/viper"

	"github.com/pdrosos/hyperledger-fabric-demo/seller/api/logger"
)

var fabricSDK *fabsdk.FabricSDK

func GetFabricSdk() (*fabsdk.FabricSDK, error) {
	if fabricSDK != nil {
		return fabricSDK, nil
	}

	// load config
	configOpt := config.FromFile("./config/network_config.yaml")

	sdk, err := fabsdk.New(configOpt)
	if err != nil {
		logger.Log.WithError(err).Error("Unable to create Fabric SDK")

		return nil, err
	}

	fabricSDK = sdk

	return fabricSDK, nil
}

func GetChannelClient() (*channel.Client, error) {
	if fabricSDK == nil {
		err := errors.New("Fabric SDK is not created. Call GetFabricSdk() first")
		logger.Log.WithError(err).Errorf("Unable to create Fabric Channel Client")

		return nil, err
	}

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

func LoadPeers() (fab.Peer, fab.Peer, error) {
	if fabricSDK == nil {
		err := errors.New("Fabric SDK is not created. Call GetFabricSdk() first")
		logger.Log.WithError(err).Errorf("Unable to get peers")

		return nil, nil, err
	}

	organization := viper.GetString("app.fabric.organization")
	courier1 := viper.GetString("app.fabric.courier1")

	organizationPeersConfig, err := fabricSDK.Config().PeersConfig(organization)
	if err != nil {
		logger.Log.WithError(err).Errorf("Unable to get organization peers config")

		return nil, nil, err
	}

	courier1PeersConfig, err := fabricSDK.Config().PeersConfig(courier1)
	if err != nil {
		logger.Log.WithError(err).Errorf("Unable to get courier1 peers config")

		return nil, nil, err
	}

	organizationPeer0, err := peer.New(
		fabricSDK.Config(),
		peer.FromPeerConfig(&core.NetworkPeer{PeerConfig: organizationPeersConfig[0]}),
	)
	if err != nil {
		logger.Log.WithError(err).Errorf("Unable to create organization peer")

		return nil, nil, err
	}

	courier1Peer0, err := peer.New(
		fabricSDK.Config(),
		peer.FromPeerConfig(&core.NetworkPeer{PeerConfig: courier1PeersConfig[0]}),
	)
	if err != nil {
		logger.Log.WithError(err).Errorf("Unable to create courier1 peer")

		return nil, nil, err
	}

	return organizationPeer0, courier1Peer0, nil
}

func RegisterChaincodeEvent(
	channelClient *channel.Client,
	chaincodeID string,
	eventID string,
	chaincodeEventHandler ChaincodeEventHandler,
) {
	reg, notifier, err := channelClient.RegisterChaincodeEvent(chaincodeID, eventID)
	if err != nil {
		logger.Log.WithFields(logrus.Fields{
			"error":       err,
			"chaincodeID": chaincodeID,
			"eventID":     eventID,
		}).Error("Unable to register chaincode event")

		return
	} else {
		logger.Log.WithFields(logrus.Fields{
			"chaincodeID": chaincodeID,
			"eventID":     eventID,
		}).Info("Registered chaincode event")
	}

	defer channelClient.UnregisterChaincodeEvent(reg)

	select {
	case ccEvent := <-notifier:
		logger.Log.WithFields(logrus.Fields{
			"chaincodeID": chaincodeID,
			"event":       ccEvent,
		}).Info("Received chaincode event")

		if chaincodeEventHandler != nil {
			chaincodeEventHandler.Handle(ccEvent)
		}
	}
}
