package main

import (
	"context"

	"github.com/Shopify/sarama"
	"github.com/slinkydeveloper/kafka-go-clients-benchmarks/common"
)

func ProducerBenchmark(stats chan<- uint64) func() {
	conf := sarama.NewConfig()

	conf.Producer.Partitioner = sarama.NewRandomPartitioner
	conf.Producer.Retry.Max = 0
	conf.Producer.Return.Successes = true
	conf.Producer.Return.Errors = false

	asyncProducer, err := sarama.NewAsyncProducer([]string{common.Broker}, conf)
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
				_ = asyncProducer.Close()
				return

			default:
				asyncProducer.Input() <- &sarama.ProducerMessage{
					Topic: common.Topic,
					Key: sarama.ByteEncoder(common.GenKey(myRand)),
					Value: sarama.ByteEncoder(common.GenPayload(myRand)),
				}
			}
		}
	}()

	go func() {
		for s := range asyncProducer.Successes(){
			stats <- uint64(s.Value.Length())
		}
	}()

	return cancel
}

type consumer struct {
	stats chan<- uint64
}

func (c consumer) Setup(sarama.ConsumerGroupSession) error {
	return nil
}

func (c consumer) Cleanup(sarama.ConsumerGroupSession) error {
	return nil
}

func (c consumer) ConsumeClaim(sess sarama.ConsumerGroupSession, claim sarama.ConsumerGroupClaim) error {
	for msg := range claim.Messages() {
		c.stats <- uint64(len(msg.Value))
	}
	return nil
}

func ConsumerBenchmark(stats chan<- uint64) func() {
	conf := sarama.NewConfig()

	conf.Version = sarama.V2_0_0_0

	consumerGroup, err := sarama.NewConsumerGroup([]string{common.Broker}, "sarama-bench", conf)
	if err != nil {
		panic(err)
	}

	go func() {
		err := consumerGroup.Consume(context.TODO(), []string{common.Topic}, consumer{stats:stats})
		if err != nil {
			panic(err)
		}
	}()

	return func() {
		err := consumerGroup.Close()
		if err != nil {
			panic(err)
		}
	}
}

func main() {
	common.Main(ProducerBenchmark, ConsumerBenchmark)
}
