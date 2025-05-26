package service

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

type DummyOrder struct {
	ID     string
	UserID string
	Amount float64
}

func TestValidateDummyOrder(t *testing.T) {
	order := DummyOrder{
		ID:     "order123",
		UserID: "user456",
		Amount: 99.99,
	}

	assert.NotEmpty(t, order.ID)
	assert.NotEmpty(t, order.UserID)
	assert.True(t, order.Amount > 0)
}
