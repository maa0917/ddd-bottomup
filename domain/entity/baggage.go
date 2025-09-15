package entity

import "github.com/google/uuid"

type BaggageID struct {
	value string
}

func NewBaggageID() *BaggageID {
	return &BaggageID{value: uuid.New().String()}
}

func (b *BaggageID) Value() string {
	return b.value
}

func (b *BaggageID) String() string {
	return b.value
}

type Baggage struct {
	id     *BaggageID
	weight int
	size   string
}

func NewBaggage(weight int, size string) *Baggage {
	return &Baggage{
		id:     NewBaggageID(),
		weight: weight,
		size:   size,
	}
}

func (b *Baggage) ID() *BaggageID {
	return b.id
}

func (b *Baggage) Weight() int {
	return b.weight
}

func (b *Baggage) Size() string {
	return b.size
}