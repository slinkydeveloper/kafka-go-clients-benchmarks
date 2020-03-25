package common

import (
	"fmt"
	"time"
)

type Record struct {
	Time      time.Time
	MsgThpt   uint64
	BytesThpt float64
}

var Records []Record

func StartMonitoringGoroutine(in <-chan uint64) {
	ticker := time.Tick(1 * time.Second)
	var beginning time.Time
	msgTotal := uint64(0)
	bytesTotal := uint64(0)
	go func() {
		for {
			select {
			case b, ok := <-in:
				if !ok {
					return
				}
				bytesTotal += b
				msgTotal += 1
			case newTime := <-ticker:
				if !beginning.IsZero() {
					duration := newTime.Sub(beginning)
					Records = append(
						Records, Record{
							Time:      newTime,
							MsgThpt:   uint64(float64(msgTotal) / duration.Seconds()),
							BytesThpt: float64(bytesTotal) / duration.Seconds(),
						},
					)
				} else if msgTotal != 0 {
					beginning = newTime
				}
			}
		}
	}()
}

func PrintRecords() {
	for _, r := range Records {
		fmt.Printf("%d,%d,%f\n", r.Time.UnixNano(), r.MsgThpt, r.BytesThpt*1024*1024)
	}
}
