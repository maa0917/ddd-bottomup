package entity

import (
	"github.com/google/uuid"
	"time"
)

type ShipmentID struct {
	value string
}

func NewShipmentID() *ShipmentID {
	return &ShipmentID{value: uuid.New().String()}
}

func (s *ShipmentID) Value() string {
	return s.value
}

func (s *ShipmentID) String() string {
	return s.value
}

type ShipmentStatus int

const (
	ShipmentStatusPreparing ShipmentStatus = iota
	ShipmentStatusShipped
	ShipmentStatusReceived
)

type Shipment struct {
	id          *ShipmentID
	fromAddress string
	toAddress   string
	status      ShipmentStatus
	shippedAt   *time.Time
	receivedAt  *time.Time
}

func NewShipment(fromAddress, toAddress string) *Shipment {
	return &Shipment{
		id:          NewShipmentID(),
		fromAddress: fromAddress,
		toAddress:   toAddress,
		status:      ShipmentStatusPreparing,
	}
}

func (s *Shipment) ID() *ShipmentID {
	return s.id
}

func (s *Shipment) Status() ShipmentStatus {
	return s.status
}

func (s *Shipment) Ship(shippedAt time.Time) {
	s.status = ShipmentStatusShipped
	s.shippedAt = &shippedAt
}

func (s *Shipment) Receive(receivedAt time.Time) {
	s.status = ShipmentStatusReceived
	s.receivedAt = &receivedAt
}

func (s *Shipment) IsShipped() bool {
	return s.status == ShipmentStatusShipped
}

func (s *Shipment) IsReceived() bool {
	return s.status == ShipmentStatusReceived
}