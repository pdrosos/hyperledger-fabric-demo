package service

import (
	"github.com/hyperledger/fabric-sdk-go/pkg/client/channel"

	"github.com/pdrosos/hyperledger-fabric-demo/seller/api/model"
)

type ShipmentService struct {
	channelClient *channel.Client
}

func NewShipmentService(channelClient *channel.Client) *ShipmentService {
	return &ShipmentService{
		channelClient: channelClient,
	}
}

func (this *ShipmentService) Create(shipment *model.Shipment) (error) {
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