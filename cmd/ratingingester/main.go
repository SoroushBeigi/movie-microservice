package main

import (
	"encoding/json"
	"fmt"
	"github.com/SoroushBeigi/movie-microservice/rating/pkg/model"
	"github.com/confluentinc/confluent-kafka-go/kafka"
	"os"
	"time"
)

const fileName = "ratingsdata.json"
const topic = "ratings"
const timeout = 10 * time.Second

func main() {
	fmt.Println("Creating a Kafka producer")
	producer, err := kafka.NewProducer(&kafka.ConfigMap{"bootstrap.servers": "localhost"})
	if err != nil {
		panic(err)
	}
	defer producer.Close()

	fmt.Println("Reading rating events from file " + fileName)
	ratingEvents, readErr := readRatingEvents(fileName)
	if readErr != nil {
		panic(readErr)
	}

	if err := produceRatingEvents(topic, producer,
		ratingEvents); err != nil {
		panic(err)
	}

	fmt.Println("Waiting " + timeout.String() + " until all events get produced")
	producer.Flush(int(timeout.Milliseconds()))
}

func readRatingEvents(fileName string) ([]model.RatingEvent, error) {
	f, err := os.Open(fileName)
	if err != nil {
		return nil, err
	}
	defer f.Close()
	var ratings []model.RatingEvent
	if err := json.NewDecoder(f).Decode(&ratings); err != nil {
		return nil, err
	}
	return ratings, nil
}

func produceRatingEvents(topic string, producer *kafka.Producer, events []model.RatingEvent) error {
	for _, event := range events {
		encodedEvent, err := json.Marshal(event)
		if err != nil {
			return err
		}

		if err := producer.Produce(&kafka.Message{
			TopicPartition: kafka.TopicPartition{Topic: &topic, Partition: kafka.PartitionAny},
			Value:          encodedEvent,
		}, nil); err != nil {
			return err
		}
	}
	return nil
}
