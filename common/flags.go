package common

import "flag"

var Mode string
var Broker string
var Topic string
var KeySize int
var PayloadSize int

func init() {
	flag.StringVar(&Broker, "broker", "localhost:9092", "broker")
	flag.StringVar(&Topic, "topic", "test", "topic")
	flag.StringVar(&Mode, "mode", "consumer", "[consumer, producer]")
	flag.IntVar(&KeySize, "key-size", 1024, "key size")
	flag.IntVar(&PayloadSize, "payload-size", 1024, "payload size")
}
