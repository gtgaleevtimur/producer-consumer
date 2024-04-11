package core

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"gitlab.simbirsoft/verify/t.galeev/internal/consumer/config"

	"github.com/IBM/sarama"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

const (
	topicName = "sync_topic"
)

var (
	messagesReceived = prometheus.NewCounter(prometheus.CounterOpts{
		Name: "sync_messages_received",
		Help: "Total number of Kafka messages received",
	})
)

type Consumer struct {
	Client   sarama.Client
	Consumer sarama.Consumer
}

// NewConsumer настраивает и создает Consumer
func NewConsumer(addresses []string) (*Consumer, error) {
	conf := sarama.NewConfig()
	conf.Consumer.Return.Errors = true
	conf.Consumer.IsolationLevel = 1

	client, err := sarama.NewClient(addresses, conf)
	if err != nil {
		return nil, fmt.Errorf("failed to create kafka client:%w", err)
	}
	consumer, err := sarama.NewConsumerFromClient(client)
	if err != nil {
		return nil, fmt.Errorf("failed to create kafka consumer:%w", err)
	}
	return &Consumer{Client: client, Consumer: consumer}, nil
}

// Run выполняет переключение между задачами
func (app *Consumer) Run() {
	conf := config.NewConfig()

	prometheus.MustRegister(messagesReceived)

	ctx, cancel := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM, syscall.SIGQUIT)
	defer cancel()

	ticker := time.NewTicker(time.Duration(conf.Interval) * time.Second)

	countChan := make(chan struct{})
	go app.CountCheck(countChan)

	go func() {
		for {
			select {
			case <-ctx.Done():
				if err := app.Stop(); err != nil {
					log.Println(err)
					os.Exit(1)
				}
				os.Exit(0)
			case <-ticker.C:
				if err := app.Consume(); err != nil {
					log.Println(err)
				}
			case <-countChan:
				if err := app.ConsumePart(); err != nil {
					log.Println(err)
				}
			}
		}
	}()
	http.Handle("/metrics", promhttp.Handler())
	log.Fatal(http.ListenAndServe(conf.PrometheusAddress, nil))
}

// Stop завершает работу
func (app *Consumer) Stop() error {
	if err := app.Consumer.Close(); err != nil {
		return fmt.Errorf("kafka consumer not closed gracefully err:%w", err)
	}
	if err := app.Client.Close(); err != nil {
		return fmt.Errorf("kafka client not closed gracefully err:%w", err)
	}
	log.Println("Successfully shutdown kafka client")
	return nil
}

// Consume вычитывает сообщение
func (app *Consumer) Consume() error {
	partitions, err := app.Consumer.Partitions(topicName)
	if err != nil {
		return fmt.Errorf("failed to get Kafka partitions:%w", err)
	}
	for _, partition := range partitions {
		pc, err := app.Consumer.ConsumePartition(topicName, partition, sarama.OffsetNewest)
		if err != nil {
			return fmt.Errorf("failed to create Kafka partition consumer:%w", err)
		}

		message := <-pc.Messages()
		log.Println("Received message:", string(message.Value))
		messagesReceived.Inc()

		if err := pc.Close(); err != nil {
			return fmt.Errorf("failed to close Kafka partition consumer: %w", err)
		}
	}
	return nil
}

// ConsumePart вычитывает блок сообщении
func (app *Consumer) ConsumePart() error {
	conf := config.NewConfig()

	partitions, err := app.Consumer.Partitions(topicName)
	if err != nil {
		return fmt.Errorf("failed to get Kafka partitions:%w", err)
	}

	for _, partition := range partitions {
		pc, err := app.Consumer.ConsumePartition(topicName, partition, sarama.OffsetNewest)
		if err != nil {
			return fmt.Errorf("failed to create Kafka partition consumer:%w", err)
		}

		result := make([]string, conf.Count)
		messageCount := 0

		for message := range pc.Messages() {
			result = append(result, string(message.Value))
			messageCount++
			if messageCount >= conf.Count {
				break
			}
		}
		log.Println("Received group of messages:", result)
		messagesReceived.Add(float64(messageCount))

		if err := pc.Close(); err != nil {
			return fmt.Errorf("failed to close Kafka partition consumer: %w", err)
		}
	}
	return nil
}

// CountCheck проверяет длинну очереди kafka
func (app *Consumer) CountCheck(c chan struct{}) {
	conf := config.NewConfig()
	for {
		partitions, err := app.Client.Partitions(topicName)
		if err != nil {
			log.Fatal(err)
		}
		for _, partition := range partitions {
			hwm, err := app.Client.GetOffset(topicName, partition, sarama.OffsetNewest)
			if err != nil {
				log.Fatal(err)
			}

			lwm, err := app.Client.GetOffset(topicName, partition, sarama.OffsetOldest)
			if err != nil {
				log.Fatal(err)
			}
			queueLength := hwm - lwm
			if queueLength >= int64(conf.Count) {
				c <- struct{}{}
				break
			}
		}
	}
}
