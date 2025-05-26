package service

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

type CartItem struct {
	ProductID string
	Quantity  int
}

func TestAddToCartItem(t *testing.T) {
	item := CartItem{
		ProductID: "product123",
		Quantity:  2,
	}

	assert.Equal(t, "product123", item.ProductID)
	assert.True(t, item.Quantity > 0)
}
