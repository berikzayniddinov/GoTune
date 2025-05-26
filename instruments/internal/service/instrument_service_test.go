package service

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

type Instrument struct {
	ID    string
	Name  string
	Brand string
	Price float64
	Stock int
}

func TestValidateInstrumentData(t *testing.T) {
	instr := Instrument{
		ID:    "instr123",
		Name:  "Acoustic Guitar",
		Brand: "Fender",
		Price: 499.99,
		Stock: 5,
	}

	assert.Equal(t, "instr123", instr.ID)
	assert.Equal(t, "Acoustic Guitar", instr.Name)
	assert.Equal(t, "Fender", instr.Brand)
	assert.True(t, instr.Price > 0)
	assert.True(t, instr.Stock >= 0)
}
