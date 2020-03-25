package main

import (
	"context"
	"time"

	"github.com/confluentinc/confluent-kafka-go/kafka"

	"github.com/slinkydeveloper/kafka-go-clients-benchmarks/common"
)

func ProducerBenchmark(stats chan<- uint64) func() {
	p, err := kafka.NewProducer(&kafka.ConfigMap{
		"bootstrap.servers": common.Broker,
		"partitioner":       "random",
	})
	if err != nil {
		panic(err)
	}

	ctx, cancel := context.WithCancel(context.Background())

	// Goroutine for input
	go func() {
		// One random per goroutine to skip all locks on rand
		myRand := common.NewRand()

		for {
			select {
			case <-ctx.Done():
				p.Close()
				return
			default:
				p.ProduceChannel() <- &kafka.Message{
					TopicPartition: kafka.TopicPartition{Topic: &common.Topic},
					Key: common.GenKey(myRand),
					Value: common.GenPayload(myRand),
				}
			}
		}
	}()

	go func() {
		for e := range p.Events() {
			switch ev := e.(type) {
			case *kafka.Message:
				m := ev
				if m.TopicPartition.Error == nil {
					stats <- uint64(len(m.Value))
				}
			}
		}
	}()

	return cancel
}

func ConsumerBenchmark(stats chan<- uint64) func() {
	c, err := kafka.NewConsumer(&kafka.ConfigMap{
		"bootstrap.servers":  common.Broker,
		"group.id":           "myGroup",
		"enable.auto.commit": "true",
		"auto.offset.reset":  "earliest",
	})
	if err != nil {
		panic(err)
	}

	err = c.Subscribe(common.Topic, nil)
	if err != nil {
		panic(err)
	}

	ctx, cancel := context.WithCancel(context.TODO())

	go func() {
		for {
			select {
			case <-ctx.Done():
				err = c.Close()
				if err != nil {
					panic(err)
				}
				return
			default:
				msg, err := c.ReadMessage(100 * time.Millisecond)
				if err == nil {
					stats <- uint64(len(msg.Value))
				}
			}
		}
	}()

	return cancel
}

func main() {
	common.Main(ProducerBenchmark, ConsumerBenchmark)
}
