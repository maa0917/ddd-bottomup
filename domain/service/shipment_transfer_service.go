package service

import (
	"ddd-bottomup/domain/entity"
	"errors"
)

type ShipmentTransferService struct{}

func NewShipmentTransferService() *ShipmentTransferService {
	return &ShipmentTransferService{}
}

func (s *ShipmentTransferService) Transport(
	from *entity.PhysicalDistributionBase,
	to *entity.PhysicalDistributionBase,
	baggage *entity.Baggage,
) error {
	if from.ID().Value() == to.ID().Value() {
		return errors.New("from and to cannot be the same")
	}

	shippedBaggage := from.Ship(baggage)

	to.Receive(shippedBaggage)

	return nil
}
