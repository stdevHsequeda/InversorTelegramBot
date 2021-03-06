package main

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestGetPrices(t *testing.T) {
	x, err := GetPrices()
	if err != nil {
		t.Error(err)
		return
	}
	assert.True(t, len(x) != 0)
	for k, v := range x {
		assert.NotNil(t, v.Fifteen, k)
		assert.NotNil(t, v.Buy, k)
		assert.NotNil(t, v.Last, k)
		assert.NotNil(t, v.Sell, k)
		assert.NotNil(t, v.Symbol, k)
	}
}

func TestGetAddress(t *testing.T) {
	// key="KehhfchD6puGoAaKKY8pIZhVUtb0tVGWPgQvh0866AE"
	if addr, err := GetAddress(); err != nil {
		t.Error(err)
	} else {
		t.Log(addr)
	}
}
