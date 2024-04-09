package main

import (
	"encoding/json"
	"fmt"
	"time"
	"v1/models"

	// "github.com/go-vgo/robotgo"
	"github.com/confluentinc/confluent-kafka-go/v2/kafka"
	hook "github.com/robotn/gohook"
)

const (
	KafkaServer          = "localhost:9092"
	KafkaTopic           = "clicks-v1-topic"
	KafkaPartitionConfig = kafka.PartitionAny

	LeftMouse  = 1
	RightMouse = 2

	MouseHold = 7
)

func main() {
	producer, err := createKafkaProducer()
	if err != nil {
		panic(err)
	}

	messageHandler(KafkaTopic, producer)

	defer producer.Close()
}

func createKafkaProducer() (*kafka.Producer, error) {
	return kafka.NewProducer(&kafka.ConfigMap{
		"bootstrap.servers": KafkaServer,
	})
}

func messageHandler(topic string, producer *kafka.Producer) {
	fmt.Println("hook start...")

	evChan := hook.Start()
	defer hook.End()

	for ev := range evChan {
		clickListener(topic, producer, ev)

		if trackerStopper(ev) {
			break
		}

	}

	fmt.Println("hook stop...")
}

func clickListener(topic string, producer *kafka.Producer, ev hook.Event) {

	if ev.Kind == MouseHold && ev.Button == LeftMouse {
		fmt.Printf("Left Mouse Clicked ! \n")
		clicks := models.NewClick(ev.When.Format(time.DateTime), ev.X, ev.Y)

		value, err := serializeMessage(clicks)
		if err != nil {
			panic(err)
		}

		err = messageProducer(topic, producer, value)

		if err != nil {
			panic(err)
		}
	}

}

func serializeMessage(clicks *models.Click) ([]byte, error) {
	// Serialize the message struct to JSON
	serialized, err := json.Marshal(clicks)
	if err != nil {
		return nil, fmt.Errorf("failed to serialize message: %w", err)
	}
	return serialized, nil
}

func trackerStopper(ev hook.Event) bool {
	if ev.Kind == MouseHold && ev.Button == RightMouse {
		fmt.Printf("Right mouse clicked and it means quitting..\n")
		return true
	}

	return false
}

func messageProducer(topic string, producer *kafka.Producer, value []byte) error {
	return producer.Produce(&kafka.Message{
		TopicPartition: kafka.TopicPartition{Topic: &topic, Partition: KafkaPartitionConfig},
		Value:          value,
	}, nil)

}
