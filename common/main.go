package common

import (
	"flag"
	"os"
	"os/signal"
	"time"
)

func Main(producerBenchmark func(chan<- uint64) func(), consumerBenchmark func(chan<- uint64) func()) {
	flag.Parse()

	sig := make(chan os.Signal, 1)
	signal.Notify(sig, os.Interrupt)

	statsChannel := make(chan uint64, 1000)
	StartMonitoringGoroutine(statsChannel)

	var closeFn func()
	if Mode == "consumer" {
		closeFn = consumerBenchmark(statsChannel)
	} else if Mode == "producer" {
		closeFn = producerBenchmark(statsChannel)
	} else {
		panic("Unrecognized mode " + Mode)
	}

	select {
	case _ = <- sig:
		closeFn()
		time.Sleep(1 * time.Second)
		close(statsChannel)
		time.Sleep(1 * time.Second)
	}

	PrintRecords()
}
