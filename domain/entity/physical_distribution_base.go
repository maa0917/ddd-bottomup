package entity

import "github.com/google/uuid"

type PhysicalDistributionBaseID struct {
	value string
}

func NewPhysicalDistributionBaseID() *PhysicalDistributionBaseID {
	return &PhysicalDistributionBaseID{value: uuid.New().String()}
}

func (p *PhysicalDistributionBaseID) Value() string {
	return p.value
}

func (p *PhysicalDistributionBaseID) String() string {
	return p.value
}

type PhysicalDistributionBase struct {
	id      *PhysicalDistributionBaseID
	name    string
	address string
}

func NewPhysicalDistributionBase(name, address string) *PhysicalDistributionBase {
	return &PhysicalDistributionBase{
		id:      NewPhysicalDistributionBaseID(),
		name:    name,
		address: address,
	}
}

func (p *PhysicalDistributionBase) ID() *PhysicalDistributionBaseID {
	return p.id
}

func (p *PhysicalDistributionBase) Name() string {
	return p.name
}

func (p *PhysicalDistributionBase) Address() string {
	return p.address
}

func (p *PhysicalDistributionBase) Ship(baggage *Baggage) *Baggage {
	return baggage
}

func (p *PhysicalDistributionBase) Receive(baggage *Baggage) {
}