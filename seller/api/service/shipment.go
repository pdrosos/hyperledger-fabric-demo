package service

import (
	"encoding/json"
	"strconv"
	"time"

	"github.com/hyperledger/fabric-sdk-go/pkg/client/channel"
	"github.com/sirupsen/logrus"
	"github.com/spf13/viper"

	"github.com/pdrosos/hyperledger-fabric-demo/seller/api/logger"
	"github.com/pdrosos/hyperledger-fabric-demo/seller/api/model"
)

type ShipmentService struct {
	channelClient *channel.Client
	chaincodeID   string
}

func NewShipmentService(channelClient *channel.Client) *ShipmentService {
	chaincodeID := viper.GetString("app.fabric.chaincodeID")

	return &ShipmentService{
		channelClient: channelClient,
		chaincodeID:   chaincodeID,
	}
}

func (this *ShipmentService) Create(shipment *model.Shipment) error {
	sender, _ := json.Marshal(shipment.Sender)
	recipient, _ := json.Marshal(shipment.Recipient)
	size, _ := json.Marshal(shipment.Size)
	lastState := "New"
	createdAt := time.Now().UTC().Format(time.RFC3339Nano)
	updatedAt := createdAt

	args := [][]byte{
		[]byte(shipment.TrackingCode),
		[]byte(shipment.Courier),
		sender,
		recipient,
		[]byte(strconv.Itoa(shipment.WeightInGrams)),
		[]byte(shipment.ShippingType),
		size,
		[]byte(shipment.Content),
		[]byte(strconv.FormatBool(shipment.IsFragile)),
		[]byte(lastState),
		[]byte(createdAt),
		[]byte(updatedAt),
	}

	response, err := this.channelClient.Execute(
		channel.Request{
			ChaincodeID: this.chaincodeID,
			Fcn:         "createShipment",
			Args:        args,
		},
	)
	if err != nil {
		logger.Log.WithFields(logrus.Fields{
			"shipment": shipment,
			"response": string(response.Payload),
		}).WithError(err).Error("Unable to create shipment")

		return err
	}

	return nil
}

func (this *ShipmentService) GetAll() ([]*model.Shipment, error) {
	return nil, nil
}

func (this *ShipmentService) GetByTrackingCode(trackingCode string) (*model.Shipment, error) {
	return nil, nil
}

func (this *ShipmentService) GetHistory(trackingCode string) ([]*model.Shipment, error) {
	return nil, nil
}
