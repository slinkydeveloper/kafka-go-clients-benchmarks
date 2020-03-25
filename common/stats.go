package common

import (
	"fmt"
	"time"
)

type Record struct {
	Time       time.Time
	Records    uint64
	TotalBytes uint64
}

var Records []Record

func StartMonitoringGoroutine(in <-chan uint64) {
	ticker := time.Tick(1 * time.Second)
	actual := Record{
		Time:       time.Now(),
		Records:    0,
		TotalBytes: 0,
	}
	go func() {
		for {
			select {
				case b, ok := <- in:
					actual.TotalBytes += b
					actual.Records += 1
					if !ok {
						return
					}
				case newTime := <- ticker:
					Records = append(Records, actual)
					actual = Record{
						Time:       newTime,
						Records:    0,
						TotalBytes: 0,
					}
			}
		}
	}()
}

func PrintRecords() {
	for _, r := range Records {
		fmt.Printf("%d,%d,%d\n", r.Time.Unix(), r.Records, r.TotalBytes)
	}
}
