package config

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNewConfig(t *testing.T) {
	t.Run("NewConfig", func(t *testing.T) {
		res := NewConfig()
		assert.NotNil(t, res)
	})
}
