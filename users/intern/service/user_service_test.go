package service

import (
	"gotune/users/pkg/hash"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestHashPassword(t *testing.T) {
	hashVal, err := hash.HashPassword("secure123")
	assert.NoError(t, err)
	assert.NotEmpty(t, hashVal)
}

func TestCheckPassword(t *testing.T) {
	hashVal, _ := hash.HashPassword("secure123")
	match := hash.CheckPasswordHash("secure123", hashVal)
	assert.True(t, match)
}
