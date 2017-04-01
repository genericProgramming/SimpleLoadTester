package main

import (
	"testing"

	"time"

	"github.com/genericProgramming/simpleLoadTester/components"
	. "github.com/smartystreets/goconvey/convey"
)

func TestMetrics(t *testing.T) {
	Convey("Hammer that service", t, func() {
		config := &components.Config{
			NumberOfBadRequestsPerTimeWindow: 10,
			Threshold: components.ResponseThreshold{
				Level:          .9,
				ResponseTimeMs: 100 * time.Millisecond,
			},
			URL:    "http://localhost:8081/hello/world",
			Window: 10 * time.Second,
		}
		monitor := NewStairCaseMonitor(config)
		RunMonitor(monitor)
		<-time.After(100 * time.Second)
		t.Log("We're done team!")
	})
}
