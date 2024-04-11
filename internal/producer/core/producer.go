package core

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"os"
	"os/signal"
	"syscall"
	"time"

	"gitlab.simbirsoft/verify/t.galeev/pkg"

	"github.com/IBM/sarama"
	ntw "moul.io/number-to-words"
)

const (
	topicName = "sync_topic"
)

type Producer struct {
	counter int // счетчик отправленных сообщений
	App     sarama.SyncProducer
}

// NewProducer настраивает и создает Producer
func NewProducer(addresses []string) (*Producer, error) {
	// Настройка producer Kafka
	config := sarama.NewConfig()
	config.Producer.Return.Successes = true
	config.Producer.Compression = sarama.CompressionSnappy
	config.Producer.Flush.Frequency = 500 * time.Millisecond

	// Создание клиента Kafka
	producer, err := sarama.NewSyncProducer(addresses, config)
	if err != nil {
		return nil, fmt.Errorf("failed to create kafka producer:%w", err)
	}
	return &Producer{App: producer}, nil
}

// Run выполняет переключение между задачами
func (app *Producer) Run() {
	ctx, cancel := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM, syscall.SIGQUIT)
	defer cancel()

	for {
		select {
		case <-ctx.Done():
			if err := app.Stop(); err != nil {
				log.Println(err)
				os.Exit(1)
			}
			os.Exit(0)
		default:
			if err := app.Send(); err != nil {
				log.Println(err)
			}
		}
	}
}

// Stop завершает работу
func (app *Producer) Stop() error {
	if err := app.App.Close(); err != nil {
		return fmt.Errorf("kafka producer not closed gracefully err:%w", err)
	}
	log.Println("Successfully shutdown kafka producer")
	return nil
}

// Send - отправляет сообщение
func (app *Producer) Send() error {
	// Генерация структуры Sync
	data := app.GenerateSync()

	// Преобразование структуры Sync в JSON
	message, err := json.Marshal(data)
	if err != nil {
		return fmt.Errorf("failed to marshal Sync struct to JSON:%w", err)
	}
	// Отправка сообщения в Kafka
	_, _, err = app.App.SendMessage(&sarama.ProducerMessage{
		Topic: topicName,
		Value: sarama.StringEncoder(message),
	})
	if err != nil {
		return fmt.Errorf("failed to send message to Kafka:%w", err)
	}

	log.Println("Sent message:", data)
	// Задержка перед отправкой следующего сообщения
	time.Sleep(1 * time.Second)
	return nil
}

// Генерация структуры Sync
func (app *Producer) GenerateSync() pkg.Sync {
	app.counter++
	roman := convertToRoman(app.counter)
	description := convertToDescription(app.counter)
	return pkg.Sync{
		Number:      app.counter,
		Roman:       roman,
		Description: description,
	}
}

// Конвертация int в римскую запись
func convertToRoman(num int) string {
	romanValues := []int{1000, 900, 500, 400, 100, 90, 50, 40, 10, 9, 5, 4, 1}
	romanSymbols := []string{"M", "CM", "D", "CD", "C", "XC", "L", "XL", "X", "IX", "V", "IV", "I"}

	result := ""
	for i := 0; i < len(romanValues); i++ {
		for num >= romanValues[i] {
			result += romanSymbols[i]
			num -= romanValues[i]
		}
	}

	return result
}

// Конвертация int в строчное описания числа на English
func convertToDescription(num int) string {
	return ntw.IntegerToEnUs(num)
}
