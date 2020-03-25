package main

import (
	"context"

	"github.com/segmentio/kafka-go"

	"github.com/slinkydeveloper/kafka-go-clients-benchmarks/common"
)

func ConsumerBenchmark(stats chan<- uint64) func() {
	r := kafka.NewReader(kafka.ReaderConfig{
		Brokers:   []string{common.Broker},
		GroupID:   "consumer-group-id",
		Topic:     common.Topic,
		QueueCapacity: int(2000),
	})

	ctx, cancel := context.WithCancel(context.TODO())

	go func() {
		for {
			select {
			case <-ctx.Done():
				err := r.Close()
				if err != nil {
					panic(err)
				}
				return
			default:
				msg, err := r.ReadMessage(ctx)
				if err == nil {
					stats <- uint64(len(msg.Value))
				}
			}
		}
	}()

	return cancel
}

func ProducerBenchmark(stats chan<- uint64) func() {
	w := kafka.NewWriter(kafka.WriterConfig{
		Brokers: []string{common.Broker},
		Topic:   common.Topic,
		Balancer: &kafka.LeastBytes{},
		Async:false,
	})

	ctx, cancel := context.WithCancel(context.Background())

	// Goroutine for input
	go func() {
		// One random per goroutine to skip all locks on rand
		myRand := common.NewRand()

		for {
			select {
			case <-ctx.Done():
				err := w.Close()
				if err != nil {
					panic(err)
				}
				return
			default:
				msg := kafka.Message{
					Key: common.GenKey(myRand),
					Value: common.GenPayload(myRand),
				}
				err := w.WriteMessages(ctx, msg)
				if err != nil {
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