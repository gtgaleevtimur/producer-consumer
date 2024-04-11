package core

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

const (
	TEST_ADDRESS = "localhost:9092"
)

// Прежде всего необходимо запустить kafka локально
func TestNewConsumer(t *testing.T) {
	t.Run("NewConsumer", func(t *testing.T) {
		result, err := NewConsumer([]string{TEST_ADDRESS})
		require.NoError(t, err)
		assert.NotNil(t, result)
	})
}

func TestConsumer_Stop(t *testing.T) {
	t.Run("ConsumerStop", func(t *testing.T) {
		result, err := NewConsumer([]string{TEST_ADDRESS})
		require.NoError(t, err)
		assert.NotNil(t, result)
		err = result.Stop()
		assert.NoError(t, err)
	})
}

func TestConsumer_Consume(t *testing.T) {
	t.Run("ConsumerConsume", func(t *testing.T) {
		result, err := NewConsumer([]string{TEST_ADDRESS})
		require.NoError(t, err)
		assert.NotNil(t, result)
		err = result.Consume()
		assert.NoError(t, err)
	})
}
