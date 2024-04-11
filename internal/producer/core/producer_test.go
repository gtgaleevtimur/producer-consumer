package core

import (
	"testing"
	"time"

	"gitlab.simbirsoft/verify/t.galeev/pkg"

	"github.com/IBM/sarama"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

const (
	TEST_ADDRESS = "localhost:9092"
)

func TestNewProducer(t *testing.T) {
	t.Run("NewProducer positive", func(t *testing.T) {
		result, err := NewProducer([]string{TEST_ADDRESS})
		require.NoError(t, err)
		assert.NotNil(t, result)
	})
}

func TestProducer_Stop(t *testing.T) {
	t.Run("ProducerStop", func(t *testing.T) {
		result, err := NewProducer([]string{TEST_ADDRESS})
		require.NoError(t, err)
		err = result.Stop()
		assert.NoError(t, err)
	})
}

func TestProducer_Send(t *testing.T) {
	t.Run("ProducerSend", func(t *testing.T) {
		config := sarama.NewConfig()
		config.Producer.Return.Successes = true
		config.Producer.Compression = sarama.CompressionSnappy
		config.Producer.Flush.Frequency = 500 * time.Millisecond

		client, err := sarama.NewClient([]string{TEST_ADDRESS}, config)
		require.NoError(t, err)
		defer client.Close()

		result, err := NewProducer([]string{TEST_ADDRESS})
		require.NoError(t, err)
		err = result.Send()
		assert.NoError(t, err)
	})
}

func TestProducer_GenerateSync(t *testing.T) {
	t.Run("GenerateSync", func(t *testing.T) {
		base, err := NewProducer([]string{TEST_ADDRESS})
		require.NoError(t, err)
		result := base.GenerateSync()
		assert.Equal(t, result, pkg.Sync{
			Number:      1,
			Roman:       "I",
			Description: "one",
		})
	})
}
