package main

import (
	"fmt"
	"testing"

	"time"

	"github.com/rcrowley/go-metrics"
)

func TestMetrics(tester *testing.T) {
	// Threadsafe registration
	t := metrics.GetOrRegisterTimer("db.get.latency", nil)
	t.Time(func() {
		time.Sleep(time.Second * 2)
	})
	// t.Update(1)

	fmt.Println(time.Duration(t.Max()))
	fmt.Println(time.Duration(t.Min()))
}
