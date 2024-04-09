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
	KafkaServer = "localhost:9092"
	KafkaTopic  = "clicks-v1-topic"
)

func main() {
	fmt.Println("------------- Hata -------------")
	producer, err := kafka.NewProducer(&kafka.ConfigMap{
		"bootstrap.servers": KafkaServer,
	})

	topic := KafkaTopic

	fmt.Println(err)
	fmt.Println("------------- Hata -------------")
	if err != nil {
		panic(err)
	}

	base(topic, producer)

	defer producer.Close()
}

func base(topic string, producer *kafka.Producer) {
	fmt.Println("hook start...")
	evChan := hook.Start()
	defer hook.End()

	for ev := range evChan {

		if ev.Kind == 7 && ev.Button == 1 {
			fmt.Printf("Mouse Tıklama = %v / X: %v Y:%v \n", ev.When.Local(), ev.X, ev.Y)
			clicks := models.Click{
				DateTime: ev.When.Local().Format(time.DateTime),
				X:        ev.X,
				Y:        ev.Y,
			}

			value, err := json.Marshal(clicks)
			if err != nil {
				panic(err)
			}
			err = producer.Produce(&kafka.Message{
				TopicPartition: kafka.TopicPartition{Topic: &topic, Partition: kafka.PartitionAny},
				Value:          value,
			}, nil)

			if err != nil {
				panic(err)
			}
		}
		if ev.Kind == 7 && ev.Button == 2 {
			fmt.Printf("Mouse Tıklama = %v / X: %v Y:%v \n", ev.When.Local(), ev.X, ev.Y)
			break
		}
	}
}
